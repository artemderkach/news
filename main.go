package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

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
		err = errors.Wrap(err, "error opening html template")
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	body, err := ioutil.ReadAll(f)
	if err != nil {
		err = errors.Wrap(err, "error reading html tempalte file")
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	t, err := template.New("").Parse(string(body))
	if err != nil {
		err = errors.Wrap(err, "error parsing tempalte")
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	t.Execute(w, feed)
}
