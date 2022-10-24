// +build !frontend

package main

import (
    "fmt"
    "net/http"
)

func serveFrontend() {
    fmt.Println("Serving without frontend")
	http.HandleFunc("/", dummyFrontendHandler)
}

func dummyFrontendHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `<h1>Dev Mode</h1>
<p>You are currently executing without frontend (e.g.  because you are
developing).</p>
<h1>API</h1>
<p>Access the API under <a href="/api/">/api/</a></p>
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
