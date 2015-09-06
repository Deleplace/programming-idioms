package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"
	"github.com/gorilla/mux"

	"appengine"
)

func randomIdiom(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	havingLang := vars["havingLang"]
	notHavingLang := vars["notHavingLang"]
	if havingLang != "" && notHavingLang != "" {
		return fmt.Errorf("Can't have both filters:", havingLang, notHavingLang)
	}

	var idiom *Idiom
	var url string
	var err error

	switch {
	case havingLang != "":
		havingLang = normLang(havingLang)
		c.Infof("Going to a random idiom having lang %v", havingLang)
		_, idiom, err = dao.randomIdiomHaving(c, havingLang)
		url = NiceIdiomURL(idiom)
		for _, impl := range idiom.Implementations {
			if impl.LanguageName == havingLang {
				url = NiceImplURL(idiom, impl.Id, havingLang)
			}
		}
	case notHavingLang != "":
		notHavingLang = normLang(notHavingLang)
		c.Infof("Going to a random idiom not having lang %v", notHavingLang)
		_, idiom, err = dao.randomIdiomNotHaving(c, notHavingLang)
		url = NiceIdiomURL(idiom)
	default:
		_, idiom, err = dao.randomIdiom(c)
		url = NiceIdiomURL(idiom)
	}

	if err != nil {
		return err
	}

	c.Infof("Picked idiom #%v: %v", idiom.Id, idiom.Title)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}
