package common

import (
	"fmt"
	"regexp"
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
func LongItemID(readerID int64) string {
	hex := readerIDToHex(readerID)

	return fmt.Sprintf(LongItemIDPrefix+"%s", hex)
}

func ReaderIDFromInput(input string) int64 {
	// If the input leads with the long-form prefix, parse the identifier as hex.
	if strings.HasPrefix(input, LongItemIDPrefix) {
		return readerIDFromHex(input[32:])
	}

	// Some clients zero-pad the ID, so trim the leading zeros.
	unpad := strings.TrimLeft(input, "0")

	// Some clients send short-form IDs as hex, so check for that.
	if hasAlpha(unpad) {
		return readerIDFromHex(unpad)
	}

	i, err := strconv.Atoi(unpad)

	if err != nil {
		panic(err)
	}

	return int64(i)
}

func readerIDFromHex(hex string) int64 {
	val, _ := strconv.ParseInt(hex, 16, 64)
	sixfour := int64(val)

	return sixfour
}

func readerIDToHex(readerID int64) string {
	return strconv.FormatInt(readerID, 16)
}

func hasAlpha(str string) bool {
	return regexp.MustCompile(`\D`).MatchString(str)
}
