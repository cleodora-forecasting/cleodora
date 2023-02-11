// Code generated by ent, DO NOT EDIT.

package estimate

import (
	"time"
)

const (
	// Label holds the string label denoting the estimate type in the database.
	Label = "estimate"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldReason holds the string denoting the reason field in the database.
	FieldReason = "reason"
	// FieldCreated holds the string denoting the created field in the database.
	FieldCreated = "created"
	// EdgeForecast holds the string denoting the forecast edge name in mutations.
	EdgeForecast = "forecast"
	// EdgeProbabilities holds the string denoting the probabilities edge name in mutations.
	EdgeProbabilities = "probabilities"
	// Table holds the table name of the estimate in the database.
	Table = "estimates"
	// ForecastTable is the table that holds the forecast relation/edge.
	ForecastTable = "estimates"
	// ForecastInverseTable is the table name for the Forecast entity.
	// It exists in this package in order to avoid circular dependency with the "forecast" package.
	ForecastInverseTable = "forecasts"
	// ForecastColumn is the table column denoting the forecast relation/edge.
	ForecastColumn = "forecast_estimates"
	// ProbabilitiesTable is the table that holds the probabilities relation/edge.
	ProbabilitiesTable = "probabilities"
	// ProbabilitiesInverseTable is the table name for the Probability entity.
	// It exists in this package in order to avoid circular dependency with the "probability" package.
	ProbabilitiesInverseTable = "probabilities"
	// ProbabilitiesColumn is the table column denoting the probabilities relation/edge.
	ProbabilitiesColumn = "estimate_probabilities"
)

// Columns holds all SQL columns for estimate fields.
var Columns = []string{
	FieldID,
	FieldReason,
	FieldCreated,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "estimates"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"forecast_estimates",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultReason holds the default value on creation for the "reason" field.
	DefaultReason string
	// DefaultCreated holds the default value on creation for the "created" field.
	DefaultCreated func() time.Time
)
