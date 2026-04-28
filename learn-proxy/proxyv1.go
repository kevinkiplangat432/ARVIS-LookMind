// problem that we aim to solve is that i eant to listen for request, look at each request, decide what to do with the request weather i should forward itmodify or block it., get the response from the target server, look at the response, decide what to do with the response weather i should forward it modify or block it. and then send the response back to the client.

package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func main(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		// read the body 
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil{
			http.Error(w, "can't read Body", http.StatusInternalServerError)
			return
		}
		// close the body when done 
		defer r.Body.Close()
		// use the data from the body
		log.Println(r.Method, r.URL.Path, "Body:", string(bodyBytes))

		switch r.Method{
		case http.MethodGet:
			w.Write([]byte("Received!"))
		}
		// restore the Body  put a fresh one in there 
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}


// an http.handlerFunc is used to register a specific function to handle incomming HTTP request for a given URL pattern 
// an http.resposewriter in go is an interface used by HTTP  handlers to construct and send an HTTP  response back to the client