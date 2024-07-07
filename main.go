package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	Text string `json:"text"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	message := Message{Text: "Hello from the API"}
	jsonResponse, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/api", apiHandler)

	loggedMux := loggingMiddleware(mux)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
