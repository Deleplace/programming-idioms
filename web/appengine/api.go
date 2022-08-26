package main

import (
	"encoding/json"
	"net/http"

	. "github.com/Deleplace/programming-idioms/idioms"
	"github.com/gorilla/mux"
	"google.golang.org/appengine/log"
)

// API JSON endpoints for idioms data.
//
// This enables a SPA client where the backend serves data
// in JSON instead of HTML.

// printJSON will be called at the end of each endpoint.
func printJSON(w http.ResponseWriter, data interface{}, pretty bool) error {
	w.Header().Set("Content-Type", "application/json")
	if pretty {
		// Advantage: output is pretty (human readable)
		// Drawback: the whole data transit through a byte buffer.
		buffer, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		_, err = w.Write(buffer)
		return err
	} else {
		// Drawback: output is ugly.
		// Drawback: the whole data still transit through a byte buffer.
		encoder := json.NewEncoder(w)
		return encoder.Encode(data)
	}
}

// Handle /api/idiom/{idiomId}
func jsonIdiom(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ctx := r.Context()

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	_, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		// TODO distinguish "not found" from "server error"
		return PiErrorf(http.StatusNotFound, "Could not find idiom %q", idiomIDStr)
	}
	// TODO cache the JSON form
	return printJSON(w, idiom, true)
}

// Handle /api/idioms/all
// jsonAllIdioms is redundant with adminImportAjax
func jsonAllIdioms(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	_, idioms, err := dao.getAllIdioms(ctx, 0, "Id")
	if err != nil {
		log.Errorf(ctx, "%v", err)
		return PiErrorf(http.StatusInternalServerError, "Could not retrieve idioms.")
	}
	// TODO cache the JSON form
	return printJSON(w, idioms, true)
}

// Handle /api/search/{q}
func jsonSearch(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	q := vars["q"]

	hits, _, err := findResults(r, q)
	if err != nil {
		return err
	}
	return printJSON(w, hits, true)
}
