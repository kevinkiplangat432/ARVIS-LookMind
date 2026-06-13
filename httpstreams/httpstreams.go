package main

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
		http.HandleFunc("/upstream-strem", mockUpstreamHandler)
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

