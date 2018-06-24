package pigae

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	gaesearch "google.golang.org/appengine/search"
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

// cheatSheetLineDoc contains some impl data, but only the field Lang is intended to be searched.
type cheatSheetLineDoc struct {
	// Lang is the language of this implementation
	Lang gaesearch.Atom
	// IdiomID is the ID of the idiom this impl belongs to.
	IdiomID gaesearch.Atom
	// IdiomTitle is the title of the idiom this impl belongs to.
	IdiomTitle gaesearch.Atom
	// IdiomLeadParagraph is the description of the idiom this impl belongs to.
	IdiomLeadParagraph gaesearch.Atom
	// ImplID is the ID of this impl.
	ImplID gaesearch.Atom
	// ImplImportsBlock is the imports block of this impl.
	ImplImportsBlock gaesearch.Atom
	// ImplCodeBlock is the code block of this impl.
	ImplCodeBlock gaesearch.Atom
	// ImplCodeBlock is the comment of this impl.
	ImplCodeBlockComment gaesearch.Atom
}

type cheatSheetLineDocs []cheatSheetLineDoc

func (lines cheatSheetLineDocs) Len() int {
	return len(lines)
}
func (lines cheatSheetLineDocs) Swap(i, j int) {
	lines[i], lines[j] = lines[j], lines[i]
}
func (lines cheatSheetLineDocs) Less(i, j int) bool {
	pad := func(a, b gaesearch.Atom) (gaesearch.Atom, gaesearch.Atom) {
		x, y := a, b
		for len(x) < len(y) {
			x = "0" + x
		}
		for len(y) < len(x) {
			y = "0" + y
		}
		return x, y
	}

	// Sort by IdiomID asc, ImplID asc
	idiomID1, idiomID2 := pad(lines[i].IdiomID, lines[j].IdiomID)
	if idiomID1 == idiomID2 {
		implID1, implID2 := pad(lines[i].ImplID, lines[j].ImplID)
		return implID1 < implID2
	}
	return idiomID1 < idiomID2
}

func indexIdiomFullText(c context.Context, idiom *Idiom, idiomKey *datastore.Key) error {
	log.Infof(c, "Reindex text of idiom %d %q", idiom.Id, idiom.Title)
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

	// Also index each impl, so we know what to highlight.
	M := len(idiom.Implementations)
	log.Infof(c, "Reindex %d impls of idiom %d", M, idiom.Id)
	indexImpl, err := gaesearch.Open("impls")
	if err != nil {
		return err
	}
	implDocIDs := make([]string, M)
	implDocs := make([]interface{}, M)
	for i, impl := range idiom.Implementations {
		implDocIDs[i] = fmt.Sprintf("%d_%d", idiom.Id, impl.Id)
		w := impl.ExtractIndexableWords()
		implDocs[i] = &searchableImplDoc{
			Lang:    impl.LanguageName,
			IdiomID: gaesearch.Atom(strconv.Itoa(idiom.Id)),
			Bulk:    strings.Join(w, " "),
		}
		// TODO: how about indexing only last edited impl?
	}
	_, err = indexImpl.PutMulti(c, implDocIDs, implDocs)

	return err
}

func indexIdiomCheatsheets(c context.Context, idiom *Idiom) error {
	log.Infof(c, "Reindex cheatsheet of idiom %d %q", idiom.Id, idiom.Title)
	index, err := gaesearch.Open("cheatsheets")
	if err != nil {
		return err
	}
	// Index each impl.
	M := len(idiom.Implementations)
	docIDs := make([]string, M)
	docs := make([]interface{}, M)
	for i, impl := range idiom.Implementations {
		docIDs[i] = fmt.Sprintf("%d_%d", idiom.Id, impl.Id)
		docs[i] = &cheatSheetLineDoc{
			Lang:                 gaesearch.Atom(impl.LanguageName),
			IdiomID:              gaesearch.Atom(strconv.Itoa(idiom.Id)),
			IdiomTitle:           gaesearch.Atom(idiom.Title),
			IdiomLeadParagraph:   gaesearch.Atom(idiom.LeadParagraph),
			ImplID:               gaesearch.Atom(strconv.Itoa(impl.Id)),
			ImplImportsBlock:     gaesearch.Atom(impl.ImportsBlock),
			ImplCodeBlock:        gaesearch.Atom(impl.CodeBlock),
			ImplCodeBlockComment: gaesearch.Atom(impl.AuthorComment),
		}
	}
	_, err = index.PutMulti(c, docIDs, docs)

	return err
}

