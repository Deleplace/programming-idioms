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
	var err error

	switch {
	case havingLang != "":
		havingLang = normLang(havingLang)
		c.Infof("Going to a random idiom having lang %v", havingLang)
		_, idiom, err = dao.randomIdiomHaving(c, havingLang)
		c.Infof("Picked idiom #%v: %v", idiom.Id, idiom.Title)
	case notHavingLang != "":
		notHavingLang = normLang(notHavingLang)
		c.Infof("Going to a random idiom not having lang %v", notHavingLang)
		_, idiom, err = dao.randomIdiomNotHaving(c, notHavingLang)
		c.Infof("Picked idiom #%v: %v", idiom.Id, idiom.Title)
	default:
		_, idiom, err = dao.randomIdiom(c)
	}

	if err != nil {
		return err
	}

	http.Redirect(w, r, NiceIdiomURL(idiom), http.StatusFound)
	return nil
}
