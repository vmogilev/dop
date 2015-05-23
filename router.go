package main

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func NewRouter(httpMount string, dopRoot string) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	//router.PathPrefix(httpMount)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	//router.PathPrefix("/").Handler(Logger(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))), "Static"))
	//router.PathPrefix("/").Handler(Logger(http.FileServer(http.Dir("./static/")), "Static"))
	router.PathPrefix(httpMount).Handler(Logger(http.FileServer(http.Dir(filepath.Join(dopRoot, "static"))), "Static"))
	//router.Methods("GET").Path(httpMount).Name("Static").Handler(Logger(http.FileServer(http.Dir(filepath.Join(dopRoot, "static"))), "Static"))

	return router
}
