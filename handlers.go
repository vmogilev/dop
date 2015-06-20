package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rigingo/dlog"
)

type JournalVars struct {
	Title      string
	Desc       string
	PubStarred bool
	Count      string
	CssLookup  map[string]string
}

func Load(f string) JournalVars {
	jv := JournalVars{}

	jf, err := ioutil.ReadFile(filepath.Join(f, "conf", "dop.json"))
	if err != nil {
		dlog.Error.Fatalf("Unable to read the data file (%s): %s", f, err)
	}
	if err := json.Unmarshal(jf, &jv); err != nil {
		dlog.Error.Fatalf("Unable to Unmarshal DOP config from data file (%s): %s", jf, err)
	}
	dlog.Trace.Println(jv)
	return jv
}

func NotFound(id string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found", Id: id}); err != nil {
		dlog.Error.Panic(err)
	}

}

func (jc *JournalConf) JsonAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	var entry string
	var journals Journals
	//var jindex JIndex

	var jv JournalVars
	jv = Load(jc.DopRoot)

	entry = vars["entryId"]
	journals, _, err = Parse(entry, "", jc, &jv)
	if (err != nil) && (err.Error() == "NotFound") {
		NotFound(entry, w)
		return
	}

	b, err := json.MarshalIndent(journals, "", "    ")
	if err != nil {
		dlog.Error.Panicf("ERROR: encoding JSON: %s\n", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (jc *JournalConf) Index(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	var entry string
	var journals Journals
	var jindex JIndex
	var current Journals
	var list bool
	var search string
	var title, desc string

	var jv JournalVars
	jv = Load(jc.DopRoot)

	search = strings.Replace(vars["term"], "+", " ", -1)
	journals, jindex, err = Parse("*", search, jc, &jv)

	entry = vars["entryId"]
	dlog.Trace.Printf("entry_POST=%s", entry)
	if entry == "" {
		list = true
		desc = jv.Desc
		title = jv.Title
		if len(journals) > 0 {
			entry = journals[0].Id
		} else {
			entry = ""
		}
	} else {
		list = false
		entry = jindex[entry]
	}

	dlog.Trace.Printf("entry_PARSE=%s", entry)
	dlog.Trace.Printf("len(jindex)=%s", len(jindex))

	current, _, err = Parse(entry, "", jc, &jv)
	if (err != nil) && (err.Error() == "NotFound") {
		NotFound(entry, w)
		return
	}

	if desc == "" {
		if desc = current[0].DopDesc; desc == "" {
			desc = jv.Desc
		}
	}

	if title == "" {
		if title = current[0].Title; title == "" {
			title = jv.Title
		}
	}

	var nextid string
	var previd string
	if currpos := journals.CurrPosition(entry); currpos != -1 {
		nextid = journals.NextId(currpos)
		previd = journals.PrevId(currpos)
	}

	page := Page{
		Title:     title,
		Desc:      desc,
		IsList:    list,
		PrevId:    previd,
		NextId:    nextid,
		HttpFQDN:  jc.HttpFQDN,
		EUrl:      jc.EUrl,
		TUrl:      jc.TUrl,
		Search:    search,
		CssLookup: jv.CssLookup,
		Navbar:    journals,
		Content:   current,
	}
	renderTemplate(w, "dop", &page)
}
