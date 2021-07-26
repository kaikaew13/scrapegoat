package scrapegoat

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Goat struct {
	URL    string
	client *http.Client
	opts   Options
	doc    *goquery.Document
	req    *http.Request
}

func NewGoat(url string, opts Options) (*Goat, error) {
	goat := Goat{
		URL:    url,
		client: new(http.Client),
		opts:   opts,
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

func (g *Goat) Scrape() []string {
	if g.req == nil {
		if err := g.newRequest(); err != nil {
			log.Panicln(err)
		}
	}

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

	data := []string{}

	g.doc.Find(".markdown-body h2").Each(func(i int, s *goquery.Selection) {
		data = append(data, s.Text())
	})

	return data
}

func (g *Goat) SetRequest(callback func(req *http.Request)) {
	callback(g.req)
}
