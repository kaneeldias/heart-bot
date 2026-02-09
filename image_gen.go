package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

func GenerateImage(description string) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-image",
		genai.Text(description),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	if len(result.Candidates) == 0 {
		log.Printf("No candidates found")
		return
	}

	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Println(part.Text)
		} else if part.InlineData != nil {
			imageBytes := part.InlineData.Data
			outputFilename := "generated_image.png"
			_ = os.WriteFile(outputFilename, imageBytes, 0644)
		}
	}
}
