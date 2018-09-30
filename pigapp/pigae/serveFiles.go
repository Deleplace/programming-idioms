package pigae

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
)

func idiomPicture(c context.Context, w http.ResponseWriter, r *http.Request) error {
	// From https://developers.google.com/appengine/docs/go/blobstore/#Complete_Sample_App
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
	return nil
}
