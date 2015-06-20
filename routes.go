package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (jc *JournalConf) Mount() Routes {
	var routes Routes
	routes = Routes{
		Route{
			"Index",
			"GET",
			"/",
			jc.Index,
		},
		Route{
			"SearchEmpty",
			"GET",
			"/s/term=",
			jc.Index,
		},
		Route{
			"Search",
			"GET",
			"/s/term={term}",
			jc.Index,
		},
		Route{
			"Entry",
			"GET",
			"/" + jc.EUrl + "/{entryId}",
			jc.Index,
		},
		Route{
			"EntryList",
			"GET",
			"/api",
			jc.JsonAPI,
		},
		Route{
			"EntryShow",
			"GET",
			"/api/e/{entryId}",
			jc.JsonAPI,
		},
	}
	return routes
}
