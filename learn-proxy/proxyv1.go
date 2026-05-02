package main

import (
    "bytes"
	"encoding/json"
    "io"
    "log"
    "net/http"
	"strings"
)

// for old completions API (/v1/completions)
type opemAICompletionsRequest struct{
	Prompt string `json:"prompt"`
	MaxTokens int `json:"max_tokens,omitempty"`
}

// for chat completions API (/v1/chat/completions)
type OpenAIChatRequest struct{
	Model string `json:"model"`
	Messages []struct{
		Role string `json:"role"`
		Content string `json:"content"`
	}`json:"messages"`
}
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // 1. READ BODY
        bodyBytes, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "can't read Body", http.StatusInternalServerError)
            return
        }
        
        // 2. RESTORE BODY (good habit)
        r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
       
        // 3. LOG REQUEST
        log.Println(r.Method, r.URL.Path, "Body:", string(bodyBytes))
        
		// 4. safety check block request conatainig hack
		blocked := false
		blockReason :=""
		// check which api endpoint 
		if strings.Contains(r.URL.Path, "/v1/completions"){
			var reqData opemAICompletionsRequest
			if err := json.Unmarshal(bodyBytes, &reqData); err == nil{
				if strings.Contains(strings.ToLower(reqData.Prompt), "hack") {
                    blocked = true
                    blockReason = "prompt contains blocked word"
				}
			}
		} else if strings.Contains(r.URL.Path, "/v1/chat/completions"){
			var reqData OpenAIChatRequest
			 if err := json.Unmarshal(bodyBytes, &reqData); err == nil {
                for _, msg := range reqData.Messages {
                    if strings.Contains(strings.ToLower(msg.Content), "hack") {
                        blocked = true
                        blockReason = "message contains blocked word"
                        break
                    }
                }
            }
		}else{
			// Generic check for other endpoints
            if strings.Contains(string(bodyBytes), "hack") {
                blocked = true
                blockReason = "contains blocked word"
            }
		}
		 // 5. BLOCK IF UNSAFE
        if blocked {
            log.Printf("BLOCKED: %s - %s", r.URL.Path, blockReason)
            http.Error(w, "Blocked: inappropriate content", http.StatusForbidden)
            return
        }
        
        // 6. CREATE FORWARD REQUEST
        proxyUrl := "https://api.openai.com" + r.URL.Path 
        newReq, err := http.NewRequest(r.Method, proxyUrl, bytes.NewReader(bodyBytes))
        if err != nil {
            http.Error(w, "failed to create request", http.StatusInternalServerError)
            return
        }
        
        // 5. COPY HEADERS 
        newReq.Header = r.Header.Clone()
        
        
        // 6. SEND REQUEST
        client := &http.Client{}
        resp, err := client.Do(newReq)
        if err != nil {
            http.Error(w, "Failed to reach destination", http.StatusBadGateway)
            return
        }
        defer resp.Body.Close()
        
        // 7. RETURN RESPONSE
        w.WriteHeader(resp.StatusCode)
        io.Copy(w, resp.Body)
    })
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}