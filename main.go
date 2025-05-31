/*
(DONE) TODO: procedure to https://ai.google.dev/gemini-api/docs/text-generation#go using genai (including stream response)
The response is streamed but if printed to screen as it comes, the markdown renderer adds an unwanted jumpline, which is tricky to debug since there are actual necessary jumplines in the response that should be shown

TODO: branch to a crypto option so don't have api key exposed in .env (even though is local)
TODO: default to flash 2.0 if none model is specified (flag for model selection, maybe use cobra)
TODO: add release so can be summoned by just typing "gemyni" after go installing it as a package
? Potentially integrating image generation with flash 2.0
? Add tests
? Add more LLMs options for api calls (maybe would break the api response and error so integrate with standarized format)
*/

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	glam "github.com/charmbracelet/glamour"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

var spinnerSet = []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●"}

func main() {
	// ANSI sequence to hide the cursor while rendering spinner (added \n for padding)
	fmt.Print("\033[?25l\n")
	s := spinner.New(
		spinnerSet, 
		200*time.Millisecond,
	)
	// Some padding to the left
	s.Prefix = "  "
	s.Start()

	// Get the query
	query := os.Args[1]
	err := godotenv.Load(".env")
	if err != nil {
		s.Stop()
		fmt.Println("Failed loading .env file")
	}
	
	ctx := context.Background()
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("LLM_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	stream := client.Models.GenerateContentStream(
		ctx,
		"gemini-2.0-flash",
		genai.Text(query),
		nil,
	)
	
	strB := &strings.Builder{}
	// Arbitrary amount of characters
	strB.Grow(500) 

	r, _ := glam.NewTermRenderer(
		glam.WithAutoStyle(),
		glam.WithTableWrap(false),
		glam.WithWordWrap(150),
	)

	for chunk := range stream {
		// Sometimes chunk is blank so this throw a nil pointer derefence 
		if chunk != nil {
			part := chunk.Candidates[0].Content.Parts[0]
			strB.WriteString(part.Text)
		}
	}

	out, err := r.Render(strB.String())
	if err != nil {
		fmt.Println("Failed rendering glamour:", err.Error())
	} else {
		s.Stop()
		fmt.Print(out)
		// Show cursor again
		fmt.Print("\033[?25h")
	}
}
