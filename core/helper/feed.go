package helper

import (
	"errors"

	"github.com/mmcdole/gofeed"
)

var ErrInvalidFeedURL = errors.New("invalid feed URL")

// Returns whether the given URL string is a parseable RSS or Atom feed.
func IsValidFeed(feed string) bool {
	parser := gofeed.NewParser()
	_, err := parser.ParseURL(feed)

	return err == nil
}

// Scrapes the given URL to look for any feed URLs in the document content.
// Currently not implemented.
func ScrapeForFeedURLs(_ string) ([]string, error) {
	return nil, nil
}

// For a given (potentially invalid,  potentially non-feed) URL, find the
// closest valid, subscribable feed.
func GatherFeed(url string) (string, error) {
	normalizedURL, err := NormalizeFeedURL(url)

	if err != nil {
		return "", ErrInvalidFeedURL
	}

	if IsValidFeed(normalizedURL) {
		return normalizedURL, nil
	}

	urls, err := ScrapeForFeedURLs(normalizedURL)

	// TODO: Handle multiple valid feeds from ScrapeForFeedURLs
	if err != nil && len(urls) > 0 && IsValidFeed(urls[0]) {
		return urls[0], nil
	}

	return "", ErrInvalidFeedURL
}
