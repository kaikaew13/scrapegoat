package scrapegoat

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Selection struct {
	gs               *goquery.Selection
	curScrapingDepth uint
	opts             options
	selectorQueue    *[]cssSelector
	reqFuncs         *[]func(req *http.Request)
}

func newSelection(opts *options, depth uint, gs *goquery.Selection) *Selection {
	return &Selection{
		gs:               gs,
		curScrapingDepth: depth,
		opts:             *opts,
		selectorQueue:    new([]cssSelector),
		reqFuncs:         new([]func(req *http.Request)),
	}
}

func (s *Selection) Scrape(url string) error {
	if s.curScrapingDepth >= s.opts.maxScrapingDepth {
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
				s.scrapeSelector(doc, css, req.URL.String())
			}(cs)
		} else {
			s.scrapeSelector(doc, cs, req.URL.String())
		}
	}

	wg.Wait()
	return nil
}

func (s *Selection) scrapeSelector(doc *goquery.Document, cs cssSelector, url string) {
	sel := doc.Find(cs.selector)

	if s.opts.enableConcurrency {
		deltas := sel.Length()

		var wg sync.WaitGroup
		var mu sync.Mutex

		wg.Add(deltas)

		sel.Each(func(i int, gs *goquery.Selection) {
			go func(gqs *goquery.Selection) {
				defer wg.Done()

				if s.opts.enableLogging {
					s.log(url, cs.selector, int(s.curScrapingDepth))
				}

				mu.Lock()
				cs.selectorFunc(*newSelection(&s.opts, s.curScrapingDepth+1, gqs))
				mu.Unlock()
			}(gs)
		})

		wg.Wait()
	} else {
		sel.Each(func(i int, gs *goquery.Selection) {
			if s.opts.enableLogging {
				s.log(url, cs.selector, int(s.curScrapingDepth))
			}

			cs.selectorFunc(*newSelection(&s.opts, s.curScrapingDepth+1, gs))
		})
	}
}

func (s *Selection) ChildrenSelector(selector string, selectorFunc func(s Selection)) {
	s.gs.ChildrenFiltered(selector).Each(func(i int, gs *goquery.Selection) {
		if s.opts.enableLogging {
			s.log("", selector, int(s.curScrapingDepth)-1)
		}

		selectorFunc(*newSelection(&s.opts, s.curScrapingDepth, gs))
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

func (s *Selection) Text() string {
	return strings.TrimSpace(s.gs.Text())
}

func (s *Selection) Attr(attr string) (val string, exist bool) {
	return s.gs.Attr(attr)
}

func (s *Selection) log(url, selector string, indent int) {
	var ind string
	for i := 0; i < indent; i++ {
		ind += "\t"
	}

	if url == "" {
		fmt.Printf("%s- child selector: %s\n", ind, selector)
		return
	}

	fmt.Printf("%surl: %s, selector: %s\n", ind, url, selector)
}
