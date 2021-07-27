package scrapegoat

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Selection struct {
	gs                *goquery.Selection
	selectorQueue     *[]cssSelector
	reqFuncs          *[]func(req *http.Request)
	maxScrapingDepth  uint
	curScrapingDepth  uint
	enableConcurrency bool
	enableLogging     bool
}

func newSelection(scraper Scraper, gs *goquery.Selection) *Selection {
	mrd, crd, ec, el := getOptions(scraper)

	return &Selection{
		gs:                gs,
		selectorQueue:     new([]cssSelector),
		reqFuncs:          new([]func(req *http.Request)),
		maxScrapingDepth:  mrd,
		curScrapingDepth:  crd,
		enableConcurrency: ec,
		enableLogging:     el,
	}
}

func (s *Selection) Scrape(url string) error {
	if s.curScrapingDepth >= s.maxScrapingDepth {
		if s.enableLogging {
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

	for _, cs := range *s.selectorQueue {
		doc.Find(cs.selector).Each(func(i int, gs *goquery.Selection) {
			if s.enableLogging {
				var indent string
				for i := 0; uint(i) < s.curScrapingDepth; i++ {
					indent += "\t"
				}

				fmt.Printf("%surl: %s, selector: %s\n", indent, req.URL, cs.selector)
			}

			cs.selectorFunc(*newSelection(s, gs))
		})
	}

	return nil
}

func (s *Selection) ChildrenSelector(selector string, selectorFunc func(s Selection)) {
	s.gs.ChildrenFiltered(selector).Each(func(i int, gs *goquery.Selection) {
		if s.enableLogging {
			var indent string
			for i := 0; uint(i) < s.curScrapingDepth-1; i++ {
				indent += "\t"
			}

			fmt.Printf("%s- child selector: %s\n", indent, selector)
		}

		selectorFunc(*newSelection(s, gs))
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
