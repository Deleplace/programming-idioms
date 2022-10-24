package main

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/idioms"
)

// IdiomCreateFacade if the Facade for the New Idiom page.
type IdiomCreateFacade struct {
	PageMeta               PageMeta
	UserProfile            UserProfile
	LanguageSingleSelector LanguageSingleSelector
}

func idiomCreate(w http.ResponseWriter, r *http.Request) error {
	myToggles := copyToggles(toggles)
	myToggles["editing"] = true

	data := &IdiomCreateFacade{
		PageMeta: PageMeta{
			PageTitle: "New Idiom",
			Toggles:   myToggles,
			ExtraCss: []string{
				hostPrefix() + themeDirectory() + "/css/edit.css",
			},
		},
		UserProfile: readUserProfile(r),
		LanguageSingleSelector: LanguageSingleSelector{
			FieldName: "impl_language",
			Selected:  "",
		},
	}

	if err := templates.ExecuteTemplate(w, "page-idiom-create", data); err != nil {
		return PiErrorf(http.StatusInternalServerError, "%v", err)
	}
	return nil
}
