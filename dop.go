package main

import (
	"encoding/json"
	"flag"
	//"fmt"
	"github.com/jpoehls/go-dayone"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
	//	"strings"
)

type Journal struct {
	Id        string      `json:"id"`
	Title     string      `json:"title"`
	Starred   bool        `json:"starred"`
	Tags      []string    `json:"tags"`
	Date      time.Time   `json:"date"`
	Photo     interface{} `json:"photo"` // interface is needed here so we can assign "nil" to entry that has no photo
	EntryText string      `json:"entrytext"`
}

// ByDate implements sort.Interface for []Journal based on
// the Date field.
type ByDate []Journal

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func main() {
	var journal string
	var entry string
	var debug bool
	flag.StringVar(&journal, "journal", "./", "Journal Directory name")
	flag.StringVar(&entry, "entry", "", "Entry UUID")
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
	var journals []Journal

	j := dayone.NewJournal(journal)

	parse := func(e *dayone.Entry, err error, gettext bool) error {
		var photo interface{}
		var etext string

		if err != nil {
			return err
		}

		log.Printf("Date: %s [%s] %s\n", e.CreationDate.Local(), e.UUID(), e.Tags)
		const layout = "Mon, 02 Jan 2006"

		p, err := j.PhotoStat(e.UUID())
		if (err == nil) && (p != nil) {
			photo = p.Name()
		} else {
			photo = nil
		}

		if gettext {
			etext = e.EntryText
		} else {
			etext = ""
		}

		journals = append(journals, Journal{
			Id:        e.UUID(),
			Title:     e.CreationDate.Local().Format(layout),
			Starred:   e.Starred,
			Tags:      e.Tags,
			Date:      e.CreationDate,
			Photo:     photo,
			EntryText: etext,
		})
		return nil
	}

	var parseall bool
	if (entry != "") && (entry != "*") {
		e, err := j.ReadEntry(entry)
		if err != nil {
			log.Panic(err)
		}
		err = parse(e, nil, true)

	} else {
		if entry == "*" {
			parseall = true
		} else {
			parseall = false
		}
		// closure to wrap the extra param
		err = j.Read(func(e *dayone.Entry, err error) error {
			return parse(e, err, parseall)
			//return err

		})
	}

	if err != nil {
		log.Panic(err)
	}

	sort.Sort(ByDate(journals))
	b, err := json.MarshalIndent(journals, "", "    ")
	if err != nil {
		log.Panicf("ERROR: encoding JSON: %s\n", err)
	}
	os.Stdout.Write(b)

}
