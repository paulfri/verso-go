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

var readerIDFromInputCases = []struct {
	in  string
	out string
}{
	// With longform prefix, reader ID returned verbatim.
	{"tag:google.com,2005:reader/item/14d367199b16429c", "14d367199b16429c"},
	// With longform prefix and invalid reader ID (too short), return padded.
	{"tag:google.com,2005:reader/item/123123", "0000000000123123"},
	// With longform prefix and invalid reader ID (too long), return all zeroes.
	{"tag:google.com,2005:reader/item/1111111111111111111111111", "0000000000000000"},
	// Valid hex string of correct length without longform prefix, returned as is.
	{"000000000000001a", "000000000000001a"},
	// Short-form identifier parsed as positive int and returned as padded hex string.
	{"123123", "000000000001e0f3"},
	// Short-form identifier parsed as negative int and returned as padded hex string.
	{"-355401917359550817", "fb115bd6d34a8e9f"},
	// Short-form identifier of unpadded hex, return padded.
	{"1a", "000000000000001a"},
	// Not a valid hex string but correct length, return all zeroes.
	{"asdf/abc123qwer1", "0000000000000000"},
	// Invalid on many levels, return all zeroes.
	{"asdf/abc123qwer9999", "0000000000000000"},
}

func TestReaderIDFromInput(t *testing.T) {
	for _, test := range readerIDFromInputCases {
		t.Run(test.in, func(t *testing.T) {
			res := ReaderIDFromInput(test.in)

			if res != test.out {
				t.Errorf("failed: got %q, want %q", res, test.out)
			}
		})
	}
}
