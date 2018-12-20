package main

import (
	"os"

	"google.golang.org/appengine"
)

// This ID is needed to create a global Datastore client.
//
// It is also used for other concerns (e.g. Cloud Tasks) but
// could be dynamically derived from http.Request.Context().
//
var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

//const projectID = "pi-go111"

func isAppengineDevServer() bool {
	return appengine.IsDevAppServer()
}
