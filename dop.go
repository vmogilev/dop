package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rigingo/dlog"
)

type Myjournal struct {
	Dir         string
	Title       string
	Count       string
	CssLookup   map[string]string
	HttpFQDN    string
	TemplateDIR string
}

var myjournal Myjournal

func Load(f string) {
	jData, err := ioutil.ReadFile(filepath.Join(f, "conf", "dop.json"))
	if err != nil {
		dlog.Error.Fatalf("Unable to read the data file (%s): %s", f, err)
	}
	if err := json.Unmarshal(jData, &myjournal); err != nil {
		dlog.Error.Fatalf("Unable to Unmarshal DOP config from data file (%s): %s", jData, err)
	}
	myjournal.TemplateDIR = filepath.Join(f, "templates")
	compileTemplate(myjournal.TemplateDIR)
	dlog.Info.Println(myjournal)

}

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

	Load(dopRoot)
	mp := MountPoint(httpMount)

	if httpHostExt != "" {
		myjournal.HttpFQDN = httpHostExt
	} else {
		if httpPort == "80" {
			myjournal.HttpFQDN = httpHost + mp
		} else {
			myjournal.HttpFQDN = httpHost + ":" + httpPort + mp
		}
	}

	var journal string
	journal = myjournal.Dir
	jd := filepath.Join(journal, "entries")
	files, err := ioutil.ReadDir(jd)
	if err != nil {
		dlog.Error.Fatalf("ERROR: Journal directory is not readable: %s\n", jd)
	}

	dlog.Info.Printf("Found %d journal entries in %s\n", len(files), journal)

	photos := filepath.Join(journal, "photos")
	router := NewRouter(httpMount, dopRoot, photos)

	dlog.Info.Fatal(http.ListenAndServe(":"+httpPort, router))
}
