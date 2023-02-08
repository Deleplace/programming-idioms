package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"google.golang.org/appengine/v2/log"
)

// The client is telling the server that it's using a feature
func using(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	str := r.URL.Path
	str = strings.TrimPrefix(str, "/using/")
	var data struct {
		Page string `json:"page"`
	}
	json.NewDecoder(r.Body).Decode(&data)
	log.Infof(ctx, "Using %s from page %s", str, data.Page)
}
