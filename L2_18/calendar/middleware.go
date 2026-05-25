package main

import (
	"log"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Method: %s, URL: %s, Time: %s", r.Method, r.URL.String(), start.Format(time.RFC3339))

		next.ServeHTTP(w, r)

		log.Printf("Completed in %v", time.Since(start))
	})
}
