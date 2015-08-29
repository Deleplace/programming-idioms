package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
	"appengine/datastore"
)

type dataAccessor interface {
	idiomGetter
	idiomSaver
	uploadProcesser
	appConfigGetter
	appConfigSaver
}

type idiomGetter interface {
	getIdiom(c appengine.Context, idiomID int) (*datastore.Key, *Idiom, error)
	getIdiomByImplID(c appengine.Context, implID int) (*datastore.Key, *Idiom, error)
	getAllIdioms(c appengine.Context, limit int, order string) ([]*datastore.Key, []*Idiom, error)
	searchIdiomsByWords(c appengine.Context, words []string, lang string, limit int) ([]*Idiom, error)
	searchIdiomsByWordsWithFavorites(c appengine.Context, words []string, favoriteLangs []string, seeNonFavorite bool, limit int) ([]*Idiom, error)
	searchIdiomsByLangs(c appengine.Context, langs []string, limit int) ([]*Idiom, error)
	recentIdioms(c appengine.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error)
	popularIdioms(c appengine.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error)
	randomIdiom(c appengine.Context) (*datastore.Key, *Idiom, error)
	randomIdiomHaving(c appengine.Context, havingLang string) (*datastore.Key, *Idiom, error)
	randomIdiomNotHaving(c appengine.Context, notHavingLang string) (*datastore.Key, *Idiom, error)
	idiomsFilterOrder(c appengine.Context, favoriteLangs []string, limitEachLang int, showOther bool, sortOrder string) ([]*Idiom, error)
	languagesHavingImpl(c appengine.Context) []string
	getIdiomHistory(c appengine.Context, idiomID int, version int) (*datastore.Key, *IdiomHistory, error)
}

type idiomSaver interface {
	saveNewIdiom(c appengine.Context, idiom *Idiom) (*datastore.Key, error)
	saveExistingIdiom(c appengine.Context, key *datastore.Key, idiom *Idiom) error
	deleteAllIdioms(c appengine.Context) error
	deleteIdiom(c appengine.Context, idiomID int) error
	deleteImpl(c appengine.Context, idiomID int, implID int) error
	nextIdiomID(c appengine.Context) (int, error)
	nextImplID(c appengine.Context) (int, error)
}

type uploadProcesser interface {
	processUploadFile(r *http.Request, name string) (string, map[string][]string, error)
	processUploadFiles(r *http.Request, names []string) ([]string, map[string][]string, error)
}

type appConfigGetter interface {
	getAppConfig(c appengine.Context) (ApplicationConfig, error)
}

type appConfigSaver interface {
	saveAppConfig(c appengine.Context, appConfig ApplicationConfig) error
}
