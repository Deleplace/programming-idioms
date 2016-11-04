package pigae

import (
	"net/http"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
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

	selectedImplID := 0
	selectedImplLang := ""
	if selectedImplIDStr := vars["implId"]; len(selectedImplIDStr) > 0 {
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

	favlangs := lookForFavoriteLanguages(r)
	includeNonFav := seeNonFavorite(r)
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

	implLangInURL := vars["implLang"]
	if implLangInURL != "" && strings.ToLower(selectedImplLang) != strings.ToLower(implLangInURL) {
		// Maybe an accident,
		// or someone is attempting a practical joke forging a funny URL ?
		properURL := NiceImplURL(idiom, selectedImplID, selectedImplLang)
		http.Redirect(w, r, properURL, 301)
		return nil
	}

	userProfile := readUserProfile(r)
	daoVotes.decorateIdiom(c, idiom, userProfile.Nickname)

	pageTitle := idiom.Title
	if selectedImplLang != "" {
		// SEO: specify the language in the HTML title, for search engine results
		if niceLang := printNiceLang(selectedImplLang); niceLang != "" {
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
	return templates.ExecuteTemplate(w, "page-idiom-detail", data)
}
