package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"

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
			return PiErrorf(http.StatusBadRequest, "Username is mandatory. No anonymous edit.")
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

	ctx := r.Context()
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
	editSummary := fmt.Sprintf("Idiom creation by user [%v] with %s implementation", username, PrintNiceLang(language))

	lead = TruncateBytes(lead, 500)
	keywords = Truncate(keywords, 250)
	imports = Truncate(imports, 200)
	code = TruncateBytes(NoCR(code), 500)
	comment = TruncateBytes(comment, 500)
	attributionURL = Truncate(attributionURL, 250)
	demoURL = Truncate(demoURL, 250)
	docURL = Truncate(docURL, 250)

	log.Infof(ctx, "[%v] is creating new idiom [%v]", username, title)

	if !StringSliceContains(AllLanguages(), language) {
		return PiErrorf(http.StatusBadRequest, "Sorry, [%v] is currently not a supported language. Supported languages are %v.", r.FormValue("impl_language"), AllNiceLangs)
	}

	// TODO put that in a transaction!
	idiomID, err := dao.nextIdiomID(ctx)
	if err != nil {
		return err
	}
	implID, err := dao.nextImplID(ctx)
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

		if u := user.Current(ctx); u != nil {
			idiom.Author = u.String()
			idiom.LastEditor = u.String()
			implementations[0].Author = u.String()
			implementations[0].LastEditor = u.String()
		}
	*/

	_, err = dao.saveNewIdiom(ctx, idiom)
	if err != nil {
		return err
	}

	htmlCacheEvict(ctx, "/about-block-all-idioms")

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
	ctx := r.Context()
	log.Infof(ctx, "[%v] is updating statement of idiom %v", username, existingIDStr)

	idiomID := String2Int(existingIDStr)
	if idiomID == -1 {
		return PiErrorf(http.StatusBadRequest, "%q is not a valid idiom id.", existingIDStr)
	}

	key, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "Could not find idiom %q", existingIDStr)
	}

	isAdmin := IsAdmin(r)
	if idiom.Protected && !isAdmin {
		return PiErrorf(http.StatusUnauthorized, "Can't edit protected idiom %q", existingIDStr)
	}
	if isAdmin {
		wasProtected := idiom.Protected
		idiom.Protected = r.FormValue("idiom_protected") != ""

		if wasProtected && !idiom.Protected {
			log.Infof(ctx, "[%v] unprotects idiom %v", username, existingIDStr)
		}
		if !wasProtected && idiom.Protected {
			log.Infof(ctx, "[%v] protects idiom %v", username, existingIDStr)
		}

		if vars := strings.Replace(r.FormValue("idiom_variables"), " ", "", -1); vars == "" {
			idiom.Variables = nil
		} else {
			idiom.Variables = strings.Split(vars, ",")
		}

		newRelatedURLs := []string{
			strings.TrimSpace(r.FormValue("related_url_1")),
			strings.TrimSpace(r.FormValue("related_url_2")),
			strings.TrimSpace(r.FormValue("related_url_3")),
		}
		newRelatedURLLabels := []string{
			strings.TrimSpace(r.FormValue("related_url_label_1")),
			strings.TrimSpace(r.FormValue("related_url_label_2")),
			strings.TrimSpace(r.FormValue("related_url_label_3")),
		}
		idiom.RelatedURLs = nil
		idiom.RelatedURLLabels = nil
		for i := range newRelatedURLs {
			url, label := newRelatedURLs[i], newRelatedURLLabels[i]
			if url != "" {
				idiom.RelatedURLs = append(idiom.RelatedURLs, url)
				if label == "" {
					label = "See also"
				}
				idiom.RelatedURLLabels = append(idiom.RelatedURLLabels, label)
			}
		}
	}

	if r.FormValue("idiom_version") != strconv.Itoa(idiom.Version) {
		return PiErrorf(http.StatusConflict, "Idiom has been concurrently modified (editing version %v, current version is %v)", r.FormValue("idiom_version"), idiom.Version)
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
	idiom.LeadParagraph = TruncateBytes(idiom.LeadParagraph, 500)
	idiom.ExtraKeywords = Truncate(idiom.ExtraKeywords, 250)
	idiom.EditSummary = Truncate(idiom.EditSummary, 120)

	err = dao.saveExistingIdiom(ctx, key, idiom)
	if err != nil {
		return err
	}

	http.Redirect(w, r, NiceIdiomURL(idiom), http.StatusFound)
	return nil
}
