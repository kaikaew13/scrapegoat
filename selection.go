package scrapegoat

import "github.com/PuerkitoBio/goquery"

type selection struct {
	selector string
	callback func(s *goquery.Selection)
}
