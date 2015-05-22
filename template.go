package main

import (
	//"fmt"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles(
	"templates/base.html",
	"templates/sidebar.html",
	"templates/content.html",
	"templates/customjs.html",
))

type Page struct {
	Title     string
	IsList    bool
	PrevId    string
	NextId    string
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
