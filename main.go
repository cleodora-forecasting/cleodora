package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed frontend/build
var embeddedFiles embed.FS

func main() {
	fmt.Println("Starting Server")
	http.HandleFunc("/api/", apiHandler)
	http.Handle("/", http.FileServer(getFileSystem()))
	http.ListenAndServe(":8080", nil)
}

func getFileSystem() http.FileSystem {

    // Get the build subdirectory as the
    // root directory so that it can be passed
    // to the http.FileServer
	fsys, err := fs.Sub(embeddedFiles, "frontend/build")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
