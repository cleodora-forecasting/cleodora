package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Outcome holds the schema definition for the Outcome entity.
type Outcome struct {
	ent.Schema
}

// Fields of the Outcome.
func (Outcome) Fields() []ent.Field {
	return []ent.Field{
		field.String("text").NotEmpty(),
		field.Bool("correct").Default(false),
	}
}

// Edges of the Outcome.
func (Outcome) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("probabilities", Probability.Type),
	}
}

func (Outcome) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
	}
}
