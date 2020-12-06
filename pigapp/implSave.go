package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"google.golang.org/appengine/log"
)

func implSave(w http.ResponseWriter, r *http.Request) error {
	idiomIDStr := r.FormValue("idiom_id")
	existingIDStr := r.FormValue("impl_id")
	username := r.FormValue("user_nickname")
	username = Truncate(username, 30)

	if !toggles["anonymousWrite"] {
		if username == "" {
			return PiErrorf(http.StatusBadRequest, "Username is mandatory. No anonymous edit.")
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

	ctx := r.Context()
	language := NormLang(r.FormValue("impl_language"))
	imports := r.FormValue("impl_imports")
	code := r.FormValue("impl_code")
	comment := r.FormValue("impl_comment")
	attributionURL := r.FormValue("impl_attribution_url")
	demoURL := r.FormValue("impl_demo_url")
	docURL := r.FormValue("impl_doc_url")
	editSummary := fmt.Sprintf("New %s implementation by user [%s]", PrintNiceLang(language), username)

	trim := strings.TrimSpace
	imports = trim(Truncate(imports, 200))
	code = TruncateBytes(NoCR(code), 500)
	comment = trim(TruncateBytes(comment, 500))
	attributionURL = trim(Truncate(attributionURL, 250))
	demoURL = trim(Truncate(demoURL, 250))
	docURL = trim(Truncate(docURL, 250))

	log.Infof(ctx, "[%s] is creating new %s impl for idiom %v", username, PrintNiceLang(language), idiomIDStr)

	if !StringSliceContains(AllLanguages(), language) {
		return PiErrorf(http.StatusBadRequest, "Sorry, [%v] is currently not a supported language. Supported languages are %v.", r.FormValue("impl_language"), AllNiceLangs)
	}

	idiomID := String2Int(idiomIDStr)
	if idiomID == -1 {
		return PiErrorf(http.StatusBadRequest, "%q is not a valid idiom id.", idiomIDStr)
	}

	key, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "Could not find idiom %q", idiomIDStr)
	}

	if err := validateURLFormatOrEmpty(attributionURL); err != nil {
		return PiErrorf(http.StatusBadRequest, "Can't accept URL [%s]", attributionURL)
	}

	if err := validateURLFormatOrEmpty(demoURL); err != nil {
		return PiErrorf(http.StatusBadRequest, "Can't accept URL [%s]", demoURL)
	}

	implID, err := dao.nextImplID(ctx)
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

	err = dao.saveExistingIdiom(ctx, key, idiom)
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

	ctx := r.Context()
	imports := r.FormValue("impl_imports")
	code := r.FormValue("impl_code")
	comment := r.FormValue("impl_comment")
	attributionURL := r.FormValue("impl_attribution_url")
	demoURL := r.FormValue("impl_demo_url")
	docURL := r.FormValue("impl_doc_url")

	trim := strings.TrimSpace
	imports = trim(Truncate(imports, 200))
	code = TruncateBytes(NoCR(code), 500)
	comment = trim(TruncateBytes(comment, 500))
	attributionURL = trim(Truncate(attributionURL, 250))
	demoURL = trim(Truncate(demoURL, 250))
	docURL = trim(Truncate(docURL, 250))

	log.Infof(ctx, "[%s] is updating impl %s of idiom %s", username, existingImplIDStr, idiomIDStr)

	idiomID := String2Int(idiomIDStr)
	if idiomID == -1 {
		return PiErrorf(http.StatusBadRequest, "%q is not a valid idiom id.", idiomIDStr)
	}

	key, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "Could not find implementation %q for idiom %q", existingImplIDStr, idiomIDStr)
	}

	implID := String2Int(existingImplIDStr)
	if implID == -1 {
		return PiErrorf(http.StatusBadRequest, "%q is not a valid implementation id.", existingImplIDStr)
	}

	_, impl, _ := idiom.FindImplInIdiom(implID)

	isAdmin := IsAdmin(r)
	if idiom.Protected && !isAdmin {
		return PiErrorf(http.StatusUnauthorized, "Can't edit protected idiom %q", idiomIDStr)
	}
	if impl.Protected && !isAdmin {
		return PiErrorf(http.StatusUnauthorized, "Can't edit protected impl %q", existingImplIDStr)
	}
	if isAdmin {
		wasProtected := impl.Protected
		impl.Protected = r.FormValue("impl_protected") != ""

		if wasProtected && !impl.Protected {
			log.Infof(ctx, "[%v] unprotects impl %v of idiom %v", username, existingImplIDStr, idiomIDStr)
		}
		if !wasProtected && impl.Protected {
			log.Infof(ctx, "[%v] protects impl %v of idiom %v", username, existingImplIDStr, idiomIDStr)
		}
	}

	if r.FormValue("impl_version") != strconv.Itoa(impl.Version) {
		return PiErrorf(http.StatusConflict, "Implementation has been concurrently modified (editing version %v, current version is %v)", r.FormValue("impl_version"), impl.Version)
	}

	if err := validateURLFormatOrEmpty(attributionURL); err != nil {
		return PiErrorf(http.StatusBadRequest, "Can't accept URL [%s]", attributionURL)
	}

	if err := validateURLFormatOrEmpty(demoURL); err != nil {
		return PiErrorf(http.StatusBadRequest, "Can't accept URL [%s]", demoURL)
	}

	idiom.EditSummary = "[" + PrintNiceLang(impl.LanguageName) + "] " + r.FormValue("edit_summary")
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

	err = dao.saveExistingIdiom(ctx, key, idiom)
	if err != nil {
		return err
	}

	http.Redirect(w, r, NiceImplURL(idiom, implID, impl.LanguageName), http.StatusFound)
	return nil
}
