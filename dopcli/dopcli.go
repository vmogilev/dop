package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rigingo/dlog"
)

func main() {

	var dopRoot string
	var httpHost string
	var httpPort string
	var httpMount string
	var httpHostExt string
	var debug bool

	flag.StringVar(&dopRoot, "dopRoot", "./", "DOP Root Directory [where the ./conf, ./static and ./templates dirs are]")
	flag.StringVar(&httpHost, "httpHost", "http://localhost", "HTTP Host Name")
	flag.StringVar(&httpPort, "httpPort", "8080", "HTTP Port")
	flag.StringVar(&httpMount, "httpMount", "/", "HTTP Mount Point [EX: /myjournal]")
	flag.StringVar(&httpHostExt, "httpHostExt", "", "Fully Qualified External Path if using Proxy [EX: http://mydomain.com/path]")
	flag.BoolVar(&debug, "debug", false, "Debug")
	flag.Parse()

	if debug {
		dlog.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else {
		dlog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}
}
