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
// 16-length hex string with a static prefix.
//
// In the original Reader API, this was stored as a 64-bit unsigned integer.
// This is an annoying conversion to do all over the place, so we just store
// it as a string.
func LongItemID(readerID string) string {
	return fmt.Sprintf(LongItemIDPrefix+"%s", readerID)
}

const InvalidReaderID = "0000000000000000"

func ReaderIDFromInput(input string) string {
	// Try parsing the reader ID as longform (with a hex string).
	if strings.HasPrefix(input, LongItemIDPrefix) {
		parts := strings.SplitAfter(input, LongItemIDPrefix)

		if len(parts) == 2 {
			id := parts[1]

			if isValidReaderID(id) {
				// Valid ID after prefix.
				return id
			} else if isValidHexString(id) && len(id) < 16 {
				// Provided ID not a valid hex string, but we can just pad it.
				return fmt.Sprintf("%016s", id)
			} else {
				// Provided ID not a valid hex string and too long (info loss).
				return InvalidReaderID
			}
		} else {
			// No ID provided after longform prefix.
			return InvalidReaderID
		}
	}

	// Check if it's a valid reader ID but not submitted with longform prefix.
	if isValidReaderID(input) {
		return input
	}

	// Parse the input as a base-10 signed int (shortform id).
	asInt, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		// Not a valid int. Try to parse it as hex, like it's a longform ID that
		// lost its prefix somehow.
		if isValidHexString(input) {
			// Valid hex. Return input as padded.
			return fmt.Sprintf("%016s", input)
		}

		return InvalidReaderID
	}

	// Valid int, convert to uint.
	asUint := uint64(asInt)

	// Return as padded hex.
	return fmt.Sprintf("%016x", asUint)
}

func ShortIDFromReaderID(readerID string) string {
	// Convert the hex string to an unsigned integer.
	val, _ := strconv.ParseUint(readerID, 16, 64)

	return fmt.Sprintf("%d", val)
}

func isValidReaderID(readerID string) bool {
	return len(readerID) == 16 && isValidHexString(readerID)
}

func isValidHexString(readerID string) bool {
	_, err := strconv.ParseUint(readerID, 16, 64)

	return err == nil
}
