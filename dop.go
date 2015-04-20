package main

import (
	"encoding/json"
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
	Id    string   `json:"Id"`
	Title string   `json:"Title"`
	Tags  []string `json:"Tags"`
	Sort  string   `json:"Sort"`
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
			Id:    e.UUID(),
			Title: e.CreationDate.Local().Format(layout),
			Tags:  e.Tags,
			Sort:  e.CreationDate.Local().Format("2006.01.02-15.04.05"),
		})
		return nil
	}

	err = j.Read(parse)
	if err != nil {
		panic(err)
	}

	//fmt.Println(l)
	b, err := json.MarshalIndent(l, "", "    ")
	if err != nil {
		fmt.Printf("ERROR: encoding JSON: %s\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(b)

}
