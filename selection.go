package scrapegoat

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Selection struct {
	gs   *goquery.Selection
	goat *Goat
}

func newSelection(opts *options, gs *goquery.Selection) *Selection {
	g := NewGoat()
	g.opts = *opts

	return &Selection{
		gs:   gs,
		goat: g,
	}
}

func (s *Selection) Scrape(url string) error {
	g := s.goat

	if g.opts.curScrapingDepth >= g.opts.maxScrapingDepth {
		if g.opts.enableLogging {
			fmt.Println("[maximum scraping depth reached]")
		}

		return nil
	}

	req, err := newRequest(s, url)
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
				scrapeSelector(s, doc, css, req.URL.String())
			}(cs)
		} else {
			scrapeSelector(s, doc, cs, req.URL.String())
		}
	}

	wg.Wait()
	return nil
}

func (s *Selection) ChildrenSelector(selector string, selectorFunc func(child Selection)) {
	g := s.goat

	s.gs.Find(selector).Each(func(i int, gs *goquery.Selection) {
		if g.opts.enableLogging {
			log(s, "", selector)
		}

		selectorFunc(*newSelection(&g.opts, gs))
	})
}

func (s *Selection) SetRequest(selectorFunc func(req *http.Request)) {
	g := s.goat
	*g.reqFuncs = append(*g.reqFuncs, selectorFunc)
}

func (s *Selection) SetSelector(selector string, selectorFunc func(ss Selection)) {
	setSelectorHelper(s, selector, selectorFunc)
}

func (s *Selection) getGoat() *Goat {
	return s.goat
}

func (s *Selection) Text() string {
	return strings.TrimSpace(s.gs.Text())
}

func (s *Selection) Attr(attr string) (val string, exist bool) {
	return s.gs.Attr(attr)
}
