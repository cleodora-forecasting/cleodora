//go:build production

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
)

//go:embed frontend_build
var embeddedFiles embed.FS

func serveFrontend(router chi.Router) {
	fmt.Println("Serving with frontend")
	router.Handle("/*", http.FileServer(getFileSystem()))
}

func getFileSystem() http.FileSystem {

	// Get the build subdirectory as the
	// root directory so that it can be passed
	// to the http.FileServer
	fsys, err := fs.Sub(embeddedFiles, "frontend_build")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

// configureCORS disables CORS for production
func configureCORS(router *chi.Mux, srv *handler.Server) {
	// We don't support CORS because the backend and fronted are bundled as a
	// single server and right now there is no plan or need to separate them.
	// Later more granular configuration can be introduced.
	// Probably it makes no sense without authentication for the backend.
	fmt.Println("Disabling CORS")
}
