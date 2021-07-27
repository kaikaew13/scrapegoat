package scrapegoat

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const defaultMaxScrapingDepth uint = 3

type Goat struct {
	MaxScrapingDepth  uint
	curScrapingDepth  uint
	EnableConcurrency bool
	EnableLogging     bool
	selectorQueue     *[]cssSelector
	reqFuncs          *[]func(req *http.Request)
}

func NewGoat() *Goat {
	return &Goat{
		MaxScrapingDepth:  defaultMaxScrapingDepth,
		EnableConcurrency: false,
		EnableLogging:     false,
		selectorQueue:     new([]cssSelector),
		reqFuncs:          new([]func(req *http.Request)),
	}
}

func (g *Goat) Scrape(url string) error {
	if g.curScrapingDepth >= g.MaxScrapingDepth {
		if g.EnableLogging {
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
		if !g.EnableConcurrency {
			g.scrapeSelector(doc, cs, req)
		} else {
			wg.Add(1)

			go func(css cssSelector) {
				defer wg.Done()
				g.scrapeSelector(doc, css, req)
			}(cs)
		}
	}

	wg.Wait()
	return nil
}

func (g *Goat) scrapeSelector(doc *goquery.Document, cs cssSelector, req *http.Request) {
	sel := doc.Find(cs.selector)
	deltas := sel.Length()

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(deltas)

	sel.Each(func(i int, gs *goquery.Selection) {
		go func(gqs *goquery.Selection) {
			defer wg.Done()

			mu.Lock()
			cs.selectorFunc(*newSelection(g, gqs))
			mu.Unlock()
		}(gs)
	})

	wg.Wait()

	// doc.Find(cs.selector).Each(func(i int, gs *goquery.Selection) {
	// 	if g.EnableLogging {
	// 		fmt.Printf("url: %s, selector: %s\n", req.URL, cs.selector)
	// 	}

	// 	cs.selectorFunc(*newSelection(g, gs))
	// })
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
