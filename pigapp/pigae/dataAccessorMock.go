package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
	"appengine/datastore"
)

// DatastoreAccessorMock should be useful for unit testing. But is is dead code for now.
type DatastoreAccessorMock struct {
	idioms []*Idiom
}

func (a DatastoreAccessorMock) getIdiom(c appengine.Context, idiomID int) (key *datastore.Key, idiom *Idiom, thiserror error) {
	return
}

func (a DatastoreAccessorMock) saveNewIdiom(c appengine.Context, idiom *Idiom) (key *datastore.Key, thiserror error) {
	return
}

func (a DatastoreAccessorMock) saveExistingIdiom(c appengine.Context, key *datastore.Key, idiom *Idiom) error {
	return nil
}

func (a DatastoreAccessorMock) getAllIdioms(c appengine.Context, limit int, order string) (keys []*datastore.Key, idioms []*Idiom, thiserror error) {
	return
}

func (a DatastoreAccessorMock) deleteAllIdioms(c appengine.Context) error {
	return nil
}

func (a DatastoreAccessorMock) deleteIdiom(c appengine.Context, idiomID int) error {
	return nil
}

func (a DatastoreAccessorMock) deleteImpl(c appengine.Context, idiomID int, implID int) error {
	return nil
}

// Language filter lang is optional.
func (a DatastoreAccessorMock) searchIdiomsByWords(c appengine.Context, words []string, lang string, limit int) ([]*Idiom, error) {
	return nil, nil
}

func (a DatastoreAccessorMock) searchIdiomsByWordsWithFavorites(c appengine.Context, words []string, favoriteLangs []string, seeNonFavorite bool, limit int) ([]*Idiom, error) {
	return nil, nil
}

func (a DatastoreAccessorMock) processUploadFile(r *http.Request, name string) (string, map[string][]string, error) {
	return "", nil, nil
}

func (a DatastoreAccessorMock) processUploadFiles(r *http.Request, names []string) ([]string, map[string][]string, error) {
	return nil, nil, nil
}

func (a DatastoreAccessorMock) nextIdiomID(c appengine.Context) (int, error) {
	return -1, nil
}

func (a DatastoreAccessorMock) nextImplID(c appengine.Context) (int, error) {
	return -1, nil
}

func (a DatastoreAccessorMock) languagesHavingImpl(c appengine.Context) []string {
	return nil
}

func (a DatastoreAccessorMock) recentIdioms(c appengine.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error) {
	return nil, nil
}

func (a DatastoreAccessorMock) popularIdioms(c appengine.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error) {
	return nil, nil
}

func (a DatastoreAccessorMock) idiomsFilterOrder(c appengine.Context, favoriteLangs []string, limitEachLang int, showOther bool, sortOrder string) ([]*Idiom, error) {
	return nil, nil
}

func (a DatastoreAccessorMock) randomIdiom(c appengine.Context) (key *datastore.Key, idiom *Idiom, thiserror error) {
	return nil, nil, nil
}
