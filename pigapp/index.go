package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"

	"context"
)

// ThemeVersion is the version of the current CSS-JS theme.
// It is the name of the folder containing the theme files.
const ThemeVersion = "default"

// ThemeDate is the prefix used for "revving" the static files and enable long-term HTTP cache.
// It MUST end with underscore _ (see app.yaml)
const ThemeDate = "20210311_"

var r = mux.NewRouter()

func init() {
	initEnv()
	initToggles()
	initRoutes()

	// We want the random results to be different even if we reboot the server. Thus, we use
	// the clock to seed the default generator.
	// See https://programming-idioms.org/idiom/70/use-clock-as-random-generator-seed/346/go
	rand.Seed(time.Now().UnixNano())
}

func initRoutes() {
	if !toggles["online"] {
		handle("/", makeWall("<i class=\"icon-wrench icon-2x\"> Under maintenance.</i>"))
		//r.HandleFunc("/", makeWall("<i class=\"icon-wrench icon-2x\"> Coming soon.</i>"))
	} else {
		//handle("/", makeWall("<i class=\"icon-wrench icon-2x\"> Coming soon.</i>"))
		handle("/", home)
		handle("/home", home)
		handle("/wall", makeWall("<i class=\"icon-wrench icon-2x\"> Coming soon.</i>"))
		handle("/about", about)
		handle("/idiom/{idiomId}", idiomDetail)
		handle("/idiom/{idiomId}/impl/{implId}", idiomDetail)
		handle("/idiom/{idiomId}/{idiomTitle}", idiomDetail)
		handle("/idiom/{idiomId}/diff/{v1}/{v2}", versionDiff)
		handle("/idiom/{idiomId}/impl/{implId}/diff/{v1}/{v2}", versionDiff)
		handle("/idiom/{idiomId}/{idiomTitle}/{implId}/{implLang}", idiomDetail)
		handle("/history/{idiomId}", idiomHistory)
		handle("/history/{idiomId}/impl/{implId}", implHistory)
		handle("/revert", revertIdiomVersion)
		handle("/admin-history-restore", restoreIdiomVersion)
		handle("/all-idioms", allIdioms)
		handle("/random-idiom/having/{havingLang}", randomIdiomHaving)
		handle("/random-idiom/not-having/{notHavingLang}", randomIdiomNotHaving)
		handle("/random-idiom", randomIdiom)
		handle("/search", searchRedirect)
		handle("/search/{q}", search)
		handle("/list-by-language/{langs}", listByLanguage)
		handle("/missing-fields/{lang}", missingList)
		handle("/idiom-picture", idiomPicture)
		handle("/rss-recently-created", rssRecentlyCreated)
		handle("/rss-recently-updated", rssRecentlyUpdated)
		handle("/rss-recent-changes", rssRecentChanges)
		handle("/my/{nickname}/{langs}", bookmarkableUserURL)
		handle("/my/{langs}", bookmarkableUserURL)
		// handle("/cheatsheet/{lang}", cheatsheet)
		handle("/cheatsheet/{langs:.+}", cheatsheetMulti)
		handleAjax("/typeahead-languages", typeaheadLanguages)
		handleAjax("/supported-languages", supportedLanguages)
		handleAjax("/ajax-other-implementations", ajaxOtherImplementations)
		handleAjax("/ajax-impl-flag/{idiomId}/{implId}", ajaxImplFlag)
		if toggles["writable"] {
			// When not in "read-only" mode
			handle("/idiom-save", idiomSave)
			handle("/idiom-edit/{idiomId}", idiomEdit)
			handle("/idiom-add-picture/{idiomId}", idiomAddPicture)
			handle("/idiom-save-picture", idiomSavePicture)
			handle("/impl-edit/{idiomId}/{implId}", implEdit)
			//handle("/fake-idiom-save", fakeIdiomSave)
			handle("/idiom-create", idiomCreate)
			handle("/impl-create/{idiomId}", implCreate)
			handle("/impl-create/{idiomId}/{lang}", implCreate)
			handle("/impl-save", implSave)
			// Ajax
			handleAjax("/ajax-idiom-vote", ajaxIdiomVote)
			handleAjax("/ajax-impl-vote", ajaxImplVote)
			handleAjax("/ajax-demo-site-suggest", ajaxDemoSiteSuggest)
			handleAjax("/ajax-user-message-box", userMessageBoxAjax)
			handleAjax("/ajax-dismiss-user-message", dismissUserMessage)
			handle("/about-block-project", ajaxAboutProject)
			handle("/about-block-all-idioms", ajaxAboutAllIdioms)
			handle("/about-block-language-coverage", ajaxAboutLanguageCoverage)
			handle("/about-block-rss", ajaxAboutRss)
			handle("/about-block-cheatsheets", ajaxAboutCheatsheets)
			handle("/about-block-see-also", ajaxAboutSeeAlso)
			handle("/about-block-contact", ajaxAboutContact)
			// Admin
			handle("/admin", admin)
			handle("/admin-data-export", adminExport)
			handle("/admin-data-import", adminImport)
			handle("/admin-resave-entities", adminResaveEntities)
			handle("/admin-flagged", adminListFlaggedContent)
			handleAjax("/admin-repair-history-versions", adminRepairHistoryVersions)
			handleAjax("/admin-data-import-ajax", adminImportAjax)
			handleAjax("/admin-reindex-ajax", adminReindexAjax)
			handleAjax("/admin-refresh-toggles-ajax", ajaxRefreshToggles)
			handleAjax("/admin-set-toggle-ajax", ajaxSetToggle)
			handleAjax("/admin-create-relation-ajax", ajaxCreateRelation)
			handleAjax("/admin-idiom-delete", idiomDelete)
			handleAjax("/admin-impl-delete", implDelete)
			handleAjax("/admin-send-message-for-user", sendMessageForUserAjax)
			handleAjax("/admin-flag-resolve", ajaxAdminFlagResolve)
			handleAjax("/admin-memcache-flush", ajaxAdminMemcacheFlush)
		}
		handleAjax("/api/idiom/{idiomId}", jsonIdiom)
		handleAjax("/api/idioms/all", jsonAllIdioms)
		handleAjax("/api/search/{q}", jsonSearch)
		r.PathPrefix("/using/").HandlerFunc(using)

		handle("/auth", handleAuth)
		handle("/_ah/login_required", handleAuth)
	}
	http.Handle("/", r)
}

