package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Estimate holds the schema definition for the Estimate entity.
type Estimate struct {
	ent.Schema
}

// Fields of the Estimate.
func (Estimate) Fields() []ent.Field {
	return []ent.Field{
		field.String("reason").Default(""),
		field.Time("created").Default(time.Now),
	}
}

// Edges of the Estimate.
func (Estimate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("forecast", Forecast.Type).Ref("estimates").Unique(),
		edge.To("probabilities", Probability.Type),
	}
}

func (Estimate) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
	}
}
