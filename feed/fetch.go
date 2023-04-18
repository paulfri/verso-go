package feed

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func Fetch(url string) []*gofeed.Item {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)

	if err != nil {
		fmt.Println(err)
	}

	return feed.Items
}
