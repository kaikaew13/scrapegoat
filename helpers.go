package scrapegoat

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNewReq = errors.New("failed to get a request")
	ErrNewDoc = errors.New("failed to get a document")

	userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	}
)

type Scraper interface {
	Scrape(url string) error
	getSelectorQueue() *[]cssSelector
	getReqFuncs() *[]func(req *http.Request)
}

type cssSelector struct {
	selector     string
	selectorFunc func(s Selection)
}

func setSelectorHelper(scraper Scraper, selector string, selectorFunc func(s Selection)) {
	sq := scraper.getSelectorQueue()

	*sq = append(*sq, cssSelector{
		selector:     selector,
		selectorFunc: selectorFunc,
	})
}

func newRequest(scraper Scraper, url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", randomUserAgent())

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

	if res.StatusCode != 200 {
		fmt.Println(res.Status)
		return nil, fmt.Errorf("got response with status: %s", res.Status)
	}

	return goquery.NewDocumentFromReader(res.Body)
}

func getOptions(scraper Scraper) (mrd, crd uint, ec, el bool) {
	switch t := scraper.(type) {
	case *Goat:
		return t.MaxScrapingDepth, t.curScrapingDepth + 1, t.EnableConcurrency, t.EnableLogging
	case *Selection:
		return t.maxScrapingDepth, t.curScrapingDepth + 1, t.enableConcurrency, t.enableLogging
	}

	return 0, 0, false, false
}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	return userAgents[rand.Int()%len(userAgents)]
}
