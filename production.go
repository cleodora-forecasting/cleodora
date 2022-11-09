//go:build production

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed frontend/build
var embeddedFiles embed.FS

func serveFrontend(router chi.Router) {
	fmt.Println("Serving with frontend")
	router.Handle("/*", http.FileServer(getFileSystem()))
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
