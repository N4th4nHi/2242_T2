package main

import (
	"log"
	"net/http"
	"os"

	"github.com/justinas/alice"
)

// loggerMiddleware is a middleware function that logs incoming HTTP requests to a file.
func loggerMiddleware(next http.Handler) http.Handler {
	// Open a log file for writing.
	logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// Create a new logger instance that writes to the log file.
	logger := log.New(logFile, "", log.LstdFlags)

	// Return a new http.Handler that logs the request and calls the next middleware in the chain.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// thirdPartyMiddleware is a third-party middleware function that does some processing and passes the request to the next middleware in the chain.
func thirdPartyMiddleware(next http.Handler) http.Handler {
	// Return a new http.Handler that does some processing and calls the next middleware in the chain.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do some processing...
		next.ServeHTTP(w, r)
	})
}

// handleRequest is the main request handler function that handles the HTTP request and sends a response back to the client.
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Handle the request...
	w.Write([]byte("Everything is functional"))
}

func main() {
	// Create a middleware chain using alice that includes our loggerMiddleware and thirdPartyMiddleware.
	chain := alice.New(
		loggerMiddleware,
		thirdPartyMiddleware,
	)

	// Create a new HTTP server with our middleware chain and handleRequest function.
	server := http.Server{
		Addr:    ":8080",
		Handler: chain.ThenFunc(handleRequest),
	}

	// Start the server and log any errors that occur.
	log.Fatal(server.ListenAndServe())
}
