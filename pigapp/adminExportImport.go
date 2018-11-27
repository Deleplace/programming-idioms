package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
)

func adminExport(w http.ResponseWriter, r *http.Request) error {
	format := "json" // TODO read FormValue

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/octet-stream")
		d := time.Now().Format("2006-01-02_15-04")
		w.Header().Set("Content-Disposition", "attachment; filename=\"programming-idioms.org."+d+".json\"")
		return exportIdiomsAsJSON(r, w, true)
	default:
		return errors.New("Not implemented: " + format)
	}

}

func adminImportAjax(w http.ResponseWriter, r *http.Request) error {
	c := r.Context()
	file, fileHeader, err := r.FormFile("importData")
	if err != nil {
		return err
	}
	// TODO import in 1 transaction
	// unless 6+ entity groups in 1 transaction is impossible
	if purge := r.FormValue("purge"); purge != "" {
		err = dao.deleteAllIdioms(c)
		if err != nil {
			return err
		}
	}
	_ = dao.deleteCache(c)
	count, err := importFile(c, file, fileHeader)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"imported": count})
	return nil
}

func importFile(c context.Context, file multipart.File, fileHeader *multipart.FileHeader) (int, error) {
	chunks := strings.Split(fileHeader.Filename, ".")
	extension := Last(chunks)
	var err error
	var idioms []*Idiom
	switch strings.ToLower(extension) {
	case "json":
		idioms, err = importFromJSON(file)
	case "csv":
		idioms, err = importFromCSV(file)
	default:
		return 0, fmt.Errorf("Unknown extension [%v]", extension)
	}
	if err != nil {
		return 0, err
	}
	n := 0
	for _, idiom := range idioms {
		if fixNewlines(idiom) {
			infof(c, "Fixed newlines in idiom #%d", idiom.Id)
		}
		if _, err = dao.saveNewIdiom(c, idiom); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// fixNewlines replaces "\r\n" with "\n", because expected newlines
// are 1 char, and having 2 chars leads to
// "API error 1 (datastore_v3: BAD_REQUEST): Property Implementations.CodeBlock is too long. Maximum length is 500.
func fixNewlines(idiom *Idiom) bool {
	touched := false
	for i := range idiom.Implementations {
		impl := &idiom.Implementations[i]
		if strings.Contains(impl.CodeBlock, "\r\n") {
			touched = true
			impl.CodeBlock = strings.Replace(impl.CodeBlock, "\r\n", "\n", -1)
		}
	}
	return touched
}

func importFromJSON(file multipart.File) ([]*Idiom, error) {
	idioms := []*Idiom{}
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&idioms)
	if err != nil {
		return nil, err
	}
	return idioms, nil
}

func importFromCSV(file multipart.File) ([]*Idiom, error) {
	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = false
	reader.TrailingComma = true
	cells, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	languages := []string{}

	headers := cells[0]
	for i, label := range headers {
		if i < 3 {
			// Id;	Title;	Description -> dummy string xxx
			languages = append(languages, "xxx")
			continue
		}
		//  aaa, aaa_comment, bbb, bbb_comment, etc.
		languages = append(languages, label)
	}

	idioms := []*Idiom{}
	// this implID works only after a purge...
	implID := 1
	for i, line := range cells {
		if i == 0 {
			// Headers
			continue
		}

		idiomID := String2Int(line[0])

		idiom := Idiom{
			Id:            idiomID,
			Title:         line[1],
			LeadParagraph: line[2],
			Author:        "programming-idioms.org",
			Version:       1,
		}

		cell := func(line []string, j int) string {
			if j >= len(line) {
				return ""
			}
			return line[j]
		}

		impls := []Impl{}
		for j := 3; j < len(line); j += 2 {
			code := cell(line, j)
			if code != "" {
				impl := Impl{
					Id:            implID,
					LanguageName:  cell(languages, j),
					CodeBlock:     code,
					AuthorComment: cell(line, j+1),
					Version:       1,
				}
				implID++
				impls = append(impls, impl)
			}
		}
		idiom.Implementations = impls
		idioms = append(idioms, &idiom)
	}
	return idioms, nil
}

func exportIdiomsAsJSON(r *http.Request, w io.Writer, pretty bool) error {
	c := r.Context()
	_, idioms, err := dao.getAllIdioms(c, 0, "Id")
	if err != nil {
		return err
	}

	if pretty {
		// Advantage: output is pretty (human readable)
		// Drawback: the whole data transit through a byte buffer.
		buffer, err := json.MarshalIndent(idioms, "", "  ")
		if err != nil {
			return err
		}
		_, err = w.Write(buffer)
		return err
	} else {
		// Advantage: encodes (potentially voluminous) data "on-the-fly"
		//   Nope: buffered anyway. See discussion link.
		// Drawback: output is ugly.
		encoder := json.NewEncoder(w)
		return encoder.Encode(idioms)
		// TODO: see if possible to pretty-print on Writer, without buffering
		// Discussion https://groups.google.com/forum/#!topic/golang-nuts/NZ0n-RUerb0
	}

	//return json.MarshalIndent(idioms, "", "  ")

	// TODO export other entities :
	// idiom votes
	// impl votes
	// app config
}

// Not used anymore. See adminImportAjax.
func adminImport(w http.ResponseWriter, r *http.Request) error {
	c := r.Context()
	var err error
	file, fileHeader, err := r.FormFile("importData")
	_, err = importFile(c, file, fileHeader)
	if err != nil {
		return err
	}
	http.Redirect(w, r, hostPrefix()+"/admin", http.StatusFound)
	return nil
}
