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
		"EntryList",
		"GET",
		"/list",
		myjournal.List,
	},
	Route{
		"EntryShow",
		"GET",
		"/entry/{entryId}",
		myjournal.List,
	},
}
