package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	Text string `json:"text"`
}
type Block struct {
	ID        int    `json:"id"`
	Data      string `json:"data"`
	Timestamp string `json:"timestamp`
}

var db *sql.DB

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
func initDB(filepath string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	createTableSQL := `CREATE TABLE IF NOT EXISTS blocks (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
        "data" TEXT,
        "timestamp" DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = database.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}
	return database, nil
}
func addBlock(data string) error {
	insertBlockSQL := `INSERT INTO blocks(data) VALUES (?)`
	statment, err := db.Prepare(insertBlockSQL)
	if err != nil {
		return err
	}
	_, err = statment.Exec(data)
	if err != nil {
		return err
	}
	return nil
}
func addBlockHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	if data == "" {
		http.Error(w, "Missing 'data' parameter", http.StatusBadRequest)
		return
	}
	err := addBlock(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Block added successfully")
}

func main() {

	var err error
	db, err = initDB("blockchain.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/api", apiHandler)
	mux.HandleFunc("/addblock", addBlockHandler)

	loggedMux := loggingMiddleware(mux)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
