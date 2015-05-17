package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		myjournal.Index,
	},
	//Route{
	//	"EntryShow",
	//	"GET",
	//	"/entry/{entryId}",
	//	myjournal.EntryShow,
	//},
}
