package main

import (
	"errors"
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/jpoehls/go-dayone"
	"github.com/russross/blackfriday"
)

type Journal struct {
	Id        string        `json:"id"`
	Title     string        `json:"title"`
	Starred   bool          `json:"starred"`
	Tags      []string      `json:"tags"`
	Date      time.Time     `json:"date"`
	Photo     interface{}   `json:"photo"` // interface is needed here so we can assign "nil" to entry that has no photo
	Count     int           `json:"count"`
	EntryText string        `json:"entrytext,omitempty"`
	EntryMD   template.HTML `json:"-"`
}

type Journals []Journal

// ByDate implements sort.Interface for []Journal based on
// the Date field.
type ByDate []Journal

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func (j Journals) CurrPosition(currId string) int {
	for p, v := range j {
		if v.Id == currId {
			return p
		}
	}
	return -1
}

func (j Journals) NextId(currPos int) string {
	l := len(j)
	if currPos+1 < l {
		return j[currPos+1].Id
	} else {
		return ""
	}
}

func (j Journals) PrevId(currPos int) string {
	if currPos > 0 {
		return j[currPos-1].Id
	} else {
		return ""
	}
}

func (myjournal *Myjournal) Parse(entry string) (Journals, error) {
	var err error
	var journals Journals

	j := dayone.NewJournal(myjournal.Dir)

	parse := func(e *dayone.Entry, err error, gettext bool) error {
		var photo interface{}
		var etext string
		var md template.HTML
		var cnt int

		if err != nil {
			return err
		}

		Trace.Printf("Date: %s [%s] %s\n", e.CreationDate.Local(), e.UUID(), e.Tags)
		const layout = "Mon, 02 Jan 2006"

		p, err := j.PhotoStat(e.UUID())
		if (err == nil) && (p != nil) {
			//photo = filepath.Join(journal, "photos", p.Name())
			//photo = filepath.Join("photos", p.Name())
			photo = p.Name()
		} else {
			photo = nil
		}

		if gettext {
			etext = e.EntryText
			md = template.HTML(blackfriday.MarkdownCommon([]byte(etext)))
			cnt = strings.Count(etext, myjournal.Count)
		} else {
			etext = ""
			md = template.HTML("")
			cnt = 0
		}

		journals = append(journals, Journal{
			Id:        e.UUID(),
			Title:     e.CreationDate.Local().Format(layout),
			Starred:   e.Starred,
			Tags:      e.Tags,
			Date:      e.CreationDate,
			Photo:     photo,
			Count:     cnt,
			EntryText: etext,
			EntryMD:   md,
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
		Error.Panic(err)
	}

	sort.Sort(ByDate(journals))
	return journals, nil

}
