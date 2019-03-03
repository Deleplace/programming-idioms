package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"google.golang.org/appengine"
)

func idiomDelete(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	idiomIDStr := r.FormValue("idiomId")
	idiomID := String2Int(idiomIDStr)

	why := r.FormValue("why")
	if why == "" {
		why = fmt.Sprintf("Admin deletes idiom %d", idiomID)
	}

	err := dao.deleteIdiom(c, idiomID, why)

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

	// Answer to the "Why?" prompt on delete
	reason := r.FormValue("reason")

	why := r.FormValue("why")
	if why == "" {
		why = fmt.Sprintf("Admin deletes impl %d: %s", implID, reason)
	}
	err := dao.deleteImpl(c, idiomID, implID, why)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		// fmt.Fprint(w, Response{"success": false, "message": err.Error()})
		return err
	}
	fmt.Fprint(w, Response{"success": true})
	return nil
}
