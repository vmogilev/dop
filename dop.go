package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rigingo/dlog"
)

type JournalConf struct {
	JDir     string
	Template string
	HttpFQDN string
	DopRoot  string
	EUrl     string
	TUrl     string
}

var jc JournalConf

func main() {
	var dopRoot string
	var httpHost string
	var httpPort string
	var httpMount string
	var httpHostExt string
	var jTemplate string
	var jDir string
	var debug bool
	var fqdn string
	var eURL string
	var tURL string

	flag.StringVar(&dopRoot, "dopRoot", "./", "DOP Root Directory [where the ./conf, ./static and ./templates dirs are]")
	flag.StringVar(&httpHost, "httpHost", "http://localhost", "HTTP Host Name")
	flag.StringVar(&httpPort, "httpPort", "8080", "HTTP Port")
	flag.StringVar(&httpMount, "httpMount", "/", "HTTP Mount Point [EX: /myjournal]")
	flag.StringVar(&httpHostExt, "httpHostExt", "", "Fully Qualified External Path if using Proxy [EX: http://mydomain.com/path]")
	flag.StringVar(&jTemplate, "jTemplate", "dop_blog", "Tempate Directory Name in ./templates/")
	flag.StringVar(&jDir, "jDir", "", "Path to Day One Journal Directory where entries subdirectory is located")
	flag.StringVar(&eURL, "eURL", "e", "Entry URL Path: http://{httpHost}/{httpMount}/{eURL}")
	flag.StringVar(&tURL, "tURL", "tag", "Tag URL Path: http://{httpHost}/{httpMount}/{tURL}")
	flag.BoolVar(&debug, "debug", false, "Debug")
	flag.Parse()

	if debug {
		dlog.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else {
		dlog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}

	mp := MountPoint(httpMount)
	if httpHostExt != "" {
		fqdn = httpHostExt
	} else {
		if httpPort == "80" {
			fqdn = httpHost + mp
		} else {
			fqdn = httpHost + ":" + httpPort + mp
		}
	}

	tdir := filepath.Join(dopRoot, "templates", jTemplate)
	compileTemplate(tdir)

	jd := filepath.Join(jDir, "entries")
	files, err := ioutil.ReadDir(jd)
	if err != nil {
		dlog.Error.Fatalf("ERROR: Journal directory is not readable: %s\n", jd)
	}

	dlog.Info.Printf("Found %d journal entries in %s\n", len(files), jDir)

	jc = JournalConf{
		JDir:     jDir,
		Template: jTemplate,
		HttpFQDN: fqdn,
		DopRoot:  dopRoot,
		EUrl:     eURL,
		TUrl:     tURL,
	}

	jv := Load(jc.DopRoot)
	dlog.Info.Println(jv)

	photos := filepath.Join(jDir, "photos")
	router := NewRouter(httpMount, dopRoot, photos)

	dlog.Info.Fatal(http.ListenAndServe(":"+httpPort, router))
}
