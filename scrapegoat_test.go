package scrapegoat

import (
	"log"
	"net/http"
	"testing"
)

func TestSetRequest(t *testing.T) {
	goat, _ := NewGoat("https://github.com/PuerkitoBio/goquery")

	goat.SetRequest(func(req *http.Request) {
		req.Header.Add("test", "abc")
	})

	want := "abc"

	if goat.req.Header.Get("test") != want {
		t.Errorf("want test header to have val of %s, got %s", want, goat.req.Header.Get("referer"))
	}
}

func TestSetSelector(t *testing.T) {
	goat, _ := NewGoat("https://github.com/PuerkitoBio/goquery")
	goat.EnableLogging = true

	data := []string{}

	goat.SetSelector(".markdown-body h2", func(s Selection) {
		data = append(data, s.Text())
	})

	goat.SetSelector(".markdown-body h1", func(s Selection) {
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

// func TestNested(t *testing.T) {
// 	goat, _ := NewGoat("https://pkg.go.dev/github.com/PuerkitoBio/goquery")

// 	goat.SetSelector(".go-Main-navDesktop .js-readmeOutline ul li", func(s Selection) {
// 		log.Println(s.Text())
// 		// s.Find("a").Each(func(i int, ss *goquery.Selection) {
// 		// 	val, _ := ss.Attr("href")
// 		// 	log.Println(val)
// 		// })

// 		s.SetSelector("a", func(ss Selection) {
// 			val, _ := ss.Attr("href")
// 			log.Println(val)
// 		})

// 		s.Scrape()
// 	})

// 	goat.Scrape()
// }

func TestNestedSetSelector(t *testing.T) {
	goat, _ := NewGoat("https://github.com/PuerkitoBio/goquery")

	data := []string{}

	goat.SetSelector(".markdown-body p:nth-child(27) ", func(s Selection) {
		data = append(data, s.Text())

		s.SetSelector("a", func(ss Selection) {
			val, _ := ss.Attr("href")
			log.Println(val)
		})

		s.Scrape()
	})

	goat.Scrape()
}
