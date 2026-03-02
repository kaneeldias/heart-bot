package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

// GenerateImage generates an image from the provided description using the genai client,
// writes it to outputPath and returns the path or an error.
func GenerateImage(description string, outputPath string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("creating genai client: %w", err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-image",
		genai.Text(description),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("generate content: %w", err)
	}

	if result == nil || len(result.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned from image generation")
	}

	cand := result.Candidates[0]
	if cand == nil || cand.Content == nil {
		return "", fmt.Errorf("candidate content is nil")
	}

	for _, part := range cand.Content.Parts {
		if part == nil {
			continue
		}
		// If the model returned inline binary data for the image, write it.
		if part.InlineData != nil && len(part.InlineData.Data) > 0 {
			imageBytes := part.InlineData.Data
			if outputPath == "" {
				outputPath = "generated_image.png"
			}
			if err := os.WriteFile(outputPath, imageBytes, 0644); err != nil {
				return "", fmt.Errorf("writing image file: %w", err)
			}
			return outputPath, nil
		}
		// If the model returned a text part describing something, print it as info.
		if part.Text != "" {
			// keep as informational output; don't treat as error
			fmt.Println(part.Text)
		}
	}

	return "", fmt.Errorf("no inline image data found in the model response")
}
