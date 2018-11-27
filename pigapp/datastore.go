package main

import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
)

var ds *datastore.Client

func init() {
	ctx := context.Background()
	projectID := "pi-go111" // TODO ?!?
	var err error
	ds, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
}
