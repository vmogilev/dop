package main

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func MountPoint(httpMount string) string {
	if httpMount == "/" {
		return ""
	} else {
		return httpMount
	}
}

func NewRouter(httpMount string, dopRoot string, photos string) *mux.Router {

	mp := MountPoint(httpMount)

	router := mux.NewRouter().StrictSlash(true)
	//router.PathPrefix(httpMount)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(mp + route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	//router.PathPrefix("/").Handler(Logger(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))), "Static"))
	//router.PathPrefix("/").Handler(Logger(http.FileServer(http.Dir("./static/")), "Static"))
	//router.PathPrefix(mp + "/").Handler(Logger(http.FileServer(http.Dir(filepath.Join(dopRoot, "static"))), "Static"))
	router.PathPrefix(mp + "/photos").Handler(http.StripPrefix(mp+"/photos", Logger(http.FileServer(http.Dir(photos)), "Photos")))
	router.PathPrefix(mp + "/").Handler(http.StripPrefix(mp+"/", Logger(http.FileServer(http.Dir(filepath.Join(dopRoot, "static"))), "Static")))

	return router
}
