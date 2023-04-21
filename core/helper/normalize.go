package helper

import (
	"errors"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/purell"
)

var ErrUnparseableUrl = errors.New("URL is not parseable")
var ErrUnnormalizeableUrl = errors.New("URL cannot be normalized")

func NormalizeFeedUrl(input string) (string, error) {
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "https://" + input
	}

	parsed, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	normalized, err := purell.NormalizeURLString(parsed.String(), purell.FlagsSafe)
	if err != nil {
		return "", err
	}

	return normalized, nil
}
