package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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
	c := appengine.NewContext(r)

	log.Debugf(c, "retrieveAllIdioms start...")
	allIdioms, err := retrieveAllIdioms(r)
	if err != nil {
		return err
	}
	log.Debugf(c, "retrieveAllIdioms end.")

	data := AboutFacade{
		PageMeta: PageMeta{
			Toggles: toggles,
		},
		UserProfile: readUserProfile(r),
		AllIdioms:   allIdioms,
	}

	log.Debugf(c, "block-about-all-idioms templating start...")
	if err := templates.ExecuteTemplate(w, "block-about-all-idioms", data); err != nil {
		return PiError{err.Error(), http.StatusInternalServerError}
	}
	log.Debugf(c, "block-about-all-idioms templating end.")
	return nil
}

func ajaxAboutLanguageCoverage(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	coverage, _ := languageCoverage(c)
	favlangs := lookForFavoriteLanguages(r)
	log.Debugf(c, "favoritesFirstWithOrder start...")
	favoritesFirstWithOrder(coverage.Languages, favlangs, coverage.LangImplCount)
	log.Debugf(c, "favoritesFirstWithOrder end.")

	data := AboutFacade{
		PageMeta: PageMeta{
			Toggles: toggles,
		},
		UserProfile: readUserProfile(r),
		Coverage:    coverage,
	}

	log.Debugf(c, "block-about-language-coverage templating start...")
	err := templates.ExecuteTemplate(w, "block-about-language-coverage", data)
	log.Debugf(c, "block-about-language-coverage templating end.")
	return err
}

func ajaxAboutRss(w http.ResponseWriter, r *http.Request) error {
	return templates.ExecuteTemplate(w, "block-about-rss", nil)
}

func languageCoverage(c context.Context) (cover CoverageFacade, err error) {
	checked := map[int]map[string]int{}
	langImplCount := map[string]int{}
	log.Debugf(c, "Loading full idiom list...")
	_, idioms, err := dao.getAllIdioms(c, 199, "-ImplCount") // TODO change 199 ?!
	if err != nil {
		return cover, err
	}
	log.Debugf(c, "Full idiom list loaded.")
	idiomIds := make([]int, len(idioms))
	idiomTitles := make([]string, len(idioms))

	log.Debugf(c, "Counting impls of each idiom...")
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
	log.Debugf(c, "Impls of each idiom counted.")

	cover = CoverageFacade{
		IdiomIds:      idiomIds,
		IdiomTitles:   idiomTitles,
		Languages:     allLanguages(),
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
