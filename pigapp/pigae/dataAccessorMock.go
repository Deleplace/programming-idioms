package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// DatastoreAccessorMock should be useful for unit testing. But is is dead code for now.
type DatastoreAccessorMock struct {
	idioms []*Idiom
}

func (a DatastoreAccessorMock) getIdiom(c context.Context, idiomID int) (key *datastore.Key, idiom *Idiom, thiserror error) {
	return
}

func (a DatastoreAccessorMock) saveNewIdiom(c context.Context, idiom *Idiom) (key *datastore.Key, thiserror error) {
	return
}

func (a DatastoreAccessorMock) saveExistingIdiom(c context.Context, key *datastore.Key, idiom *Idiom) error {
	return nil
}

func (a DatastoreAccessorMock) getAllIdioms(c context.Context, limit int, order string) (keys []*datastore.Key, idioms []*Idiom, thiserror error) {
	return
}

func (a DatastoreAccessorMock) deleteAllIdioms(c context.Context) error {
	return nil
}

func (a DatastoreAccessorMock) deleteIdiom(c context.Context, idiomID int, why string) error {
	return nil
}

func (a DatastoreAccessorMock) deleteImpl(c context.Context, idiomID int, implID int, why string) error {
	return nil
}

// Language filter lang is optional.
func (a DatastoreAccessorMock) searchIdiomsByWords(c context.Context, words []string, lang string, limit int) ([]*Idiom, error) {
	return nil, nil
}

func (a DatastoreAccessorMock) searchIdiomsByWordsWithFavorites(c context.Context, words []string, favoriteLangs []string, seeNonFavorite bool, limit int) ([]*Idiom, error) {
	return nil, nil
}

func (a DatastoreAccessorMock) processUploadFile(r *http.Request, name string) (string, map[string][]string, error) {
	return "", nil, nil
}

func (a DatastoreAccessorMock) processUploadFiles(r *http.Request, names []string) ([]string, map[string][]string, error) {
	return nil, nil, nil
}

func (a DatastoreAccessorMock) nextIdiomID(c context.Context) (int, error) {
	return -1, nil
}

func (a DatastoreAccessorMock) nextImplID(c context.Context) (int, error) {
	return -1, nil
}

func (a DatastoreAccessorMock) languagesHavingImpl(c context.Context) []string {
	return nil
}

func (a DatastoreAccessorMock) recentIdioms(c context.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error) {
	return nil, nil
}

func (a DatastoreAccessorMock) popularIdioms(c context.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error) {
	return nil, nil
}

func (a DatastoreAccessorMock) idiomsFilterOrder(c context.Context, favoriteLangs []string, limitEachLang int, showOther bool, sortOrder string) ([]*Idiom, error) {
	return nil, nil
}

func (a DatastoreAccessorMock) randomIdiom(c context.Context) (key *datastore.Key, idiom *Idiom, thiserror error) {
	return nil, nil, nil
}
