package scrapegoat

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const defaultMaxRecursionDepth uint = 3

type Goat struct {
	MaxRecursionDepth uint
	curRecursionDepth uint
	EnableConcurrency bool
	EnableLogging     bool
	selectorQueue     *[]cssSelector
	reqFuncs          *[]func(req *http.Request)
}

func NewGoat() *Goat {
	return &Goat{
		MaxRecursionDepth: defaultMaxRecursionDepth,
		EnableConcurrency: false,
		EnableLogging:     false,
		selectorQueue:     new([]cssSelector),
		reqFuncs:          new([]func(req *http.Request)),
	}
}

func (g *Goat) Scrape(url string) error {
	if g.curRecursionDepth >= g.MaxRecursionDepth {
		if g.EnableLogging {
			fmt.Println("[maximum recursion depth reached]")
		}

		return nil
	}

	req, err := newRequest(g, url)
	if err != nil {
		return ErrNewReq
	}

	doc, err := getDocumentFromRequest(req)
	if err != nil {
		return ErrNewDoc
	}

	for _, cs := range *g.selectorQueue {
		doc.Find(cs.selector).Each(func(i int, gs *goquery.Selection) {
			if g.EnableLogging {
				fmt.Printf("url: %s, selector: %s\n", req.URL, cs.selector)
			}

			cs.selectorFunc(*newSelection(g, gs))
		})
	}

	return nil
}

func (g *Goat) SetRequest(selectorFunc func(req *http.Request)) {
	*g.reqFuncs = append(*g.reqFuncs, selectorFunc)
}

func (g *Goat) SetSelector(selector string, selectorFunc func(s Selection)) {
	setSelectorHelper(g, selector, selectorFunc)
}

func (g *Goat) getSelectorQueue() *[]cssSelector {
	return g.selectorQueue
}

func (g *Goat) getReqFuncs() *[]func(req *http.Request) {
	return g.reqFuncs
}
