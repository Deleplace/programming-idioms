package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"context"

	"google.golang.org/appengine"
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
	// ImplDemoURL is the demo URL of this impl.
	ImplDemoURL gaesearch.Atom
	// ImplDocURL is the documentation URL of this impl.
	ImplDocURL gaesearch.Atom
	// Protected when "only admin can edit"
	Protected gaesearch.Atom
	// IdiomVersion is the number of this idiom revision
	IdiomVersion gaesearch.Atom
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

func indexIdiomFullText(ctx context.Context, idiom *Idiom, idiomKey *datastore.Key) error {
	log.Infof(ctx, "Reindex text of idiom %d %q", idiom.Id, idiom.Title)
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
	_, err = index.Put(ctx, docID, doc)
	if err != nil {
		return err
	}

	// Also index each impl, so we know what to highlight.
	M := len(idiom.Implementations)
	log.Infof(ctx, "Reindex %d impls of idiom %d", M, idiom.Id)
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
	_, err = indexImpl.PutMulti(ctx, implDocIDs, implDocs)

	return err
}

func indexIdiomCheatsheets(ctx context.Context, idiom *Idiom) error {
	log.Infof(ctx, "Reindex cheatsheet of idiom %d %q", idiom.Id, idiom.Title)
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
			ImplDocURL:           gaesearch.Atom(impl.DocumentationURL),
			ImplDemoURL:          gaesearch.Atom(impl.DemoURL),
			IdiomVersion:         gaesearch.Atom(strconv.Itoa(idiom.Version)),
			Protected:            gaesearch.Atom(strconv.FormatBool(idiom.Protected || impl.Protected)),
		}
	}
	_, err = index.PutMulti(ctx, docIDs, docs)

	if err != nil {
		if multierr, ok := err.(appengine.MultiError); ok {
			log.Warningf(ctx, "PutMulti returned %d errors", len(multierr))
			for i, singleerr := range multierr {
				log.Warningf(ctx, "  error %d: %v", i, singleerr)
			}
		} else {
			log.Warningf(ctx, "Can't convert PutMulti error into []error")
		}
	}

	return err
}

func (a *GaeDatastoreAccessor) unindexAll(ctx context.Context) error {
	log.Infof(ctx, "Unindexing everything (from the text search indexes)")
	start := time.Now()

	for _, indexName := range []string{
		"idioms",
		"impls",
		"cheatsheets",
	} {
		log.Infof(ctx, "Unindexing items of [%v]", indexName)
		index, err := gaesearch.Open(indexName)
		if err != nil {
			return err
		}
		it := index.List(ctx, &gaesearch.ListOptions{IDsOnly: true})
		docIDs := make([]string, 0, 100)
		for {
			docID, err := it.Next(nil)
			if err == gaesearch.Done {
				break
			}
			if err != nil {
				log.Errorf(ctx, "Error getting next indexed object to unindex: %v", err)
				return err
			}
			docIDs = append(docIDs, docID)
			if len(docIDs) == 100 {
				// It seems that we can't make a huge batch call, hard limit is ~200
				err = index.DeleteMulti(ctx, docIDs)
				if err != nil {
					return err
				}
				docIDs = docIDs[:0]
			}
		}
		err = index.DeleteMulti(ctx, docIDs)
		if err != nil {
			return err
		}
	}

	log.Infof(ctx, "Unindexed everything in %v", time.Since(start))
	return nil
}

func (a *GaeDatastoreAccessor) unindex(ctx context.Context, idiom *Idiom) error {
	log.Infof(ctx, "Unindexing idiom %d", idiom.Id)

	docID := strconv.Itoa(idiom.Id)
	index, err := gaesearch.Open("idioms")
	if err != nil {
		return err
	}
	err = index.Delete(ctx, docID)
	if err != nil {
		return err
	}

	log.Infof(ctx, "Unindexing %d implementations from idiom %d", len(idiom.Implementations), idiom.Id)
	for _, impl := range idiom.Implementations {
		log.Infof(ctx, "Unindexing idiom %d impl %d", idiom.Id, impl.Id)
		err := unindexImpl(ctx, idiom.Id, impl.Id)
		if err != nil {
			log.Errorf(ctx, "Unindexing idiom %d impl %d: %v", idiom.Id, impl.Id, err)
			// Keep going though
		}
	}
	return nil
}

