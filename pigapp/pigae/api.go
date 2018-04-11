package pigae

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
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
	c := appengine.NewContext(r)

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	_, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		// TODO distinguish "not found" from "server error"
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
	}
	// TODO cache the JSON form
	return printJSON(w, idiom, true)
}

// Handle /api/idioms/all
// jsonAllIdioms is redundant with adminImportAjax
func jsonAllIdioms(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	_, idioms, err := dao.getAllIdioms(c, 0, "Id")
	if err != nil {
		log.Errorf(c, "%v", err)
		return PiError{"Could not retrieve idioms.", http.StatusInternalServerError}
	}

	// Serialize only relevant fields, to minimize response payload.
	cidioms := make([]CompactIdiom, len(idioms))
	for i, idiom := range idioms {
		cidioms[i] = compactIdiom(idiom)
	}

	// TODO cache the JSON form
	return printJSON(w, cidioms, false)
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

// Handle /api/random-id
func jsonRandomID(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	url, err := randomIdiomURL(c)
	if err != nil {
		return err
	}
	// url looks like "/idiom/16/depth-first-traversing-of-a-binary-tree"
	parts := strings.Split(url, "/")
	idiomIDstr := parts[2]
	// idiomIDstr looks like "16"
	idiomID, err := strconv.Atoi(idiomIDstr)
	if err != nil {
		return err
	}
	return printJSON(w, idiomID, true)
}
