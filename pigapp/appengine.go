package main

// This ID is needed to create a global Datastore client.
//
// It is also used for other concerns (e.g. Cloud Tasks) but
// could be dynamically derived from http.Request.Context().
//
const projectID = "pi-go111"
