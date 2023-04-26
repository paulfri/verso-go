package common

import (
	"fmt"
	"strconv"
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

const (
	LongItemIDPrefix = "tag:google.com,2005:reader/item/"
)

func StreamIDType(streamID string) string {
	switch norm := strings.ReplaceAll(streamID, " ", ""); norm {
	case StreamIDReadingList:
		return StreamIDReadingList
	case StreamIDBroadcastFriends:
		return StreamIDBroadcastFriends
	case StreamIDStarred:
		return StreamIDStarred
	case StreamIDRead:
		return StreamIDRead
	case StreamIDKeptUnread:
		return StreamIDKeptUnread
	case StreamIDBroadcast:
		return StreamIDBroadcast
	case StreamIDLiked:
		return StreamIDLiked
	default:
		if FeedURLFromReaderStreamID(streamID) != "" {
			return StreamIDFormatFeed
		} else if UserLabelFromReaderStreamID(streamID) != "" {
			return StreamIDFormatLabel
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

func UserLabelFromReaderStreamID(label string) string {
	var feedURL string
	if _, err := fmt.Sscanf(label, StreamIDFormatLabel, &feedURL); err != nil {
		return ""
	}

	return feedURL
}

func ReaderStreamIDFromUserLabel(label string) string {
	return fmt.Sprintf(StreamIDFormatLabel, label)
}

// https://github.com/mihaip/google-reader-api/blob/master/wiki/ItemId.wiki?plain=1

// Returns the long-form item ID for the given UUID. This is a zero-padded,
// 16-length unsigned hex string with a static prefix.
func LongItemID(readerID string) string {
	return fmt.Sprintf(LongItemIDPrefix+"%s", readerID)
}

func ReaderIDFromInput(input string) string {
	// If the input leads with the long-form prefix, parse the identifier as hex.
	if strings.HasPrefix(input, LongItemIDPrefix) {
		return input[32:]
	}

	// Otherwise, parse the input as an integer, and convert it to hex.
	val, _ := strconv.ParseInt(input, 16, 64)

	return strconv.FormatInt(int64(val), 16)
}

func ShortIDFromReaderID(readerID string) string {
	// Convert the hex string to an unsigned integer.
	val, _ := strconv.ParseUint(readerID, 16, 64)

	return fmt.Sprintf("%d", val)
}