func unindexImpl(ctx context.Context, idiomID, implID int) error {
	var err error
	for _, indexName := range []string{
		"impls",
		"cheatsheets",
	} {
		index, err := gaesearch.Open(indexName)
		if err != nil {
			return err
		}
		docID := fmt.Sprintf("%d_%d", idiomID, implID)
		err2 := index.Delete(ctx, docID)
		if err2 != nil {
			err = err2
		}
	}
	return err

	// Index "idioms":
	// Reindexing of the Idiom itself from index "idioms", doc "id", is handled elsewhere,
	// async via indexDelayer.
}

// retriever returns a list of Idiom Key strings
type retriever func() ([]string, error)

// searchIdiomsByWordsWithFavorites must return idioms that contain *all* the searched words.
// If seeNonFavorite==false, it must only return idioms that have at least 1 implementation in 1 of the user favoriteLangs.
// If seeNonFavorite==true, it must return the same list but extended with idioms that contain all the searched words but no implementation in a user favoriteLang.
func (a *GaeDatastoreAccessor) searchIdiomsByWordsWithFavorites(ctx context.Context, typedWords, typedLangs []string, favoriteLangs []string, seeNonFavorite bool, limit int) ([]*Idiom, error) {
	terms := append(append([]string(nil), typedWords...), typedLangs...)

	var retrievers []retriever
	idiomKeyStrings := make([]string, 0, limit)
	seenIdiomKeyStrings := make(map[string]bool, limit)

	var idiomQueryRetriever = func(q string) retriever {
		return func() ([]string, error) {
			return executeIdiomKeyTextSearchQuery(ctx, q, limit)
		}
	}

	if len(typedLangs) == 1 {
		// Exactly 1 term is a lang: assume user really wants this lang
		lang := typedLangs[0]
		log.Debugf(ctx, "User is looking for results in [%v]", lang)
		// 1) Impls in lang, containing all words
		implRetriever := func() ([]string, error) {
			var keystrings []string
			implQuery := "Bulk:(~" + strings.Join(terms, " AND ~") + ") AND Lang:" + lang
			implIdiomIDs, _, err := executeImplTextSearchQuery(ctx, implQuery, limit)
			if err != nil {
				return nil, err
			}
			for _, idiomID := range implIdiomIDs {
				idiomKey := newIdiomKey(ctx, idiomID)
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
				log.Errorf(ctx, "problem fetching search results: %v", err)
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
					log.Debugf(ctx, "%d new results, %d dupes, stopping here.", m, dupes)
					break harvestloop
				}
			}
		}
		log.Debugf(ctx, "%d new results, %d dupes.", m, dupes)
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
	err = datastore.GetMulti(ctx, idiomKeys, buffer)
	// Convert []Idiom to []*Idiom
	idioms := make([]*Idiom, len(buffer))
	for i := range buffer {
		// Do not take the address of the 2nd range variable, it would make a copy.
		// Better take the address in the existing buffer.
		idioms[i] = &buffer[i]
	}
	return idioms, err
}

