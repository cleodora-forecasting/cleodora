package main

import (
	"fmt"
	"net/http"
)

var VERSION = "dev"

func main() {
	fmt.Printf(
		"Starting Cleodora (version: %s) http://localhost:8080\n",
		VERSION,
	)
	http.HandleFunc("/api/", apiHandler)
	serveFrontend()
	http.ListenAndServe(":8080", nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Cleodora (version: %s)", VERSION)
}
