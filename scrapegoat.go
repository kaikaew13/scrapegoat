package scrapegoat

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Goat struct {
	// URL               string
	MaxRecursionDepth int
	EnableConcurrency bool
	EnableLogging     bool
	doc               *goquery.Document
	reqFuncs          *[]func(req *http.Request)
	selectorQueue     *[]cssSelector
}

func NewGoat() (*Goat, error) {
	goat := Goat{
		// URL:               url,
		MaxRecursionDepth: 3,
		EnableConcurrency: false,
		EnableLogging:     false,
		selectorQueue:     new([]cssSelector),
	}

	// if err := goat.newRequest(); err != nil {
	// 	return nil, err
	// }

	if goat.EnableLogging {
		goat.SetRequest(func(req *http.Request) {
			log.Println(req.URL)
		})
	}

	return &goat, nil
}

func (g *Goat) Scrape(url string) {
	req, err := newRequest(g, url)
	if err != nil {
		log.Panicln(err)
	}

	client := new(http.Client)
	res, err := client.Do(req)
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

	for _, cs := range *g.selectorQueue {
		g.doc.Find(cs.selector).Each(func(i int, gs *goquery.Selection) {
			if g.EnableLogging {
				log.Printf("url: %s, selector: %s\n", req.URL, cs.selector)
			}

			cs.callback(Selection{
				gs:            gs,
				selectorQueue: new([]cssSelector),
			})
		})
	}
}

func (g *Goat) SetRequest(callback func(req *http.Request)) {
	*g.reqFuncs = append(*g.reqFuncs, callback)
}

func (g *Goat) SetSelector(selector string, callback func(s Selection)) {
	setSelectorHelper(g, selector, callback)
}

func (g *Goat) getSelectorQueue() *[]cssSelector {
	return g.selectorQueue
}

func (g *Goat) getReqFuncs() *[]func(req *http.Request) {
	return g.reqFuncs
}
