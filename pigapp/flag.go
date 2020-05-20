package main

import (
	"net/http"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"
	"github.com/gorilla/mux"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Let visitors "flag" an inappropriate content, i.e. notify admins.

type FlaggedContent struct {
	IdiomID int

	// IdiomVersion is needed because the flagged content existed in a specific version of the
	// Idiom, which may not be the current version anymore.
	IdiomVersion int

	// ImplID is optional: if the flagged content is in a specific impl.
	ImplID int

	// Timestamp of the report.
	Timestamp time.Time

	// Rationale provided by the user: why are they flagging this content.
	Rationale string

	// UserNickname is optional. Anonymous reports are fine.
	UserNickname string
}

func ajaxImplFlag(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ctx := r.Context()

	idiomIDStr := vars["idiomId"]
	implIDStr := vars["implId"]
	idiomVersionStr := r.FormValue("idiomVersion")
	rationale := r.FormValue("rationale")
	nickname := lookForNickname(r)

	log.Infof(ctx, "Flagging idiom %v version %v impl %v because: %q", idiomIDStr, idiomVersionStr, implIDStr, rationale)

	flag := FlaggedContent{
		IdiomID:      String2Int(idiomIDStr),
		IdiomVersion: String2Int(idiomVersionStr),
		ImplID:       String2Int(implIDStr),
		Timestamp:    time.Now(),
		Rationale:    rationale,
		UserNickname: nickname,
	}
	ikey := datastore.NewIncompleteKey(ctx, "FlaggedContent", nil)
	key, err := datastore.Put(ctx, ikey, &flag)
	if err != nil {
		log.Errorf(ctx, "saving FlagContent: %v", err)
		return PiError{"Could not save flagged content data", http.StatusInternalServerError}
	}
	log.Infof(ctx, "Saved content flag %s", key.Encode())
	return nil
}
