package scrapegoat

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// func TestScrape(t *testing.T) {
// 	goat, _ := NewGoat("https://github.com/PuerkitoBio/goquery", DefaultOptions)

// 	want := []string{
// 		"Table of Contents",
// 		"Installation",
// 		"Changelog",
// 		"API",
// 		"Examples",
// 		"Related Projects",
// 		"Support",
// 		"License",
// 	}

// 	goat.SetRequest(func(req *http.Request) {
// 		req.Header.Add("test", "abc")
// 	})

// 	data := goat.Scrape()

// 	log.Println(goat.req.Header.Get("test"))

// 	if len(data) != len(want) {
// 		t.Errorf("want slice of data with length of %d, got %d", len(want), len(data))
// 	}

// 	for i := range want {
// 		if data[i] != want[i] {
// 			t.Errorf("want data at index %d to be %s, got %s", i, want[i], data[i])
// 		}
// 	}
// }

func TestSetRequest(t *testing.T) {
	goat, _ := NewGoat("https://github.com/PuerkitoBio/goquery", DefaultOptions)

	goat.SetRequest(func(req *http.Request) {
		req.Header.Add("test", "abc")
	})

	want := "abc"

	if goat.req.Header.Get("test") != want {
		t.Errorf("want test header to have val of %s, got %s", want, goat.req.Header.Get("referer"))
	}
}

func TestSetTags(t *testing.T) {
	goat, _ := NewGoat("https://github.com/PuerkitoBio/goquery", Options{
		EnableLogging: true,
	})

	data := []string{}

	goat.SetTags(".markdown-body h2", func(s *goquery.Selection) {
		data = append(data, s.Text())
	})

	goat.SetTags(".markdown-body h1", func(s *goquery.Selection) {
		data = append(data, s.Text())
	})

	goat.Scrape()

	want := []string{
		"Table of Contents",
		"Installation",
		"Changelog",
		"API",
		"Examples",
		"Related Projects",
		"Support",
		"License",
		"goquery - a little like that j-thing, only in Go",
	}

	if len(data) != len(want) {
		t.Errorf("want slice of data with length of %d, got %d", len(want), len(data))
	}

	for i := range want {
		if data[i] != want[i] {
			t.Errorf("want data at index %d to be %s, got %s", i, want[i], data[i])
		}
	}
}
