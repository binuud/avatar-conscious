package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/sirupsen/logrus"
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

func openAiHandler(w http.ResponseWriter, r *http.Request) {

	logrus.Infof("openAiHandler %v", r)
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

	// Accept the Chat Completions request body from the caller.
	var params responses.ResponseNewParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Force streaming, regardless of what the caller sent.
	// optParams := &option.RequestOption{
	// 	openai.ChatCompletionStreamOptionsParam{
	// 		IncludeUsage: openai.Bool(true),
	// 		IncludeObfuscation: openai.Bool(false),
	// 	}
	// }

	client := openai.NewClient(
		option.WithBaseURL("http://host.docker.internal:11434/v1/"),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	logrus.Infof("Send data to llm server %s", params.User.Value)

	stream := client.Responses.NewStreaming(r.Context(), params)
	defer func() {
		_ = stream.Close()
	}()

	// Initialize the streaming session using the standard Responses endpoint
	fmt.Print("AI Response: ")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	logrus.Info("\n\n\n Starting streaming")
	// Iterate over incoming chunks as they arrive from the API

	for stream.Next() {
		data := stream.Current()
		println(data.Delta)
		if data.JSON.Text.Valid() {
			println()
			println("Finished Content")
			// completeText = data.Text
			break
		}
		flusher.Flush()
	}
	// for stream.Next() {
	// 	data := stream.Current()

	// 	// Print the token/text delta immediately to the terminal without a newline
	// 	fmt.Println(data)
	// 	if _, err := fmt.Fprint(w, data); err != nil {
	// 		return // client disconnected
	// 	}
	// 	flusher.Flush()
	// }

	// Always check for errors after the stream loop terminates
	if err := stream.Err(); err != nil {
		fmt.Printf("\nStream encountered an error: %v\n", err)
	}

	fmt.Println() // Print a final trailing newline when complete
}

func newHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	client := openai.NewClient(
		option.WithBaseURL("http://host.docker.internal:11434/v1/"),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

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

	logrus.Infof("Send data to llm server %s", prompt)
	// Initialize the streaming session using the standard Responses endpoint
	stream := client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(prompt),
		},
		Model: "gemma4:latest", // Or any preferred model variant
	})
	defer func() {
		_ = stream.Close()
	}()

	fmt.Print("AI Response: ")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Iterate over incoming chunks as they arrive from the API
	for stream.Next() {
		data := stream.Current()

		// Print the token/text delta immediately to the terminal without a newline
		fmt.Print(data.Delta)
		if _, err := fmt.Fprint(w, data); err != nil {
			return // client disconnected
		}
		flusher.Flush()
	}

	// Always check for errors after the stream loop terminates
	if err := stream.Err(); err != nil {
		fmt.Printf("\nStream encountered an error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println() // Print a final trailing newline when complete
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
		Model:  "gemma4:latest",
		Prompt: prompt,
		Stream: true,
	}
	jsonData, _ := json.Marshal(llmReq)

	// 3. Connect to LLM Server
	logrus.Infof("Send data to llm server %s", llmReq.Prompt)
	resp, err := http.Post("http://host.docker.internal:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to connect to LLM", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 2. FIX: Read JSON chunks sequentially using json.Decoder instead of bufio.Scanner
	decoder := json.NewDecoder(resp.Body)
	for decoder.More() {
		var llmResp string
		if err := decoder.Decode(&llmResp); err != nil {
			break
		}

		if llmResp != "" {
			// Format correctly as an SSE data payload without altering embedded newlines
			logrus.Infof("%s", llmResp)
			fmt.Fprintf(w, "%s", llmResp)
			flusher.Flush()
		}
	}
}

func main() {
	http.HandleFunc("/api/oldchat", chatStreamHandler)
	http.HandleFunc("/api/chat", newHandler)
	http.HandleFunc("/v1/chat/completions", openAiHandler)

	log.Println("Server starting on :9080...")
	log.Fatal(http.ListenAndServe(":9080", nil))
}
