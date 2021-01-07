package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"context"

	"github.com/gorilla/mux"

	"google.golang.org/appengine/log"
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

	// #112 Auto add favorite languages
	// Nah, let's do this in the frontend instead
	// if implLangInURL := vars["implLang"]; implLangInURL != "" {
	// 	if favlangs := lookForFavoriteLanguages(r); !StringSliceContainsCaseInsensitive(favlangs, implLangInURL) {
	// 		log.Infof(ctx, "Adding welcome lang in cookie: %q", implLangInURL)
	// 		favlangs = append(favlangs, implLangInURL)
	// 		langsList := strings.Join(favlangs, "_") + "_"
	// 		setLanguagesCookie(w, langsList)
	// 	}
	// }

	if userProfile.Empty() {
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
				http.Redirect(w, r, string(properURL), 302)
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

	idiomTitleInURL := vars["idiomTitle"]
	if idiomTitleInURL != "" && uriNormalize(idiom.Title) != idiomTitleInURL {
		// Maybe the title has changed recently,
		// or someone is attempting a practical joke forging a funny URL ?
		properURL := NiceIdiomURL(idiom)
		http.Redirect(w, r, properURL, 301)
		return nil
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

	implLangInURL := vars["implLang"]
	if implLangInURL != "" && strings.ToLower(selectedImplLang) != strings.ToLower(implLangInURL) {
		// Maybe an accident,
		// or someone is attempting a practical joke forging a funny URL ?
		properURL := NiceImplURL(idiom, selectedImplID, selectedImplLang)
		http.Redirect(w, r, properURL, 301)
		return nil
	}

	if toggles.Any("idiomVotingUp", "implVotingUp") {
		log.Debugf(ctx, "Decorate with votes start...")
		daoVotes.decorateIdiom(ctx, idiom, userProfile.Nickname)
		log.Debugf(ctx, "Decorate with votes end.")
	}

	pageTitle := idiom.Title
	extraKeywords := idiom.ExtraKeywords
	if selectedImplLang != "" {
		// SEO: specify the language in the HTML title, and in meta keywords, for search engine results
		langAliases := selectedImplLang
		if niceLang := PrintNiceLang(selectedImplLang); niceLang != "" {
			pageTitle += ", in " + niceLang
			if selectedImplLang != niceLang {
				langAliases = niceLang + " " + langAliases
			}
		}
		langAliases += " " + strings.Join(LanguageExtraKeywords(selectedImplLang), " ")
		extraKeywords = langAliases + " " + extraKeywords
	}

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
			PageKeywords: extraKeywords,
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
		// SEO: specify the language in the HTML title, for search engine results
		if niceLang := PrintNiceLang(selectedImplLang); niceLang != "" {
			pageTitle += ", in " + niceLang
		}
	}

	myToggles := copyToggles(toggles)
	myToggles["actionEditIdiom"] = !idiom.Protected
	myToggles["actionIdiomHistory"] = true
	myToggles["actionAddImpl"] = !idiom.Protected
	data := &IdiomDetailFacade{
		PageMeta: PageMeta{
			PageTitle:    pageTitle,
			PageKeywords: idiom.ExtraKeywords,
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
