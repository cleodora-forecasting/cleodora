package main

import (
	"fmt"
	"net/http"
)

func main() {
    fmt.Println("Starting Server http://localhost:8080")
	http.HandleFunc("/api/", apiHandler)
    serveFrontend()
	http.ListenAndServe(":8080", nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
