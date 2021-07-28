package scrapegoat

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
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
	getGoat() *Goat
}

type cssSelector struct {
	selector     string
	selectorFunc func(s Selection)
}

func setSelectorHelper(scraper Scraper, selector string, selectorFunc func(s Selection)) {
	g := scraper.getGoat()
	sq := g.selectorQueue

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

	req.Header.Add("User-Agent", getRandomUserAgent())

	g := scraper.getGoat()
	reqFuncs := g.reqFuncs
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

func scrapeSelector(scraper Scraper, doc *goquery.Document, cs cssSelector, url string) {
	sel := doc.Find(cs.selector)

	g := scraper.getGoat()
	opts := g.opts
	newOpts := opts
	newOpts.curScrapingDepth++

	if opts.enableConcurrency {
		deltas := sel.Length()

		var wg sync.WaitGroup
		var mu sync.Mutex

		wg.Add(deltas)

		sel.Each(func(i int, gs *goquery.Selection) {
			go func(gqs *goquery.Selection) {
				defer wg.Done()

				if opts.enableLogging {
					log(scraper, url, cs.selector)
				}

				mu.Lock()
				cs.selectorFunc(*newSelection(&newOpts, gqs))
				mu.Unlock()
			}(gs)
		})

		wg.Wait()
	} else {
		sel.Each(func(i int, gs *goquery.Selection) {
			if opts.enableLogging {
				log(scraper, url, cs.selector)
			}

			cs.selectorFunc(*newSelection(&newOpts, gs))
		})
	}
}

func getRandomUserAgent() string {
	rand.Seed(time.Now().Unix())
	return userAgents[rand.Int()%len(userAgents)]
}

func log(scraper Scraper, url, selector string) {
	g := scraper.getGoat()
	opts := g.opts

	var indent string
	for i := 0; i < int(opts.curScrapingDepth); i++ {
		indent += "\t"
	}

	if url == "" {
		fmt.Printf("%s- child selector: %s\n", indent[:len(indent)-1], selector)
		return
	}

	fmt.Printf("%surl: %s, selector: %s\n", indent, url, selector)
}
