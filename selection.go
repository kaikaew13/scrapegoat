package scrapegoat

import (
	"log"
	"net/http"

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

func (s *Selection) Scrape(url string) {
	req, err := newRequest(s, url)
	if err != nil {
		log.Panicln(err)
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panicln(err)
	}

	for _, cs := range *s.selectorQueue {
		doc.Find(cs.selector).Each(func(i int, gs *goquery.Selection) {
			cs.callback(Selection{
				gs:            gs,
				selectorQueue: new([]cssSelector),
			})
		})
	}

	// for _, each := range *s.selectorQueue {
	// 	s.gs.ChildrenFiltered(each.selector).Each(func(i int, gs *goquery.Selection) {
	// 		each.callback(Selection{
	// 			gs:            gs,
	// 			selectorQueue: new([]cssSelector),
	// 		})
	// 	})
	// }
}

// func (s *Selection) SetChildrenSelector(selector string,callback func(s Selection)) {
// 	s.gs.ChildrenFiltered(selector).Each()
// }

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
	return s.gs.Text()
}

func (s *Selection) Attr(attr string) (val string, exist bool) {
	return s.gs.Attr(attr)
}
