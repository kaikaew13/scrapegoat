package scrapegoat

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type cssSelector struct {
	selector string
	callback func(s Selection)
}

type Selection struct {
	gs            *goquery.Selection
	selectorQueue *[]cssSelector
	reqFuncs      *[]func(req *http.Request)
}

func newSelection(gs *goquery.Selection) *Selection {
	return &Selection{
		gs:            gs,
		selectorQueue: new([]cssSelector),
		reqFuncs:      new([]func(req *http.Request)),
	}
}

func (s *Selection) Scrape(url string) error {
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
			cs.callback(*newSelection(gs))
		})
	}

	return nil
}

func (s *Selection) SetChildrenSelector(selector string, callback func(s Selection)) {
	s.gs.ChildrenFiltered(selector).Each(func(i int, gs *goquery.Selection) {
		callback(*newSelection(gs))
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
