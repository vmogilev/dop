package main

import (
	//"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var templates = template.Must(template.ParseFiles(
	filepath.Join(myjournal.TemplateDIR, "base.html"),
	filepath.Join(myjournal.TemplateDIR, "sidebar.html"),
	filepath.Join(myjournal.TemplateDIR, "content.html"),
	filepath.Join(myjournal.TemplateDIR, "customjs.html"),
))

type Page struct {
	Title     string
	IsList    bool
	PrevId    string
	NextId    string
	HttpFQDN  string
	Search    string
	CssLookup map[string]string
	Navbar    interface{}
	Content   interface{}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.Execute(w, p)
	//err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
