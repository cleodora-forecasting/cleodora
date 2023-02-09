// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/estimate"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/forecast"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/predicate"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/probability"
)

// EstimateUpdate is the builder for updating Estimate entities.
type EstimateUpdate struct {
	config
	hooks    []Hook
	mutation *EstimateMutation
}

// Where appends a list predicates to the EstimateUpdate builder.
func (eu *EstimateUpdate) Where(ps ...predicate.Estimate) *EstimateUpdate {
	eu.mutation.Where(ps...)
	return eu
}

// SetReason sets the "reason" field.
func (eu *EstimateUpdate) SetReason(s string) *EstimateUpdate {
	eu.mutation.SetReason(s)
	return eu
}

// SetNillableReason sets the "reason" field if the given value is not nil.
func (eu *EstimateUpdate) SetNillableReason(s *string) *EstimateUpdate {
	if s != nil {
		eu.SetReason(*s)
	}
	return eu
}

// SetCreated sets the "created" field.
func (eu *EstimateUpdate) SetCreated(t time.Time) *EstimateUpdate {
	eu.mutation.SetCreated(t)
	return eu
}

// SetNillableCreated sets the "created" field if the given value is not nil.
func (eu *EstimateUpdate) SetNillableCreated(t *time.Time) *EstimateUpdate {
	if t != nil {
		eu.SetCreated(*t)
	}
	return eu
}

// SetForecastID sets the "forecast" edge to the Forecast entity by ID.
func (eu *EstimateUpdate) SetForecastID(id int) *EstimateUpdate {
	eu.mutation.SetForecastID(id)
	return eu
}

// SetNillableForecastID sets the "forecast" edge to the Forecast entity by ID if the given value is not nil.
func (eu *EstimateUpdate) SetNillableForecastID(id *int) *EstimateUpdate {
	if id != nil {
		eu = eu.SetForecastID(*id)
	}
	return eu
}

// SetForecast sets the "forecast" edge to the Forecast entity.
func (eu *EstimateUpdate) SetForecast(f *Forecast) *EstimateUpdate {
	return eu.SetForecastID(f.ID)
}

// AddProbabilityIDs adds the "probabilities" edge to the Probability entity by IDs.
func (eu *EstimateUpdate) AddProbabilityIDs(ids ...int) *EstimateUpdate {
	eu.mutation.AddProbabilityIDs(ids...)
	return eu
}

// AddProbabilities adds the "probabilities" edges to the Probability entity.
func (eu *EstimateUpdate) AddProbabilities(p ...*Probability) *EstimateUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eu.AddProbabilityIDs(ids...)
}

// Mutation returns the EstimateMutation object of the builder.
func (eu *EstimateUpdate) Mutation() *EstimateMutation {
	return eu.mutation
}

// ClearForecast clears the "forecast" edge to the Forecast entity.
func (eu *EstimateUpdate) ClearForecast() *EstimateUpdate {
	eu.mutation.ClearForecast()
	return eu
}

// ClearProbabilities clears all "probabilities" edges to the Probability entity.
func (eu *EstimateUpdate) ClearProbabilities() *EstimateUpdate {
	eu.mutation.ClearProbabilities()
	return eu
}

// RemoveProbabilityIDs removes the "probabilities" edge to Probability entities by IDs.
func (eu *EstimateUpdate) RemoveProbabilityIDs(ids ...int) *EstimateUpdate {
	eu.mutation.RemoveProbabilityIDs(ids...)
	return eu
}

