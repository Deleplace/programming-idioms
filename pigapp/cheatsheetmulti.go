package main

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"
)

// CheatSheetMultipleFacade is the Facade for the Cheat Sheets with multiple languages.
type CheatSheetMultipleFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Langs       []string
	Lines       []cheatSheetLineMulti
}

type cheatSheetLineMulti struct {
	IdiomID            int
	IdiomTitle         string
	IdiomLeadParagraph string

	// Depth represents the max number of impl for one language in ByLanguage
	Depth []struct{}

	ByLanguage []cheatSheetLineDocs
}

func (chsh *cheatSheetLineMulti) hasCell(langIndex, alt int) bool {
	return alt < len(chsh.Depth) &&
		langIndex < len(chsh.ByLanguage) &&
		alt < len(chsh.ByLanguage[langIndex])
}

func (chsh *cheatSheetLineMulti) computeDepth() {
	n := 0
	for _, lineDocs := range chsh.ByLanguage {
		if len(lineDocs) > n {
			n = len(lineDocs)
		}
	}
	chsh.Depth = make([]struct{}, n)
}

func cheatsheetMulti(w http.ResponseWriter, r *http.Request) error {
	subpath := strings.TrimPrefix(r.URL.Path, "/cheatsheet/")
	subpath = strings.TrimSuffix(subpath, "/")
	langs := strings.Split(subpath, "/")
	return cheatsheetMultiLangs(w, r, langs)
}

func cheatsheetMultiLangs(w http.ResponseWriter, r *http.Request, langs []string) error {
	ctx := r.Context()

	for i := range langs {
		lang := NormLang(langs[i])
		if !StringSliceContains(AllLanguages(), lang) {
			return PiErrorf(http.StatusBadRequest, "Sorry, [%v] is currently not a supported language. Supported languages are %v.", langs[i], AllNiceLangs)
		}
		langs[i] = lang
	}

	// Security belt. Might be changed if needed.
	limit := 1000

	var lines []cheatSheetLineMulti

	byIdiomID := map[int][]cheatSheetLineDocs{}
	idiomTitles := map[int]string{}
	idiomLeadParagraphs := map[int]string{}
	for langIndex, lang := range langs {
		cheatsheetLines, err := dao.getCheatSheet(ctx, lang, limit)
		if err != nil {
			return PiErrorf(http.StatusInternalServerError, "%v", err)
		}
		for _, line := range cheatsheetLines {
			idiomID, err := strconv.Atoi(string(line.IdiomID))
			if err != nil {
				return err
			}
			if _, exists := byIdiomID[idiomID]; exists {
				byIdiomID[idiomID][langIndex] = append(byIdiomID[idiomID][langIndex], line)
			} else {
				byIdiomID[idiomID] = make([]cheatSheetLineDocs, len(langs))
				byIdiomID[idiomID][langIndex] = cheatSheetLineDocs{line}
				idiomTitles[idiomID] = string(line.IdiomTitle)
				idiomLeadParagraphs[idiomID] = string(line.IdiomLeadParagraph)
			}

		}
	}
	idiomIDs := make([]int, 0, len(byIdiomID))
	for id := range byIdiomID {
		idiomIDs = append(idiomIDs, id)
	}
	sort.Ints(idiomIDs)
	for _, idiomID := range idiomIDs {
		line := cheatSheetLineMulti{
			IdiomID:            idiomID,
			IdiomTitle:         idiomTitles[idiomID],
			IdiomLeadParagraph: idiomLeadParagraphs[idiomID],
			// Depth:              make([]struct{}, 2),
			ByLanguage: byIdiomID[idiomID],
		}
		line.computeDepth()
		lines = append(lines, line)
	}

	pageTitle := ""
	glue := ""
	for _, lang := range langs {
		pageTitle += glue + PrintNiceLang(lang)
		glue = ", "
	}
	pageTitle += " cheat sheet"
	data := CheatSheetMultipleFacade{
		PageMeta: PageMeta{
			PageTitle: pageTitle,
			Toggles:   toggles,
			ExtraCss:  []string{hostPrefix() + themeDirectory() + "/css/pages/cheatsheetmulti.css"},
			ExtraJs:   []string{hostPrefix() + themeDirectory() + "/js/pages/cheatsheetmulti.js"},
		},
		UserProfile: readUserProfile(r),
		Langs:       langs,
		Lines:       lines,
	}

	if err := templates.ExecuteTemplate(w, "page-cheatsheet-multi", data); err != nil {
		return PiErrorf(http.StatusInternalServerError, "%v", err)
	}
	return nil
}

// In case someone manually tries to reach /cheatsheet or /cheatsheet/ .
func cheatsheetPageRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/cheatsheets", http.StatusFound)
}
