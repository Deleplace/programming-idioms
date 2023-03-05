package main

import (
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	. "github.com/Deleplace/programming-idioms/idioms"
	"github.com/gorilla/mux"
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

	// Resolved by the admin.
	Resolved bool

	// ResolveDate is when the admin marked the flag as "Resolved".
	ResolveDate time.Time
}

func ajaxImplFlag(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ctx := r.Context()

	idiomIDStr := vars["idiomId"]
	implIDStr := vars["implId"]
	idiomVersionStr := r.FormValue("idiomVersion")
	rationale := r.FormValue("rationale")
	nickname := lookForNickname(r)

	logf(ctx, "Flagging idiom %v version %v impl %v because: %q", idiomIDStr, idiomVersionStr, implIDStr, rationale)

	flag := FlaggedContent{
		IdiomID:      String2Int(idiomIDStr),
		IdiomVersion: String2Int(idiomVersionStr),
		ImplID:       String2Int(implIDStr),
		Timestamp:    time.Now(),
		Rationale:    rationale,
		UserNickname: nickname,
	}
	ikey := datastore.IncompleteKey("FlaggedContent", nil)
	key, err := dao.dsClient.Put(ctx, ikey, &flag)
	if err != nil {
		errf(ctx, "saving FlaggedContent: %v", err)
		return PiErrorf(http.StatusInternalServerError, "Could not save flagged content data")
	}
	logf(ctx, "Saved content flag %s", key.Encode())
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
	Key          *datastore.Key
	IdiomHistory *IdiomHistory
	Impl         *Impl
}

func adminListFlaggedContent(w http.ResponseWriter, r *http.Request) error {
	// reports are raw user input: type FlaggedContent.
	// table contains decorated flagged contents: type FlaggedContentFacade.
	ctx := r.Context()

	var reports []FlaggedContent
	q := datastore.NewQuery("FlaggedContent").
		Order("-Timestamp").
		Limit(100)
	keys, err := dao.dsClient.GetAll(ctx, q, &reports)
	if err != nil {
		return err
	}

	var table []FlaggedContentFacade
	for i, report := range reports {
		line := FlaggedContentFacade{
			FlaggedContent: report,
			Key:            keys[i],
		}

		_, idiomHistory, err := dao.getIdiomHistory(ctx, report.IdiomID, report.IdiomVersion)
		if err == nil {
			line.IdiomHistory = idiomHistory
			_, impl, found := idiomHistory.FindImplInIdiom(report.ImplID)
			if found {
				line.Impl = impl
			} else {
				errf(ctx, "impl %d not found in idiom %d v%d", report.ImplID, idiomHistory.Id, idiomHistory.Version)
			}
		} else {
			errf(ctx, "loading idiom %d history v%d for impl %d: %v", report.IdiomID, report.IdiomVersion, report.ImplID, err)
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

func ajaxAdminFlagResolve(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return PiErrorf(http.StatusBadRequest, "POST only")
	}
	flagKeyStr := r.FormValue("flagkey")
	flagKey, err := datastore.DecodeKey(flagKeyStr)
	if err != nil {
		return PiErrorf(http.StatusBadRequest, "Could not decode key %q", flagKeyStr)
	}

	ctx := r.Context()
	var flag FlaggedContent
	err = dao.dsClient.Get(ctx, flagKey, &flag)
	if err == datastore.ErrNoSuchEntity {
		return PiErrorf(http.StatusNotFound, "Flagged contents %q no longer exists", flagKeyStr)
	}
	if err != nil {
		errf(ctx, "retrieving FlaggegContent: %v", err)
		return PiErrorf(http.StatusInternalServerError, "Could not retrieve Flagged Contents entry :(")
	}
	flag.Resolved = true
	flag.ResolveDate = time.Now()
	_, err = dao.dsClient.Put(ctx, flagKey, &flag)
	if err != nil {
		errf(ctx, "saving FlaggedContent: %v", err)
		return PiErrorf(http.StatusInternalServerError, "Could not save flagged content data")
	}
	logf(ctx, "Saved content flag %s", flagKey.Encode())

	return nil
}
