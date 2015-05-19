package main

import (
	"html/template"
	//"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles(
	"templates/base.html",
	"templates/index.html",
))

type Page struct {
	Title   string
	Navbar  interface{}
	Content interface{}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	//err := templates.ExecuteTemplate(w, tmpl+".html", j)
	err := templates.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
