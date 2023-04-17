package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
)

func main() {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://www.sounderatheart.com/rss/current.xml")

	fmt.Println(feed.Title)

	items := lo.Map(feed.Items, func(item *gofeed.Item, index int) string {
		return item.Title + "\n"
	})

	fmt.Printf("%v", items)
}
