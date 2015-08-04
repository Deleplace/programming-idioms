package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
)

func idiomDelete(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	idiomIDStr := r.FormValue("idiomId")
	idiomID := String2Int(idiomIDStr)

	err := dao.deleteIdiom(c, idiomID)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		// fmt.Fprint(w, Response{"success": false, "message": err.Error()})
		return err
	}
	fmt.Fprint(w, Response{"success": true})
	return nil
}

func implDelete(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	idiomIDStr := r.FormValue("idiomId")
	idiomID := String2Int(idiomIDStr)

	implIDStr := r.FormValue("implId")
	implID := String2Int(implIDStr)

	err := dao.deleteImpl(c, idiomID, implID)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		// fmt.Fprint(w, Response{"success": false, "message": err.Error()})
		return err
	}
	fmt.Fprint(w, Response{"success": true})
	return nil
}
