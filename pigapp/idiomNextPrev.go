package main

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"
)

//
// #192 Keyboard shortcuts
// p Go to "previous" idiom
// n Go to "next" idiom
//
// Sometimes the exact (n-1) or (n+1) target ID doesn't exist, then
// try the next existing one, or loop.

func nextIdiom(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	vars := mux.Vars(r)
	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	// This is is bit fragile and inefficient.
	// Consider a memcached full list instead...
	for nextID := idiomID + 1; nextID <= idiomID+20; nextID++ {
		_, nextHead, err := dao.getIdiom(ctx, nextID)
		if err == nil {
			url := NiceIdiomRelativeURL(nextHead)
			http.Redirect(w, r, url, http.StatusFound)
			return nil
		}
	}
	// Loop: after the very last idiom, we may return... idiom 1
	_, minIdiomHead, err := dao.getIdiom(ctx, 1)
	if err != nil {
		return nil
	}
	url := NiceIdiomRelativeURL(minIdiomHead)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

func previousIdiom(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	vars := mux.Vars(r)
	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	// This is is bit fragile and inefficient.
	// Consider a memcached full list instead...
	for prevID := idiomID - 1; prevID >= idiomID-20 && prevID >= 1; prevID-- {
		_, prevHead, err := dao.getIdiom(ctx, prevID)
		if err == nil {
			url := NiceIdiomRelativeURL(prevHead)
			http.Redirect(w, r, url, http.StatusFound)
			return nil
		}
	}
	// Loop: before idiom 1, we may return... the very max idiom
	maxIdiomHead, err := dao.getMaxIdiomIDTitle(ctx)
	if err != nil {
		return nil
	}
	url := NiceIdiomRelativeURL(maxIdiomHead)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}
