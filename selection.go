package scrapegoat

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Selection struct {
	gs            *goquery.Selection
	opts          options
	selectorQueue *[]cssSelector
	reqFuncs      *[]func(req *http.Request)
}

func newSelection(opts *options, gs *goquery.Selection) *Selection {
	return &Selection{
		gs:            gs,
		opts:          *opts,
		selectorQueue: new([]cssSelector),
		reqFuncs:      new([]func(req *http.Request)),
	}
}

func (s *Selection) Scrape(url string) error {
	if s.opts.curScrapingDepth >= s.opts.maxScrapingDepth {
		if s.opts.enableLogging {
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

	for _, cs := range *s.selectorQueue {
		if s.opts.enableConcurrency {
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

func (s *Selection) ChildrenSelector(selector string, selectorFunc func(s Selection)) {
	s.gs.ChildrenFiltered(selector).Each(func(i int, gs *goquery.Selection) {
		if s.opts.enableLogging {
			log(s, "", selector)
		}

		selectorFunc(*newSelection(&s.opts, gs))
	})
}

func (s *Selection) SetRequest(selectorFunc func(req *http.Request)) {
	*s.reqFuncs = append(*s.reqFuncs, selectorFunc)
}

func (s *Selection) SetSelector(selector string, selectorFunc func(s Selection)) {
	setSelectorHelper(s, selector, selectorFunc)
}

func (s *Selection) getSelectorQueue() *[]cssSelector {
	return s.selectorQueue
}

func (s *Selection) getReqFuncs() *[]func(req *http.Request) {
	return s.reqFuncs
}

func (s *Selection) getOptions() *options {
	return &s.opts
}

func (s *Selection) Text() string {
	return strings.TrimSpace(s.gs.Text())
}

func (s *Selection) Attr(attr string) (val string, exist bool) {
	return s.gs.Attr(attr)
}
