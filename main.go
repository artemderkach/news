package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// some comments
// blah
// blah
func main() {
	http.HandleFunc("/", news)
	fmt.Println("listenign on localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func news(w http.ResponseWriter, r *http.Request) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://news.radio-t.com/rss")

	f, err := os.Open("./news.html")
	if err != nil {
		sendError(w, errors.Wrap(err, "error opening html template"))
		return
	}

	body, err := ioutil.ReadAll(f)
	if err != nil {
		sendError(w, errors.Wrap(err, "error reading html tempalte file"))
		return
	}

	t, err := template.New("").Parse(string(body))
	if err != nil {
		sendError(w, errors.Wrap(err, "error parsing tempalte"))
		return
	}

	err = changeDateFormat(feed)
	if err != nil {
		err = errors.Wrap(err, "error changing date format")
		fmt.Println(err)
	}

	err = t.Execute(w, feed)
	if err != nil {
		err = errors.Wrap(err, "error applying parsed tempalte")
		fmt.Println(err)
	}
}

func changeDateFormat(feed *gofeed.Feed) error {
	// Mon, 05 Aug 2019 21:34:31 +0000
	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	t, err := time.Parse(layout, feed.Published)
	if err != nil {
		return errors.Wrap(err, "error parsing time")
	}
	fmt.Printf("%01d.%d.%v\n", t.Year(), t.Month(), t.Day())

	return nil
}

func sendError(w http.ResponseWriter, err error) {
	fmt.Println(err)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		err = errors.Wrap(err, "error sending response")
		fmt.Println(err)
	}
}
