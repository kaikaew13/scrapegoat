package scrapegoat

import "github.com/PuerkitoBio/goquery"

type cssSelector struct {
	selector string
	callback func(s Selection)
}

type Selection struct {
	*goquery.Selection
}
