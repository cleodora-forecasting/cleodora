// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/forecast"
)

// Forecast is the model entity for the Forecast schema.
type Forecast struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Created holds the value of the "created" field.
	Created time.Time `json:"created,omitempty"`
	// Resolves holds the value of the "resolves" field.
	Resolves time.Time `json:"resolves,omitempty"`
	// Closes holds the value of the "closes" field.
	Closes *time.Time `json:"closes,omitempty"`
	// Resolution holds the value of the "resolution" field.
	Resolution forecast.Resolution `json:"resolution,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ForecastQuery when eager-loading is set.
	Edges ForecastEdges `json:"edges"`
}

// ForecastEdges holds the relations/edges for other nodes in the graph.
type ForecastEdges struct {
	// Estimates holds the value of the estimates edge.
	Estimates []*Estimate `json:"estimates,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int

	namedEstimates map[string][]*Estimate
}

// EstimatesOrErr returns the Estimates value or an error if the edge
// was not loaded in eager-loading.
func (e ForecastEdges) EstimatesOrErr() ([]*Estimate, error) {
	if e.loadedTypes[0] {
		return e.Estimates, nil
	}
	return nil, &NotLoadedError{edge: "estimates"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Forecast) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case forecast.FieldID:
			values[i] = new(sql.NullInt64)
		case forecast.FieldTitle, forecast.FieldDescription, forecast.FieldResolution:
			values[i] = new(sql.NullString)
		case forecast.FieldCreated, forecast.FieldResolves, forecast.FieldCloses:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Forecast", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Forecast fields.
func (f *Forecast) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case forecast.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			f.ID = int(value.Int64)
		case forecast.FieldTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[i])
			} else if value.Valid {
				f.Title = value.String
			}
		case forecast.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				f.Description = value.String
			}
		case forecast.FieldCreated:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created", values[i])
			} else if value.Valid {
				f.Created = value.Time
			}
		case forecast.FieldResolves:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field resolves", values[i])
			} else if value.Valid {
				f.Resolves = value.Time
			}
		case forecast.FieldCloses:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field closes", values[i])
			} else if value.Valid {
				f.Closes = new(time.Time)
				*f.Closes = value.Time
			}
		case forecast.FieldResolution:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field resolution", values[i])
			} else if value.Valid {
				f.Resolution = forecast.Resolution(value.String)
			}
		}
	}
	return nil
}

// QueryEstimates queries the "estimates" edge of the Forecast entity.
func (f *Forecast) QueryEstimates() *EstimateQuery {
	return NewForecastClient(f.config).QueryEstimates(f)
}

// Update returns a builder for updating this Forecast.
// Note that you need to call Forecast.Unwrap() before calling this method if this Forecast
// was returned from a transaction, and the transaction was committed or rolled back.
func (f *Forecast) Update() *ForecastUpdateOne {
	return NewForecastClient(f.config).UpdateOne(f)
}

// Unwrap unwraps the Forecast entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (f *Forecast) Unwrap() *Forecast {
	_tx, ok := f.config.driver.(*txDriver)
	if !ok {
		panic("ent: Forecast is not a transactional entity")
	}
	f.config.driver = _tx.drv
	return f
}

// String implements the fmt.Stringer.
func (f *Forecast) String() string {
	var builder strings.Builder
	builder.WriteString("Forecast(")
	builder.WriteString(fmt.Sprintf("id=%v, ", f.ID))
	builder.WriteString("title=")
	builder.WriteString(f.Title)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(f.Description)
	builder.WriteString(", ")
	builder.WriteString("created=")
	builder.WriteString(f.Created.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("resolves=")
	builder.WriteString(f.Resolves.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := f.Closes; v != nil {
		builder.WriteString("closes=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("resolution=")
	builder.WriteString(fmt.Sprintf("%v", f.Resolution))
	builder.WriteByte(')')
	return builder.String()
}

// NamedEstimates returns the Estimates named value or an error if the edge was not
// loaded in eager-loading with this name.
func (f *Forecast) NamedEstimates(name string) ([]*Estimate, error) {
	if f.Edges.namedEstimates == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := f.Edges.namedEstimates[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (f *Forecast) appendNamedEstimates(name string, edges ...*Estimate) {
	if f.Edges.namedEstimates == nil {
		f.Edges.namedEstimates = make(map[string][]*Estimate)
	}
	if len(edges) == 0 {
		f.Edges.namedEstimates[name] = []*Estimate{}
	} else {
		f.Edges.namedEstimates[name] = append(f.Edges.namedEstimates[name], edges...)
	}
}

// Forecasts is a parsable slice of Forecast.
type Forecasts []*Forecast

func (f Forecasts) config(cfg config) {
	for _i := range f {
		f[_i].config = cfg
	}
}
