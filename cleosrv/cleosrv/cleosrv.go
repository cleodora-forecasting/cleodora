package cleosrv

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"

	"github.com/cleodora-forecasting/cleodora/cleosrv/graph"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/generated"
	"github.com/cleodora-forecasting/cleodora/cleoutils"
)

func Start() error {
	fmt.Printf(
		"Starting Cleodora (version: %s) http://localhost:8080\n",
		cleoutils.Version,
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

	return http.ListenAndServe(":8080", router)
}
