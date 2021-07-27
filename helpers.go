package scrapegoat

import (
	"net/http"
)

type Scraper interface {
	Scrape(url string)
	getSelectorQueue() *[]cssSelector
	getReqFuncs() *[]func(req *http.Request)
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
