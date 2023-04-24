package common

import (
	"testing"
)

var streamIDTypeTestCases = []struct {
	streamID     string
	streamIDType string
}{
	{"user/-/state/com.google/reading-list", StreamIDReadingList},
	{"user/-/state/com.google/broadcast-friends", StreamIDBroadcastFriends},
	{"user/-/state/com.google/starred", StreamIDStarred},
	{"user/-/state/com.google/read", StreamIDRead},
	{"user/-/state/com.google/kept-unread", StreamIDKeptUnread},
	{"user/-/state/com.google/broadcast", StreamIDBroadcast},
	{"user/-/state/com.google/like", StreamIDLiked},
	{"feed/https://blog.paulfri.xyz/atom.xml", StreamIDFormatFeed},
	{"user/-/label/Sounders", StreamIDFormatLabel},
	{"garbage in", ""},
}

func TestStreamIDType(t *testing.T) {
	for _, test := range streamIDTypeTestCases {
		t.Run(test.streamID, func(t *testing.T) {
			streamIDType := StreamIDType(test.streamID)

			if streamIDType != test.streamIDType {
				t.Errorf("failed: got %q, want %q", streamIDType, test.streamIDType)
			}
		})
	}
}

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
				t.Errorf("failed: got %q, want %q", streamID, test.streamID)
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
				t.Errorf("failed: got %q, want %q", url, test.url)
			}
		})
	}
}
