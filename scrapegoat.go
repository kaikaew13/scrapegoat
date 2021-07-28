package scrapegoat

import (
	"fmt"
	"net/http"
	"sync"
)

type Goat struct {
	opts          options
	selectorQueue *[]cssSelector
	reqFuncs      *[]func(req *http.Request)
}

func NewGoat(opts ...optionFunc) *Goat {
	g := &Goat{
		opts:          defaultOptions,
		selectorQueue: new([]cssSelector),
		reqFuncs:      new([]func(req *http.Request)),
	}

	for _, opt := range opts {
		opt(&g.opts)
	}

	return g
}

func (g *Goat) Scrape(url string) error {
	if g.opts.curScrapingDepth >= g.opts.maxScrapingDepth {
		if g.opts.enableLogging {
			fmt.Println("[maximum scraping depth reached]")
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

	var wg sync.WaitGroup

	for _, cs := range *g.selectorQueue {
		if g.opts.enableConcurrency {
			wg.Add(1)

			go func(css cssSelector) {
				defer wg.Done()
				scrapeSelector(g, doc, css, req.URL.String())
			}(cs)
		} else {
			scrapeSelector(g, doc, cs, req.URL.String())
		}
	}

	wg.Wait()
	return nil
}

func (g *Goat) SetRequest(selectorFunc func(req *http.Request)) {
	*g.reqFuncs = append(*g.reqFuncs, selectorFunc)
}

func (g *Goat) SetSelector(selector string, selectorFunc func(s Selection)) {
	setSelectorHelper(g, selector, selectorFunc)
}

func (g *Goat) getGoat() *Goat {
	return g
}
