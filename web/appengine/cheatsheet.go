package main

import (
	"net/http"
	"strconv"

	. "github.com/Deleplace/programming-idioms/idioms"
	"github.com/gorilla/mux"
	gaesearch "google.golang.org/appengine/search"
)

// LEGACY
// This code is currently not used. See cheatsheetmulti.go, even for 1 single language.

//
// A Cheat Sheet is a single page containing the implementations of all idioms, for a
// single language.
//

// CheatSheetFacade is the Facade for the Cheat Sheets.
type CheatSheetFacade struct {
	PageMeta        PageMeta
	UserProfile     UserProfile
	Lang            string
	CheatsheetLines []cheatSheetLineDoc
}

func cheatsheet(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	lang := vars["lang"]
	ctx := r.Context()

	// Security belt. Might be changed if needed.
	limit := 1000

	// This uses the Search API to retrieve just the data we need.
	cheatsheetLines, err := dao.getCheatSheet(ctx, lang, limit)
	if err != nil {
		return PiErrorf(http.StatusInternalServerError, "%v", err)
	}

	// Don't repeat idiom ID, title, lead on consecutive rows.
	// Note: this *may* not play well with the JS text filter.
	for i := 1; i < len(cheatsheetLines); i++ {
		if cheatsheetLines[i].IdiomID == cheatsheetLines[i-1].IdiomID {
			cheatsheetLines[i].IdiomID = ""
			cheatsheetLines[i].IdiomTitle = ""
			cheatsheetLines[i].IdiomLeadParagraph = "Alternative implementation"
		}
	}

	data := CheatSheetFacade{
		PageMeta: PageMeta{
			PageTitle: PrintNiceLang(lang) + " cheat sheet",
			Toggles:   toggles,
		},
		UserProfile:     readUserProfile(r),
		Lang:            lang,
		CheatsheetLines: cheatsheetLines,
	}

	if err := templates.ExecuteTemplate(w, "page-cheatsheet", data); err != nil {
		return PiErrorf(http.StatusInternalServerError, "%v", err)
	}
	return nil
}

// useful for calling markup2CSS on cheatSheetLineDoc fields
func atom2string(atom gaesearch.Atom) string {
	return string(atom)
}

// useful for calling markup2CSS on cheatSheetLineDoc fields
func atom2int(atom gaesearch.Atom) int {
	i, _ := strconv.Atoi(string(atom))
	// In case of an invalid argument, the return value is 0
	return i
}
