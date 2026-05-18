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
	"log/slog"
	"os"
)

// Middleware defines a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// statusResponseWriter wraps the standard http.ResponsWriter to cappture the status code.
type statusResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// writeHeader intercepts the status code before  sending it to the cliet
func (sw *statusResponseWriter) WriteHeader(statusCode int){
	sw.StatusCode =  statusCode // record the status code for our logs
	sw.ResponseWriter.WriteHeader(statusCode) // forward it to the real client response
}


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
	// worth checking if we need to upgrade this log
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
// this function currently handles the core proxy logging using a starndard text formatting .
// i need to remove this text logger because my new middleware will hande this is JSON format
func ProxyHandler(proxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// forward the request to the target server
		proxy.ServeHTTP(w, r)
	}
}

// // simpleLogger logs before and after the request is processed
// func simpleLogger(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		id := GetRequestID(r.Context())
// 		// instead of printing separate "before" and "after "  text lines  this must be refactored into a single JSON logging powerhouse
// 		// track time, inject custom reposnse write wrapper into next.ServeHTTP, Extract the requestID , calculate duration, use slog to write single JSON log containing all five required fileds 
// 		log.Println("[MIDDLEWARE] before request")
// 		log.Printf("[MIDDLEWARE] request_id=%s before request", id)


// 		// continue the chain
// 		next.ServeHTTP(w, r)

// 		log.Println("[MIDDLEWARE] after request")
// 	})
// }

// simpleBlocker blocks access to specific routes
func simpleBlocker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// block a specific path
		// when a request is blocked currently it just print unstructured text log and immdiately returns 403 Forbidden Status. upgrade this to use a new JSON slog layout so that denials much my starndard formatting

		if r.URL.Path == "/blocked" {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			logger.Warn("Acess Denied", "path", r.URL.Path, "status_code", http.StatusForbidden, "request_id", GetRequestID(r.Context()))
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

func jsonLoggerMiddleware(next http.Handler) http.Handler{

	// initialize the slog JSON handler writing to starndard output
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// caprture the start time before the request execution starts
		startTime := time.Now()

		// wrap the original response writer in our status capture
		// default the status to 200 because if the upstream handler succeeds
		// it often omits calling WriteHeader() explicitly
		wrappedWriter := &statusResponseWriter{
			ResponseWriter: w,
			StatusCode: http.StatusOK,
		}

		// pass control down the chain  using the wrapped  writer 
		next.ServeHTTP(wrappedWriter, r)

		// calculate execution during after the request finishes 
		duration := time.Since(startTime)

		// fetch the request Id from the context 
		reqID := GetRequestID(r.Context())

		// output a single jSON log containing all 5 required fields
		logger.Info("HTTP Request Processed",
			"request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", wrappedWriter.StatusCode,
			"duration_ms", duration.Milliseconds(), // logged as an easily indexable number
	)

	})
}

func main() {

	//set up default logger 
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil))) 
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
	// SL := simpleLogger(Sb)
	JL := jsonLoggerMiddleware(Sb)
	// Request Id runs first setting up the context for the logger
	Rid := requestIDMiddleware(JL)

	
	// register route
	http.Handle("/", Rid)


	// upgrade this logger massively change this to slog.error or slog.log

	slog.Info("server started on :8080")

	// upgrade this logger massively(Js onified)
	// i can either instanciate it so i can use it all through the algorithm or initialize it as a global logger.
	// start server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("failed to start server", err)
		os.Exit(1)
	}
}