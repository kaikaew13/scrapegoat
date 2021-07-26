package scrapegoat

import (
	"github.com/PuerkitoBio/goquery"
)

type cssSelector struct {
	selector string
	callback func(s Selection)
}

type Selection struct {
	*goquery.Selection
	selectorQueue []cssSelector
}

func (s *Selection) SetSelector(selector string, callback func(s Selection)) {
	s.selectorQueue = append(s.selectorQueue, cssSelector{
		selector: selector,
		callback: callback,
	})
}

func (s *Selection) Scrape() {
	for _, each := range s.selectorQueue {
		s.ChildrenFiltered(each.selector).Each(func(i int, gs *goquery.Selection) {
			each.callback(Selection{
				gs,
				[]cssSelector{},
			})
		})
	}
}
