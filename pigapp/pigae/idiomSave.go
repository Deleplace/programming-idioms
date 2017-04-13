package pigae

import (
	"fmt"
	"net/http"
	"strconv"

	. "github.com/Deleplace/programming-idioms/pig"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// Save an new idiom OR an existing idiom, depending on
// parameter "idiom_id"
func idiomSave(w http.ResponseWriter, r *http.Request) error {
	existingIDStr := r.FormValue("idiom_id")
	title := r.FormValue("idiom_title")
	username := r.FormValue("user_nickname")
	username = Truncate(username, 30)

	if !toggles["anonymousWrite"] {
		if username == "" {
			return PiError{"Username is mandatory. No anonymous edit.", http.StatusBadRequest}
		}
	}

	setNicknameCookie(w, username)

	if existingIDStr == "" {
		return newIdiomSave(w, r, username, title)
	}
	return existingIdiomSave(w, r, username, existingIDStr, title)
}

func newIdiomSave(w http.ResponseWriter, r *http.Request, username string, title string) error {
	if err := togglesMissing(w, r, "idiomCreation"); err != nil {
		return err
	}
	if err := parametersMissing(w, r, "impl_language", "impl_code"); err != nil {
		return err
	}

	c := appengine.NewContext(r)
	lead := r.FormValue("idiom_lead")
	keywords := r.FormValue("idiom_keywords")
	picture := r.FormValue("idiom_picture") /* TODO upload file ?! */
	language := NormLang(r.FormValue("impl_language"))
	imports := r.FormValue("impl_imports")
	code := r.FormValue("impl_code")
	comment := r.FormValue("impl_comment")
	attributionURL := r.FormValue("impl_attribution_url")
	demoURL := r.FormValue("impl_demo_url")
	docURL := r.FormValue("impl_doc_url")
	editSummary := fmt.Sprintf("Idiom creation by user [%v]", username)

	lead = Truncate(lead, 500)
	keywords = Truncate(keywords, 250)
	imports = Truncate(imports, 200)
	code = Truncate(NoCR(code), 500)
	comment = Truncate(comment, 500)
	attributionURL = Truncate(attributionURL, 250)
	demoURL = Truncate(demoURL, 250)
	docURL = Truncate(docURL, 250)

	log.Infof(c, "[%v] is creating new idiom [%v]", username, title)

	if !StringSliceContains(AllLanguages(), language) {
		return PiError{fmt.Sprintf("Sorry, [%v] is currently not a supported language. Supported languages are %v.", r.FormValue("impl_language"), AllNiceLangs), http.StatusBadRequest}
	}

	// TODO put that in a transaction!
	idiomID, err := dao.nextIdiomID(c)
	if err != nil {
		return err
	}
	implID, err := dao.nextImplID(c)
	if err != nil {
		return err
	}

	implementations := []Impl{
		Impl{
			Id:                     implID,
			OrigId:                 -1,
			Author:                 username,
			LastEditor:             username,
			LanguageName:           language,
			ImportsBlock:           imports,
			CodeBlock:              code,
			AuthorComment:          comment,
			OriginalAttributionURL: attributionURL,
			DemoURL:                demoURL,
			DocumentationURL:       docURL,
		},
	}
	idiom := &Idiom{
		Id:               idiomID,
		Title:            title,
		LeadParagraph:    lead,
		ExtraKeywords:    keywords,
		Picture:          picture, /* TODO upload file ?! */
		Author:           username,
		LastEditor:       username,
		LastEditedImplID: implID,
		EditSummary:      editSummary,
		Rating:           0,
		Implementations:  implementations,
	}
	/*
		Authenticated user name not needed here, as of 2015.
		Especially not for the Admin.

		if u := user.Current(c); u != nil {
			idiom.Author = u.String()
			idiom.LastEditor = u.String()
			implementations[0].Author = u.String()
			implementations[0].LastEditor = u.String()
		}
	*/

	_, err = dao.saveNewIdiom(c, idiom)
	if err != nil {
		return err
	}

	http.Redirect(w, r, NiceIdiomURL(idiom), http.StatusFound)
	return nil
}

func existingIdiomSave(w http.ResponseWriter, r *http.Request, username string, existingIDStr string, title string) error {
	if err := togglesMissing(w, r, "idiomEditing"); err != nil {
		return err
	}
	if err := parametersMissing(w, r, "idiom_version"); err != nil {
		return err
	}
	c := appengine.NewContext(r)
	log.Infof(c, "[%v] is updating statement of idiom %v", username, existingIDStr)

	idiomID := String2Int(existingIDStr)
	if idiomID == -1 {
		return PiError{existingIDStr + " is not a valid idiom id.", http.StatusBadRequest}
	}

	key, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + existingIDStr, http.StatusNotFound}
	}

	isAdmin := IsAdmin(r)
	if idiom.Protected && !isAdmin {
		return PiError{"Can't edit protected idiom " + existingIDStr, http.StatusUnauthorized}
	}
	if isAdmin {
		wasProtected := idiom.Protected
		idiom.Protected = r.FormValue("idiom_protected") != ""

		if wasProtected && !idiom.Protected {
			log.Infof(c, "[%v] unprotects idiom %v", username, existingIDStr)
		}
		if !wasProtected && idiom.Protected {
			log.Infof(c, "[%v] protects idiom %v", username, existingIDStr)
		}
	}

	if r.FormValue("idiom_version") != strconv.Itoa(idiom.Version) {
		return PiError{fmt.Sprintf("Idiom has been concurrently modified (editing version %v, current version is %v)", r.FormValue("idiom_version"), idiom.Version), http.StatusConflict}
	}

	idiom.LastEditor = username
	idiom.LastEditedImplID = 0
	idiom.Title = title
	idiom.LeadParagraph = r.FormValue("idiom_lead")
	idiom.ExtraKeywords = r.FormValue("idiom_keywords")
	idiom.EditSummary = r.FormValue("edit_summary")
	/* idiomPicture.go
	idiom.Picture, err = processUploadFile(r, "idiom_picture")
	if err != nil {
		return err
	}
	*/
	idiom.LeadParagraph = Truncate(idiom.LeadParagraph, 500)
	idiom.ExtraKeywords = Truncate(idiom.ExtraKeywords, 250)
	idiom.EditSummary = Truncate(idiom.EditSummary, 120)

	err = dao.saveExistingIdiom(c, key, idiom)
	if err != nil {
		return err
	}

	http.Redirect(w, r, NiceIdiomURL(idiom), http.StatusFound)
	return nil
}
