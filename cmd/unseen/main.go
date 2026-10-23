package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", healthLive)

	log.Print("starting server at http://127.0.0.1:8080")

	err := http.ListenAndServe("127.0.0.1:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func healthLive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload := map[string]string{"status": "ok"}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write liveness response: %v", err)
	}
}
