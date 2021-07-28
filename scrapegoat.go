package scrapegoat

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Goat struct {
	curScrapingDepth uint
	opts             options
	selectorQueue    *[]cssSelector
	reqFuncs         *[]func(req *http.Request)
}

func NewGoat(opts ...optionFunc) *Goat {
	goat := &Goat{
		selectorQueue: new([]cssSelector),
		reqFuncs:      new([]func(req *http.Request)),
		opts:          defaultOptions,
	}

	for _, opt := range opts {
		opt(&goat.opts)
	}

	return goat
}

func (g *Goat) Scrape(url string) error {
	if g.curScrapingDepth >= g.opts.maxScrapingDepth {
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
				g.scrapeSelector(doc, css, req.URL.String())
			}(cs)
		} else {
			g.scrapeSelector(doc, cs, req.URL.String())
		}
	}

	wg.Wait()
	return nil
}

func (g *Goat) scrapeSelector(doc *goquery.Document, cs cssSelector, url string) {
	sel := doc.Find(cs.selector)

	if g.opts.enableConcurrency {
		deltas := sel.Length()

		var wg sync.WaitGroup
		var mu sync.Mutex

		wg.Add(deltas)

		sel.Each(func(i int, gs *goquery.Selection) {
			go func(gqs *goquery.Selection) {
				defer wg.Done()

				if g.opts.enableLogging {
					g.log(url, cs.selector)
				}

				mu.Lock()
				cs.selectorFunc(*newSelection(&g.opts, g.curScrapingDepth+1, gqs))
				mu.Unlock()
			}(gs)
		})

		wg.Wait()
	} else {
		sel.Each(func(i int, gs *goquery.Selection) {
			if g.opts.enableLogging {
				g.log(url, cs.selector)
			}

			cs.selectorFunc(*newSelection(&g.opts, g.curScrapingDepth+1, gs))
		})
	}
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

func (g *Goat) log(url, selector string) {
	fmt.Printf("url: %s, selector: %s\n", url, selector)
}
