package main

import (
	"context"
	"crypto/rand"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	"fmt"
)

// Middleware defines a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

type contextKey string

const requestIDKey contextKey = "requestID"

func SetRequestID(ctx context.Context, ID string) context.Context {
	return context.WithValue(ctx, requestIDKey, ID )
}

func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey) .(string)
	return id
}

func generateId() string{
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "req-" + time.Now().Format("05.000")
	}
	return fmt.Sprintf("%x", b)
}

func requestIDMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		//generate id
		id := generateId()
		//store in context
		newCtx := SetRequestID(r.Context(), id)
		// attach new context to request
		r = r.WithContext(newCtx)
		// set reponse header
		w.Header().Set("X-Request-ID", id)

		// call next
		next.ServeHTTP(w, r)
	})

}

// ProxyHandler handles incoming requests and forwards them using the reverse proxy
func ProxyHandler(proxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"[PROXY] METHOD %s | PATH %s | USER-AGENT %s",
			r.Method,
			r.URL.Path,
			r.Header.Get("User-Agent"),
		)

		// forward the request to the target server
		proxy.ServeHTTP(w, r)
	}
}

// simpleLogger logs before and after the request is processed
func simpleLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestID(r.Context())
		log.Println("[MIDDLEWARE] before request")
		log.Printf("[MIDDLEWARE] request_id=%s before request", id)


		// continue the chain
		next.ServeHTTP(w, r)

		log.Println("[MIDDLEWARE] after request")
	})
}

// simpleBlocker blocks access to specific routes
func simpleBlocker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// block a specific path
		if r.URL.Path == "/blocked" {
			log.Printf("[BLOCKER] denied access to %s", r.URL.Path)

			http.Error(w, "you are not allowed here", http.StatusForbidden)
			return // stop the chain
		}

		// continue if not blocked
		next.ServeHTTP(w, r)
	})
}

// headerModifier adds a custom header before forwarding the request
func headerModifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// add a custom header to the outgoing request
		r.Header.Set("X-Learning-Proxy", "Active")

		// continue the chain
		next.ServeHTTP(w, r)
	})
}

func main() {
	// parse the target URL
	target, err := url.Parse("http://httpbin.org")
	if err != nil {
		log.Fatalf("failed to parse the target url: %v", err)
	}

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// configure timeout for upstream response
	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 5 * time.Second,
	}



	// build the handler chain
	Ph := ProxyHandler(proxy)

	// middleware wrapping (last added runs first)
	Hm := headerModifier(Ph)
	Sb := simpleBlocker(Hm)
	SL := simpleLogger(Sb)
	Rid := requestIDMiddleware(SL)

	
	// register route
	http.Handle("/", Rid)

	log.Println("server started on :8080")

	// start server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}