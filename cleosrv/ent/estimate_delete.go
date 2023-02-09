// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/estimate"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/predicate"
)

// EstimateDelete is the builder for deleting a Estimate entity.
type EstimateDelete struct {
	config
	hooks    []Hook
	mutation *EstimateMutation
}

// Where appends a list predicates to the EstimateDelete builder.
func (ed *EstimateDelete) Where(ps ...predicate.Estimate) *EstimateDelete {
	ed.mutation.Where(ps...)
	return ed
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ed *EstimateDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, EstimateMutation](ctx, ed.sqlExec, ed.mutation, ed.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ed *EstimateDelete) ExecX(ctx context.Context) int {
	n, err := ed.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ed *EstimateDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: estimate.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: estimate.FieldID,
			},
		},
	}
	if ps := ed.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ed.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ed.mutation.done = true
	return affected, err
}

// EstimateDeleteOne is the builder for deleting a single Estimate entity.
type EstimateDeleteOne struct {
	ed *EstimateDelete
}

// Where appends a list predicates to the EstimateDelete builder.
func (edo *EstimateDeleteOne) Where(ps ...predicate.Estimate) *EstimateDeleteOne {
	edo.ed.mutation.Where(ps...)
	return edo
}

// Exec executes the deletion query.
func (edo *EstimateDeleteOne) Exec(ctx context.Context) error {
	n, err := edo.ed.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{estimate.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (edo *EstimateDeleteOne) ExecX(ctx context.Context) {
	if err := edo.Exec(ctx); err != nil {
		panic(err)
	}
}