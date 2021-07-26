package scrapegoat

import (
	"log"
	"net/http"
	"testing"
)

func TestScrape(t *testing.T) {
	goat, _ := NewGoat("https://github.com/PuerkitoBio/goquery", DefaultOptions)

	want := []string{
		"Table of Contents",
		"Installation",
		"Changelog",
		"API",
		"Examples",
		"Related Projects",
		"Support",
		"License",
	}

	goat.SetRequest(func(req *http.Request) {
		req.Header.Add("test", "abc")
	})

	data := goat.Scrape()

	log.Println(goat.req.Header.Get("test"))

	if len(data) != len(want) {
		t.Errorf("want slice of data with length of %d, got %d", len(want), len(data))
	}

	for i := range want {
		if data[i] != want[i] {
			t.Errorf("want data at index %d to be %s, got %s", i, want[i], data[i])
		}
	}
}

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
