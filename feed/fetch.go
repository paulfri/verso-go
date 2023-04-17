package feed

import (
	"github.com/mmcdole/gofeed"
)

func Fetch(url string) []*gofeed.Item {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)

	return feed.Items
}
