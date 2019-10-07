package main

import "time"

type EventResponse struct {
	Message string `json:"message"`
}

type SNSEvent struct {
	S3Bucket    string   `json:"s3Bucket"`
	S3ObjectKey []string `json:"s3ObjectKey"`
}

type Event struct {
	EventRecord []EventRecord `json:"records" binding:"required"`
}

type EventRecord struct {
	EventVersion string `yaml:"eventVersion"`
	UserIdentity struct {
		Type        string `yaml:"type"`
		UserName    string `yaml:"userName"`
		PrincipalID string `yaml:"principalId"`
		Arn         string `yaml:"arn"`
		AccountID   string `yaml:"accountId"`
		AccessKeyID string `yaml:"accessKeyId"`
	} `yaml:"userIdentity"`
	EventTime         time.Time `yaml:"eventTime"`
	EventSource       string    `yaml:"eventSource"`
	EventName         string    `yaml:"eventName"`
	AwsRegion         string    `yaml:"awsRegion"`
	SourceIPAddress   string    `yaml:"sourceIPAddress"`
	UserAgent         string    `yaml:"userAgent"`
	RequestParameters struct {
		NextToken string `yaml:"nextToken"`
	} `yaml:"requestParameters"`
	ResponseElements   interface{} `yaml:"responseElements"`
	RequestID          string      `yaml:"requestID"`
	EventID            string      `yaml:"eventID"`
	EventType          string      `yaml:"eventType"`
	RecipientAccountID string      `yaml:"recipientAccountId"`
}
