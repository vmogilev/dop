package main

import (
	//"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var templates *template.Template

func compileTemplate(p string) {
	templates = template.Must(template.ParseFiles(
		filepath.Join(p, "base.html"),
		filepath.Join(p, "sidebar.html"),
		filepath.Join(p, "content.html"),
		filepath.Join(p, "customjs.html"),
	))
}

type Page struct {
	SiteTitle string
	Title     string
	Desc      string
	IsList    bool
	PrevId    string
	NextId    string
	HttpFQDN  string
	EUrl      string
	TUrl      string
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
