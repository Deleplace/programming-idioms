package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	. "github.com/Deleplace/programming-idioms/idioms"

	"context"

	"github.com/gorilla/mux"

	"google.golang.org/appengine/v2/log"
)

// IdiomDetailFacade is the Facade for the Idiom Detail page.
type IdiomDetailFacade struct {
	PageMeta         PageMeta
	UserProfile      UserProfile
	Idiom            *Idiom
	SelectedImplID   int
	SelectedImplLang string
}

func idiomDetail(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ctx := r.Context()
	userProfile := readUserProfile(r)
	favlangs := userProfile.FavoriteLanguages

	pushResources := func() {
		if _, err := r.Cookie("v"); err == nil {
			log.Infof(ctx, "Returning visitor: no resource server push needed.")
		} else {
			// TODO enable when page-idiom-detail-minimal is default, and cookie "v" is in use
			//
			// log.Infof(ctx, "New visitor: server push resources!")
			// prefix := hostPrefix() + themeDirectory()
			// var links bytes.Buffer
			// fmt.Fprintf(&links, "<%s>; rel=preload; as=%s, ", prefix+"/js/pages/idiom-detail-minimal.js", "script")
			// fmt.Fprintf(&links, "<%s>; rel=preload; as=%s, ", prefix+"/css/pages/idiom-detail-minimal.css", "style")
			// fmt.Fprintf(&links, "<%s>; rel=preload; as=%s", "/default_20200205_/img/wheel.svg", "image")
			// w.Header().Set("Link", links.String())
		}
	}

	userProfileEmpty := userProfile.Empty()

	// #185 Take "Language of current impl" into account for the right bar Cheatsheet links,
	// even it this language is "not yet" in the user's favorites.
	implLangInURL := vars["implLang"]
	if implLangInURL != "" {
		if !StringSliceContainsCaseInsensitive(userProfile.FavoriteLanguages, implLangInURL) {
			// userProfile.FavoriteLanguages = append(userProfile.FavoriteLanguages, implLangInURL)
			// But that would mess up a bit with the favlanag bar, so... let's do this in the frontend instead
		}
		// #112 Auto add favorite languages by setting cookie
		// Nah, let's do this in the frontend instead
	}

	if userProfileEmpty {
		//
		// Zero-preference ≡ anonymous visit ≡ server cache enabled
		//

		path := r.URL.RequestURI()
		if cachedPage := htmlCacheRead(ctx, path); cachedPage != nil {
			// Using the whole HTML block from Memcache
			log.Debugf(ctx, "%s from memcache!", path)
			_, err := w.Write(cachedPage)
			return err
		}
		log.Debugf(ctx, "%s not in memcache.", path)

		var buffer bytes.Buffer
		err := generateIdiomDetailPage(ctx, &buffer, vars)
		if err != nil {
			if properURL, ok := err.(needRedirectError); ok {
				http.Redirect(w, r, string(properURL), http.StatusFound) // 302
				return nil
			}
			return err
		}
		pushResources()
		_, err = w.Write(buffer.Bytes())
		if err != nil {
			return err
		}

		htmlCacheWrite(ctx, path, buffer.Bytes(), 24*time.Hour)
		// Note that this cache entry must be later invalidated in case
		// of any modification in this idiom.

		// Here we just cached 1 HTML page for 1 day.
		// We tried previously to agressively trigger htmlRecacheNowAndTomorrow,
		// but it didn't lead to great results in shared memcache.

		return nil
	}

	//
	// Below: for users having custom profile
	//
	// WARNING the code below is currently very redundant with the second part of generateIdiomDetail.
	// Please try to not diverge.
	//

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)
	var canonicalURL string

	_, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "Could not find idiom %q", idiomIDStr)
	}

	var selectedImplID int
	var selectedImplLang string
	if selectedImplIDStr := vars["implId"]; selectedImplIDStr != "" {
		selectedImplID = String2Int(selectedImplIDStr)
		for _, impl := range idiom.Implementations {
			if selectedImplID == impl.Id {
				selectedImplLang = impl.LanguageName
				break
			}
		}
		if selectedImplLang == "" {
			// The requested implementation was not found.
			properURL := NiceIdiomURL(idiom)
			http.Redirect(w, r, properURL, 302)
			return nil
		}
		canonicalURL = host() + NiceImplRelativeURL(idiom, selectedImplID, selectedImplLang)
	} else {
		// Just the idiom, no specific impl
		canonicalURL = host() + NiceIdiomRelativeURL(idiom)
	}

	idiomTitleInURL := vars["idiomTitle"]
	if idiomTitleInURL != "" && uriNormalize(idiom.Title) != idiomTitleInURL {
		// Maybe the title has changed recently,
		// or someone is attempting a practical joke forging a funny URL ?
		properURL := canonicalURL
		http.Redirect(w, r, properURL, http.StatusMovedPermanently) // 301
		return nil
	}

	includeNonFav := seeNonFavorite(r)
	log.Debugf(ctx, "Reorder impls start...")
	implFavoriteLanguagesFirstWithOrder(idiom, favlangs, selectedImplLang, includeNonFav)

	// Selected impl as very first element
	for i := range idiom.Implementations {
		if idiom.Implementations[i].LanguageName != selectedImplLang {
			break
		}
		if idiom.Implementations[i].Id == selectedImplID {
			idiom.Implementations[0], idiom.Implementations[i] = idiom.Implementations[i], idiom.Implementations[0]
			break
		}
	}
	log.Debugf(ctx, "Reorder impls end.")

	if implLangInURL != "" && strings.ToLower(selectedImplLang) != strings.ToLower(implLangInURL) {
		// Maybe an accident,
		// or someone is attempting a practical joke forging a funny URL ?
		properURL := NiceImplURL(idiom, selectedImplID, selectedImplLang)
		http.Redirect(w, r, properURL, http.StatusMovedPermanently) // 301
		return nil
	}

	if toggles.Any("idiomVotingUp", "implVotingUp") {
		log.Debugf(ctx, "Decorate with votes start...")
		daoVotes.decorateIdiom(ctx, idiom, userProfile.Nickname)
		log.Debugf(ctx, "Decorate with votes end.")
	}

	pageTitle := idiom.Title
	if selectedImplLang != "" {
		// SEO: specify the language in the HTML title
		pageTitle += ", in " + PrintNiceLang(selectedImplLang)
	}

	pageKeywords := preparePageKeywords(idiom, selectedImplLang)

	extraJS := []string{}
	if IsAdmin(r) {
		extraJS = append(extraJS, hostPrefix()+themeDirectory()+"/js/programming-idioms-admin.js")
	}

	myToggles := copyToggles(toggles)
	myToggles["actionEditIdiom"] = !idiom.Protected || IsAdmin(r)
	myToggles["actionIdiomHistory"] = true
	myToggles["actionAddImpl"] = !idiom.Protected || IsAdmin(r)
	data := &IdiomDetailFacade{
		PageMeta: PageMeta{
			PageTitle:    pageTitle,
			PageKeywords: pageKeywords,
			CanonicalURL: canonicalURL,
			Toggles:      myToggles,
			ExtraJs:      extraJS,
		},
		UserProfile:      userProfile,
		Idiom:            idiom,
		SelectedImplID:   selectedImplID,
		SelectedImplLang: selectedImplLang,
	}

	pushResources()
	log.Debugf(ctx, "ExecuteTemplate start...")
	err = templates.ExecuteTemplate(w, "page-idiom-detail", data)
	// err = templates.ExecuteTemplate(w, "page-idiom-detail-minimal", data)
	log.Debugf(ctx, "ExecuteTemplate end.")
	return err
}

