package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// CustomItem just a bit expanded version of default *gofeed.Item
type CustomItem struct {
	GlobalTitle string // gofeed.Feed.Title
	Item        *gofeed.Item
}

// CustomFeed is needed for storing all Items(from different urls) is one place.
// in case of default feed we will have []*gofeed.Fead, in thes case there is no possibility
// to order all items in all feeds by name.
// CustomFeed struct solve the problem of storing items for futere filtering by date
type CustomFeed struct {
	Items []*CustomItem
}

// urls represents path with sites rss feeds
var urls []string = []string{
	"https://news.radio-t.com/rss",
	"https://news.ycombinator.com/rss",
}

const layout string = "Mon, 2 Jan 2006 15:04:05 -0700"

func main() {
	http.HandleFunc("/", news)
	fmt.Println("listenign on localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func news(w http.ResponseWriter, r *http.Request) {
	// create template form .html file
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

	// put all data from urls to custom feed struct
	customFeed := &CustomFeed{}
	fp := gofeed.NewParser()
	for _, url := range urls {
		feed, err := fp.ParseURL(url)
		if err != nil {
			fmt.Println(errors.Wrap(err, "error retrieving rss data"))
			continue
		}

		for _, item := range feed.Items {
			customItem := &CustomItem{
				GlobalTitle: feed.Title,
				Item:        item,
			}
			customItem.GlobalTitle = feed.Title
			customFeed.Items = append(customFeed.Items, customItem)
		}
	}

	orderByDate(customFeed)

	for _, customItem := range customFeed.Items {
		err = changeDateFormat(customItem.Item)
		if err != nil {
			sendError(w, errors.Wrap(err, "error changing date format"))
			return
		}
	}

	err = t.Execute(w, customFeed.Items)
	if err != nil {
		sendError(w, errors.Wrap(err, "error applying parsed tempalte"))
		return
	}
}

func orderByDate(feed *CustomFeed) {
	for i := 0; i < len(feed.Items)-1; i += 1 {
		for j := 0; j < len(feed.Items)-1-i; j += 1 {
			t1, err := time.Parse(layout, feed.Items[j].Item.Published)
			if err != nil {
				fmt.Println(errors.Wrapf(err, "error parsing time %s", feed.Items[j].Item.Published))
				continue
			}
			t2, err := time.Parse(layout, feed.Items[j+1].Item.Published)
			if err != nil {
				fmt.Println(errors.Wrapf(err, "error parsing time %s", feed.Items[j+1].Item.Published))
				continue
			}

			if t1.Before(t2) {
				feed.Items[j], feed.Items[j+1] = feed.Items[j+1], feed.Items[j]
			}
		}
	}
}

// changeDateFormat changes date format to setisfy simpler rule "DD.MM"
func changeDateFormat(item *gofeed.Item) error {
	t, err := time.Parse(layout, item.Published)
	if err != nil {
		errors.Wrap(err, "error parsing time")
	}
	item.Published = fmt.Sprintf("%02d.%02d.%s", t.Day(), t.Month(), strconv.Itoa(t.Year())[2:])

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
