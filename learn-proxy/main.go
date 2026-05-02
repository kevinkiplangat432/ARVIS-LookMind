package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Middleware func(http.Handler) http.Handler

func ProxyHandlerfunc(proxy *httputil.ReverseProxy) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		log.Printf(
			"[PROXY] METHOD %s PATH %s USER-AGENT %s",
			r.Method,
			r.URL.Path,
			r.Header.Get("User-Agent"),
		)
		proxy.ServeHTTP(w, r)
	}
}

func simplelogger(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			log.Printf("[Middleware] before it runs")
			next.ServeHTTP(w, r)
			log.Print("[middleware] after the middleware runs")
		})
	}


func main(){

	// parse the target url string 
	target, err := url.Parse("http://httpbin.org")
	if err != nil{
		log.Fatalf("Failed to parse the target url %v", err)
	}
	// create a new proxy 
	proxy := httputil.NewSingleHostReverseProxy(target)
	// set a proxy timeout
	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 5* time.Second,
	}

	PHF:= ProxyHandlerfunc(proxy)

	wrapped :=simplelogger(PHF)
	
	http.Handle("/home", wrapped)
	log.Print("server started successfully == .... ====....===...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil{
		log.Fatalf("failed to start up the server %v", err)
	}

}