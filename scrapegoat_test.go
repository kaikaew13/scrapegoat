package scrapegoat

import (
	"fmt"
	"testing"
)

const testingURL string = "https://github.com/PuerkitoBio/goquery"

// func TestSetRequest(t *testing.T) {
// 	goat := NewGoat()

// 	goat.SetRequest(func(req *http.Request) {
// 		req.Header.Add("test", "abc")
// 	})

// 	req, err := newRequest(goat, testingURL)
// 	if err != nil {
// 		t.Error("erororor")
// 	}

// 	want := "abc"

// 	if req.Header.Get("test") != want {
// 		t.Errorf("want test header to have val of %s, got %s", want, req.Header.Get("test"))
// 	}
// }

// func TestSetSelector(t *testing.T) {
// 	goat := NewGoat()
// 	goat.EnableLogging = true

// 	data := []string{}

// 	goat.SetSelector(".markdown-body h2", func(s Selection) {
// 		data = append(data, s.Text())
// 	})

// 	goat.SetSelector(".markdown-body h1", func(s Selection) {
// 		data = append(data, s.Text())
// 	})

// 	if err := goat.Scrape(testingURL); err != nil {
// 		t.Error(err)
// 	}

// 	want := []string{
// 		"Table of Contents",
// 		"Installation",
// 		"Changelog",
// 		"API",
// 		"Examples",
// 		"Related Projects",
// 		"Support",
// 		"License",
// 		"goquery - a little like that j-thing, only in Go",
// 	}

// 	compareSliceHelper(t, want, data)
// }

// func TestSetChildrenSelector(t *testing.T) {
// 	goat := NewGoat()
// 	goat.EnableLogging = true

// 	data := []string{}

// 	goat.SetSelector(".markdown-body", func(s Selection) {
// 		s.ChildrenSelector("h2", func(ss Selection) {
// 			data = append(data, ss.Text())
// 		})
// 	})

// 	if err := goat.Scrape(testingURL); err != nil {
// 		t.Error(err)
// 	}

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

// 	compareSliceHelper(t, want, data)
// }

// func TestNestedSetSelector(t *testing.T) {
// 	tests := []struct {
// 		desc string
// 		mrd  uint
// 		want []string
// 	}{
// 		{
// 			desc: "set MaxRecursionDepth to 2 - should be able to fetch data",
// 			mrd:  2,
// 			want: []string{
// 				"Handle Non-UTF8 html Pages",
// 				"Handle Javascript-based Pages",
// 				"For Loop",
// 			},
// 		},
// 		{
// 			desc: "set MaxRecursionDepth to 1 - should not be able to fetch data",
// 			mrd:  1,
// 			want: []string{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.desc, func(t *testing.T) {
// 			goat := NewGoat()
// 			goat.EnableLogging = true
// 			goat.MaxScrapingDepth = tt.mrd

// 			data := []string{}

// 			goat.SetSelector(".markdown-body p:nth-child(27) a", func(s Selection) {
// 				val, _ := s.Attr("href")
// 				s.SetSelector(".markdown-body h2", func(ss Selection) {
// 					data = append(data, ss.Text())
// 				})

// 				if err := s.Scrape(val); err != nil {
// 					t.Error(err)
// 				}

// 			})

// 			if err := goat.Scrape(testingURL); err != nil {
// 				t.Error(err)
// 			}

// 			compareSliceHelper(t, tt.want, data)
// 		})
// 	}
// }

func TestEnableConcurrency(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(i, "----------------------------------------------")
		goat := NewGoat()
		goat.EnableConcurrency = true
		// goat.EnableLogging = true

		data := []string{}

		goat.SetSelector(".markdown-body h2", func(s Selection) {
			data = append(data, s.Text())
		})

		goat.SetSelector(".markdown-body h1", func(s Selection) {
			data = append(data, s.Text())
		})

		goat.Scrape(testingURL)

		if len(data) != 9 {
			t.Fatal(len(data))
		}
		for _, v := range data {
			fmt.Println(v)
		}
	}
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
