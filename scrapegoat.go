package scrapegoat

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Goat struct {
	URL               string
	MaxRecursionDepth int
	EnableConcurrency bool
	EnableLogging     bool
	doc               *goquery.Document
	req               *http.Request
	selectorQueue     []cssSelector
}

func NewGoat(url string) (*Goat, error) {
	goat := Goat{
		URL:               url,
		MaxRecursionDepth: 3,
		EnableConcurrency: false,
		EnableLogging:     false,
		req:               nil,
		selectorQueue:     []cssSelector{},
	}

	if err := goat.newRequest(); err != nil {
		return nil, err
	}

	if goat.EnableLogging {
		goat.SetRequest(func(req *http.Request) {
			log.Println(req.URL)
		})
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
	client := new(http.Client)
	res, err := client.Do(g.req)
	if err != nil {
		log.Panicln(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Panicf("got a response with response code of %d, want %d\n", res.StatusCode, http.StatusOK)
	}

	g.doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panicln(err)
	}

	for _, s := range g.selectorQueue {
		g.doc.Find(s.selector).Each(func(i int, gs *goquery.Selection) {
			if g.EnableLogging {
				log.Printf("url: %s, selector: %s\n", g.req.URL, s.selector)
			}

			s.callback(Selection{
				gs,
				[]cssSelector{},
			})
		})
	}
}

func (g *Goat) SetRequest(callback func(req *http.Request)) {
	callback(g.req)
}

func (g *Goat) SetSelector(selector string, callback func(s Selection)) {
	g.selectorQueue = append(g.selectorQueue, cssSelector{
		selector: selector,
		callback: callback,
	})
}
