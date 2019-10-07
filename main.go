package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// Info logger
	Info *log.Logger
	// Warning logger
	Warning *log.Logger
	// Error logger
	Error *log.Logger
)

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

func main() {
	initLog(os.Stdout, os.Stdout, os.Stderr)

	lambda.Start(handler)
}

func handler(context context.Context, request SNSEvent) (EventResponse, error) {
	Info.Printf("Input: %+v", request)
	return EventResponse{Message: "Done"}, nil
}
