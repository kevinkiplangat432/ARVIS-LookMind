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


// Handler for the arvis Proxy server - forward data 
func arvisProxyHandler(w http.ResponseWriter, r *http.Request) {
	//assert that our proxy response writer can flush to the end client

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "proxy streaming not supported", http.StatusInternalServerError)
		return
	}


	// set downstream response headers
	w.Header().Set("Content-Type", "text/event-stream")

	// call the upstream streaming endpoint
	resp, err := http.Get("http://Localhost:8081/upstream-stream")
	if err != nil {
		http.Error(w, "Failed to connect to uostream", http.StatusBadGateway)
		return
	}

	
	defer resp.Body.Close()

	// stream the data chunk by chunk writer to writer
	//insted of loading the resp body using io.readall we use bufio.Reader
	reader := bufio.NewReader(resp.Body)

	for {
		// read up until each newline chunk sent by upstream 
		line, err := reader.ReadBytes('\n')

		if len(line) > 0 {
			// write the chunk to our client response writer
			w.Write(line)
			flusher.Flush() // push directly to the client rn
			fmt.Printf("[Proxy log] Forwarding chunk: %s", string(line))

		}

		if err !=nil {
			if err == io.EOF {
				fmt.Println("[Proxy log] Upstream stream ended")
				break
			}
			fmt.Printf("[proxy Log] stream reading errors: %v\n", err)
			break
		}
	}
}
