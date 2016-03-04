package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
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
	getIdiom(c context.Context, idiomID int) (*datastore.Key, *Idiom, error)
	getIdiomByImplID(c context.Context, implID int) (*datastore.Key, *Idiom, error)
	getAllIdioms(c context.Context, limit int, order string) ([]*datastore.Key, []*Idiom, error)
	searchIdiomsByWordsWithFavorites(c context.Context, words, typedLangs []string, favoriteLangs []string, seeNonFavorite bool, limit int) ([]*Idiom, error)
	searchIdiomsByLangs(c context.Context, langs []string, limit int) ([]*Idiom, error)
	searchImplIDs(c context.Context, words, langs []string) (map[string]bool, error)
	recentIdioms(c context.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error)
	popularIdioms(c context.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error)
	randomIdiom(c context.Context) (*datastore.Key, *Idiom, error)
	randomIdiomHaving(c context.Context, havingLang string) (*datastore.Key, *Idiom, error)
	randomIdiomNotHaving(c context.Context, notHavingLang string) (*datastore.Key, *Idiom, error)
	idiomsFilterOrder(c context.Context, favoriteLangs []string, limitEachLang int, showOther bool, sortOrder string) ([]*Idiom, error)
	getIdiomHistory(c context.Context, idiomID int, version int) (*datastore.Key, *IdiomHistory, error)
	getIdiomHistoryList(c context.Context, idiomID int) ([]*datastore.Key, []*IdiomHistory, error)
}

type idiomSaver interface {
	saveNewIdiom(c context.Context, idiom *Idiom) (*datastore.Key, error)
	saveExistingIdiom(c context.Context, key *datastore.Key, idiom *Idiom) error
	deleteAllIdioms(c context.Context) error
	unindexAll(c context.Context) error
	unindex(c context.Context, idiomId int) error
	deleteIdiom(c context.Context, idiomID int, why string) error
	deleteImpl(c context.Context, idiomID int, implID int, why string) error
	nextIdiomID(c context.Context) (int, error)
	nextImplID(c context.Context) (int, error)
	revert(c context.Context, idiomID int, version int) (*Idiom, error)
	historyRestore(c context.Context, idiomID int, version int) (*Idiom, error)
}

type uploadProcesser interface {
	processUploadFile(r *http.Request, name string) (string, map[string][]string, error)
	processUploadFiles(r *http.Request, names []string) ([]string, map[string][]string, error)
}

type appConfigGetter interface {
	getAppConfig(c context.Context) (ApplicationConfig, error)
}

type appConfigSaver interface {
	saveAppConfig(c context.Context, appConfig ApplicationConfig) error
	saveAppConfigProperty(c context.Context, prop AppConfigProperty) error
}

type messenger interface {
	saveNewMessage(c context.Context, msg *MessageForUser) (*datastore.Key, error)
	getMessagesForUser(c context.Context, username string) ([]*datastore.Key, []*MessageForUser, error)
	dismissMessage(c context.Context, key *datastore.Key) (*MessageForUser, error)
}
