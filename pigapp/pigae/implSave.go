package pigae

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func implSave(w http.ResponseWriter, r *http.Request) error {
	idiomIDStr := r.FormValue("idiom_id")
	existingIDStr := r.FormValue("impl_id")
	username := r.FormValue("user_nickname")
	username = Truncate(username, 30)

	if !toggles["anonymousWrite"] {
		if username == "" {
			return PiError{"Username is mandatory. No anonymous edit.", http.StatusBadRequest}
		}
	}

	setNicknameCookie(w, username)

	if existingIDStr == "" {
		return newImplSave(w, r, username, idiomIDStr)
	}
	return existingImplSave(w, r, username, idiomIDStr, existingIDStr)
}

func newImplSave(w http.ResponseWriter, r *http.Request, username string, idiomIDStr string) error {
	if err := togglesMissing(w, r, "implAddition"); err != nil {
		return err
	}
	if err := parametersMissing(w, r, "impl_language"); err != nil {
		return err
	}

	c := appengine.NewContext(r)
	language := normLang(r.FormValue("impl_language"))
	imports := r.FormValue("impl_imports")
	code := r.FormValue("impl_code")
	comment := r.FormValue("impl_comment")
	attributionURL := r.FormValue("impl_attribution_url")
	demoURL := r.FormValue("impl_demo_url")
	docURL := r.FormValue("impl_doc_url")
	editSummary := fmt.Sprintf("New %v implementation by user [%v]", language, username)

	imports = Truncate(imports, 200)
	code = Truncate(code, 500)
	comment = Truncate(comment, 500)
	attributionURL = Truncate(attributionURL, 250)
	demoURL = Truncate(demoURL, 250)
	docURL = Truncate(docURL, 250)

	log.Infof(c, "[%v] is creating new %v impl for idiom %v", username, language, idiomIDStr)

	if !StringSliceContains(allLanguages(), language) {
		return PiError{fmt.Sprintf("Sorry, [%v] is currently not a supported language. Supported languages are %v.", r.FormValue("impl_language"), allNiceLangs), http.StatusBadRequest}
	}

	idiomID := String2Int(idiomIDStr)
	if idiomID == -1 {
		return PiError{idiomIDStr + " is not a valid idiom id.", http.StatusBadRequest}
	}

	key, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
	}

	if err := validateURLFormatOrEmpty(attributionURL); err != nil {
		return PiError{"Can't accept URL [" + attributionURL + "]", http.StatusBadRequest}
	}

	if err := validateURLFormatOrEmpty(demoURL); err != nil {
		return PiError{"Can't accept URL [" + demoURL + "]", http.StatusBadRequest}
	}

	implID, err := dao.nextImplID(c)
	if err != nil {
		return err
	}
	now := time.Now()
	newImpl := Impl{
		Id:                     implID,
		OrigId:                 implID,
		Author:                 username,
		CreationDate:           now,
		LastEditor:             username,
		LanguageName:           language,
		ImportsBlock:           imports,
		CodeBlock:              code,
		AuthorComment:          comment,
		OriginalAttributionURL: attributionURL,
		DemoURL:                demoURL,
		DocumentationURL:       docURL,
		Version:                1,
		VersionDate:            now,
	}

	if IsAdmin(r) {
		// 2016-10: only Admin may set an impl picture
		newImpl.PictureURL = r.FormValue("impl_picture_url")
	}

	idiom.Implementations = append(idiom.Implementations, newImpl)
	idiom.EditSummary = editSummary
	idiom.LastEditedImplID = implID

	err = dao.saveExistingIdiom(c, key, idiom)
	if err != nil {
		return err
	}

	http.Redirect(w, r, NiceImplURL(idiom, implID, language), http.StatusFound)
	return nil
}

func existingImplSave(w http.ResponseWriter, r *http.Request, username string, idiomIDStr string, existingImplIDStr string) error {
	if err := togglesMissing(w, r, "implEditing"); err != nil {
		return err
	}
	if err := parametersMissing(w, r, "impl_version"); err != nil {
		return err
	}

	c := appengine.NewContext(r)
	imports := r.FormValue("impl_imports")
	code := r.FormValue("impl_code")
	comment := r.FormValue("impl_comment")
	attributionURL := r.FormValue("impl_attribution_url")
	demoURL := r.FormValue("impl_demo_url")
	docURL := r.FormValue("impl_doc_url")

	imports = Truncate(imports, 200)
	code = Truncate(code, 500)
	comment = Truncate(comment, 500)
	attributionURL = Truncate(attributionURL, 250)
	demoURL = Truncate(demoURL, 250)
	docURL = Truncate(docURL, 250)

	log.Infof(c, "[%v] is updating impl %v of idiom %v", username, existingImplIDStr, idiomIDStr)

	idiomID := String2Int(idiomIDStr)
	if idiomID == -1 {
		return PiError{idiomIDStr + " is not a valid idiom id.", http.StatusBadRequest}
	}

	key, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find implementation " + existingImplIDStr + " for idiom " + idiomIDStr, http.StatusNotFound}
	}

	implID := String2Int(existingImplIDStr)
	if implID == -1 {
		return PiError{existingImplIDStr + " is not a valid implementation id.", http.StatusBadRequest}
	}

	_, impl, _ := idiom.FindImplInIdiom(implID)

	isAdmin := IsAdmin(r)
	if idiom.Protected && !isAdmin {
		return PiError{"Can't edit protected idiom " + idiomIDStr, http.StatusUnauthorized}
	}
	if impl.Protected && !isAdmin {
		return PiError{"Can't edit protected impl " + existingImplIDStr, http.StatusUnauthorized}
	}
	if isAdmin {
		wasProtected := impl.Protected
		impl.Protected = r.FormValue("impl_protected") != ""

		if wasProtected && !impl.Protected {
			log.Infof(c, "[%v] unprotects impl %v of idiom %v", username, existingImplIDStr, idiomIDStr)
		}
		if !wasProtected && impl.Protected {
			log.Infof(c, "[%v] protects impl %v of idiom %v", username, existingImplIDStr, idiomIDStr)
		}
	}

	if r.FormValue("impl_version") != strconv.Itoa(impl.Version) {
		return PiError{fmt.Sprintf("Implementation has been concurrently modified (editing version %v, current version is %v)", r.FormValue("impl_version"), impl.Version), http.StatusConflict}
	}

	if err := validateURLFormatOrEmpty(attributionURL); err != nil {
		return PiError{"Can't accept URL [" + attributionURL + "]", http.StatusBadRequest}
	}

	if err := validateURLFormatOrEmpty(demoURL); err != nil {
		return PiError{"Can't accept URL [" + demoURL + "]", http.StatusBadRequest}
	}

	idiom.EditSummary = "[" + impl.LanguageName + "] " + r.FormValue("edit_summary")
	idiom.LastEditedImplID = implID
	impl.ImportsBlock = imports
	impl.CodeBlock = code
	impl.AuthorComment = comment
	impl.LastEditor = username
	impl.OriginalAttributionURL = attributionURL
	impl.DemoURL = demoURL
	impl.DocumentationURL = docURL
	impl.Version = impl.Version + 1
	impl.VersionDate = time.Now()

	if isAdmin {
		// 2016-10: only Admin may set an impl picture
		impl.PictureURL = r.FormValue("impl_picture_url")
	}

	err = dao.saveExistingIdiom(c, key, idiom)
	if err != nil {
		return err
	}

	http.Redirect(w, r, NiceImplURL(idiom, implID, impl.LanguageName), http.StatusFound)
	return nil
}
