package scrapegoat

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Selection struct {
	gs                *goquery.Selection
	selectorQueue     *[]cssSelector
	reqFuncs          *[]func(req *http.Request)
	maxRecursionDepth uint
	curRecursionDepth uint
	enableConcurrency bool
	enableLogging     bool
}

func newSelection(scraper Scraper, gs *goquery.Selection) *Selection {
	mrd, crd, ec, el := getOptions(scraper)

	return &Selection{
		gs:                gs,
		selectorQueue:     new([]cssSelector),
		reqFuncs:          new([]func(req *http.Request)),
		maxRecursionDepth: mrd,
		curRecursionDepth: crd,
		enableConcurrency: ec,
		enableLogging:     el,
	}
}

func (s *Selection) Scrape(url string) error {
	if s.curRecursionDepth >= s.maxRecursionDepth {
		if s.enableLogging {
			log.Println("[maximum recursion depth reached]")
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
				for i := 0; uint(i) < s.curRecursionDepth; i++ {
					indent += "\t"
				}

				log.Printf("%surl: %s, selector: %s\n", indent, req.URL, cs.selector)
			}

			cs.callback(*newSelection(s, gs))
		})
	}

	return nil
}

func (s *Selection) ChildrenSelector(selector string, callback func(s Selection)) {
	s.gs.ChildrenFiltered(selector).Each(func(i int, gs *goquery.Selection) {
		if s.enableLogging {
			var indent string
			for i := 0; uint(i) < s.curRecursionDepth-1; i++ {
				indent += "\t"
			}

			log.Printf("%s- child selector: %s\n", indent, selector)
		}

		callback(*newSelection(s, gs))
	})
}

func (s *Selection) SetRequest(callback func(req *http.Request)) {
	*s.reqFuncs = append(*s.reqFuncs, callback)
}

func (s *Selection) SetSelector(selector string, callback func(s Selection)) {
	setSelectorHelper(s, selector, callback)
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
