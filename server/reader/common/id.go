package common

import (
	"fmt"
	"strings"

	"github.com/versolabs/verso/core/helper"
)

const (
	// Subscription-level
	StreamIDReadingList      = "user/-/state/com.google/reading-list"
	StreamIDBroadcastFriends = "user/-/state/com.google/broadcast-friends"

	// Item-level
	StreamIDStarred    = "user/-/state/com.google/starred"
	StreamIDRead       = "user/-/state/com.google/read"
	StreamIDKeptUnread = "user/-/state/com.google/kept-unread"
	StreamIDBroadcast  = "user/-/state/com.google/broadcast"
	StreamIDLiked      = "user/-/state/com.google/like"
)

const (
	StreamIDFormatFeed  = "feed/%s"
	StreamIDFormatLabel = "user/-/label/%s"
)

func StreamIDType(streamID string) string {
	switch streamID {
	case StreamIDReadingList:
		return StreamIDReadingList
	case StreamIDBroadcastFriends:
		return StreamIDBroadcastFriends
	case StreamIDStarred:
		return StreamIDStarred
	default:
		if FeedURLFromReaderStreamID(streamID) != "" {
			return StreamIDFormatFeed
		}

		return ""
	}
}

func ReaderStreamIDFromFeedURL(url string) string {
	normalized, err := helper.NormalizeFeedURL(url)

	if err != nil {
		return ""
	}

	return fmt.Sprintf(StreamIDFormatFeed, normalized)
}

func FeedURLFromReaderStreamID(streamID string) string {
	var feedURL string
	if _, err := fmt.Sscanf(streamID, StreamIDFormatFeed, &feedURL); err != nil {
		return ""
	}

	// Sanity check that it's an http/s URL.
	if !strings.HasPrefix(feedURL, "http") {
		return ""
	}

	return feedURL
}
