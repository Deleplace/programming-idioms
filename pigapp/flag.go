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

// FlaggedContent is a user report of an inappropriate content.
// It contains an Idiom ID, but not the Idiom itself.
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

// AdminFlaggedFacade is the Facade for the Admin Flagged Contents page.
type AdminFlaggedFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Flagged     []FlaggedContentFacade
}

// FlaggedContentFacade is the Facade for 1 line of the Flagged Contents table.
type FlaggedContentFacade struct {
	FlaggedContent
	Idiom *Idiom
	Impl  *Impl
}

func adminListFlaggedContent(w http.ResponseWriter, r *http.Request) error {
	// reports are raw user input: type FlaggedContent.
	// table contains decorated flagged contents: type FlaggedContentFacade.
	ctx := r.Context()

	var reports []FlaggedContent
	_, err := datastore.NewQuery("FlaggedContent").
		Order("-Timestamp").
		Limit(100).
		GetAll(ctx, &reports)
	if err != nil {
		return err
	}

	var table []FlaggedContentFacade
	for _, report := range reports {
		line := FlaggedContentFacade{
			FlaggedContent: report,
		}

		_, idiom, err := dao.getIdiom(ctx, report.IdiomID)
		if err == nil {
			line.Idiom = idiom
			_, impl, found := idiom.FindImplInIdiom(report.ImplID)
			if found {
				line.Impl = impl
			} else {
				log.Errorf(ctx, "impl %d not found in idiom %d", report.ImplID, idiom.Id)
			}
		} else {
			log.Errorf(ctx, "loading idiom for impl %d: %v", report.ImplID, err)
		}

		table = append(table, line)
	}

	data := &AdminFlaggedFacade{
		PageMeta: PageMeta{
			PageTitle: "Flagged Contents",
			ExtraCss:  []string{hostPrefix() + themeDirectory() + "/css/admin.css"},
			ExtraJs:   []string{hostPrefix() + themeDirectory() + "/js/programming-idioms-admin.js"},
			Toggles:   toggles,
		},
		Flagged: table,
	}

	return templates.ExecuteTemplate(w, "page-admin-list-flagged", data)
}
