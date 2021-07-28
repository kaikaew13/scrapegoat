package scrapegoat

import (
	"net/http"
	"testing"
)

const testingURL string = "https://github.com/PuerkitoBio/goquery"

func TestSetRequest(t *testing.T) {
	goat := NewGoat()

	goat.SetRequest(func(req *http.Request) {
		req.Header.Add("test", "abc")
	})

	req, err := newRequest(goat, testingURL)
	if err != nil {
		t.Fatal(err)
	}

	want := "abc"

	if req.Header.Get("test") != want {
		t.Errorf("want test header to have val of %s, got %s", want, req.Header.Get("test"))
	}
}

func TestSetSelector(t *testing.T) {
	goat := NewGoat(
		EnableLogging(true),
	)

	data := []string{}

	goat.SetSelector(".markdown-body h2", func(s Selection) {
		data = append(data, s.Text())
	})

	goat.SetSelector(".markdown-body h1", func(s Selection) {
		data = append(data, s.Text())
	})

	if err := goat.Scrape(testingURL); err != nil {
		t.Fatal(err)
	}

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

func TestSetChildrenSelector(t *testing.T) {
	goat := NewGoat(
		EnableLogging(true),
	)

	data := []string{}

	goat.SetSelector(".markdown-body", func(s Selection) {
		s.ChildrenSelector("h2", func(child Selection) {
			data = append(data, child.Text())
		})
	})

	if err := goat.Scrape(testingURL); err != nil {
		t.Fatal(err)
	}

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

	compareSliceHelper(t, want, data)
}

func TestNestedSetSelector(t *testing.T) {
	tests := []struct {
		desc string
		msd  uint
		ec   bool
		want []string
	}{
		{
			desc: "set MaxScrapingDepth to 2 - should be able to fetch data",
			msd:  2,
			ec:   false,
			want: []string{
				"Handle Non-UTF8 html Pages",
				"Handle Javascript-based Pages",
				"For Loop",
			},
		},
		{
			desc: "set MaxScrapingDepth to 1 - should not be able to fetch data",
			msd:  1,
			ec:   false,
			want: []string{},
		},
		{
			desc: "concurrent with nested scraping - should be able to fetch data",
			msd:  2,
			ec:   true,
			want: []string{
				"Handle Non-UTF8 html Pages",
				"Handle Javascript-based Pages",
				"For Loop",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			it := 1
			if tt.ec {
				it = 10
			}

			for i := 0; i < it; i++ {
				goat := NewGoat(
					EnableConcurrency(tt.ec),
					EnableLogging(true),
					MaxScrapingDepth(tt.msd),
				)

				dataMap := map[string]bool{}
				dataSlice := []string{}

				goat.SetSelector(".markdown-body p:nth-child(27) a", func(s Selection) {
					val, _ := s.Attr("href")
					s.SetSelector(".markdown-body h2", func(ss Selection) {
						if tt.ec {
							dataMap[ss.Text()] = true
						} else {

							dataSlice = append(dataSlice, ss.Text())
						}
					})

					if err := s.Scrape(val); err != nil {
						t.Fatal(err)
					}

				})

				if err := goat.Scrape(testingURL); err != nil {
					t.Fatal(err)
				}

				if tt.ec {
					compareMapHelper(t, tt.want, dataMap)
				} else {
					compareSliceHelper(t, tt.want, dataSlice)
				}
			}
		})
	}
}

func TestEnableConcurrency(t *testing.T) {
	for i := 0; i < 10; i++ {
		goat := NewGoat(
			EnableConcurrency(true),
			EnableLogging(true),
		)

		data := map[string]bool{}

		goat.SetSelector(".markdown-body h2", func(s Selection) {
			data[s.Text()] = true
		})

		goat.SetSelector(".markdown-body h1", func(s Selection) {
			data[s.Text()] = true
		})

		if err := goat.Scrape(testingURL); err != nil {
			t.Fatal(err)
		}

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

		compareMapHelper(t, want, data)
	}
}

func compareMapHelper(t testing.TB, want []string, got map[string]bool) {
	if len(got) != len(want) {
		t.Errorf("want map of data with length of %d, got %d", len(want), len(got))
	}

	for _, w := range want {
		if !got[w] {
			t.Errorf("want data map to have %s as true, got %t", w, got[w])
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
