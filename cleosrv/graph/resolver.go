package graph

import (
	"github.com/99designs/gqlgen/graphql"

	"github.com/cleodora-forecasting/cleodora/cleosrv/ent"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/generated"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the resolver root.
type Resolver struct{ client *ent.Client }

// NewSchema creates a graphql executable schema.
func NewSchema(client *ent.Client) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &Resolver{client},
		},
	)
}
