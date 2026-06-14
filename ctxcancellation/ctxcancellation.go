package main
// This example demonstrates how to implement a simple HTTP proxy that forwards streaming data from an upstream server to a downstream client.
import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	// start the mock upstream server sim OpenAI or a slow db
	go func() {
		http.HandleFunc("/upstream-stream", mockUpstreamHandler)
		fmt.Println("[Upstream] Server running on :8081...")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			panic(err)
		}
	}()

	// give the upstream server a split second to spin up
	time.Sleep(100 * time.Millisecond)

	// start the arvis Proxy server 
	http.HandleFunc("/proxy-stream", arvisProxyHandler)
	fmt.Println("[ARVIS Proxy] Listening on :8080...")
	fmt.Println("[ARVIS Proxy] Test your stream using: curl -N http://localhost:8080/proxy-stream")


	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)

	}



}



// hander for the mock upstream server -- generate data slowly
func mockUpstreamHandler(w http.ResponseWriter, r *http.Request) {
	//set headers for streaming / chunked transfer Encoding
	w.Header().Set("Content-type", "text/event-stream")
	w.Header().Set("cache-control", "no-cache")
	w.Header().Set("connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	for i := 1; i <= 5; i++ {
		//write a single chunk
		fmt.Fprintf(w, "chunk #%d from upstream server\n", i)
		flusher.Flush() //force the bytes down the wire immediately

		time.Sleep(500* time.Millisecond) // simulate slow data generation
	}

}




// Upgraded Handler with Context Cancellation
func arvisProxyHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Proxy streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")

	// EXTRACT THE CLIENT'S CONTEXT
	// This context represents the client's connection lifespan.
	clientCtx := r.Context()

	// CREATE A NEW REQUEST LINKED TO THE CLIENT'S CONTEXT
	// By using NewRequestWithContext, if clientCtx cancels, Go instantly 
	// tears down this upstream HTTP request automatically!
	req, err := http.NewRequestWithContext(clientCtx, "GET", "http://localhost:8081/upstream-stream", nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Execute the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// If the client disconnected before we could even connect upstream, handle it gracefully
		fmt.Println("[Proxy Log] Upstream request failed or aborted:", err)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		// Read chunk-by-chunk
		line, err := reader.ReadBytes('\n')
		
		if len(line) > 0 {
			w.Write(line)
			flusher.Flush()
			fmt.Printf("[Proxy Log] Forwarded chunk: %s", string(line))
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("[Proxy Log] Upstream finished stream safely.")
				break
			}
			
			// CHECK IF THE ERROR WAS CAUSED BY CLIENT CANCELLATION
			select {
			case <-clientCtx.Done():
				fmt.Println("[Proxy Log] Client disconnected! Cleaning up resources and exiting loop.")
			default:
				fmt.Printf("[Proxy Log] Network error reading upstream: %v\n", err)
			}
			break
		}
	}
}
