package main

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cleodora-forecasting/cleodora/graph"
	"github.com/cleodora-forecasting/cleodora/graph/generated"
)

var VERSION = "dev"

func main() {
	fmt.Printf(
		"Starting Cleodora (version: %s) http://localhost:8080\n",
		VERSION,
	)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/playground/",
		playground.Handler("GraphQL playground", "/query"),
	)
	http.Handle("/query", srv)

	serveFrontend()

	http.ListenAndServe(":8080", nil)
}
