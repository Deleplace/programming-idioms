package main

import (
	"bytes"
	"net/http"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
)

// AboutFacade is the Facade for the About page.
type AboutFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	AllIdioms   []*Idiom
	Coverage    CoverageFacade
}

// CoverageFacade is the facade for the Language Cover block of the About page.
type CoverageFacade struct {
	IdiomIds    []int
	IdiomTitles []string
	Languages   []string
	// Ex: Checked[123]["java"] == 957
	Checked map[int]map[string]int
	// How many impls does each language have?
	// Ex : langImplCount["java"] == 52
	LangImplCount map[string]int
	// What is the cumulated "score" for this language?
	// Including number of impls, doc URLs, and demo URLs.
	// Ex : langImplScore["java"] == 81
	LangImplScore map[string]int
}

func about(w http.ResponseWriter, r *http.Request) error {
	data := AboutFacade{
		PageMeta: PageMeta{
			PageTitle: "About Programming-Idioms",
			Toggles:   toggles,
			ExtraCss:  []string{hostPrefix() + themeDirectory() + "/css/docs.css"},
			ExtraJs:   []string{hostPrefix() + themeDirectory() + "/js/pages/about.js"},
		},
		UserProfile: readUserProfile(r),
	}

	if err := templates.ExecuteTemplate(w, "page-about", data); err != nil {
		return PiError{err.Error(), http.StatusInternalServerError}
	}
	return nil
}

func ajaxAboutProject(w http.ResponseWriter, r *http.Request) error {
	return templates.ExecuteTemplate(w, "block-about-project", nil)
}

func ajaxAboutSeeAlso(w http.ResponseWriter, r *http.Request) error {
	return templates.ExecuteTemplate(w, "block-about-see-also", nil)
}

func ajaxAboutContact(w http.ResponseWriter, r *http.Request) error {
	return templates.ExecuteTemplate(w, "block-about-contact", nil)
}

func ajaxAboutAllIdioms(w http.ResponseWriter, r *http.Request) error {
	c := r.Context()

	debugf(c, "retrieveAllIdioms start...")
	allIdioms, err := retrieveAllIdioms(r)
	if err != nil {
		return err
	}
	debugf(c, "retrieveAllIdioms end.")

	data := AboutFacade{
		PageMeta: PageMeta{
			Toggles: toggles,
		},
		UserProfile: readUserProfile(r),
		AllIdioms:   allIdioms,
	}

	debugf(c, "block-about-all-idioms templating start...")
	if err := templates.ExecuteTemplate(w, "block-about-all-idioms", data); err != nil {
		return PiError{err.Error(), http.StatusInternalServerError}
	}
	debugf(c, "block-about-all-idioms templating end.")
	return nil
}

func ajaxAboutLanguageCoverage(w http.ResponseWriter, r *http.Request) error {
	c := r.Context()
	favlangs := lookForFavoriteLanguages(r)

	if len(favlangs) == 0 {
		if coverageHtml := htmlCacheZipRead(c, "about-block-language-coverage"); coverageHtml != nil {
			// Using the whole HTML block from cache
			debugf(c, "block-about-language-coverage from cache!")
			_, err := w.Write(coverageHtml)
			return err
		}
	}

	coverage, err := languageCoverage(c)
	if err != nil {
		errorf(c, "Error generating language coverage: %v", err)
		return PiError{"Couldn't generate language coverage", 500}
	}
	favoritesFirstWithOrder(coverage.Languages, favlangs, coverage.LangImplCount, coverage.LangImplScore)

	data := AboutFacade{
		PageMeta: PageMeta{
			Toggles: toggles,
		},
		UserProfile: readUserProfile(r),
		Coverage:    coverage,
	}

	var buffer bytes.Buffer
	debugf(c, "block-about-language-coverage templating start...")
	err = templates.ExecuteTemplate(&buffer, "block-about-language-coverage", data)
	debugf(c, "block-about-language-coverage templating end.")
	if err != nil {
		return err
	}

	if len(favlangs) == 0 {
		// Caching may also be done in own goroutine, or defered as a task.
		htmlCacheZipWrite(c, "about-block-language-coverage", buffer.Bytes(), 24*time.Hour)
	}

	_, err = w.Write(buffer.Bytes())
	return err
}

func ajaxAboutRss(w http.ResponseWriter, r *http.Request) error {
	return templates.ExecuteTemplate(w, "block-about-rss", nil)
}

type AboutCheatsheetsFacade struct {
	UserProfile UserProfile
	Langs       []string
}

func ajaxAboutCheatsheets(w http.ResponseWriter, r *http.Request) error {
	data := AboutCheatsheetsFacade{
		UserProfile: readUserProfile(r),
		Langs:       AllLanguages(),
	}
	return templates.ExecuteTemplate(w, "block-about-cheatsheets", data)
}

func languageCoverage(c context.Context) (cover CoverageFacade, err error) {
	checked := map[int]map[string]int{}
	langImplCount := map[string]int{}
	langImplScore := map[string]int{}
	debugf(c, "Loading full idiom list...")
	_, idioms, err := dao.getAllIdioms(c, 299, "-ImplCount") // TODO change 299 ?!
	if err != nil {
		return cover, err
	}
	debugf(c, "Full idiom list loaded.")
	idiomIds := make([]int, len(idioms))
	idiomTitles := make([]string, len(idioms))

	debugf(c, "Counting impls of each idiom...")
	for i, idiom := range idioms {
		idiomIds[i] = idiom.Id
		idiomTitles[i] = idiom.Title

		for _, impl := range idiom.Implementations {
			if checked[idiom.Id] == nil {
				checked[idiom.Id] = map[string]int{impl.LanguageName: impl.Id}
			} else {
				checked[idiom.Id][impl.LanguageName] = impl.Id
			}
			langImplCount[impl.LanguageName]++
			// 1 point for existing
			score := 1
			langImplScore[impl.LanguageName]++
			if len(impl.DocumentationURL) > 6 {
				// 1 point for documentation URL
				score++
			}
			if len(impl.DemoURL) > 6 {
				// 1 point for online demo
				score++
			}
			langImplScore[impl.LanguageName] += score
		}
	}
	debugf(c, "Impls of each idiom counted.")

	cover = CoverageFacade{
		IdiomIds:      idiomIds,
		IdiomTitles:   idiomTitles,
		Languages:     AllLanguages(),
		Checked:       checked,
		LangImplCount: langImplCount,
		LangImplScore: langImplScore,
	}
	return cover, nil
}

func hiddenizeExtraColumns(j int, classBefore, classExtra string) string {
	const maxCols = 11 // 12 lang columns 0..11 are shown
	if j > maxCols {
		return classExtra
	}
	return classBefore
}
