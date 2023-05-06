package ai

import (
	"context"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

func Summarize(content string) string {
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
		log.Printf("ChatCompletion error: %v\n", err)
	}

	return resp.Choices[0].Text
}
