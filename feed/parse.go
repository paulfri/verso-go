package feed

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func Parse(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return feed, nil
}