// Request will fail if path parameters are missing
var neededPathVariables = map[string][]string{
	"/idiom/{idiomId}":                                  {"idiomId"},
	"/idiom/{idiomId}/impl/{implId}":                    {"idiomId"},
	"/idiom/{idiomId}/{idiomTitle}":                     {"idiomId"},
	"/idiom/{idiomId}/{idiomTitle}/{implId}/{implLang}": {"idiomId"},
	"/search/{q}":                                       {"q"},
	"/my/{nickname}/{langs}":                            {"nickname", "langs"},
	"/idiom-edit/{idiomId}":                             {"idiomId"},
	"/idiom-add-picture/{idiomId}":                      {"idiomId"},
	"/impl-edit/{idiomId}/{implId}":                     {"idiomId", "implId"},
	"/impl-create/{idiomId}":                            {"idiomId"},
	"/impl-create/{idiomId}/{lang}":                     {"idiomId"},
	"/cheatsheet/{lang}":                                {"lang"},
}

// Request will fail if it doesn't provide the required GET or POST parameters
var neededParameters = map[string][]string{
	"/typeahead-languages":         { /*todo*/ },
	"/idiom-save":                  {"idiom_title"},
	"/idiom-save-picture":          { /*todo*/ },
	"/impl-save":                   {"idiom_id", "impl_code"},
	"/revert":                      {"idiomId", "version"},
	"/ajax-idiom-vote":             {"idiomId", "choice"},
	"/ajax-impl-vote":              {"implId", "choice"},
	"/ajax-demo-site-suggest":      { /*todo*/ },
	"/ajax-dismiss-user-message":   {"key"},
	"/admin-data-export":           { /*todo*/ },
	"/admin-data-import":           { /*todo*/ },
	"/admin-data-import-ajax":      { /*todo*/ },
	"/admin-set-toggle-ajax":       {"toggle", "value"},
	"/admin-create-relation-ajax":  {"idiomAId", "idiomBId"},
	"/admin-idiom-delete":          {"idiomId"},
	"/admin-impl-delete":           {"idiomId", "implId"},
	"/admin-send-message-for-user": {"username", "message"},
	"/admin-flag-resolve":          {"flagkey"},
	"/api/idiom":                   {"idiomId"},
}