func generateIdiomDetailPage(ctx context.Context, w io.Writer, vars map[string]string) error {
	//
	// WARNING this code is currently very redundant with the second part of idiomDetail.
	// Please try to not diverge.
	//

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)
	var canonicalURL string

	_, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "Could not find idiom %q", idiomIDStr)
	}

	idiomTitleInURL := vars["idiomTitle"]
	if idiomTitleInURL != "" && uriNormalize(idiom.Title) != idiomTitleInURL {
		// Maybe the title has changed recently,
		// or someone is attempting a practical joke forging a funny URL ?
		properURL := NiceIdiomURL(idiom)
		return needRedirectError(properURL)
	}

	var selectedImplID int
	var selectedImplLang string
	if selectedImplIDStr := vars["implId"]; selectedImplIDStr != "" {
		selectedImplID = String2Int(selectedImplIDStr)
		for _, impl := range idiom.Implementations {
			if selectedImplID == impl.Id {
				selectedImplLang = impl.LanguageName
				break
			}
		}
		if selectedImplLang == "" {
			// The requested implementation was not found.
			properURL := NiceIdiomURL(idiom)
			return needRedirectError(properURL)
		}
		canonicalURL = host() + NiceImplRelativeURL(idiom, selectedImplID, selectedImplLang)
	} else {
		// Just the idiom, no specific impl
		canonicalURL = host() + NiceIdiomRelativeURL(idiom)
	}

	implFavoriteLanguagesFirstWithOrder(idiom, nil, selectedImplLang, true)

	// Selected impl as very first element
	for i := range idiom.Implementations {
		if idiom.Implementations[i].LanguageName != selectedImplLang {
			break
		}
		if idiom.Implementations[i].Id == selectedImplID {
			idiom.Implementations[0], idiom.Implementations[i] = idiom.Implementations[i], idiom.Implementations[0]
			break
		}
	}

	implLangInURL := vars["implLang"]
	if implLangInURL != "" && strings.ToLower(selectedImplLang) != strings.ToLower(implLangInURL) {
		// Maybe an accident,
		// or someone is attempting a practical joke forging a funny URL ?
		properURL := NiceImplURL(idiom, selectedImplID, selectedImplLang)
		return needRedirectError(properURL)
	}

	pageTitle := idiom.Title
	if selectedImplLang != "" {
		// SEO: specify the language in the HTML title
		pageTitle += ", in " + PrintNiceLang(selectedImplLang)
	}

	pageKeywords := preparePageKeywords(idiom, selectedImplLang)

	myToggles := copyToggles(toggles)
	myToggles["actionEditIdiom"] = !idiom.Protected
	myToggles["actionIdiomHistory"] = true
	myToggles["actionAddImpl"] = !idiom.Protected
	data := &IdiomDetailFacade{
		PageMeta: PageMeta{
			PageTitle:    pageTitle,
			PageKeywords: pageKeywords,
			CanonicalURL: canonicalURL,
			Toggles:      myToggles,
		},
		UserProfile:      EmptyUserProfile(),
		Idiom:            idiom,
		SelectedImplID:   selectedImplID,
		SelectedImplLang: selectedImplLang,
	}

	log.Debugf(ctx, "ExecuteTemplate start...")
	err = templates.ExecuteTemplate(w, "page-idiom-detail", data)
	// err = templates.ExecuteTemplate(w, "page-idiom-detail-minimal", data)
	log.Debugf(ctx, "ExecuteTemplate end.")
	return err
}

type needRedirectError string

func (err needRedirectError) Error() string {
	return string(err)
}

// SEO meta keywords
func preparePageKeywords(idiom *Idiom, selectedImplLang string) string {
	// Keywords from the Idiom itself
	keywords := idiom.ExtraKeywords
	// #144 most idiom.ExtraKeywords have space separators, but wee need commas.
	// Let's replace with commas. Drawback: a multiple word keyword don't survive.
	keywords = strings.Replace(keywords, " ", ", ", -1)
	keywords = strings.Replace(keywords, ",,", ",", -1)

	// Extra keywords for the languqge name of the selected impl
	if selectedImplLang != "" {
		langAliases := selectedImplLang
		niceLang := PrintNiceLang(selectedImplLang)
		if selectedImplLang != niceLang {
			langAliases = niceLang + ", " + langAliases
		}
		langAliases += ", " + strings.Join(LanguageExtraKeywords(selectedImplLang), ", ")
		keywords = langAliases + ", " + keywords
	}

	return keywords
}
