// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/outcome"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/predicate"
)

// OutcomeDelete is the builder for deleting a Outcome entity.
type OutcomeDelete struct {
	config
	hooks    []Hook
	mutation *OutcomeMutation
}

// Where appends a list predicates to the OutcomeDelete builder.
func (od *OutcomeDelete) Where(ps ...predicate.Outcome) *OutcomeDelete {
	od.mutation.Where(ps...)
	return od
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (od *OutcomeDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, OutcomeMutation](ctx, od.sqlExec, od.mutation, od.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (od *OutcomeDelete) ExecX(ctx context.Context) int {
	n, err := od.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (od *OutcomeDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: outcome.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: outcome.FieldID,
			},
		},
	}
	if ps := od.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, od.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	od.mutation.done = true
	return affected, err
}

// OutcomeDeleteOne is the builder for deleting a single Outcome entity.
type OutcomeDeleteOne struct {
	od *OutcomeDelete
}

// Where appends a list predicates to the OutcomeDelete builder.
func (odo *OutcomeDeleteOne) Where(ps ...predicate.Outcome) *OutcomeDeleteOne {
	odo.od.mutation.Where(ps...)
	return odo
}

// Exec executes the deletion query.
func (odo *OutcomeDeleteOne) Exec(ctx context.Context) error {
	n, err := odo.od.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{outcome.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (odo *OutcomeDeleteOne) ExecX(ctx context.Context) {
	if err := odo.Exec(ctx); err != nil {
		panic(err)
	}
}
