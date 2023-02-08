package main

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/idioms"
)

func idiomDelete(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	idiomIDStr := r.FormValue("idiomId")
	idiomID := String2Int(idiomIDStr)

	why := r.FormValue("why")
	if why == "" {
		why = fmt.Sprintf("Admin deletes idiom %d", idiomID)
	}

	err := dao.deleteIdiom(ctx, idiomID, why)

	htmlCacheEvict(ctx, "/about-block-all-idioms")

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		// fmt.Fprint(w, Response{"success": false, "message": err.Error()})
		return err
	}
	fmt.Fprint(w, Response{"success": true})
	return nil
}

func implDelete(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

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
	err := dao.deleteImpl(ctx, idiomID, implID, why)

	err2 := unindexImpl(ctx, idiomID, implID)
	if err2 != nil {
		errf(ctx, "Unindexing impl %d from idiom %d: %v", implID, idiomID, err2)
		// But keep going
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		// fmt.Fprint(w, Response{"success": false, "message": err.Error()})
		return err
	}
	fmt.Fprint(w, Response{"success": true})
	return nil
}
