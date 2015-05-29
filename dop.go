package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func Load(f string) {
	jData, err := ioutil.ReadFile(filepath.Join(f, "conf", "dop.json"))
	if err != nil {
		Error.Fatalf("Unable to read the data file (%s): %s", f, err)
	}
	if err := json.Unmarshal(jData, &myjournal); err != nil {
		Error.Fatalf("Unable to Unmarshal DOP config from data file (%s): %s", jData, err)
	}
	myjournal.TemplateDIR = filepath.Join(f, "templates")
	compileTemplate(myjournal.TemplateDIR)
	Info.Println(myjournal)

}

func main() {
	var dopRoot string
	var httpHost string
	var httpPort string
	var httpMount string
	var debug bool

	flag.StringVar(&dopRoot, "dopRoot", "./", "DOP Root Directory [where the ./conf, ./static and ./templates dirs are]")
	flag.StringVar(&httpHost, "httpHost", "http://localhost", "HTTP Host Name")
	flag.StringVar(&httpPort, "httpPort", "8080", "HTTP Port")
	flag.StringVar(&httpMount, "httpMount", "/", "HTTP Mount Point [EX: /myjournal]")
	flag.BoolVar(&debug, "debug", false, "Debug")
	flag.Parse()

	if debug {
		Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else {
		Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}

	Load(dopRoot)
	mp := MountPoint(httpMount)
	if httpPort == "80" {
		myjournal.HttpFQDN = httpHost + mp
	} else {
		myjournal.HttpFQDN = httpHost + ":" + httpPort + mp
	}

	var journal string
	journal = myjournal.Dir
	jd := filepath.Join(journal, "entries")
	files, err := ioutil.ReadDir(jd)
	if err != nil {
		Error.Fatalf("ERROR: Journal directory is not readable: %s\n", jd)
	}

	Info.Printf("Found %d journal entries in %s\n", len(files), journal)

	photos := filepath.Join(journal, "photos")
	router := NewRouter(httpMount, dopRoot, photos)

	Info.Fatal(http.ListenAndServe(":"+httpPort, router))
}
