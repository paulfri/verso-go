package common

import (
	"fmt"
	"strings"

	"github.com/versolabs/verso/core/helper"
)

func ReaderStreamIDFromFeedURL(url string) string {
	normalized, err := helper.NormalizeFeedURL(url)

	if err != nil {
		return ""
	}

	return fmt.Sprintf("feed/%v", normalized)
}

func FeedURLFromReaderStreamID(streamID string) string {
	parts := strings.Split(streamID, "feed/")

	if len(parts) > 1 && strings.HasPrefix(parts[1], "http") {
		return parts[1]
	}

	return ""
}
