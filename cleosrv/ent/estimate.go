// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/estimate"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/forecast"
)

// Estimate is the model entity for the Estimate schema.
type Estimate struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Reason holds the value of the "reason" field.
	Reason string `json:"reason,omitempty"`
	// Created holds the value of the "created" field.
	Created time.Time `json:"created,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EstimateQuery when eager-loading is set.
	Edges              EstimateEdges `json:"edges"`
	forecast_estimates *int
}

// EstimateEdges holds the relations/edges for other nodes in the graph.
type EstimateEdges struct {
	// Forecast holds the value of the forecast edge.
	Forecast *Forecast `json:"forecast,omitempty"`
	// Probabilities holds the value of the probabilities edge.
	Probabilities []*Probability `json:"probabilities,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
	// totalCount holds the count of the edges above.
	totalCount [2]map[string]int

	namedProbabilities map[string][]*Probability
}

// ForecastOrErr returns the Forecast value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EstimateEdges) ForecastOrErr() (*Forecast, error) {
	if e.loadedTypes[0] {
		if e.Forecast == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: forecast.Label}
		}
		return e.Forecast, nil
	}
	return nil, &NotLoadedError{edge: "forecast"}
}

// ProbabilitiesOrErr returns the Probabilities value or an error if the edge
// was not loaded in eager-loading.
func (e EstimateEdges) ProbabilitiesOrErr() ([]*Probability, error) {
	if e.loadedTypes[1] {
		return e.Probabilities, nil
	}
	return nil, &NotLoadedError{edge: "probabilities"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Estimate) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case estimate.FieldID:
			values[i] = new(sql.NullInt64)
		case estimate.FieldReason:
			values[i] = new(sql.NullString)
		case estimate.FieldCreated:
			values[i] = new(sql.NullTime)
		case estimate.ForeignKeys[0]: // forecast_estimates
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Estimate", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Estimate fields.
func (e *Estimate) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case estimate.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			e.ID = int(value.Int64)
		case estimate.FieldReason:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field reason", values[i])
			} else if value.Valid {
				e.Reason = value.String
			}
		case estimate.FieldCreated:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created", values[i])
			} else if value.Valid {
				e.Created = value.Time
			}
		case estimate.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field forecast_estimates", value)
			} else if value.Valid {
				e.forecast_estimates = new(int)
				*e.forecast_estimates = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryForecast queries the "forecast" edge of the Estimate entity.
func (e *Estimate) QueryForecast() *ForecastQuery {
	return NewEstimateClient(e.config).QueryForecast(e)
}

// QueryProbabilities queries the "probabilities" edge of the Estimate entity.
func (e *Estimate) QueryProbabilities() *ProbabilityQuery {
	return NewEstimateClient(e.config).QueryProbabilities(e)
}

// Update returns a builder for updating this Estimate.
// Note that you need to call Estimate.Unwrap() before calling this method if this Estimate
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Estimate) Update() *EstimateUpdateOne {
	return NewEstimateClient(e.config).UpdateOne(e)
}

// Unwrap unwraps the Estimate entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Estimate) Unwrap() *Estimate {
	_tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("ent: Estimate is not a transactional entity")
	}
	e.config.driver = _tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Estimate) String() string {
	var builder strings.Builder
	builder.WriteString("Estimate(")
	builder.WriteString(fmt.Sprintf("id=%v, ", e.ID))
	builder.WriteString("reason=")
	builder.WriteString(e.Reason)
	builder.WriteString(", ")
	builder.WriteString("created=")
	builder.WriteString(e.Created.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// NamedProbabilities returns the Probabilities named value or an error if the edge was not
// loaded in eager-loading with this name.
func (e *Estimate) NamedProbabilities(name string) ([]*Probability, error) {
	if e.Edges.namedProbabilities == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := e.Edges.namedProbabilities[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (e *Estimate) appendNamedProbabilities(name string, edges ...*Probability) {
	if e.Edges.namedProbabilities == nil {
		e.Edges.namedProbabilities = make(map[string][]*Probability)
	}
	if len(edges) == 0 {
		e.Edges.namedProbabilities[name] = []*Probability{}
	} else {
		e.Edges.namedProbabilities[name] = append(e.Edges.namedProbabilities[name], edges...)
	}
}

// Estimates is a parsable slice of Estimate.
type Estimates []*Estimate

func (e Estimates) config(cfg config) {
	for _i := range e {
		e[_i].config = cfg
	}
}