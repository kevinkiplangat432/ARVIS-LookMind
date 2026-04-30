package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)


func ProxyHandlerFunction(proxy *httputil.ReverseProxy) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		log.Printf(
			"[PROXY]  %s  %s | UA: %s",
			r.Method,
			r.URL.Path,
			r.Header.Get("User-Agent"),
		)
		r.Header.Set("x-Forwarding-By", "Go-Reverse-Proxy")
		
		proxy.ServeHTTP(w, r)
	}

}

func main(){
	// create a target url
	target, err := url.Parse("https://httpbin.org")
	if err != nil{
		log.Fatal("Failed to parse the target url", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 5* time.Second,
	}

	http.HandleFunc("/", ProxyHandlerFunction(proxy))

	log.Println("Server running at port 8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil{
		log.Fatalf("server failed: %v", err)
	}
}