package datageneratorgemini

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/vertexai/genai"
)

// DataGeneratorGemini handles HTTP requests to generate data using Gemini AI
func DataGeneratorGemini(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming request
	req := new(RequestModel)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		SendError(w, err, http.StatusBadRequest)

		return
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(contextTimeoutS)*time.Second)
	defer func() {
		cancel()
		log.Printf("Done, context closed due to: %v", ctx.Err())
	}()

	// Get a client from the pool
	client := clientPool.Get().(*genai.Client)
	defer func() {
		if client != nil {
			clientPool.Put(client)
			log.Print("Client returned to pool")
		}
	}()
	log.Print("Client retrieved from pool")

	// Process the request using textsToTexts function
	resp := textsToTexts(ctx, client, req)

	// Send the successful response
	SendSuccess(w, resp)
}
