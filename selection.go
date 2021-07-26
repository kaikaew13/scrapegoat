package scrapegoat

import "github.com/PuerkitoBio/goquery"

type Selection struct {
	selector string
	callback func(s *goquery.Selection)
}