func (a *GaeDatastoreAccessor) searchImplIDs(ctx context.Context, words, langs []string) (map[string]bool, error) {
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
	it := index.Search(ctx, query, &gaesearch.SearchOptions{
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

func executeImplTextSearchQuery(ctx context.Context, query string, limit int) (idiomIDs, implIDs []int, err error) {
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
	it := index.Search(ctx, query, &gaesearch.SearchOptions{
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

func executeIdiomKeyTextSearchQuery(ctx context.Context, query string, limit int) (keystrings []string, err error) {
	// log.Infof(ctx, query)
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
	it := index.Search(ctx, query, &gaesearch.SearchOptions{
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
		idiomKeyString := newIdiomKey(ctx, idiomID).Encode()
		idiomKeyStrings = append(idiomKeyStrings, idiomKeyString)
	}
	//log.Debugf(ctx, "Query [%v] yields %d results.", query, len(idiomKeyStrings))
	return idiomKeyStrings, nil
}

func executeIdiomTextSearchQuery(ctx context.Context, query string, limit int) ([]*Idiom, error) {
	// log.Infof(ctx, query)
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
	it := index.Search(ctx, query, &gaesearch.SearchOptions{
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
		key := newIdiomKey(ctx, idiomID)
		idiomKeys = append(idiomKeys, key)
	}
	// Fetch Idioms in a []Idiom
	buffer := make([]Idiom, len(idiomKeys))
	err = datastore.GetMulti(ctx, idiomKeys, buffer)
	// Convert []Idiom to []*Idiom
	idioms := make([]*Idiom, len(buffer))
	for i := range buffer {
		// Do not take the address of the 2nd range variable, it would make a copy.
		// Better take the address in the existing buffer.
		idioms[i] = &buffer[i]
	}
	return idioms, err
}

func (a *GaeDatastoreAccessor) searchIdiomsByLangs(ctx context.Context, langs []string, limit int) ([]*Idiom, error) {
	dsq := datastore.NewQuery("Idiom")
	dsq = dsq.Filter("Implementations.LanguageName = ", langs[0])
	if len(langs) >= 2 {
		return nil, fmt.Errorf("Not yet implemented: list for more than 1 language")
	}
	dsq = dsq.Order("-Rating").Limit(limit)
	hits := make([]*Idiom, 0, 10)
	if _, err := dsq.GetAll(ctx, &hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (a *GaeDatastoreAccessor) getCheatSheet(ctx context.Context, lang string, limit int) ([]cheatSheetLineDoc, error) {
	index, err := gaesearch.Open("cheatsheets")
	if err != nil {
		return nil, err
	}
	if limit == 0 {
		// Limit is not optional. 0 means zero result.
		return nil, nil
	}
	cheatLines := make([]cheatSheetLineDoc, 0, 500)
	query := "Lang:" + lang
	it := index.Search(ctx, query, &gaesearch.SearchOptions{
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

func searchRandomImplsForLang(ctx context.Context, lang string, n int) ([]*IdiomSingleton, error) {
	index, err := gaesearch.Open("cheatsheets")
	if err != nil {
		return nil, err
	}
	query := "Lang:" + lang
	limit := 1000
	// This is an *IDsOnly* search, where docID == idiomID_implID
	it := index.Search(ctx, query, &gaesearch.SearchOptions{
		Limit: limit,
		//IDsOnly: true,
	})
	var results []*IdiomSingleton
	for {
		var doc cheatSheetLineDoc
		_, err := it.Next(&doc)
		if err == gaesearch.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		singleton := &IdiomSingleton{
			Id:            String2Int(string(doc.IdiomID)),
			Title:         string(doc.IdiomTitle),
			LeadParagraph: string(doc.IdiomLeadParagraph),
			Implementations: []Impl{
				{
					Id:               String2Int(string(doc.ImplID)),
					CodeBlock:        string(doc.ImplCodeBlock),
					DemoURL:          string(doc.ImplDemoURL),
					DocumentationURL: string(doc.ImplDocURL),
					AuthorComment:    string(doc.ImplCodeBlockComment),
				},
			},
			Version: String2Int(string(doc.IdiomVersion)),
			// All other fields left blank for Backlog display
			// All other impls ignored for Backlog display of 1 impl
		}
		results = append(results, singleton)
	}
	log.Debugf(ctx, "searchRandomImplsForLang handling %d %s impls", len(results), lang)
	rand.Shuffle(len(results), func(i, j int) {
		results[i], results[j] = results[j], results[i]
	})
	if len(results) < n {
		return results, nil
	}
	return results[:n], nil
}

type backlogMissingDocDemo struct {
	// Lang is the language these impls.
	Lang string

	// MissingDoc holds a small subset of the Implementations in Lang that don't have a
	// DocumentationURL.
	// They are idiom singletons: only Id, Title, LeadParagraph, and one single implementation.
	MissingDoc []*IdiomSingleton

	// MissingDemo holds a small subset of the Implementations in Lang that don't have a
	// DemoURL.
	// They are idiom singletons: only Id, Title, LeadParagraph, and one single implementation.
	MissingDemo []*IdiomSingleton

	Stats struct {
		// CountImplsMissingDoc is the total database number of Implementations in Lang that
		// don't have a DocumentationURL.
		CountImplsMissingDoc int

		// CountImplsMissingDemo is the total database number of Implementations in Lang that
		// don't have a DemoURL.
		CountImplsMissingDemo int

		// CountImplsLangTotal is the total number of Impls for this Lang, in the database.
		CountImplsLangTotal int
	}
}

func (bmdd backlogMissingDocDemo) MissingDocRatio() float64 {
	return float64(bmdd.Stats.CountImplsMissingDoc) / float64(bmdd.Stats.CountImplsLangTotal)
}

func (bmdd backlogMissingDocDemo) MissingDemoRatio() float64 {
	return float64(bmdd.Stats.CountImplsMissingDemo) / float64(bmdd.Stats.CountImplsLangTotal)
}

func (bmdd backlogMissingDocDemo) MissingDocPercent() string {
	ratio := bmdd.MissingDocRatio()
	return fmt.Sprintf("%.0f%%", 100*ratio)
}

func (bmdd backlogMissingDocDemo) MissingDemoPercent() string {
	ratio := bmdd.MissingDemoRatio()
	return fmt.Sprintf("%.0f%%", 100*ratio)
}

func (bmdd backlogMissingDocDemo) HavingDocPercent() string {
	ratio := 1 - bmdd.MissingDocRatio()
	return fmt.Sprintf("%.0f%%", 100*ratio)
}

func (bmdd backlogMissingDocDemo) HavingDemoPercent() string {
	ratio := 1 - bmdd.MissingDemoRatio()
	return fmt.Sprintf("%.0f%%", 100*ratio)
}

func searchMissingDocDemoForLang(ctx context.Context, lang string, n int) (bmdd backlogMissingDocDemo, err error) {
	bmdd.Lang = lang
	index, err := gaesearch.Open("cheatsheets")
	if err != nil {
		return bmdd, err
	}
	query := "Lang:" + lang
	limit := 1000 // TODO 1000 here is... somewhat arbitrary
	it := index.Search(ctx, query, &gaesearch.SearchOptions{
		Limit: limit,
	})
	for {
		var doc cheatSheetLineDoc
		_, err := it.Next(&doc)
		if err == gaesearch.Done {
			break
		}
		if err != nil {
			return bmdd, err
		}
		bmdd.Stats.CountImplsLangTotal++
		if doc.IdiomID == "149" {
			// "Rescue the princess" is not open to contributions.
			continue
			// Should also be caught by the Protected test below
		}
		if doc.Protected == "true" {
			// Some idioms e.g. #149 "Rescue the princess" are not open to contributions.
			log.Infof(ctx, "Skipping protected idiom %s impl %s", doc.IdiomID, doc.ImplID)
			continue
		}
		singleton := &IdiomSingleton{
			Id:            String2Int(string(doc.IdiomID)),
			Title:         string(doc.IdiomTitle),
			LeadParagraph: string(doc.IdiomLeadParagraph),
			Implementations: []Impl{
				{
					Id:               String2Int(string(doc.ImplID)),
					CodeBlock:        string(doc.ImplCodeBlock),
					DemoURL:          string(doc.ImplDemoURL),
					DocumentationURL: string(doc.ImplDocURL),
					AuthorComment:    string(doc.ImplCodeBlockComment),
				},
			},
			// All other fields left blank for Backlog display
			// All other impls ignored for Backlog display of 1 impl
		}
		if strings.TrimSpace(singleton.Implementations[0].DocumentationURL) == "" {
			bmdd.MissingDoc = append(bmdd.MissingDoc, singleton)
		}
		if strings.TrimSpace(singleton.Implementations[0].DemoURL) == "" {
			bmdd.MissingDemo = append(bmdd.MissingDemo, singleton)
		}
	}
	rand.Shuffle(len(bmdd.MissingDoc), func(i, j int) {
		bmdd.MissingDoc[i], bmdd.MissingDoc[j] = bmdd.MissingDoc[j], bmdd.MissingDoc[i]
	})
	bmdd.Stats.CountImplsMissingDoc = len(bmdd.MissingDoc)
	if len(bmdd.MissingDoc) > n {
		bmdd.MissingDoc = bmdd.MissingDoc[:n]
	}
	rand.Shuffle(len(bmdd.MissingDemo), func(i, j int) {
		bmdd.MissingDemo[i], bmdd.MissingDemo[j] = bmdd.MissingDemo[j], bmdd.MissingDemo[i]
	})
	bmdd.Stats.CountImplsMissingDemo = len(bmdd.MissingDemo)
	if len(bmdd.MissingDemo) > n {
		bmdd.MissingDemo = bmdd.MissingDemo[:n]
	}
	// Note we resliced to a subset of n, but we're not releasing the extra memory at all.
	// TODO copy in new small slices, let the GC reclaim the rest.
	return bmdd, nil
}

type backlogMissingImpl struct {
	// Lang is the language for which the current struct shows missing impls.
	Lang string

	// Stubs holds a small subset of the Idioms that have zero impls in Lang.
	// They are stubs: only Id, Title, LeadParagraph. Does not implementations.
	Stubs []*IdiomStub

	Stats struct {
		// CountIdiomsMissingImpl is the total database number of Idioms having zero.
		// impls in Lang
		CountIdiomsMissingImpl int

		// CountIdiomsTotal is the total number of Idioms in the database.
		CountIdiomsTotal int
	}
}

func (bmi backlogMissingImpl) MissingImplRatio() float64 {
	return float64(bmi.Stats.CountIdiomsMissingImpl) / float64(bmi.Stats.CountIdiomsTotal)
}

func (bmi backlogMissingImpl) MissingImplPercent() string {
	ratio := bmi.MissingImplRatio()
	return fmt.Sprintf("%.0f%%", 100*ratio)
}

func (bmi backlogMissingImpl) HavingImplPercent() string {
	ratio := 1 - bmi.MissingImplRatio()
	return fmt.Sprintf("%.0f%%", 100*ratio)
}

func searchMissingImplForLang(ctx context.Context, lang string, n int) (bmi backlogMissingImpl, err error) {
	bmi.Lang = lang

	// Indirect way of finding out idioms where lang is missing:
	// 1) Get the set of all idiomIDs having at least 1 impl in lang.
	// 2) Iterate through all idiomIDs, keep only those not in the set.
	// 3) Shuffle.
	// 4) Keep n first idioms.

	// 1)
	haveLang := map[string]bool{}
	index, err := gaesearch.Open("cheatsheets")
	if err != nil {
		return bmi, err
	}
	query := "Lang:" + lang
	limit := 1000
	it := index.Search(ctx, query, &gaesearch.SearchOptions{
		Limit: limit,
	})
	for {
		var doc cheatSheetLineDoc
		_, err := it.Next(&doc)
		if err == gaesearch.Done {
			break
		}
		if err != nil {
			return bmi, err
		}
		haveLang[string(doc.IdiomID)] = true
	}
	//log.Infof(ctx, "Found idioms having %s: %v", lang, haveLang)

	// 2)
	idiomIDsWithoutLang := make([]string, 0, 300)
	seen := map[string]bool{}
	index, err = gaesearch.Open("idioms")
	if err != nil {
		return bmi, err
	}
	// This is an *IDsOnly* search, where docID == Idiom.Id
	it = index.List(ctx, &gaesearch.ListOptions{
		IDsOnly: true,
	})
	for {
		docID, err := it.Next(nil)
		if err == gaesearch.Done {
			break
		}
		if err != nil {
			return bmi, err
		}
		idiomID := docID
		if seen[idiomID] {
			continue
		}
		seen[idiomID] = true
		if haveLang[idiomID] {
			continue
		}
		if idiomID == "149" {
			// "Rescue the princess" is not open to contributions.
			continue
			// TODO find a way to skip Protected idioms, in general.
		}
		idiomIDsWithoutLang = append(idiomIDsWithoutLang, idiomID)
	}
	//log.Infof(ctx, "Found idioms without %s: %v", lang, idiomIDsWithoutLang)
	bmi.Stats.CountIdiomsMissingImpl = len(idiomIDsWithoutLang)
	bmi.Stats.CountIdiomsTotal = len(seen)

	// 3)
	rand.Shuffle(len(idiomIDsWithoutLang), func(i, j int) {
		idiomIDsWithoutLang[i], idiomIDsWithoutLang[j] = idiomIDsWithoutLang[j], idiomIDsWithoutLang[i]
	})

	// 4)
	if len(idiomIDsWithoutLang) > n {
		idiomIDsWithoutLang = idiomIDsWithoutLang[:n]
	}

	// 5) Collect idiom stub data to be displayed
	bmi.Stubs = make([]*IdiomStub, len(idiomIDsWithoutLang))
	index, err = gaesearch.Open("cheatsheets")
	if err != nil {
		return bmi, err
	}
	for i, idiomID := range idiomIDsWithoutLang {
		query := "IdiomID:" + idiomID
		it := index.Search(ctx, query, &gaesearch.SearchOptions{
			Limit: 1,
		})
		var doc cheatSheetLineDoc
		_, err := it.Next(&doc)
		if err == gaesearch.Done {
			log.Errorf(ctx, "couldn't find any cheatsheet data for idiom %s ??", idiomID)
			continue
		}
		if err != nil {
			return bmi, err
		}
		bmi.Stubs[i] = &IdiomStub{
			Id:            String2Int(string(doc.IdiomID)),
			Title:         string(doc.IdiomTitle),
			LeadParagraph: string(doc.IdiomLeadParagraph),
		}
	}
	log.Infof(ctx, "Found idioms without %s: %v", lang, idiomIDsWithoutLang)
	return bmi, nil
}