// RemoveProbabilities removes "probabilities" edges to Probability entities.
func (eu *EstimateUpdate) RemoveProbabilities(p ...*Probability) *EstimateUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eu.RemoveProbabilityIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (eu *EstimateUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, EstimateMutation](ctx, eu.sqlSave, eu.mutation, eu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (eu *EstimateUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *EstimateUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *EstimateUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eu *EstimateUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   estimate.Table,
			Columns: estimate.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: estimate.FieldID,
			},
		},
	}
	if ps := eu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eu.mutation.Reason(); ok {
		_spec.SetField(estimate.FieldReason, field.TypeString, value)
	}
	if value, ok := eu.mutation.Created(); ok {
		_spec.SetField(estimate.FieldCreated, field.TypeTime, value)
	}
	if eu.mutation.ForecastCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   estimate.ForecastTable,
			Columns: []string{estimate.ForecastColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: forecast.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.ForecastIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   estimate.ForecastTable,
			Columns: []string{estimate.ForecastColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: forecast.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.mutation.ProbabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   estimate.ProbabilitiesTable,
			Columns: []string{estimate.ProbabilitiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: probability.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.RemovedProbabilitiesIDs(); len(nodes) > 0 && !eu.mutation.ProbabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   estimate.ProbabilitiesTable,
			Columns: []string{estimate.ProbabilitiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: probability.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.ProbabilitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   estimate.ProbabilitiesTable,
			Columns: []string{estimate.ProbabilitiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: probability.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, eu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{estimate.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	eu.mutation.done = true
	return n, nil
}

// EstimateUpdateOne is the builder for updating a single Estimate entity.
type EstimateUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *EstimateMutation
}

// SetReason sets the "reason" field.
func (euo *EstimateUpdateOne) SetReason(s string) *EstimateUpdateOne {
	euo.mutation.SetReason(s)
	return euo
}

// SetNillableReason sets the "reason" field if the given value is not nil.
func (euo *EstimateUpdateOne) SetNillableReason(s *string) *EstimateUpdateOne {
	if s != nil {
		euo.SetReason(*s)
	}
	return euo
}

// SetCreated sets the "created" field.
func (euo *EstimateUpdateOne) SetCreated(t time.Time) *EstimateUpdateOne {
	euo.mutation.SetCreated(t)
	return euo
}

// SetNillableCreated sets the "created" field if the given value is not nil.
func (euo *EstimateUpdateOne) SetNillableCreated(t *time.Time) *EstimateUpdateOne {
	if t != nil {
		euo.SetCreated(*t)
	}
	return euo
}

// SetForecastID sets the "forecast" edge to the Forecast entity by ID.
func (euo *EstimateUpdateOne) SetForecastID(id int) *EstimateUpdateOne {
	euo.mutation.SetForecastID(id)
	return euo
}

// SetNillableForecastID sets the "forecast" edge to the Forecast entity by ID if the given value is not nil.
func (euo *EstimateUpdateOne) SetNillableForecastID(id *int) *EstimateUpdateOne {
	if id != nil {
		euo = euo.SetForecastID(*id)
	}
	return euo
}

// SetForecast sets the "forecast" edge to the Forecast entity.
func (euo *EstimateUpdateOne) SetForecast(f *Forecast) *EstimateUpdateOne {
	return euo.SetForecastID(f.ID)
}

// AddProbabilityIDs adds the "probabilities" edge to the Probability entity by IDs.
func (euo *EstimateUpdateOne) AddProbabilityIDs(ids ...int) *EstimateUpdateOne {
	euo.mutation.AddProbabilityIDs(ids...)
	return euo
}

// AddProbabilities adds the "probabilities" edges to the Probability entity.
func (euo *EstimateUpdateOne) AddProbabilities(p ...*Probability) *EstimateUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return euo.AddProbabilityIDs(ids...)
}

// Mutation returns the EstimateMutation object of the builder.
func (euo *EstimateUpdateOne) Mutation() *EstimateMutation {
	return euo.mutation
}

// ClearForecast clears the "forecast" edge to the Forecast entity.
func (euo *EstimateUpdateOne) ClearForecast() *EstimateUpdateOne {
	euo.mutation.ClearForecast()
	return euo
}

// ClearProbabilities clears all "probabilities" edges to the Probability entity.
func (euo *EstimateUpdateOne) ClearProbabilities() *EstimateUpdateOne {
	euo.mutation.ClearProbabilities()
	return euo
}

// RemoveProbabilityIDs removes the "probabilities" edge to Probability entities by IDs.
func (euo *EstimateUpdateOne) RemoveProbabilityIDs(ids ...int) *EstimateUpdateOne {
	euo.mutation.RemoveProbabilityIDs(ids...)
	return euo
}

// RemoveProbabilities removes "probabilities" edges to Probability entities.
func (euo *EstimateUpdateOne) RemoveProbabilities(p ...*Probability) *EstimateUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return euo.RemoveProbabilityIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (euo *EstimateUpdateOne) Select(field string, fields ...string) *EstimateUpdateOne {
	euo.fields = append([]string{field}, fields...)
	return euo
}

// Save executes the query and returns the updated Estimate entity.
func (euo *EstimateUpdateOne) Save(ctx context.Context) (*Estimate, error) {
	return withHooks[*Estimate, EstimateMutation](ctx, euo.sqlSave, euo.mutation, euo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (euo *EstimateUpdateOne) SaveX(ctx context.Context) *Estimate {
	node, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (euo *EstimateUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *EstimateUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (euo *EstimateUpdateOne) sqlSave(ctx context.Context) (_node *Estimate, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   estimate.Table,
			Columns: estimate.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: estimate.FieldID,
			},
		},
	}
	id, ok := euo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Estimate.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := euo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, estimate.FieldID)
		for _, f := range fields {
			if !estimate.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != estimate.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := euo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := euo.mutation.Reason(); ok {
		_spec.SetField(estimate.FieldReason, field.TypeString, value)
	}
	if value, ok := euo.mutation.Created(); ok {
		_spec.SetField(estimate.FieldCreated, field.TypeTime, value)
	}
	if euo.mutation.ForecastCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   estimate.ForecastTable,
			Columns: []string{estimate.ForecastColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: forecast.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.ForecastIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   estimate.ForecastTable,
			Columns: []string{estimate.ForecastColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: forecast.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.mutation.ProbabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   estimate.ProbabilitiesTable,
			Columns: []string{estimate.ProbabilitiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: probability.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.RemovedProbabilitiesIDs(); len(nodes) > 0 && !euo.mutation.ProbabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   estimate.ProbabilitiesTable,
			Columns: []string{estimate.ProbabilitiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: probability.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.ProbabilitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   estimate.ProbabilitiesTable,
			Columns: []string{estimate.ProbabilitiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: probability.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Estimate{config: euo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, euo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{estimate.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	euo.mutation.done = true
	return _node, nil
}