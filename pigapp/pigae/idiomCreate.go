package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"
	"golang.org/x/net/context"
)

// IdiomCreateFacade if the Facade for the New Idiom page.
type IdiomCreateFacade struct {
	PageMeta               PageMeta
	UserProfile            UserProfile
	LanguageSingleSelector LanguageSingleSelector
}

func idiomCreate(c context.Context, w http.ResponseWriter, r *http.Request) error {
	myToggles := copyToggles(toggles)
	myToggles["editing"] = true

	data := &IdiomCreateFacade{
		PageMeta: PageMeta{
			PageTitle: "New Idiom",
			Toggles:   myToggles,
		},
		UserProfile: readUserProfile(c, r),
		LanguageSingleSelector: LanguageSingleSelector{
			FieldName: "impl_language",
			Selected:  "",
		},
	}

	if err := templates.ExecuteTemplate(w, "page-idiom-create", data); err != nil {
		return PiError{err.Error(), http.StatusInternalServerError}
	}
	return nil
}
