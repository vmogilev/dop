package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rigingo/dlog"
)

func NotFound(id string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found", Id: id}); err != nil {
		dlog.Error.Panic(err)
	}

}

func (myjournal *Myjournal) List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	var entry string
	var journals Journals
	//var jindex JIndex

	entry = vars["entryId"]
	journals, _, err = myjournal.Parse(entry, "")
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

func (mj *Myjournal) Index(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	var entry string
	var journals Journals
	var jindex JIndex
	var current Journals
	var list bool
	var search string
	var desc string

	search = strings.Replace(vars["term"], "+", " ", -1)
	journals, jindex, err = myjournal.Parse("*", search)

	entry = vars["entryId"]
	dlog.Trace.Printf("entry_POST=%s", entry)
	if entry == "" {
		list = true
		desc = mj.Desc
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

	current, _, err = myjournal.Parse(entry, "")
	if (err != nil) && (err.Error() == "NotFound") {
		NotFound(entry, w)
		return
	}

	if desc == "" {
		if desc = current[0].DopDesc; desc == "" {
			desc = mj.Desc
		}
	}

	var nextid string
	var previd string
	if currpos := journals.CurrPosition(entry); currpos != -1 {
		nextid = journals.NextId(currpos)
		previd = journals.PrevId(currpos)
	}

	page := Page{
		Title:     mj.Title,
		Desc:      desc,
		IsList:    list,
		PrevId:    previd,
		NextId:    nextid,
		HttpFQDN:  mj.HttpFQDN,
		Search:    search,
		CssLookup: mj.CssLookup,
		Navbar:    journals,
		Content:   current,
	}
	renderTemplate(w, "dop", &page)
}
