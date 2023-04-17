package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"github.com/sashabaranov/go-openai"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("OPENAI_API_KEY")

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://www.sounderatheart.com/rss/current.xml")

	fmt.Println(feed.Title)

	items := lo.Map(feed.Items, func(item *gofeed.Item, index int) string {
		return item.Title + "\n"
	})

	fmt.Printf("%v", items)

	item := feed.Items[0]

	client := openai.NewClient(key)
	resp, err := client.CreateCompletion(
		context.Background(),
		openai.CompletionRequest{
			Model:            openai.GPT3TextDavinci003,
			Prompt:           item.Content + "\n\nTl;dr",
			MaxTokens:        100,
			Temperature:      0.7,
			TopP:             1.0,
			FrequencyPenalty: 0.0,
			PresencePenalty:  1,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Text)
}
