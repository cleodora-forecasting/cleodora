package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Probability holds the schema definition for the Probability entity.
type Probability struct {
	ent.Schema
}

// Fields of the Probability.
func (Probability) Fields() []ent.Field {
	return []ent.Field{
		field.Int("value").Range(0, 100),
	}
}

// Edges of the Probability.
func (Probability) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("estimate", Estimate.Type).Ref("probabilities").Unique(),
		edge.From("outcome", Outcome.Type).Ref("probabilities").Unique(),
	}
}

func (Probability) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
	}
}
