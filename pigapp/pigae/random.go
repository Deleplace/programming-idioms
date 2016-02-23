package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"
	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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
		if havingLang == "" {
			return fmt.Errorf("Invalid language [%s]", vars["havingLang"])
		}
		log.Infof(c, "Going to a random idiom having lang %v", havingLang)
		_, idiom, err = dao.randomIdiomHaving(c, havingLang)
		if err != nil {
			return err
		}
		url = NiceIdiomURL(idiom)
		for _, impl := range idiom.Implementations {
			if impl.LanguageName == havingLang {
				url = NiceImplURL(idiom, impl.Id, havingLang)
			}
		}
	case notHavingLang != "":
		notHavingLang = normLang(notHavingLang)
		if notHavingLang == "" {
			return fmt.Errorf("Invalid language [%s]", vars["notHavingLang"])
		}
		log.Infof(c, "Going to a random idiom not having lang %v", notHavingLang)
		_, idiom, err = dao.randomIdiomNotHaving(c, notHavingLang)
		if err != nil {
			return err
		}
		url = NiceIdiomURL(idiom)
	default:
		_, idiom, err = dao.randomIdiom(c)
		if err != nil {
			return err
		}
		url = NiceIdiomURL(idiom)
	}

	log.Infof(c, "Picked idiom #%v: %v", idiom.Id, idiom.Title)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}
