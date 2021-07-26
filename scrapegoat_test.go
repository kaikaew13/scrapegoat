package scrapegoat

import "testing"

func TestScrape(t *testing.T) {
	goat := NewGoat("https://github.com/PuerkitoBio/goquery")

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

	data := goat.Scrape()

	if len(data) != len(want) {
		t.Errorf("want slice of data with length of %d, got %d", len(want), len(data))
	}

	for i := range want {
		if data[i] != want[i] {
			t.Errorf("want data at index %d to be %s, got %s", i, want[i], data[i])
		}
	}
}
