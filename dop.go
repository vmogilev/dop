package main

import (
	"flag"
	"fmt"
	"github.com/jpoehls/go-dayone"
	"io/ioutil"
	"os"
	"path/filepath"
//	"time"
//	"strings"
)

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

	j := dayone.NewJournal(journal)

	parse := func(e *dayone.Entry, err error) error {
		if err != nil {
			return err
		}

		// Do something with the entry,
		// or return dayone.ErrStopRead to break.
		fmt.Printf("Date: %s [%s] %s\n", e.CreationDate.Local(), e.UUID(), e.Tags)
		return nil
	}

	err = j.Read(parse)
	if err != nil {
		panic(err)
	}
}
