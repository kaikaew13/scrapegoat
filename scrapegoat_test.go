package scrapegoat

import (
	"net/http"
	"testing"
)

func TestSetRequest(t *testing.T) {
	goat := NewGoat()

	goat.SetRequest(func(req *http.Request) {
		req.Header.Add("test", "abc")
	})

	req, err := newRequest(goat, "https://github.com/PuerkitoBio/goquery")
	if err != nil {
		t.Error("erororor")
	}

	want := "abc"

	if req.Header.Get("test") != want {
		t.Errorf("want test header to have val of %s, got %s", want, req.Header.Get("test"))
	}
}

func TestSetSelector(t *testing.T) {
	goat := NewGoat()
	goat.EnableLogging = true

	data := []string{}

	goat.SetSelector(".markdown-body h2", func(s Selection) {
		data = append(data, s.Text())
	})

	goat.SetSelector(".markdown-body h1", func(s Selection) {
		data = append(data, s.Text())
	})

	goat.Scrape("https://github.com/PuerkitoBio/goquery")

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

	compareSliceHelper(t, want, data)
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
	goat := NewGoat()

	data := []string{}

	goat.SetSelector(".markdown-body p:nth-child(27) a ", func(s Selection) {

		// s.SetSelector("a", func(ss Selection) {
		// 	val, _ := ss.Attr("href")
		// 	log.Println(val)
		// })

		val, _ := s.Attr("href")
		s.SetSelector(".markdown-body h2", func(ss Selection) {
			// log.Println(ss.Text())
			data = append(data, ss.Text())
		})

		s.Scrape(val)

	})

	goat.Scrape("https://github.com/PuerkitoBio/goquery")

	want := []string{
		"Handle Non-UTF8 html Pages",
		"Handle Javascript-based Pages",
		"For Loop",
	}

	compareSliceHelper(t, want, data)
}

func compareSliceHelper(t testing.TB, want, got []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("want slice of data with length of %d, got %d", len(want), len(got))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("want data at index %d to be %s, got %s", i, want[i], got[i])
		}
	}
}
