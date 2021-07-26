package scrapegoat

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Goat struct {
	URL        string
	client     *http.Client
	opts       Options
	doc        *goquery.Document
	req        *http.Request
	selections []Selection
}

func NewGoat(url string, opts Options) (*Goat, error) {
	goat := Goat{
		URL:        url,
		client:     new(http.Client),
		opts:       opts,
		selections: []Selection{},
	}

	if err := goat.newRequest(); err != nil {
		return nil, err
	}

	return &goat, nil
}

func (g *Goat) newRequest() error {
	req, err := http.NewRequest(http.MethodGet, g.URL, nil)
	if err != nil {
		return err
	}

	g.req = req
	return nil
}

func (g *Goat) Scrape() {
	res, err := g.client.Do(g.req)
	if err != nil {
		log.Panicln(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Panicf("got a response with response code of %d, want %d", res.StatusCode, http.StatusOK)
	}

	g.doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panicln(err)
	}

	for _, q := range g.selections {
		g.doc.Find(q.selector).Each(func(i int, s *goquery.Selection) {
			q.callback(s)
		})
	}
}

func (g *Goat) SetRequest(callback func(req *http.Request)) {
	callback(g.req)
}

func (g *Goat) SetTags(selector string, callback func(s *goquery.Selection)) {
	g.selections = append(g.selections, Selection{
		selector: selector,
		callback: callback,
	})
}
