package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

 type Middleware func(http.Handler) http.Handler

func ProxyHandler(proxy *httputil.ReverseProxy) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request ){
		log.Printf(
			"[proxy] METHOD %s | PATH %s | USER-AGENT %s",
			r.Method,
			r.URL.Path,
			r.Header.Get("User-Agent"),
		)
		proxy.ServeHTTP(w, r)
	}
}

func simpleLogger(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log.Printf("before the middleware")

		next.ServeHTTP(w,r)

		log.Printf("After the middleware")
	})
}


func main(){
	target, err := url.Parse("http://httpbin.org")
	if err != nil{
		log.Fatalf("failed to parse the target url %v ", err)
	}
	

	// create the forwarding

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 5 * time.Second,
	}

	handler := ProxyHandler(proxy)
	wrapped := simpleLogger(handler)
	//register the func

	http.Handle("/home", wrapped)
	err = http.ListenAndServe(":8080", nil)
	if err != nil{
		log.Fatalf("server did not start %v", err)
	}
} 