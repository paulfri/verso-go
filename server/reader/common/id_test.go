package common

import (
	"testing"
)

var readerStreamIDToURLCases = []struct {
	url      string
	streamID string
}{
	{"https://blog.paulfri.xyz/atom.xml", "feed/https://blog.paulfri.xyz/atom.xml"},
	{"blog.paulfri.xyz/atom.xml", "feed/https://blog.paulfri.xyz/atom.xml"},
}

func TestReaderStreamIDFromURL(t *testing.T) {
	for _, test := range readerStreamIDToURLCases {
		t.Run(test.url, func(t *testing.T) {
			streamID := ReaderStreamIDFromFeedURL(test.url)

			if streamID != test.streamID {
				t.Errorf("normalize failed: got %q, want %q", streamID, test.streamID)
			}
		})
	}
}

var feedURLToStreamIDCases = []struct {
	streamID string
	url      string
}{
	{"feed/https://blog.paulfri.xyz/atom.xml", "https://blog.paulfri.xyz/atom.xml"},
	{"https://blog.paulfri.xyz/atom.xml", ""},
	{"feed/feed/https://blog.paulfri.xyz/atom.xml", ""},
}

func TestFeedURLFromReaderStreamID(t *testing.T) {
	for _, test := range feedURLToStreamIDCases {
		t.Run(test.url, func(t *testing.T) {
			url := FeedURLFromReaderStreamID(test.streamID)

			if url != test.url {
				t.Errorf("normalize failed: got %q, want %q", url, test.url)
			}
		})
	}
}
