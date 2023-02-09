// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/predicate"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/probability"
)

// ProbabilityDelete is the builder for deleting a Probability entity.
type ProbabilityDelete struct {
	config
	hooks    []Hook
	mutation *ProbabilityMutation
}

// Where appends a list predicates to the ProbabilityDelete builder.
func (pd *ProbabilityDelete) Where(ps ...predicate.Probability) *ProbabilityDelete {
	pd.mutation.Where(ps...)
	return pd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pd *ProbabilityDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, ProbabilityMutation](ctx, pd.sqlExec, pd.mutation, pd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (pd *ProbabilityDelete) ExecX(ctx context.Context) int {
	n, err := pd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pd *ProbabilityDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: probability.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: probability.FieldID,
			},
		},
	}
	if ps := pd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	pd.mutation.done = true
	return affected, err
}

// ProbabilityDeleteOne is the builder for deleting a single Probability entity.
type ProbabilityDeleteOne struct {
	pd *ProbabilityDelete
}

// Where appends a list predicates to the ProbabilityDelete builder.
func (pdo *ProbabilityDeleteOne) Where(ps ...predicate.Probability) *ProbabilityDeleteOne {
	pdo.pd.mutation.Where(ps...)
	return pdo
}

// Exec executes the deletion query.
func (pdo *ProbabilityDeleteOne) Exec(ctx context.Context) error {
	n, err := pdo.pd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{probability.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pdo *ProbabilityDeleteOne) ExecX(ctx context.Context) {
	if err := pdo.Exec(ctx); err != nil {
		panic(err)
	}
}
