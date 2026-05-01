package main

import (
	"log"
	"net/http"
	"time"
)

type timeHandler struct{
	format string
}

func (th timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	tm := time.Now().Format(th.format)
	w.Write([]byte("the time is:" + tm))
}

func main(){
	// Use the http.NewServeMux() function to create an empty servemux
	mux := http.NewServeMux()

	// use the http.RedirectHandler() function to create a handler which 307
	// redirects all the request it receives to http://example.org
	rh := http.RedirectHandler("http://example.org", 307)

	//next we use the mux.handl() functions to register this with our new
	//servemux, so it acts as the handler for all incomming request with the URL
	// path /foo
	mux.Handle("/foo", rh)

	log.Print("Listening...")

	// Then we create a new server and start listening for incoming requests
	// with the http.ListenAndServe() function, passing in our servemux for it to
	// match requests against as the second argument.
	http.ListenAndServe(":3000", mux)
}


// type Handler interface{
// 	ServeHTTP(w http.ResponseWriter, r *http.Request)
// }