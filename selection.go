package scrapegoat

import (
	"github.com/PuerkitoBio/goquery"
)

type cssSelector struct {
	selector string
	callback func(s Selection)
}

type Selection struct {
	gs            *goquery.Selection
	selectorQueue *[]cssSelector
}

func (s *Selection) Scrape() {
	for _, each := range *s.selectorQueue {
		s.gs.ChildrenFiltered(each.selector).Each(func(i int, gs *goquery.Selection) {
			each.callback(Selection{
				gs:            gs,
				selectorQueue: new([]cssSelector),
			})
		})
	}
}

func (s *Selection) SetSelector(selector string, callback func(s Selection)) {
	setSelectorHelper(s, s.selectorQueue, selector, callback)
}

func (s *Selection) Text() string {
	return s.gs.Text()
}

func (s *Selection) Attr(attr string) (val string, exist bool) {
	return s.gs.Attr(attr)
}
