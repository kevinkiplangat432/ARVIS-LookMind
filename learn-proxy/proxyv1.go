// problem that we aim to solve is that i eant to listen for request, look at each request, decide what to do with the request weather i should forward itmodify or block it., get the response from the target server, look at the response, decide what to do with the response weather i should forward it modify or block it. and then send the response back to the client.

package main

import (
	"log"
	"net/http"
)

func main(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("hello"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}


// an http.handlerFunc is used to register a specific function to handle incomming HTTP request for a given URL pattern 
// an http.resposewriter in go is an interface used by HTTP  handlers to construct and send an HTTP  response back to the client