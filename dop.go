package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jpoehls/go-dayone"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
	//	"strings"
)

type Journal struct {
	Id      string      `json:"id"`
	Title   string      `json:"title"`
	Starred bool        `json:"starred"`
	Tags    []string    `json:"tags"`
	Date    time.Time   `json:"date"`
	Photo   interface{} `json:"photo"` // interface is needed here so we can assign "nil" to entry that has no photo
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
	flag.StringVar(&journal, "journal", "./", "Journal Directory name")
	flag.StringVar(&entry, "entry", "", "Entry UUID")
	flag.Parse()

	jd := filepath.Join(journal, "entries")
	files, err := ioutil.ReadDir(jd)
	if err != nil {
		fmt.Printf("ERROR: Journal directory is not readable: %s\n", jd)
		os.Exit(1)
	}

	fmt.Printf("Found %d journal entries in %s\n", len(files), journal)
	var journals []Journal

	j := dayone.NewJournal(journal)

	parse := func(e *dayone.Entry, err error) error {
		var photo interface{}
		if err != nil {
			return err
		}

		fmt.Printf("Date: %s [%s] %s\n", e.CreationDate.Local(), e.UUID(), e.Tags)
		const layout = "Mon, 02 Jan 2006"

		p, err := j.PhotoStat(e.UUID())
		if (err == nil) && (p != nil) {
			photo = p.Name()
		} else {
			photo = nil
		}

		journals = append(journals, Journal{
			Id:      e.UUID(),
			Title:   e.CreationDate.Local().Format(layout),
			Starred: e.Starred,
			Tags:    e.Tags,
			Date:    e.CreationDate,
			Photo:   photo,
		})
		return nil
	}

	if entry != "" {
		e, err := j.ReadEntry(entry)
		if err != nil {
			panic(err)
		}
		err = parse(e, nil)
		if err != nil {
			panic(err)
		}

	} else {
		err = j.Read(parse)
		if err != nil {
			panic(err)
		}
	}

	sort.Sort(ByDate(journals))
	b, err := json.MarshalIndent(journals, "", "    ")
	if err != nil {
		fmt.Printf("ERROR: encoding JSON: %s\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(b)

}
