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
						Description: "Rewritten prompt with modified expected number of rows to be generated and follow the system instructions",
					},
					"expectedTotalRows": {
						Type:        genai.TypeInteger,
						Description: "The original number of rows requested in the prompts",
					},
					"newTotalRows": {
						Type:        genai.TypeInteger,
						Description: "The new number of rows given the limit that set in the system instructions",
					},
				},
				Required: []string{"prompt", "expectedTotalRows", "newTotalRows"},
			},
		}},
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, os.Getenv("PROJECT_ID"), os.Getenv("LOCATION"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	properties := map[string]int{
		"expectedTotalRows": 101,
		"rowsPerCall":       20,
	}
	systemInstruction := fmt.Sprintf(`You are a built-in tool to parse the user's prompt to generate a dummy data with given JSON schema. Your task is to create an array of prompts, each of which will be used to generate a chunk of dummy data (maximum of %d rows per call) to avoid running into the token limit for the model. Return an error if the prompt is not related to data generation request or dummy data generation request doesn't include any JSON schema.

	For example:
	User: "Generate %d rows of dummy data for a table that has id, name, and email columns using this JSON schema: Data = {"id": "integer", "name": "string", "email": "string"} Return: Array<Data>"
	You: "{"prompt": "Generate %d rows of dummy data for a table that has id, name, and email columns using this JSON schema: Data = {"id":"integer","name":"string","email":"string"} Return = Array<Data>","expectedTotalRows":%d,"newTotalRows":%d}
	`, properties["rowsPerCall"], properties["expectedTotalRows"], properties["rowsPerCall"], properties["expectedTotalRows"], properties["rowsPerCall"])
	fmt.Println(systemInstruction)

	model := client.GenerativeModel("gemini-1.5-flash-002")
	model.Tools = []*genai.Tool{promptParser}
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(systemInstruction),
		},
	}
	session := model.StartChat()

	prompt := `Create 50 dummy game statistics data with this JSON schema: GameDetails = {'player_name': string, 'accuracy_percentage': integer, 'device': string} Return: Array<GameDetails>`
	// prompt := `list 5 cloud run region`
	resp, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
	}
	fmt.Printf("Received request:\n%q\n\n", prompt)

	part := resp.Candidates[0].Content.Parts[0]
	funcall, ok := part.(genai.FunctionCall)
	if !ok {
		//FIXME: somehow it return this error: Expected type FunctionCall, got genai.Text
		log.Fatalf("Expected type FunctionCall, got %T", part)
	}
	if g, e := funcall.Name, promptParser.FunctionDeclarations[0].Name; g != e {
		log.Fatalf("Expected FunctionCall.Name %q, got %q", e, g)
	}
	fmt.Printf("Received function call response:\n%q\n\n", part)

	fmt.Printf("Calling actual function:\n")
	apiResult := processPrompt(ctx, client, funcall.Args)
	if len(apiResult["errors"].([]error)) > 0 {
		log.Fatalf("Error setting prompt: %v", apiResult["errors"])
	}
	fmt.Printf("Length of result: %d\n", len(apiResult["results"].([]*genai.GenerateContentResponse)))

	fmt.Print("Sample result: \n")
	results := apiResult["results"].([]*genai.GenerateContentResponse)[0]
	printResponse(results)
}

func processPrompt(ctx context.Context, client *genai.Client, input map[string]any) map[string]any {
	iter := int(math.Ceil(input["expectedTotalRows"].(float64) / input["newTotalRows"].(float64)))
	resps := make([]*genai.GenerateContentResponse, iter)
	errs := make([]error, 0)

	model := client.GenerativeModel("gemini-1.5-flash-002")
	model.ResponseMIMEType = "application/json"

	for i := 0; i < int(iter); i++ {
		fmt.Println("iter", i)
		// Generate content using the configured model
		prompt := fmt.Sprintf("%v", input["prompt"])
		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			errs = append(errs, err)
		}

		resps[i] = resp
	}

	return map[string]any{
		"results": resps,
		"errors":  errs,
		"status":  "success",
	}
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
