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
	messenger
}

type idiomGetter interface {
	getIdiom(c appengine.Context, idiomID int) (*datastore.Key, *Idiom, error)
	getIdiomByImplID(c appengine.Context, implID int) (*datastore.Key, *Idiom, error)
	getAllIdioms(c appengine.Context, limit int, order string) ([]*datastore.Key, []*Idiom, error)
	searchIdiomsByWordsWithFavorites(c appengine.Context, words, typedLangs []string, favoriteLangs []string, seeNonFavorite bool, limit int) ([]*Idiom, error)
	searchIdiomsByLangs(c appengine.Context, langs []string, limit int) ([]*Idiom, error)
	searchImplIDs(c appengine.Context, words, langs []string) (map[string]bool, error)
	recentIdioms(c appengine.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error)
	popularIdioms(c appengine.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error)
	randomIdiom(c appengine.Context) (*datastore.Key, *Idiom, error)
	randomIdiomHaving(c appengine.Context, havingLang string) (*datastore.Key, *Idiom, error)
	randomIdiomNotHaving(c appengine.Context, notHavingLang string) (*datastore.Key, *Idiom, error)
	idiomsFilterOrder(c appengine.Context, favoriteLangs []string, limitEachLang int, showOther bool, sortOrder string) ([]*Idiom, error)
	languagesHavingImpl(c appengine.Context) []string
	getIdiomHistory(c appengine.Context, idiomID int, version int) (*datastore.Key, *IdiomHistory, error)
	getIdiomHistoryList(c appengine.Context, idiomID int) ([]*datastore.Key, []*IdiomHistory, error)
}

type idiomSaver interface {
	saveNewIdiom(c appengine.Context, idiom *Idiom) (*datastore.Key, error)
	saveExistingIdiom(c appengine.Context, key *datastore.Key, idiom *Idiom) error
	deleteAllIdioms(c appengine.Context) error
	unindexAll(c appengine.Context) error
	unindex(c appengine.Context, idiomId int) error
	deleteIdiom(c appengine.Context, idiomID int, why string) error
	deleteImpl(c appengine.Context, idiomID int, implID int, why string) error
	nextIdiomID(c appengine.Context) (int, error)
	nextImplID(c appengine.Context) (int, error)
	revert(c appengine.Context, idiomID int, version int) (*Idiom, error)
	historyRestore(c appengine.Context, idiomID int, version int) (*Idiom, error)
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
	saveAppConfigProperty(c appengine.Context, prop AppConfigProperty) error
}

type messenger interface {
	saveNewMessage(c appengine.Context, msg *MessageForUser) (*datastore.Key, error)
	getMessagesForUser(c appengine.Context, username string) ([]*datastore.Key, []*MessageForUser, error)
	dismissMessage(c appengine.Context, key *datastore.Key) (*MessageForUser, error)
}
