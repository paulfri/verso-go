package helper

import (
	"testing"
)

var cases = []struct {
	in  string
	out string
}{
	{"example.com", "https://example.com"},
	{"example.com/feed.rss", "https://example.com/feed.rss"},
	{"EXAMPLE.COM/FEED.RSS", "https://example.com/FEED.RSS"},
	{"http://www.example.com/feed.rss", "http://www.example.com/feed.rss"},
}

func TestNormalizeFeedURL(t *testing.T) {
	for _, test := range cases {
		t.Run(test.in, func(t *testing.T) {
			result, err := NormalizeFeedURL(test.in)

			if err != nil {
				t.Errorf("normalize failed: %v", err)
			}

			if result != test.out {
				t.Errorf("normalize failed: got %q, want %q", result, test.out)
			}
		})
	}
}
