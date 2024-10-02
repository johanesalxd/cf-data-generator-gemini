package main

import (
	"context"
	"fmt"
	"log"
	"math"
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
					"counter": {
						Type:        genai.TypeInteger,
						Description: "How many times the prompt should be executed, at minimum, to process all the rows (+1 if needed to cover the remainder)",
					},
				},
				Required: []string{"prompt", "counter"},
			},
		}},
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, os.Getenv("PROJECT_ID"), os.Getenv("LOCATION"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	properties := map[string]int{
		"totalRows":   101,
		"rowsPerCall": 20,
		"counter":     0,
	}
	properties["counter"] = int(math.Ceil(float64(properties["totalRows"]) / float64(properties["rowsPerCall"])))
	systemInstruction := fmt.Sprintf(`You are a built-in tool to parse the user's prompt to generate a dummy data. Your task is to create an array of prompts, each of which will be used to generate a chunk of dummy data (maximum of %d rows per call) to avoid running into the token limit for the model.
	
	For example:
	User: "Generate %d rows of dummy data for a table that has id, name, and email columns using this JSON schema: Data = {"id": "integer", "name": "string", "email": "string"} Return: Array<Data>"
	
	You: "{\"prompt\": \"Generate %d rows of dummy data for a table that has id, name, and email columns using this JSON schema: Data = {\"id\": \"integer\", \"name\": \"string\", \"email\": \"string\"} Return: Array<Data>\", \"counter\": %d}
	`, properties["rowsPerCall"], properties["totalRows"], properties["totalRows"], properties["counter"])
	model := client.GenerativeModel("gemini-1.5-flash-002")
	model.Tools = []*genai.Tool{promptParser}
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(systemInstruction),
		},
	}
	session := model.StartChat()

	prompt := `Create 149 dummy cookie recipes using this JSON schema: Recipe = {'recipeName': string} Return: Array<Recipe>`
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

	apiResult, err := setPrompt(ctx, client, funcall.Args)
	if err != nil {
		log.Fatalf("Error setting prompt: %v", err)
	}
	fmt.Printf("Calling actual function for one time only:\n")
	printResponse(apiResult)
}

func setPrompt(ctx context.Context, client *genai.Client, input map[string]any) (*genai.GenerateContentResponse, error) {
	model := client.GenerativeModel("gemini-1.5-flash-002")
	model.ResponseMIMEType = "application/json"

	// Generate content using the configured model
	prompt := fmt.Sprintf("%v", input["prompt"])
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("error generating text: %w", err)
	}

	return resp, nil
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}