// Request will fail if corresponding toggle is off
var neededToggles = map[string][]string{
	"/home":                         {"online"},
	"/search":                       {"searchable"},
	"/search/{q}":                   {"searchable"},
	"/idiom-save":                   {"writable"},
	"/idiom-edit/{idiomId}":         {"writable", "writable", "idiomEditing"},
	"/idiom-add-picture/{idiomId}":  {"writable", "idiomEditing"},
	"/idiom-save-picture":           {"writable", "idiomEditing"},
	"/impl-edit/{idiomId}/{implId}": {"writable", "implEditing"},
	"/idiom-create":                 {"writable"},
	"/impl-create/{idiomId}":        {"writable", "implAddition"},
	"/impl-create/{idiomId}/{lang}": {"writable", "implAddition"},
	"/impl-save":                    {"writable"},
	"/ajax-idiom-vote":              {"writable"},
	"/ajax-impl-vote":               {"writable"},
	"/admin":                        {"administrable"},
	"/admin-data-export":            {"administrable"},
	"/admin-data-import":            {"administrable"},
	"/admin-data-import-ajax":       {"administrable"},
	"/admin-set-toggle-ajax":        {"administrable"},
	"/admin-create-relation-ajax":   {"administrable"},
	"/admin-idiom-delete":           {"administrable"},
	"/admin-impl-delete":            {"administrable"},
}

type standardHandler func(w http.ResponseWriter, r *http.Request)
type betterHandler func(w http.ResponseWriter, r *http.Request) error

// Wrap HandleFunc with
// - error handling
// - mandatory path variables check
// - mandatory parameters check
// - toggles check
func handle(path string, h betterHandler) {
	r.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			if isSpam(w, r) {
				return
			}

			defer func() {
				if msg := recover(); msg != nil {
					errorPage(w, r, PiErrorf(http.StatusInternalServerError, "%v", msg))
					return
				}
			}()
			if configTime == "0" {
				ctx := r.Context()
				_ = refreshToggles(ctx)
				// If it fails... well, ignore for now and continue with non-fresh toggles.
			}

			if err := muxVarsMissing(w, r, neededPathVariables[path]...); err != nil {
				errorPage(w, r, err)
				return
			}
			if err := togglesMissing(w, r, neededToggles[path]...); err != nil {
				errorPage(w, r, err)
				return
			}
			if err := parametersMissing(w, r, neededParameters[path]...); err != nil {
				errorPage(w, r, err)
				return
			}
			err := h(w, r)
			if err != nil {
				errorPage(w, r, err)
			}
		})
}

func handleAjax(path string, h betterHandler) {
	r.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			if isSpam(w, r) {
				return
			}

			defer func() {
				if msg := recover(); msg != nil {
					errorJSON(w, r, PiErrorf(http.StatusInternalServerError, "%v", msg))
					return
				}
			}()
			if configTime == "0" {
				ctx := r.Context()
				_ = refreshToggles(ctx)
				// If it fails... well, ignore for now and continue with non-fresh toggles.
			}

			if err := muxVarsMissing(w, r, neededPathVariables[path]...); err != nil {
				errorJSON(w, r, err)
				return
			}
			if err := togglesMissing(w, r, neededToggles[path]...); err != nil {
				errorJSON(w, r, err)
				return
			}
			if err := parametersMissing(w, r, neededParameters[path]...); err != nil {
				errorJSON(w, r, err)
				return
			}

			// TODO: maybe we never want any other Content-Type so...?
			// w.Header().Set("Content-Type", "application/json")
			err := h(w, r)
			if err != nil {
				errorJSON(w, r, err)
			}
		})
}

var datastoreDao = GaeDatastoreAccessor{}
var memcachedDao = MemcacheDatastoreAccessor{datastoreDao}
var dao = memcachedDao

var daoVotes = GaeVotesAccessor{}

func parametersMissing(w http.ResponseWriter, r *http.Request, params ...string) error {
	missing := []string{}
	for _, param := range params {
		if r.FormValue(param) == "" {
			missing = append(missing, param)
		}
	}
	if len(missing) > 0 {
		return PiErrorf(http.StatusBadRequest, "Missing parameters : %s", missing)
	}
	return nil
}

// Looks in gorilla mux populated variables
func muxVarsMissing(w http.ResponseWriter, r *http.Request, params ...string) error {
	missing := []string{}
	muxvars := mux.Vars(r)
	for _, param := range params {
		if muxvars[param] == "" {
			missing = append(missing, param)
		}
	}
	if len(missing) > 0 {
		return PiErrorf(http.StatusBadRequest, "Missing parameters : %s", missing)
	}
	return nil
}

func validateURLFormat(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return fmt.Errorf("Requires an absolute URL")
	}
	return nil
}

func validateURLFormatOrEmpty(urlStr string) error {
	if urlStr == "" {
		return nil
	}
	return validateURLFormat(urlStr)
}

/*
func logIf(err error, logfunc func(format string, args ...interface{}), when string) {
	if err != nil {
		logfunc("Problem on %v: %v", when, err.Error())
	}
}
*/

func logIf(err error, logfunc func(ctx context.Context, format string, args ...interface{}), ctx context.Context, when string) {
	if err != nil {
		logfunc(ctx, "Problem on %v: %v", when, err.Error())
	}
}
