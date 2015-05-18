package main

import (
	"errors"
	"log"
	"path/filepath"
	"sort"
	"time"

	"github.com/jpoehls/go-dayone"
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

type Journals []Journal

// ByDate implements sort.Interface for []Journal based on
// the Date field.
type ByDate []Journal

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func (myjournal *Myjournal) JournalParser(entry string) (Journals, error) {
	var err error
	var journals Journals

	j := dayone.NewJournal(myjournal.dir)

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
			//photo = filepath.Join(journal, "photos", p.Name())
			photo = filepath.Join("photos", p.Name())
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
			return nil, errors.New("NotFound")
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
	return journals, nil

}
