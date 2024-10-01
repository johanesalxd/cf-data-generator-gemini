package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	datageneratorgemini "github.com/johanesalxd/cf-data-generator-gemini"
)

// main starts the function framework server on the specified port
func main() {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	funcframework.RegisterHTTPFunction("/", datageneratorgemini.DataGeneratorGemini)
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