func (a *GaeDatastoreAccessor) unindexAll(c context.Context) error {
	log.Infof(c, "Unindexing everything (from the text search indexes)")

	// Must remove 1 by 1 (Index has no batch methods)
	for _, indexName := range []string{
		"idioms",
		"impls",
		"cheatsheets",
	} {
		log.Infof(c, "Unindexing items of [%v]", indexName)
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
				log.Errorf(c, "Error getting next indexed object to unindex: %v", err)
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

func (a *GaeDatastoreAccessor) unindex(c context.Context, idiomID int) error {
	log.Infof(c, "Unindexing idiom %d", idiomID)

	docID := strconv.Itoa(idiomID)
	index, err := gaesearch.Open("idioms")
	if err != nil {
		return err
	}
	return index.Delete(c, docID)
}

// retriever returns a list of Idiom Key strings
type retriever func() ([]string, error)

// searchIdiomsByWordsWithFavorites must return idioms that contain *all* the searched words.
// If seeNonFavorite==false, it must only return idioms that have at least 1 implementation in 1 of the user favoriteLangs.
// If seeNonFavorite==true, it must return the same list but extended with idioms that contain all the searched words but no implementation in a user favoriteLang.
func (a *GaeDatastoreAccessor) searchIdiomsByWordsWithFavorites(c context.Context, typedWords, typedLangs []string, favoriteLangs []string, seeNonFavorite bool, limit int) ([]*Idiom, error) {
	terms := append(append([]string(nil), typedWords...), typedLangs...)

	var retrievers []retriever
	idiomKeyStrings := make([]string, 0, limit)
	seenIdiomKeyStrings := make(map[string]bool, limit)

	var idiomQueryRetriever = func(q string) retriever {
		return func() ([]string, error) {
			return executeIdiomKeyTextSearchQuery(c, q, limit)
		}
	}

	if len(typedLangs) == 1 {
		// Exactly 1 term is a lang: assume user really wants this lang
		lang := typedLangs[0]
		log.Debugf(c, "User is looking for results in [%v]", lang)
		// 1) Impls in lang, containing all words
		implRetriever := func() ([]string, error) {
			var keystrings []string
			implQuery := "Bulk:(~" + strings.Join(terms, " AND ~") + ") AND Lang:" + lang
			implIdiomIDs, _, err := executeImplTextSearchQuery(c, implQuery, limit)
			if err != nil {
				return nil, err
			}
			for _, idiomID := range implIdiomIDs {
				idiomKey := newIdiomKey(c, idiomID)
				idiomKeyString := idiomKey.Encode()
				keystrings = append(keystrings, idiomKeyString)
			}
			return keystrings, nil
		}
		retrievers = []retriever{
			// 1) Idioms with words in title, having an impl in lang
			idiomQueryRetriever("TitleWords:(~" + strings.Join(typedWords, " AND ~") + ") AND Langs:(" + lang + ")"),
			// 2) Implementations in lang, containing all terms
			implRetriever,
			// 3) Idioms with words in lead paragraph (or title), having an impl in lang
			idiomQueryRetriever("TitleOrLeadWords:(~" + strings.Join(typedWords, " AND ~") + ") AND Langs:(" + lang + ")"),
			// 4) Just all the terms
			idiomQueryRetriever("Bulk:(~" + strings.Join(terms, " AND ~") + ")"),
		}

	} else {
		// Either 0 or many langs. Just make sure all terms are respected.
		retrievers = append(retrievers,
			// 1) Words in idiom title, having all the langs implemented
			idiomQueryRetriever("TitleWords:(~"+strings.Join(typedWords, " AND ~")+") AND Bulk:(~"+strings.Join(terms, " AND ~")+")"),
			// 2) Words in idiom lead paragraph (or title), having all the langs implemented
			idiomQueryRetriever("TitleOrLeadWords:(~"+strings.Join(typedWords, " AND ~")+") AND Bulk:(~"+strings.Join(terms, " AND ~")+")"),
			// 3) Terms (words and langs) somewhere in idiom
			idiomQueryRetriever("Bulk:(~"+strings.Join(terms, " AND ~")+")"),
		)
	}

	// Each retriever will send 1 slice in 1 channel. So we can harvest them in right order.
	promises := make([]chan []string, len(retrievers))
	for i := range retrievers {
		retriever := retrievers[i]
		promises[i] = make(chan []string, 1)
		ch := promises[i]
		go func() {
			keyStrings, err := retriever()
			if err != nil {
				log.Errorf(c, "problem fetching search results: %v", err)
				ch <- nil
			} else {
				ch <- keyStrings
			}
			close(ch)
		}()
	}
harvestloop:
	for _, promise := range promises {
		kstrChunk := <-promise
		m := 0
		dupes := 0
		for _, kstr := range kstrChunk {
			if seenIdiomKeyStrings[kstr] {
				dupes++
			} else {
				m++
				idiomKeyStrings = append(idiomKeyStrings, kstr)
				seenIdiomKeyStrings[kstr] = true
				if len(idiomKeyStrings) == limit {
					log.Debugf(c, "%d new results, %d dupes, stopping here.", m, dupes)
					break harvestloop
				}
			}
		}
		log.Debugf(c, "%d new results, %d dupes.", m, dupes)
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

func (a *GaeDatastoreAccessor) searchImplIDs(c context.Context, words, langs []string) (map[string]bool, error) {
	index, err := gaesearch.Open("impls")
	if err != nil {
		return nil, err
	}
	hits := map[string]bool{}
	query := "Bulk:(~" + strings.Join(words, " AND ~") + ")"
	if len(langs) > 0 {
		query += " AND Lang:(" + strings.Join(langs, " OR ") + ")"
	}
	// Beware of INVALID_REQUEST: The limit 1000000 must be between 1 and 1000
	limit := 1000
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

func executeImplTextSearchQuery(c context.Context, query string, limit int) (idiomIDs, implIDs []int, err error) {
	index, err := gaesearch.Open("impls")
	if err != nil {
		return nil, nil, err
	}
	if limit == 0 {
		// Limit is not optional. 0 means zero result.
		return nil, nil, nil
	}
	idiomIDs = make([]int, 0, limit)
	implIDs = make([]int, 0, limit)
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
			return nil, nil, err
		}
		ids := strings.Split(docID, "_")
		idiomID, err := strconv.Atoi(ids[0])
		if err != nil {
			return nil, nil, err
		}
		implID, err := strconv.Atoi(ids[0])
		if err != nil {
			return nil, nil, err
		}
		idiomIDs = append(idiomIDs, idiomID)
		implIDs = append(implIDs, implID)
	}
	return idiomIDs, implIDs, nil
}

func executeIdiomKeyTextSearchQuery(c context.Context, query string, limit int) (keystrings []string, err error) {
	// log.Infof(c, query)
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
	//log.Debugf(c, "Query [%v] yields %d results.", query, len(idiomKeyStrings))
	return idiomKeyStrings, nil
}

func executeIdiomTextSearchQuery(c context.Context, query string, limit int) ([]*Idiom, error) {
	// log.Infof(c, query)
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

func (a *GaeDatastoreAccessor) searchIdiomsByLangs(c context.Context, langs []string, limit int) ([]*Idiom, error) {
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

func (a *GaeDatastoreAccessor) getCheatSheet(c context.Context, lang string, limit int) ([]cheatSheetLineDoc, error) {
	index, err := gaesearch.Open("cheatsheets")
	if err != nil {
		return nil, err
	}
	if limit == 0 {
		// Limit is not optional. 0 means zero result.
		return nil, nil
	}
	cheatLines := make([]cheatSheetLineDoc, 0, 200)
	query := "Lang:" + lang
	it := index.Search(c, query, &gaesearch.SearchOptions{
		Limit: limit,
	})
	for {
		var cheatLine cheatSheetLineDoc
		_, err := it.Next(&cheatLine)
		if err == gaesearch.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		cheatLines = append(cheatLines, cheatLine)
	}
	// Sort by IdiomID asc, ImplID asc
	sort.Sort(cheatSheetLineDocs(cheatLines))
	return cheatLines, err
}
