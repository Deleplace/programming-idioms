package main

import (
	"context"
	"log"
)

func init() {
	// Don't prefix each message with (redundant) datetime
	log.SetFlags(0)
}

func logf(c context.Context, level, msg string, args ...interface{}) {
	// "cloud.google.com/go/logging" doesn't work yet?
	// Actually it logs under "Google Project", not under "GAE Application"

	// Sets your Google Cloud Platform project ID.
	// projectID := appengine.AppID(c)

	// Creates a client.
	// client, err := logging.NewClient(c, projectID)
	// if err != nil {
	// 	log.Fatalf("Failed to create stackdriver logging client: %v", err)
	// }
	// defer client.Close()

	// Sets the name of the log to write to.
	// logName := "my-log"
	// logger := client.Logger(logName).StandardLogger(logging.Info)
	// logger.Printf(msg, args...)

	log.Printf(level+" "+msg, args...)
}

func debugf(c context.Context, msg string, args ...interface{}) {
	logf(c, "DEBUG", msg, args...)
}

func infof(c context.Context, msg string, args ...interface{}) {
	logf(c, "INFO", msg, args...)
}

func warningf(c context.Context, msg string, args ...interface{}) {
	logf(c, "WARN", msg, args...)
}

func errorf(c context.Context, msg string, args ...interface{}) {
	logf(c, "ERROR", msg, args...)
}
