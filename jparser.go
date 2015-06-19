package main

import (
	"errors"
	"html/template"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jpoehls/go-dayone"
	"github.com/rigingo/dlog"
	"github.com/russross/blackfriday"
)

type Journal struct {
	Id        string        `json:"id"`
	Title     string        `json:"title"`
	Starred   bool          `json:"starred"`
	Tags      []string      `json:"tags"`
	Date      time.Time     `json:"date"`
	CharDate  string        `json:"chardate"`
	Photo     interface{}   `json:"photo"` // interface is needed here so we can assign "nil" to entry that has no photo
	Thumb     interface{}   `json:"thumb"` // interface is needed here so we can assign "nil" to entry that has no photo
	Small     interface{}   `json:"small"` // interface is needed here so we can assign "nil" to entry that has no photo
	Count     int           `json:"count"`
	EntryText string        `json:"entrytext,omitempty"`
	EntryMD   template.HTML `json:"-"`
	DopTitle  string        `json:"doptitle"`
	DopDesc   string        `json:"dopdesc"`
	DopLink   string        `json:"doplink"`
}

type Journals []Journal

type JIndex map[string]string

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
		return j[currPos+1].DopLink
	} else {
		return ""
	}
}

// Note: the sort is now reversed that why Prev is Next and Next is Prev
func (j Journals) NextId(currPos int) string {
	if currPos > 0 {
		return j[currPos-1].DopLink
	} else {
		return ""
	}
}

func GrepI(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

type DopTokens struct {
	Title string
	Desc  string
	Link  string
}

func ParseToken(s, t string) string {
	rp := regexp.MustCompile(t)

	if rp.MatchString(s) {
		dlog.Trace.Printf("t=%s | s=%s", t, s)
		return s[len(t)-1:]
	}
	return ""
}

// GetTokens gets values for three dop tokens from lines 1,2,3
// line 1: Title	- should start with "# "
// line 2: Description	- should start with "//dop:desc "
// line 3: Link		- should start with "//dop:link "
func GetTokens(s string) DopTokens {
	var r DopTokens
	r = DopTokens{}
	t := "^# "
	d := "^//dop:desc "
	l := "^//dop:link "

	lines := strings.Split(s, "\n")

	if len(lines) >= 3 {
		dlog.Trace.Printf("lines[0]=%s", lines[0])
		dlog.Trace.Printf("lines[1]=%s", lines[1])
		dlog.Trace.Printf("lines[2]=%s", lines[2])
		r.Title = ParseToken(lines[0], t)
		r.Desc = ParseToken(lines[1], d)
		r.Link = ParseToken(lines[2], l)
	}
	return r
}

func (myjournal *Myjournal) Parse(entry string, s string) (Journals, JIndex, error) {
	var err error
	var journals Journals
	var jindex JIndex
	jindex = JIndex{}

	j := dayone.NewJournal(myjournal.Dir)

	parse := func(e *dayone.Entry, err error, gettext bool, search string, buildIndex bool) error {
		var photo interface{}
		var thumb interface{}
		var small interface{}
		var etext string
		var md template.HTML
		var cnt int
		var serp bool
		var tokens DopTokens
		var uuid string
		var pub bool

		uuid = e.UUID()

		if err != nil {
			return err
		}

		dlog.Trace.Printf("Date: %s [%s] %s\n", e.CreationDate.Local(), uuid, e.Tags)
		const layout = "Mon, 02 Jan 2006"

		p, err := j.PhotoStat(uuid)
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
			tokens = GetTokens(etext)
			if buildIndex {
				if tokens.Link != "" {
					jindex[tokens.Link] = uuid
				} else {
					jindex[uuid] = uuid
				}
			}
		} else {
			etext = ""
			md = template.HTML("")
			cnt = 0
			tokens = DopTokens{}
		}

		dlog.Trace.Printf("1: tokens.Title=%s", tokens.Title)
		dlog.Trace.Printf("1: tokens.Desc=%s", tokens.Desc)
		dlog.Trace.Printf("1: tokens.Link=%s", tokens.Link)

		if (search != "") && (gettext) {
			serp = GrepI(etext, search)
		} else {
			serp = true
		}

		var title string
		charDate := e.CreationDate.Local().Format(layout)

		if tokens.Title != "" {
			title = tokens.Title
		} else {
			title = charDate
		}

		if tokens.Title == "" {
			tokens.Title = title
		}
		if tokens.Desc == "" {
			tokens.Desc = title
		}
		if tokens.Link == "" {
			tokens.Link = uuid
		}

		dlog.Trace.Printf("2: tokens.Title=%s", tokens.Title)
		dlog.Trace.Printf("2: tokens.Desc=%s", tokens.Desc)
		dlog.Trace.Printf("2: tokens.Link=%s", tokens.Link)

		pub = true
		if myjournal.PubStarred {
			pub = false
			if e.Starred {
				pub = true
			}
		}

		if (serp) && (pub) {
			journals = append(journals, Journal{
				Id:        uuid,
				Title:     title,
				Starred:   e.Starred,
				Tags:      e.Tags,
				Date:      e.CreationDate,
				CharDate:  charDate,
				Photo:     photo,
				Thumb:     thumb,
				Small:     small,
				Count:     cnt,
				EntryText: etext,
				EntryMD:   md,
				DopTitle:  tokens.Title,
				DopDesc:   tokens.Desc,
				DopLink:   tokens.Link,
			})
		}
		return nil
	}

	var parseall bool
	var buildIndex bool

	if (entry != "") && (entry != "*") {
		e, err := j.ReadEntry(entry)
		if err != nil {
			return nil, nil, errors.New("NotFound")
		}
		parseall = true
		buildIndex = false
		err = parse(e, nil, parseall, "", buildIndex)

	} else {
		if entry == "*" {
			parseall = true
			buildIndex = true
		} else {
			parseall = false
			buildIndex = false
		}
		// closure to wrap the extra param
		err = j.Read(func(e *dayone.Entry, err error) error {
			return parse(e, err, parseall, s, buildIndex)
			//return err

		})
	}

	if err != nil {
		dlog.Error.Panic(err)
	}

	sort.Sort(ByDate(journals))
	return journals, jindex, nil

}
