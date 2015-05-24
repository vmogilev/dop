package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NotFound(id string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found", Id: id}); err != nil {
		Error.Panic(err)
	}

}

func (myjournal *Myjournal) List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	var entry string
	var journals Journals

	entry = vars["entryId"]
	journals, err = myjournal.Parse(entry, "")
	if (err != nil) && (err.Error() == "NotFound") {
		NotFound(entry, w)
		return
	}

	b, err := json.MarshalIndent(journals, "", "    ")
	if err != nil {
		Error.Panicf("ERROR: encoding JSON: %s\n", err)
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
	var current Journals
	var list bool
	var search string

	search = vars["term"]
	journals, err = myjournal.Parse("*", search)

	entry = vars["entryId"]
	if entry == "" {
		list = true
		if len(journals) > 0 {
			entry = journals[0].Id
		} else {
			entry = ""
		}
	} else {
		list = false
	}

	current, err = myjournal.Parse(entry, "")
	if (err != nil) && (err.Error() == "NotFound") {
		NotFound(entry, w)
		return
	}

	var nextid string
	var previd string
	if currpos := journals.CurrPosition(entry); currpos != -1 {
		nextid = journals.NextId(currpos)
		previd = journals.PrevId(currpos)
	}

	page := Page{
		Title:     mj.Title,
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

func (mj *Myjournal) Test(w http.ResponseWriter, r *http.Request) {
	var csslookup = make(map[string]string)
	csslookup["horrible"] = "danger"
	csslookup["Tarlov"] = "danger"
	csslookup["IN_experiment"] = "warning"

	var journals Journals
	page := Page{Title: "Vitaliy's Food Journal",
		IsList:    false,
		PrevId:    "PREVID",
		NextId:    "NEXTID",
		CssLookup: csslookup,
		Navbar:    journals,
		Content:   journals,
	}
	b, err := json.MarshalIndent(page, "", "    ")
	if err != nil {
		Error.Panicf("ERROR: encoding JSON: %s\n", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
