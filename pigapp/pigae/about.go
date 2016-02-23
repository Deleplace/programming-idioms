package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
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
	// Ex : langImplCount["java"] == 52
	LangImplCount map[string]int
}

func about(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	allIdioms, err := retrieveAllIdioms(r)
	if err != nil {
		return err
	}

	coverage, _ := languageCoverage(c)
	favlangs := lookForFavoriteLanguages(r)
	favoritesFirstWithOrder(coverage.Languages, favlangs, coverage.LangImplCount)

	aboutToggles := copyToggles(toggles)

	data := AboutFacade{
		PageMeta: PageMeta{
			PageTitle: "About Programming-Idioms",
			Toggles:   aboutToggles,
			ExtraCss:  []string{hostPrefix() + themeDirectory() + "/css/docs.css"},
			ExtraJs:   []string{hostPrefix() + themeDirectory() + "/js/pages/about.js"},
		},
		UserProfile: readUserProfile(r),
		AllIdioms:   allIdioms,
		Coverage:    coverage,
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
	allIdioms, err := retrieveAllIdioms(r)
	if err != nil {
		return err
	}

	data := AboutFacade{
		PageMeta: PageMeta{
			Toggles: toggles,
		},
		UserProfile: readUserProfile(r),
		AllIdioms:   allIdioms,
	}

	if err := templates.ExecuteTemplate(w, "block-about-all-idioms", data); err != nil {
		return PiError{err.Error(), http.StatusInternalServerError}
	}
	return nil
}

func ajaxAboutLanguageCoverage(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	coverage, _ := languageCoverage(c)
	favlangs := lookForFavoriteLanguages(r)
	favoritesFirstWithOrder(coverage.Languages, favlangs, coverage.LangImplCount)

	data := AboutFacade{
		PageMeta: PageMeta{
			Toggles: toggles,
		},
		UserProfile: readUserProfile(r),
		Coverage:    coverage,
	}

	return templates.ExecuteTemplate(w, "block-about-language-coverage", data)
}

func ajaxAboutRss(w http.ResponseWriter, r *http.Request) error {
	return templates.ExecuteTemplate(w, "block-about-rss", nil)
}

func languageCoverage(c context.Context) (cover CoverageFacade, err error) {
	checked := map[int]map[string]int{}
	langImplCount := map[string]int{}
	_, idioms, err := dao.getAllIdioms(c, 199, "-ImplCount") // TODO change 199 ?!
	if err != nil {
		return cover, err
	}
	idiomIds := make([]int, len(idioms))
	idiomTitles := make([]string, len(idioms))

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
		}
	}

	cover = CoverageFacade{
		IdiomIds:      idiomIds,
		IdiomTitles:   idiomTitles,
		Languages:     dao.languagesHavingImpl(c),
		Checked:       checked,
		LangImplCount: langImplCount,
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
