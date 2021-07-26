package scrapegoat

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Goat struct {
	URL               string
	Client            *http.Client
	MaxRecursionDepth int
	EnableConcurrency bool
}

func NewGoat(url string) *Goat {
	goat := Goat{
		URL:               url,
		Client:            new(http.Client),
		MaxRecursionDepth: 0,
		EnableConcurrency: false,
	}

	return &goat
}

func (g *Goat) Scrape() []string {
	req, err := http.NewRequest(http.MethodGet, g.URL, nil)
	if err != nil {
		log.Panicln(err)
	}

	res, err := g.Client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Panicf("got a response with response code of %d, want %d", res.StatusCode, http.StatusOK)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panicln(err)
	}

	data := []string{}

	doc.Find(".markdown-body h2").Each(func(i int, s *goquery.Selection) {
		data = append(data, s.Text())
	})

	return data
}
