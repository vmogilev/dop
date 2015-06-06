package main

import (
	"errors"
	"html/template"
	"path/filepath"
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
	Thumb     interface{}   `json:"photo"` // interface is needed here so we can assign "nil" to entry that has no photo
	Small     interface{}   `json:"photo"` // interface is needed here so we can assign "nil" to entry that has no photo
	Count     int           `json:"count"`
	EntryText string        `json:"entrytext,omitempty"`
	EntryMD   template.HTML `json:"-"`
}

type Journals []Journal

// ByDate implements sort.Interface for []Journal based on
// the Date field.
type ByDate []Journal

func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

//func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }
func (a ByDate) Less(i, j int) bool { return a[i].Date.After(a[j].Date) }

func (j Journals) CurrPosition(currId string) int {
	for p, v := range j {
		if v.Id == currId {
			return p
		}
	}
	return -1
}

// Note: the sort is now reversed that why Prev is Next and Next is Prev
func (j Journals) PrevId(currPos int) string {
	l := len(j)
	if currPos+1 < l {
		return j[currPos+1].Id
	} else {
		return ""
	}
}

// Note: the sort is now reversed that why Prev is Next and Next is Prev
func (j Journals) NextId(currPos int) string {
	if currPos > 0 {
		return j[currPos-1].Id
	} else {
		return ""
	}
}

func GrepI(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func (myjournal *Myjournal) Parse(entry string, s string) (Journals, error) {
	var err error
	var journals Journals

	j := dayone.NewJournal(myjournal.Dir)

	parse := func(e *dayone.Entry, err error, gettext bool, search string) error {
		var photo interface{}
		var thumb interface{}
		var small interface{}
		var etext string
		var md template.HTML
		var cnt int
		var serp bool

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
			//thumb = MakeThumbnail(filepath.Join(myjournal.Dir, "photos"), photo.(string), 28, 28)
			//small = MakeThumbnail(filepath.Join(myjournal.Dir, "photos"), photo.(string), 640, 0)
			thumb = MakeThumbnailVIPS(filepath.Join(myjournal.Dir, "photos"), photo.(string), 28, 28)
			small = MakeThumbnailVIPS(filepath.Join(myjournal.Dir, "photos"), photo.(string), 640, 0)
		} else {
			photo = nil
			thumb = nil
			small = nil
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

		if (search != "") && (gettext) {
			serp = GrepI(etext, search)
		} else {
			serp = true
		}

		if serp {
			journals = append(journals, Journal{
				Id:        e.UUID(),
				Title:     e.CreationDate.Local().Format(layout),
				Starred:   e.Starred,
				Tags:      e.Tags,
				Date:      e.CreationDate,
				Photo:     photo,
				Thumb:     thumb,
				Small:     small,
				Count:     cnt,
				EntryText: etext,
				EntryMD:   md,
			})
		}
		return nil
	}

	var parseall bool
	if (entry != "") && (entry != "*") {
		e, err := j.ReadEntry(entry)
		if err != nil {
			return nil, errors.New("NotFound")
		}
		err = parse(e, nil, true, "")

	} else {
		if entry == "*" {
			parseall = true
		} else {
			parseall = false
		}
		// closure to wrap the extra param
		err = j.Read(func(e *dayone.Entry, err error) error {
			return parse(e, err, parseall, s)
			//return err

		})
	}

	if err != nil {
		Error.Panic(err)
	}

	sort.Sort(ByDate(journals))
	return journals, nil

}
