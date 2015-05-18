package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Myjournal struct {
	dir string
}

var myjournal Myjournal

func main() {
	var journal string
	var debug bool
	flag.StringVar(&journal, "journal", "./", "Journal Directory name")
	flag.BoolVar(&debug, "debug", false, "Debug")
	flag.Parse()

	log.SetFlags(log.LstdFlags)
	if debug {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	jd := filepath.Join(journal, "entries")
	files, err := ioutil.ReadDir(jd)
	if err != nil {
		log.Fatalf("ERROR: Journal directory is not readable: %s\n", jd)
	}

	log.Printf("Found %d journal entries in %s\n", len(files), journal)

	myjournal = Myjournal{dir: journal}

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
