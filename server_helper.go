package datageneratorgemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"cloud.google.com/go/vertexai/genai"
)

var (
	clientPool      *sync.Pool
	initOnce        sync.Once
	contextTimeoutS int
)

// SendError sends an error response with the given error message and HTTP status code
func SendError(w http.ResponseWriter, err error, code int) {
	resp := new(ResponseModel)
	resp.ErrorMessage = fmt.Sprintf("Got error with details: %v", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

// SendSuccess sends a successful response with the given ResponseModel
func SendSuccess(w http.ResponseWriter, resp *ResponseModel) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func initAll() {
	var err error

	clientPool = &sync.Pool{
		New: func() interface{} {
			client, err := genai.NewClient(context.Background(), os.Getenv("PROJECT_ID"), os.Getenv("LOCATION"))
			if err != nil {
				log.Fatalf("Failed to create client: %v", err)
			}
			log.Print("Client created")

			return client
		},
	}

	// Parse the context timeout from the environment variable
	contextTimeoutS, err = strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT_S"))
	if err != nil {
		log.Printf("Failed to parse CONTEXT_TIMEOUT_S, using default value of 30 seconds: %v", err)
		contextTimeoutS = 30
	}
}
