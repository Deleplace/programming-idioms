package main

import (
	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/blobstore"
	"net/http"
)

func idiomPicture(w http.ResponseWriter, r *http.Request) error {
	// From https://developers.google.com/appengine/docs/go/blobstore/#Complete_Sample_App
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
	return nil
}
