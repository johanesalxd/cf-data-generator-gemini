package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/vertexai/genai"
)

func main() {
	lightControlTool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{{
			Name:        "controlLight",
			Description: "Set the brightness and color temperature of a room light.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"brightness": {
						Type: genai.TypeString,
						Description: "Light level from 0 to 100. Zero is off and" +
							" 100 is full brightness.",
					},
					"colorTemperature": {
						Type: genai.TypeString,
						Description: "Color temperature of the light fixture which" +
							" can be `daylight`, `cool` or `warm`.",
					},
				},
				Required: []string{"brightness", "colorTemperature"},
			},
		}},
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, os.Getenv("PROJECT_ID"), os.Getenv("LOCATION"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Use a model that supports function calling, like a Gemini 1.5 model
	model := client.GenerativeModel("gemini-1.5-flash")

	// Specify the function declaration.
	model.Tools = []*genai.Tool{lightControlTool}

	// Start new chat session.
	session := model.StartChat()

	prompt := "Dim the lights so the room feels cozy and warm."

	// Send the message to the generative model.
	resp, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
	}

	// Check that you got the expected function call back.
	part := resp.Candidates[0].Content.Parts[0]
	funcall, ok := part.(genai.FunctionCall)
	if !ok {
		log.Fatalf("Expected type FunctionCall, got %T", part)
	}
	if g, e := funcall.Name, lightControlTool.FunctionDeclarations[0].Name; g != e {
		log.Fatalf("Expected FunctionCall.Name %q, got %q", e, g)
	}
	fmt.Printf("Received function call response:\n%q\n\n", part)

	fmt.Printf("%v\n", funcall.Args)

	apiResult := setLightValues(funcall.Args)

	// Send the hypothetical API result back to the generative model.
	fmt.Printf("Sending API result:\n%q\n\n", apiResult)
	resp, err = session.SendMessage(ctx, genai.FunctionResponse{
		Name:     lightControlTool.FunctionDeclarations[0].Name,
		Response: apiResult,
	})
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
	}

	// Show the model's response, which is expected to be text.
	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("%v\n", part)
	}
}

func setLightValues(input map[string]any) map[string]any {
	// This mock API returns the requested lighting values
	return map[string]any{
		"brightness":       input["brightness"],
		"colorTemperature": input["colorTemperature"]}
}
