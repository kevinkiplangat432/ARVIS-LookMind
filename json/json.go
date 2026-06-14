package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 1. Define the nested structure matching an OpenAI response chunk
type ChatCompletionChunk struct {
	ID      string `json:"id"`
	Choices []struct {
		Delta struct {
			Content string `json:"content,omitempty"` // The text token we want
		} `json:"delta"`
	} `json:"choices"`
	// Use RawMessage to defer parsing any other complex metadata fields we don't care about
	Usage json.RawMessage `json:"usage,omitempty"` 
}

func main() {
	// A mock array representing raw stream chunks coming straight from an AI engine over network bytes
	mockSSEStream := []string{
		`data: {"id": "chat-1", "choices": [{"delta": {"content": "Hello"}}], "usage": null}`,
		`data: {"id": "chat-1", "choices": [{"delta": {"content": " there"}}], "usage": null}`,
		`data: {"id": "chat-1", "choices": [{"delta": {"content": " world"}}], "usage": null}`,
		`data: {"id": "chat-1", "choices": [{"delta": {"content": "!"}}], "usage": {"prompt_tokens": 10, "completion_tokens": 4}}`,
	}

	fmt.Println("Processing AI Stream Chunks:")

	for _, rawLine := range mockSSEStream {
		// Clean the SSE protocol prefix
		if !strings.HasPrefix(rawLine, "data: ") {
			continue
		}
		jsonStr := strings.TrimPrefix(rawLine, "data: ")

		// 2. Decode the JSON string into our targeted struct
		var chunk ChatCompletionChunk
		err := json.Unmarshal([]byte(jsonStr), &chunk)
		if err != nil {
			fmt.Printf("Error decoding chunk: %v\n", err)
			continue
		}

		// 3. Extract and print the content field if it exists
		if len(chunk.Choices) > 0 {
			token := chunk.Choices[0].Delta.Content
			if token != "" {
				fmt.Print(token) // Print tokens side-by-side to simulate typing
			}
		}

		// 4. If we hit the final chunk containing the deferred raw usage data, print it out
		if chunk.Usage != nil && string(chunk.Usage) != "null" {
			fmt.Printf("\n\n[Metadata Raw Bytes]: %s\n", string(chunk.Usage))
		}
	}
	fmt.Println("Stream parsing complete!")
}
