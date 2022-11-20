//go:build !production

package cleosrv

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

func serveFrontend(router chi.Router) {
	fmt.Println("Serving without frontend")
	router.HandleFunc("/", dummyFrontendHandler)
}

func dummyFrontendHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<h1>Dev Mode</h1>
<p>You are currently executing without frontend (e.g.  because you are
developing).</p>
<h1>GraphQL API</h1>
<p>Access the API under <a href="/query">/query</a></p>
<p>Access the playground under <a href="/playground/">/playground/</a></p>
<h1>Start Dev Frontend</h1>
To start the frontend during development as a separate process:
<pre>
cd frontend
npm start
</pre>
Then access it under <a href="http://localhost:3000">http://localhost:3000</a></p>
<h1>Embed Frontend</h1>
<p>To build and embed the frontend execute:
<pre>
cd frontend
npm run build
cd ..
</pre>
Then start the server with frontend:
<pre>
go run -tags production .
</pre>
Note that if you make any changes to the frontend code you will have to rebuild
to make those changes visible.
</p>
`,
	)
}

// configureCORS allows everything during development/testing.
func configureCORS(router *chi.Mux, srv *handler.Server) {
	fmt.Println("Enabling very permissive CORS")
	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		// It's against the CORS spec to allow credentialed requests and
		// asterisk * for Access-Control-Allow-Origin,
		// Access-Control-Allow-Headers or Access-Control-Allow-Methods .
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
		AllowCredentials: false,
		Debug:            true,
	}).Handler)

	// The following code was copied and slightly modified from the following
	// link, but it's not clear if or why it's needed since we are not using
	// websockets right now. Just in case, I'll leave it.
	// https://gqlgen.com/recipes/cors/
	// See also: https://github.com/99designs/gqlgen/issues/1250
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})
}
