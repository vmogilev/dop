package main

import (
	"flag"
	"fmt"
	"github.com/jpoehls/go-dayone"
	"io/ioutil"
	"os"
	"path/filepath"
	//"time"
	//	"strings"
)

type Journal struct {
	Key   string   `json:"Key"`
	Title string   `json:"Title"`
	Tags  []string `json:"Tags"`
}

func main() {
	var journal string
	flag.StringVar(&journal, "journal", "./", "Journal Directory name")
	flag.Parse()

	//jd := strings.Join([]string{journal, "entries"}, string(os.PathSeparator))
	jd := filepath.Join(journal, "entries")
	files, err := ioutil.ReadDir(jd)
	if err != nil {
		fmt.Printf("ERROR: Journal directory is not readable: %s\n", jd)
		os.Exit(1)
	}

	fmt.Printf("Found %d journal entries in %s\n", len(files), journal)
	//l := make([]Journal, len(files))
	var l []Journal

	j := dayone.NewJournal(journal)

	parse := func(e *dayone.Entry, err error) error {
		if err != nil {
			return err
		}

		// Do something with the entry,
		// or return dayone.ErrStopRead to break.
		fmt.Printf("Date: %s [%s] %s\n", e.CreationDate.Local(), e.UUID(), e.Tags)
		//const layout = time.RubyDate
		const layout = "Mon, 02 Jan 2006"
		l = append(l, Journal{
			Key:   e.UUID(),
			Title: e.CreationDate.Local().Format(layout),
			Tags:  e.Tags,
		})
		return nil
	}

	err = j.Read(parse)
	fmt.Println(l)
	if err != nil {
		panic(err)
	}
}
