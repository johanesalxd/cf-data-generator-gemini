package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/vertexai/genai"
)

func main() {
	promptParser := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{{
			Name:        "promptParser",
			Description: "Prompt input parser to separate between the actual prompt and number of rows to generate",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"prompt": {
						Type:        genai.TypeString,
						Description: "Actual prompt to generate a single answer (and omitting the number of rows to generate)",
					},
					"rows": {
						Type:        genai.TypeInteger,
						Description: "Number of rows to generate",
					},
				},
				Required: []string{"prompt", "rows"},
			},
		}},
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, os.Getenv("PROJECT_ID"), os.Getenv("LOCATION"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash-002")
	model.Tools = []*genai.Tool{promptParser}
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(`You are a built-in tool to parse the user's prompt to generate a dummy data. So you need to separate between the actual prompt and number of rows to generate.

			For example:
			User: "Generate 100 rows of dummy data for a table that has id, name, and email columns."
			You: "{\"prompt\": \"Generate a table with id, name, and email columns.\", \"rows\": 100}"
			`),
		},
	}
	session := model.StartChat()

	prompt := "Generate 100 rows of dummy data for a table that has id, name, and email columns."
	resp, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
	}
	fmt.Printf("Received request:\n%q\n\n", prompt)

	part := resp.Candidates[0].Content.Parts[0]
	funcall, ok := part.(genai.FunctionCall)
	if !ok {
		log.Fatalf("Expected type FunctionCall, got %T", part)
	}
	if g, e := funcall.Name, promptParser.FunctionDeclarations[0].Name; g != e {
		log.Fatalf("Expected FunctionCall.Name %q, got %q", e, g)
	}
	fmt.Printf("Received function call response:\n%q\n\n", part)

	apiResult := setPrompt(funcall.Args)
	fmt.Printf("Executing function call response:\n%q\n\n", apiResult)
	resp, err = session.SendMessage(ctx, genai.FunctionResponse{
		Name:     promptParser.FunctionDeclarations[0].Name,
		Response: apiResult,
	})
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("Received response:\n%v\n\n", part)
	}
}

func setPrompt(input map[string]any) map[string]any {
	return map[string]any{
		"prompt": input["prompt"],
		"rows":   input["rows"],
	}
}
