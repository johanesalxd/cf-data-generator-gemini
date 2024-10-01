package datageneratorgemini

import (
	"encoding/json"
	"log"
)

type PromptRequest struct {
	PromptInput string      `json:"promptInput"`
	Model       string      `json:"model"`
	ModelConfig ModelConfig `json:"modelConfig"`
}

type ModelConfig struct {
	Temperature     float32 `json:"temperature"`
	MaxOutputTokens int32   `json:"maxOutputTokens"`
	TopP            float32 `json:"topP"`
	TopK            int32   `json:"topK"`
}

func newModelConfig() ModelConfig {
	return ModelConfig{
		Temperature:     1,
		MaxOutputTokens: 1000,
		TopP:            0.95,
		TopK:            1,
	}
}

// parseModelConfig takes a JSON raw message and returns a ModelConfig
// If the input can't be unmarshaled, it returns a default ModelConfig
func parseModelConfig(input json.RawMessage) ModelConfig {
	// Start with default configuration
	config := newModelConfig()

	// Attempt to unmarshal the input into the config
	if err := json.Unmarshal(input, &config); err != nil {
		// Log the error and return the default config if unmarshaling fails
		log.Printf("Default value used due to error unmarshaling model config: %v", err)
		return config
	}

	// Return the parsed (or default) configuration
	return config
}
