package scrapegoat

import (
	"errors"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNewReq = errors.New("failed to get a request")
	ErrNewDoc = errors.New("failed to get a document")
)

type Scraper interface {
	Scrape(url string) error
	getSelectorQueue() *[]cssSelector
	getReqFuncs() *[]func(req *http.Request)
}

type cssSelector struct {
	selector string
	callback func(s Selection)
}

func setSelectorHelper(scraper Scraper, selector string, callback func(s Selection)) {
	sq := scraper.getSelectorQueue()

	*sq = append(*sq, cssSelector{
		selector: selector,
		callback: callback,
	})
}

func newRequest(scraper Scraper, url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	reqFuncs := scraper.getReqFuncs()
	if reqFuncs != nil {
		for _, fn := range *reqFuncs {
			fn(req)
		}
	}

	return req, nil
}

func getDocumentFromRequest(req *http.Request) (*goquery.Document, error) {
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return goquery.NewDocumentFromReader(res.Body)
}

func getOptions(scraper Scraper) (mrd, crd int, ec, el bool) {
	switch t := scraper.(type) {
	case *Goat:
		return t.MaxRecursionDepth, t.curRecursionDepth + 1, t.EnableConcurrency, t.EnableLogging
	case *Selection:
		return t.maxRecursionDepth, t.curRecursionDepth + 1, t.enableConcurrency, t.enableLogging
	}

	return 0, 0, false, false
}
