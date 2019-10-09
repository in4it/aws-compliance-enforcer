package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mitchellh/mapstructure"
)

var (
	// Info logger
	Info *log.Logger
	// Warning logger
	Warning *log.Logger
	// Error logger
	Error *log.Logger
)

//Handler bundles the different functions with the AWS session
type Handler struct {
	sess *session.Session
}

// initialize the log formatting
func initLog(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// initialize new handler with an AWS SDK session
func newHandler() *Handler {
	sess, _ := session.NewSession()
	return &Handler{sess: sess}
}

// main function
func main() {
	initLog(os.Stdout, os.Stdout, os.Stderr)

	handler := newHandler()

	lambda.Start(handler.start)
}

// handle snsEvent
func (h *Handler) start(context context.Context, snsEvent events.SNSEvent) (EventResponse, error) {

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		var snsMessage SNSEventMessage

		Info.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)

		err := json.Unmarshal([]byte(snsRecord.Message), &snsMessage)
		if err != nil {
			return EventResponse{Message: "Error"}, err
		}

		// download s3 file
		for _, s3ObjectKey := range snsMessage.S3ObjectKey {
			s3Contents, err := h.downloadGzippedS3Object(snsMessage.S3Bucket, s3ObjectKey)
			if err != nil {
				return EventResponse{Message: "Error"}, err
			}

			var cloudTrailEvent CloudTrailEvent
			err = json.Unmarshal(s3Contents, &cloudTrailEvent)
			if err != nil {
				return EventResponse{Message: "Error"}, err
			}

			if err = h.processCloudTrailEvent(cloudTrailEvent); err != nil {
				return EventResponse{Message: "Error"}, err
			}

		}
	}

	return EventResponse{Message: "Done"}, nil
}

// process a multi-record cloudtrail event
func (h *Handler) processCloudTrailEvent(cloudTrailEvent CloudTrailEvent) error {
	for _, record := range cloudTrailEvent.Records {
		err := h.processCloudTrailEventRecord(record)
		if err != nil {
			return err
		}
	}
	return nil
}

// process a single cloudtrail record
func (h *Handler) processCloudTrailEventRecord(record CloudTrailEventRecord) error {
	var err error
	switch record.EventName {
	case "AuthorizeSecurityGroupIngress":
		Info.Printf("Found AuthorizeSecurityGroupIngress event")
		var securityGroupEvent CloudTrailSecurityGroupEvent
		if err = mapstructure.Decode(record.RequestParameters, &securityGroupEvent); err != nil {
			return err
		}
		if err = h.securityGroupRuleHandler(securityGroupEvent); err != nil {
			return err
		}
	}
	return nil
}

// download an s3 object and uncompress it
func (h *Handler) downloadGzippedS3Object(bucket, key string) ([]byte, error) {
	contents := aws.NewWriteAtBuffer([]byte{})
	downloader := s3manager.NewDownloader(h.sess)
	_, err := downloader.Download(contents,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return nil, err
	}
	gr, err := gzip.NewReader(bytes.NewBuffer(contents.Bytes()))
	defer gr.Close()
	data, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// handle a security group event
func (h *Handler) securityGroupRuleHandler(requestParameters CloudTrailSecurityGroupEvent) error {
	Info.Printf("RequestParameter: %+v", requestParameters)

	for _, ipPermission := range requestParameters.IPPermissions.Items {
		for _, ipRange := range ipPermission.IPRanges.Items {
			if ipRange.CidrIP == "0.0.0.0/0" {
				if err := h.revokeSecurityGroupIngress(requestParameters.GroupID, ipRange.CidrIP, ipPermission.IPProtocol, ipPermission.FromPort, ipPermission.ToPort); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// revoke a security group rule
func (h *Handler) revokeSecurityGroupIngress(securityGroupID, ipRangeCidrIP, protocol string, fromPort, toPort int64) error {
	svc := ec2.New(h.sess)
	input := &ec2.RevokeSecurityGroupIngressInput{
		CidrIp:     aws.String(ipRangeCidrIP),
		GroupId:    aws.String(securityGroupID),
		IpProtocol: aws.String(protocol),
	}

	if protocol != "-1" {
		input.SetFromPort(fromPort)
		input.SetToPort(toPort)
	}

	Info.Printf("Revoking security group ingress with cidr %s, protocol %s (%s)", ipRangeCidrIP, protocol, securityGroupID)

	_, err := svc.RevokeSecurityGroupIngress(input)
	if err != nil {
		return err
	}

	return nil
}
