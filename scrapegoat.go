package scrapegoat

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Goat struct {
	MaxRecursionDepth int
	EnableConcurrency bool
	EnableLogging     bool
	doc               *goquery.Document
	reqFuncs          *[]func(req *http.Request)
	selectorQueue     *[]cssSelector
}

func NewGoat() *Goat {
	return &Goat{
		MaxRecursionDepth: 3,
		EnableConcurrency: false,
		EnableLogging:     false,
		selectorQueue:     new([]cssSelector),
		reqFuncs:          new([]func(req *http.Request)),
	}
}

func (g *Goat) Scrape(url string) {
	req, err := newRequest(g, url)
	if err != nil {
		log.Panicln(ErrNewRequest, err)
	}

	g.doc, err = getDocumentFromRequest(req)
	if err != nil {
		log.Panicln(ErrNewDoc, err)
	}

	for _, cs := range *g.selectorQueue {
		g.doc.Find(cs.selector).Each(func(i int, gs *goquery.Selection) {
			if g.EnableLogging {
				log.Printf("url: %s, selector: %s\n", req.URL, cs.selector)
			}

			cs.callback(*newSelection(gs))
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
