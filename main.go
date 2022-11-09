package main

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cleodora-forecasting/cleodora/graph"
	"github.com/cleodora-forecasting/cleodora/graph/generated"
	"github.com/go-chi/chi/v5"
)

var VERSION = "dev"

func main() {
	fmt.Printf(
		"Starting Cleodora (version: %s) http://localhost:8080\n",
		VERSION,
	)

	router := chi.NewRouter()

	resolver := graph.Resolver{}
	resolver.AddDummyData()

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}),
	)

	configureCORS(router, srv)

	router.Handle("/playground/",
		playground.Handler("GraphQL playground", "/query"),
	)
	router.Handle("/query", srv)

	serveFrontend(router)

	http.ListenAndServe(":8080", router)
}
