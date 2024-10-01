package datageneratorgemini

import "encoding/json"

type RequestModel struct {
	RequestID   string          `json:"requestId,omitempty"`
	PromptInput string          `json:"promptInput"`
	Model       string          `json:"model"`
	ModelConfig json.RawMessage `json:"modelConfig"`
}

type ResponseModel struct {
	Data         json.RawMessage `json:"data"`
	ErrorMessage string          `json:"errorMessage"`
}
