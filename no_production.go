//go:build !production

package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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
go run -tags frontend .
</pre>
Note that if you make any changes to the frontend code you will have to rebuild
to make those changes visible.
</p>
`,
	)
}
