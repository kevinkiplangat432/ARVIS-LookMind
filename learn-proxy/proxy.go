// main.go - HTTP Reverse Proxy with inspection capability
package main

import (
    "bytes" // manipulates bytes slices(raw data)
    "encoding/json" // converts JSON <> Go structs/mpas
    "io" // input output
    "log" // print with timestamps and proper logging 
    "net/http" // HTTP client and server very needed
    "net/http/httputil" // HTTP utilities including reverse proxy
    "net/url" //  Parse and build URLs
    "strings"
    "time"
)
// vreate a new custom type 
// 
type SafetyProxy struct { // struct to hold my proxy and target info
    targetURL  *url.URL  // the pointer to a url.URL object 
    proxy      *httputil.ReverseProxy //pointer to Go built in reverse proxy 
}
//
func NewSafetyProxy(target string) *SafetyProxy {
    targetURL, _ := url.Parse(target) // converts string into a url obkect and returns two things a parsed url and an error but we ignore the error for now with _
    
    sp := &SafetyProxy{
        targetURL: targetURL,
    } // sp is a variable that holds a pointer to a new SafetyProxy struct with the targetURL field set to the parsed URL. The proxy field will be initialized later.
    
    // Create reverse proxy
    sp.proxy = httputil.NewSingleHostReverseProxy(targetURL)
    
    // Customize request handling
	// director is a function that modifies the incoming request before it is sent to the target server. We save the original director function provided by the reverse proxy and then wrap it with our own function that calls the original director and then adds our custom logic to modify the request.
    originalDirector := sp.proxy.Director
    sp.proxy.Director = func(req *http.Request) {
        originalDirector(req)
        sp.modifyRequest(req)
    }
    
    // Customize response handling
    sp.proxy.ModifyResponse = sp.modifyResponse
    
    // Error handler
    sp.proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
        log.Printf("Proxy error: %v", err)
        http.Error(rw, "Proxy error", http.StatusBadGateway)
    }
    
    return sp
}

func (sp *SafetyProxy) modifyRequest(req *http.Request) {
    body, err := io.ReadAll(req.Body)
    if err != nil {
        log.Printf("Error reading body: %v", err)
        return
    }
    req.Body = io.NopCloser(bytes.NewReader(body))
    log.Printf("→ [%s] %s", req.Method, req.URL.Path)
    // Inspect for AI-specific patterns
    if strings.Contains(req.Header.Get("Content-Type"), "application/json") {
        var data map[string]interface{}
        if json.Unmarshal(body, &data) == nil {
            // Check for prompt injection patterns
            if prompt, ok := data["prompt"].(string); ok {
                if strings.Contains(strings.ToLower(prompt), "ignore previous instructions") {
                    log.Printf("Potential prompt injection detected!")
                }
            }
            // Log token count or other metrics
            log.Printf("Request body size: %d bytes", len(body))
        }
    }
    
    // Add safety headers
    req.Header.Set("X-Proxy-Version", "1.0")
    req.Header.Set("X-Safety-Check", "enabled")
}

func (sp *SafetyProxy) modifyResponse(resp *http.Response) error {
    log.Printf("← [%d] %s", resp.StatusCode, resp.Request.URL.Path)
    
    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    resp.Body = io.NopCloser(bytes.NewReader(body))
    
    // Inspect AI response for content safety
    if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
        var data map[string]interface{}
        if json.Unmarshal(body, &data) == nil {
            // Check for unsafe content (simplified)
            if text, ok := data["choices"].([]interface{}); ok {
                log.Printf("📊 Response contains %d choices", len(text))
            }
        }
    }
    
    // Add response headers
    resp.Header.Set("X-Proxy-Processed", time.Now().Format(time.RFC3339))
    
    return nil
}

// func main() {
//     // Target AI provider (change to your endpoint)
//     target := "https://api.openai.com"
    
//     proxy := NewSafetyProxy(target)
    
//     // Add rate limiting, auth, etc. middleware
//     handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         // Add your safety rules here
//         proxy.proxy.ServeHTTP(w, r)
//     })
    
//     port := ":8080"
//     log.Printf("🚀 AI Safety Proxy starting on %s", port)
//     log.Printf("🎯 Forwarding to %s", target)
    
//     if err := http.ListenAndServe(port, handler); err != nil {
//         log.Fatal(err)
//     }
// }

// go run main.go

// # Test it:
// curl -X POST http://localhost:8080/v1/completions \
//   -H "Content-Type: application/json" \
//   -d '{"prompt": "Hello AI", "max_tokens": 10}'




//ask about the constructor function NewSafetyProxy