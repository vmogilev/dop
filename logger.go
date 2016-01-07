package main

import (
	"net/http"
	"time"

	"github.com/vmogilev/dlog"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		var client string
		client = r.Header.Get("X-Forwarded-For")
		if client == "" {
			client = r.RemoteAddr
		}

		dlog.Info.Printf(
			"%s\t%s\t%s\t%s\t%s\t%s",
			client,
			time.Since(start),
			r.Method,
			name,
			r.RequestURI,
			r.Referer(),
		)
	})
}
