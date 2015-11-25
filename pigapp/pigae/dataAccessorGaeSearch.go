package pigae

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
	"appengine/datastore"
	gaesearch "appengine/search"
)

// searchableIdiomDoc is the searchable unit for 1 idiom.
// We keep only some references (id and key) in the indexed "documents", not
// the whole idioms+implementations data.
// By choosing the idiom ID as docID we don't need to retrieve
// the searchableIdiomDoc contents when searching, because IDs suffice
// to constuct datastore Keys.
// See https://cloud.google.com/appengine/docs/go/search/
type searchableIdiomDoc struct {

	// IdiomKeyString is the idiom datastore key string.
	IdiomKeyString gaesearch.Atom

	// IdiomID is the idiom ID.
	IdiomID gaesearch.Atom

	// Bulk is a simple concatenation of (normalized) words, space-separated
	Bulk string

	// Langs is a concatenation of implemented languages, for filtering
	Langs string

	// TitleWords is a concatenation of (normalized) title words, space-separated
	TitleWords string

	// LeadWords is a concatenation of (normalized) lead paragraph words, space-separated
	LeadWords string

	// TitleOrLeadWords is a concatenation of (normalized) idiom description words, space-separated
	TitleOrLeadWords string

	// + displayable data for result list?
	//idiomTitle
	//implLanguages
	//implIDs
	//...
}

// We choose "idiomID_implID" as docID
type searchableImplDoc struct {
	// Lang is the language of this implementation
	Lang string
	// IdiomID is the ID of the idiom this impl belongs to.
	IdiomID gaesearch.Atom
	// Bulk is a simple concatenation of (normalized) words, space-separated
	Bulk string
}

