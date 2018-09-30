package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"
	"golang.org/x/net/context"
)

func idiomDelete(c context.Context, w http.ResponseWriter, r *http.Request) error {
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

func implDelete(c context.Context, w http.ResponseWriter, r *http.Request) error {
	idiomIDStr := r.FormValue("idiomId")
	idiomID := String2Int(idiomIDStr)

	implIDStr := r.FormValue("implId")
	implID := String2Int(implIDStr)

	why := r.FormValue("why")
	if why == "" {
		why = fmt.Sprintf("Admin deletes impl %d", implID)
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
