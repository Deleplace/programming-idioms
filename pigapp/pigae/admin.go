package pigae

import (
	"fmt"
	"net/http"
	"strconv"

	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
	"appengine/user"
)

// IsAdmin determines whether the current user is regarded as Admin by the Google auth provider.
func IsAdmin(r *http.Request) bool {
	c := appengine.NewContext(r) // TODOD check if NewContext is expensive
	u := user.Current(c)
	return u != nil && u.Admin
}

// AdminFacade is the Facade for the Admin page.
type AdminFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
}

func admin(w http.ResponseWriter, r *http.Request) error {
	data := &AdminFacade{
		PageMeta: PageMeta{
			ExtraJs: []string{hostPrefix() + themeDirectory() + "/js/programming-idioms-admin.js"},
			Toggles: toggles,
		},
	}

	return templates.ExecuteTemplate(w, "page-admin", data)
}

func ajaxRefreshToggles(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	dao.deleteCache(c)
	return refreshToggles(c)
}

func ajaxSetToggle(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	name := r.FormValue("toggle")
	valueAsString := r.FormValue("value")

	value, err := strconv.ParseBool(valueAsString)
	if err != nil {
		return err
	}
	toggles[name] = value

	// Save config in distributed Datastore and Memcached
	err = dao.saveAppConfigProperty(c, AppConfigProperty{
		AppConfigId: 0, // TODO meaningful AppConfigId
		Name:        name,
		Value:       value,
	})
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": true})
	return nil
}

func ajaxCreateRelation(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	idiomAIdStr := r.FormValue("idiomAId")
	idiomAId := String2Int(idiomAIdStr)
	idiomBIdStr := r.FormValue("idiomBId")
	idiomBId := String2Int(idiomBIdStr)

	keyA, idiomA, err := dao.getIdiom(c, idiomAId)
	if err != nil {
		return PiError{err.Error(), http.StatusNotFound}
	}

	keyB, idiomB, err := dao.getIdiom(c, idiomBId)
	if err != nil {
		return PiError{err.Error(), http.StatusNotFound}
	}

	idiomA.AddRelation(idiomB)
	if err := dao.saveExistingIdiom(c, keyA, idiomA); err != nil {
		return PiError{err.Error(), http.StatusNotFound}
	}
	if err := dao.saveExistingIdiom(c, keyB, idiomB); err != nil {
		return PiError{err.Error(), http.StatusNotFound}
	}
	return nil
}
