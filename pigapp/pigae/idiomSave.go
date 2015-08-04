package pigae

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
	"appengine/user"
)

// Save an new idiom OR an existing idiom, depending on
// parameter "idiom_id"
func idiomSave(w http.ResponseWriter, r *http.Request) error {
	existingIDStr := r.FormValue("idiom_id")
	title := r.FormValue("idiom_title")
	username := r.FormValue("user_nickname")

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
	picture := r.FormValue("idiom_picture") /* TODO upload file ?! */
	language := normLang(r.FormValue("impl_language"))
	imports := r.FormValue("impl_imports")
	code := r.FormValue("impl_code")
	comment := r.FormValue("impl_comment")
	attributionURL := r.FormValue("impl_attribution_url")
	demoURL := r.FormValue("impl_demo_url")

	if !StringSliceContains(allLanguages(), language) {
		return PiError{fmt.Sprintf("Sorry, [%v] is currently not a supported language. Supported languages are %v.", r.FormValue("impl_language"), allNiceLangs), http.StatusBadRequest}
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
			Version:                1,
		},
	}
	idiom := &Idiom{
		Id:              idiomID,
		Title:           title,
		LeadParagraph:   lead,
		Picture:         picture, /* TODO upload file ?! */
		Author:          username,
		LastEditor:      username,
		Rating:          0,
		VersionDate:     time.Now(),
		Implementations: implementations,
	}
	if u := user.Current(c); u != nil {
		idiom.Author = u.String()
		idiom.LastEditor = u.String()
		implementations[0].Author = u.String()
		implementations[0].LastEditor = u.String()
	}

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

	idiomID := String2Int(existingIDStr)
	if idiomID == -1 {
		return PiError{existingIDStr + " is not a valid idiom id.", http.StatusBadRequest}
	}

	key, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + existingIDStr, http.StatusNotFound}
	}
	if r.FormValue("idiom_version") != strconv.Itoa(idiom.Version) {
		return PiError{fmt.Sprintf("Idiom has been concurrently modified (editing version %v, current version is %v)", r.FormValue("idiom_version"), idiom.Version), http.StatusConflict}
	}

	idiom.LastEditor = username
	idiom.Title = title
	idiom.LeadParagraph = r.FormValue("idiom_lead")
	/* idiomPicture.go
	idiom.Picture, err = processUploadFile(r, "idiom_picture")
	if err != nil {
		return err
	}
	*/

	err = dao.saveExistingIdiom(c, key, idiom)
	if err != nil {
		return err
	}

	http.Redirect(w, r, NiceIdiomURL(idiom), http.StatusFound)
	return nil
}
