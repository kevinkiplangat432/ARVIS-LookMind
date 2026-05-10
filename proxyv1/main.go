package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)
func ProxyHandler(proxy *httputil.ReverseProxy) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		log.Printf(
			"[PROXY] METHOD %s | PATH %s | USER-AGENT %s",
			r.Method,
			r.URL.Path,
			r.Header.Get("User-Agent"),
		)
		proxy.ServeHTTP(w,r)
	}

}

func main(){
	// parse the url
	target, err := url.Parse("http://httpbin.org")
	if err != nil{
		log.Fatal("Failed to parse the target url")
	}
	// create tge proxy 
	proxy := httputil.NewSingleHostReverseProxy(target)

	// add a timeout 
	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 5 * time.Second,
	}

	ph := ProxyHandler(proxy)

	http.HandleFunc("/", ph )
	http.ListenAndServe(":8080", nil)

}