package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Adjust to match your LLM server's expected payload (e.g., Ollama, OpenAI)
type LLMRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type LLMResponse struct {
	Response string `json:"response"` // Adjust key based on provider
}

func chatStreamHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Set CORS and SSE Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Handle preflight CORS request
	if r.Method == http.MethodOptions {
		return
	}

	prompt := r.URL.Query().Get("prompt")
	if prompt == "" {
		prompt = "Hello"
	}

	// 2. Prepare payload for the LLM Server (e.g., local Ollama instance)
	llmReq := LLMRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: true,
	}
	jsonData, _ := json.Marshal(llmReq)

	// 3. Connect to LLM Server
	resp, err := http.Post("http://host.docker.internal:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to connect to LLM", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 4. Stream from LLM to Client
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var llmResp LLMResponse
		line := scanner.Bytes()

		if err := json.Unmarshal(line, &llmResp); err == nil {
			// SSE format requires starting with "data: " and ending with double newlines
			fmt.Fprintf(w, "data: %s\n\n", llmResp.Response)
			flusher.Flush()
		}
	}
}

func main() {
	http.HandleFunc("/api/chat", chatStreamHandler)
	log.Println("Server starting on :9080...")
	log.Fatal(http.ListenAndServe(":9080", nil))
}
