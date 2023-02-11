package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Forecast holds the schema definition for the Forecast entity.
type Forecast struct {
	ent.Schema
}

// Fields of the Forecast.
func (Forecast) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("description").Default(""),
		field.Time("created").Default(time.Now),
		field.Time("resolves"),
		field.Time("closes").Optional().Nillable(),
		field.Enum("resolution").
			Values("UNRESOLVED", "RESOLVED", "NOT_APPLICABLE").
			Default("UNRESOLVED"),
	}
}

// Edges of the Forecast.
func (Forecast) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("estimates", Estimate.Type),
	}
}

func (Forecast) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
	}
}