func indexIdiomFullText(c appengine.Context, idiom *Idiom, idiomKey *datastore.Key) error {
	index, err := gaesearch.Open("idioms")
	if err != nil {
		return err
	}
	// By using directly the idiom Key as docID,
	// we can leverage faster ID-only search later.
	docID := strconv.Itoa(idiom.Id)
	w, wTitle, wLead := idiom.ExtractIndexableWords()
	doc := &searchableIdiomDoc{
		IdiomKeyString: gaesearch.Atom(idiomKey.Encode()),
		IdiomID:        gaesearch.Atom(strconv.Itoa(idiom.Id)),
		Bulk:           strings.Join(w, " "),
		Langs:          implementedLanguagesConcat(idiom),
		TitleWords:     strings.Join(wTitle, " "),
		LeadWords:      strings.Join(wLead, " "),
	}
	doc.TitleOrLeadWords = doc.TitleWords + " " + doc.LeadWords
	_, err = index.Put(c, docID, doc)
	if err != nil {
		return err
	}

	// Also index each impl individually,
	// so we know what to highlight.
	indexImpl, err := gaesearch.Open("impls")
	if err != nil {
		return err
	}
	for _, impl := range idiom.Implementations {
		implDocID := fmt.Sprintf("%d_%d", idiom.Id, impl.Id)
		w := impl.ExtractIndexableWords()
		implDoc := &searchableImplDoc{
			Lang:    impl.LanguageName,
			IdiomID: gaesearch.Atom(strconv.Itoa(idiom.Id)),
			Bulk:    strings.Join(w, " "),
		}
		// Weird that the search API doesn't have batch queries.
		// TODO: index each impl concurrently?
		// TODO: index only last edited impl?
		_, err = indexImpl.Put(c, implDocID, implDoc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *GaeDatastoreAccessor) unindexAll(c appengine.Context) error {
	c.Infof("Unindexing everything (from the text search indexes)")

	// Must remove 1 by 1 (Index has no batch methods)
	for _, indexName := range []string{
		"idioms",
		"impls",
	} {
		c.Infof("Unindexing items of [%v]", indexName)
		index, err := gaesearch.Open(indexName)
		if err != nil {
			return err
		}
		it := index.List(c, &gaesearch.ListOptions{IDsOnly: true})
		for {
			docID, err := it.Next(nil)
			if err == gaesearch.Done {
				break
			}
			if err != nil {
				c.Errorf("Error getting next indexed object to unindex: %v", err)
				return err
			}
			err = index.Delete(c, docID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *GaeDatastoreAccessor) unindex(c appengine.Context, idiomID int) error {
	c.Infof("Unindexing idiom %d", idiomID)

	docID := strconv.Itoa(idiomID)
	index, err := gaesearch.Open("idioms")
	if err != nil {
		return err
	}
	return index.Delete(c, docID)
}

// searchIdiomsByWordsWithFavorites must return idioms that contain *all* the searched words.
// If seeNonFavorite==false, it must only return idioms that have at least 1 implementation in 1 of the user favoriteLangs.
// If seeNonFavorite==true, it must return the same list but extended with idioms that contain all the searched words but no implementation in a user favoriteLang.
func (a *GaeDatastoreAccessor) searchIdiomsByWordsWithFavorites(c appengine.Context, typedWords, typedLangs []string, favoriteLangs []string, seeNonFavorite bool, limit int) ([]*Idiom, error) {
	var queries []string
	terms := append(append([]string(nil), typedWords...), typedLangs...)

	if len(typedLangs) == 1 {
		// Exactly 1 term is a lang: assume user really wants this lang
		lang := typedLangs[0]
		queries = []string{
			// 1) Impls in lang, containing all words
			// TODO
			// 2) Idioms with words in title, having an impl in lang
			"TitleWords:(~" + strings.Join(typedWords, " AND ~") + ") AND Langs:(" + lang + ")",
			// 3) Idioms with words in lead paragraph (or title), having an impl in lang
			"TitleOrLeadWords:(~" + strings.Join(typedWords, " AND ~") + ") AND Langs:(" + lang + ")",
			// 4) Just all the terms
			"Bulk:(~" + strings.Join(terms, " AND ~") + ")",
		}
		_ = lang

	} else {
		// Either 0 or many langs. Just make sure all terms are respected.
		queries = []string{
			// 1) Words in idiom title, having all the langs implemented
			"TitleWords:(~" + strings.Join(typedWords, " AND ~") + ") AND Bulk:(~" + strings.Join(terms, " AND ~") + ")",
			// 2) Words in idiom lead paragraph (or title), having all the langs implemented
			"TitleOrLeadWords:(~" + strings.Join(typedWords, " AND ~") + ") AND Bulk:(~" + strings.Join(terms, " AND ~") + ")",
			// 3) Terms (words and langs) somewhere in idiom
			"Bulk:(~" + strings.Join(terms, " AND ~") + ")",
		}
	}

	idiomKeyStrings := make([]string, 0, limit)
	seenIdiomKeyStrings := make(map[string]bool, limit)

queryloop:
	for _, query := range queries {
		// TODO measure that, then parallelize
		keyStrings, err := executeIdiomKeyTextSearchQuery(c, query, limit)
		if err != nil {
			return nil, err
		}
		m := 0
		dupes := 0
		for _, kstr := range keyStrings {
			if seenIdiomKeyStrings[kstr] {
				dupes++
			} else {
				m++
				idiomKeyStrings = append(idiomKeyStrings, kstr)
				seenIdiomKeyStrings[kstr] = true
				if len(idiomKeyStrings) == limit {
					c.Infof("[%v] -> %d new results, %d dupes, stopping here.", query, m, dupes)
					break queryloop
				}
			}
		}
		c.Infof("[%v] -> %d new results, %d dupes.", query, m, dupes)
	}

	// TODO use favoriteLangs
	// TODO use seeNonFavorite (or not)

	var err error
	idiomKeys := make([]*datastore.Key, len(idiomKeyStrings))
	for i, kstr := range idiomKeyStrings {
		idiomKeys[i], err = datastore.DecodeKey(kstr)
		if err != nil {
			return nil, err
		}
	}

	// Fetch Idioms in a []Idiom
	buffer := make([]Idiom, len(idiomKeys))
	err = datastore.GetMulti(c, idiomKeys, buffer)
	// Convert []Idiom to []*Idiom
	idioms := make([]*Idiom, len(buffer))
	for i := range buffer {
		// Do not take the address of the 2nd range variable, it would make a copy.
		// Better take the address in the existing buffer.
		idioms[i] = &buffer[i]
	}
	return idioms, err
}

func (a *GaeDatastoreAccessor) searchImplIDs(c appengine.Context, words []string) (map[string]bool, error) {
	index, err := gaesearch.Open("impls")
	if err != nil {
		return nil, err
	}
	hits := map[string]bool{}
	query := "Bulk:(~" + strings.Join(words, " AND ~") + ")"
	// No real limit for now, those are just IDs and we don't have 1M impls.
	limit := 1000000
	// This is an *IDsOnly* search, where docID == idiomID_implID
	it := index.Search(c, query, &gaesearch.SearchOptions{
		Limit:   limit,
		IDsOnly: true,
	})
	for {
		docID, err := it.Next(nil)
		if err == gaesearch.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		implIDStr := strings.Split(docID, "_")[1]
		hits[implIDStr] = true
	}
	return hits, nil
}

func executeIdiomKeyTextSearchQuery(c appengine.Context, query string, limit int) (keystrings []string, err error) {
	// c.Infof(query)
	index, err := gaesearch.Open("idioms")
	if err != nil {
		return nil, err
	}
	if limit == 0 {
		// Limit is not optional. 0 means zero result.
		return nil, nil
	}
	idiomKeyStrings := make([]string, 0, limit)
	// This is an *IDsOnly* search, where docID == Idiom.Id
	it := index.Search(c, query, &gaesearch.SearchOptions{
		Limit:   limit,
		IDsOnly: true,
	})
	for {
		docID, err := it.Next(nil)
		if err == gaesearch.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		idiomID, err := strconv.Atoi(docID)
		if err != nil {
			return nil, err
		}
		idiomKeyString := newIdiomKey(c, idiomID).Encode()
		idiomKeyStrings = append(idiomKeyStrings, idiomKeyString)
	}
	return idiomKeyStrings, nil
}

func executeIdiomTextSearchQuery(c appengine.Context, query string, limit int) ([]*Idiom, error) {
	// c.Infof(query)
	index, err := gaesearch.Open("idioms")
	if err != nil {
		return nil, err
	}
	if limit == 0 {
		// Limit is not optional. 0 means zero result.
		return nil, nil
	}
	idiomKeys := make([]*datastore.Key, 0, limit)
	// This is an *IDsOnly* search, where docID == Idiom.Id
	it := index.Search(c, query, &gaesearch.SearchOptions{
		Limit:   limit,
		IDsOnly: true,
	})
	for {
		docID, err := it.Next(nil)
		if err == gaesearch.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		idiomID, err := strconv.Atoi(docID)
		if err != nil {
			return nil, err
		}
		key := newIdiomKey(c, idiomID)
		idiomKeys = append(idiomKeys, key)
	}
	// Fetch Idioms in a []Idiom
	buffer := make([]Idiom, len(idiomKeys))
	err = datastore.GetMulti(c, idiomKeys, buffer)
	// Convert []Idiom to []*Idiom
	idioms := make([]*Idiom, len(buffer))
	for i := range buffer {
		// Do not take the address of the 2nd range variable, it would make a copy.
		// Better take the address in the existing buffer.
		idioms[i] = &buffer[i]
	}
	return idioms, err
}

func (a *GaeDatastoreAccessor) searchIdiomsByLangs(c appengine.Context, langs []string, limit int) ([]*Idiom, error) {
	dsq := datastore.NewQuery("Idiom")
	dsq = dsq.Filter("Implementations.LanguageName = ", langs[0])
	if len(langs) >= 2 {
		return nil, fmt.Errorf("Not yet implemented: list for more than 1 language")
	}
	dsq = dsq.Order("-Rating").Limit(limit)
	hits := make([]*Idiom, 0, 10)
	if _, err := dsq.GetAll(c, &hits); err != nil {
		return nil, err
	}
	return hits, nil
}
