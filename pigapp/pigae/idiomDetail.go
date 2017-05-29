package pigae

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"google.golang.org/appengine"
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
	c := appengine.NewContext(r)
	userProfile := readUserProfile(r)
	favlangs := userProfile.FavoriteLanguages

	if userProfile.Empty() {
		//
		// Zero-preference ≡ anonymous visit ≡ cache enabled
		//
		path := r.URL.RequestURI()
		if cachedPage := htmlCacheRead(c, path); cachedPage != nil {
			// Using the whole HTML block from Memcache
			log.Debugf(c, "%s from memcache!", path)
			_, err := w.Write(cachedPage)
			return err
		}
		log.Debugf(c, "%s not in memcache.", path)

		var buffer bytes.Buffer
		err := generateIdiomDetailPage(c, &buffer, vars)
		if err != nil {
			if properURL, ok := err.(needRedirectError); ok {
				http.Redirect(w, r, string(properURL), 302)
				return nil
			}
			return err
		}
		_, err = w.Write(buffer.Bytes())
		if err != nil {
			return err
		}

		htmlCacheWrite(c, path, buffer.Bytes(), 24*time.Hour)
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

	_, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
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
	}

	includeNonFav := seeNonFavorite(r)
	log.Debugf(c, "Reorder impls start...")
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
	log.Debugf(c, "Reorder impls end.")

	implLangInURL := vars["implLang"]
	if implLangInURL != "" && strings.ToLower(selectedImplLang) != strings.ToLower(implLangInURL) {
		// Maybe an accident,
		// or someone is attempting a practical joke forging a funny URL ?
		properURL := NiceImplURL(idiom, selectedImplID, selectedImplLang)
		http.Redirect(w, r, properURL, 301)
		return nil
	}

	log.Debugf(c, "Decorate with votes start...")
	daoVotes.decorateIdiom(c, idiom, userProfile.Nickname)
	log.Debugf(c, "Decorate with votes end.")

	pageTitle := idiom.Title
	if selectedImplLang != "" {
		// SEO: specify the language in the HTML title, for search engine results
		if niceLang := PrintNiceLang(selectedImplLang); niceLang != "" {
			pageTitle += ", in " + niceLang
		}
	}

	myToggles := copyToggles(toggles)
	myToggles["actionEditIdiom"] = !idiom.Protected || IsAdmin(r)
	myToggles["actionIdiomHistory"] = true
	myToggles["actionAddImpl"] = !idiom.Protected || IsAdmin(r)
	data := &IdiomDetailFacade{
		PageMeta: PageMeta{
			PageTitle:    pageTitle,
			PageKeywords: idiom.ExtraKeywords,
			Toggles:      myToggles,
		},
		UserProfile:      userProfile,
		Idiom:            idiom,
		SelectedImplID:   selectedImplID,
		SelectedImplLang: selectedImplLang,
	}

	log.Debugf(c, "ExecuteTemplate start...")
	err = templates.ExecuteTemplate(w, "page-idiom-detail", data)
	log.Debugf(c, "ExecuteTemplate end.")
	return err
}

func generateIdiomDetailPage(c context.Context, w io.Writer, vars map[string]string) error {
	//
	// WARNING this code is currently very redundant with the second part of idiomDetail.
	// Please try to not diverge.
	//

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	_, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
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
			Toggles:      myToggles,
		},
		UserProfile:      EmptyUserProfile(),
		Idiom:            idiom,
		SelectedImplID:   selectedImplID,
		SelectedImplLang: selectedImplLang,
	}

	log.Debugf(c, "ExecuteTemplate start...")
	err = templates.ExecuteTemplate(w, "page-idiom-detail", data)
	log.Debugf(c, "ExecuteTemplate end.")
	return err
}

type needRedirectError string

func (err needRedirectError) Error() string {
	return string(err)
}
