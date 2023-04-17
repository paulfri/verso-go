package ai

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func Summarize(content string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(key)

	resp, err := client.CreateCompletion(
		context.Background(),
		openai.CompletionRequest{
			// Model: openai.GPT3TextDavinci003,
			Model:            openai.GPT3TextBabbage001,
			Prompt:           content + "\n\nTl;dr",
			MaxTokens:        100,
			Temperature:      0.7,
			TopP:             1.0,
			FrequencyPenalty: 0.0,
			PresencePenalty:  1,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
	}

	return resp.Choices[0].Text
}
