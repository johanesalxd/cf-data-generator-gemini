package datageneratorgemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

// textsToTexts processes a text-to-text request using the Vertex AI Generative AI model
func textsToTexts(ctx context.Context, client *genai.Client, req *RequestModel) *ResponseModel {
	// Validate input: ensure promptInput and model are not empty
	if req.PromptInput == "" || req.Model == "" {
		return &ResponseModel{
			ErrorMessage: "Invalid input: check your promptInput and model again",
		}
	}

	// Prepare the input for the generative model
	input := PromptRequest{
		PromptInput: req.PromptInput,
		Model:       req.Model,
		ModelConfig: parseModelConfig(req.ModelConfig),
	}

	// Call the textToText function to generate content
	texts, err := textToText(ctx, client, &input)
	if err != nil {
		return &ResponseModel{
			ErrorMessage: err.Error(),
		}
	}

	// Return the generated content as a successful response
	return &ResponseModel{
		Data: texts,
	}
}

// textToText generates content using the Vertex AI Generative AI model
func textToText(ctx context.Context, client *genai.Client, input *PromptRequest) (json.RawMessage, error) {
	// Create a new generative model instance
	mdl := client.GenerativeModel(input.Model)

	// Configure the model parameters based on the input
	mdl.SetMaxOutputTokens(input.ModelConfig.MaxOutputTokens)
	mdl.SetTemperature(input.ModelConfig.Temperature)
	mdl.SetTopP(input.ModelConfig.TopP)
	mdl.SetTopK(input.ModelConfig.TopK)

	// Set the response MIME type to JSON
	mdl.ResponseMIMEType = "application/json"

	// Generate content using the configured model
	resp, err := mdl.GenerateContent(ctx, genai.Text(input.PromptInput))
	if err != nil {
		return nil, fmt.Errorf("error generating text: %w", err)
	}

	// Process the response and return it as JSON
	return GenerateJSONResponse(resp)
}

// As per example here: https://github.com/google/generative-ai-go/blob/2ec7e23d0c921b95b2ef733030715a298972724d/genai/internal/samples/docs-snippets_test.go#L1752
func GenerateJSONResponse(resp *genai.GenerateContentResponse) (json.RawMessage, error) {
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("invalid response: no candidates found")
	}

	if resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("invalid response: no content found")
	}

	jsonBytes := []byte(fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]))

	if !json.Valid(jsonBytes) {
		return nil, fmt.Errorf("invalid JSON response")
	}

	trimmed := bytes.TrimSpace(jsonBytes)
	if len(trimmed) > 0 && trimmed[0] == '[' {
		return json.RawMessage(jsonBytes), nil
	}

	return nil, fmt.Errorf("invalid response: not an array")
}
