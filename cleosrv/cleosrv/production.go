//go:build production

package cleosrv

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/adrg/xdg"
	"github.com/go-chi/chi/v5"
)

//go:embed frontend_build
var embeddedFiles embed.FS

func serveFrontend(router chi.Router, frontendFooterText string) {
	fmt.Println("Serving with frontend")
	router.HandleFunc("/config.json", getFrontendConfig(frontendFooterText))
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

type frontendConfig struct {
	FooterText string `json:"FOOTER_TEXT"`
}

// getFrontendConfig returns a function that generates the content of the
// config.json file used to configure the frontend.
func getFrontendConfig(frontendFooterText string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := frontendConfig{
			FooterText: frontendFooterText,
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			fmt.Println("ERROR while encoding frontendConfig", err) // TODO log
		}
	}
}

func DefaultDatabasePath() string {
	return filepath.Join(xdg.DataHome, "cleosrv", "cleosrv.db")
}
