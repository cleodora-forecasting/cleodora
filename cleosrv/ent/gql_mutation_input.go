// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/forecast"
)

// CreateEstimateInput represents a mutation input for creating estimates.
type CreateEstimateInput struct {
	Reason         *string
	Created        *time.Time
	ForecastID     *int
	ProbabilityIDs []int
}

// Mutate applies the CreateEstimateInput on the EstimateMutation builder.
func (i *CreateEstimateInput) Mutate(m *EstimateMutation) {
	if v := i.Reason; v != nil {
		m.SetReason(*v)
	}
	if v := i.Created; v != nil {
		m.SetCreated(*v)
	}
	if v := i.ForecastID; v != nil {
		m.SetForecastID(*v)
	}
	if v := i.ProbabilityIDs; len(v) > 0 {
		m.AddProbabilityIDs(v...)
	}
}

// SetInput applies the change-set in the CreateEstimateInput on the EstimateCreate builder.
func (c *EstimateCreate) SetInput(i CreateEstimateInput) *EstimateCreate {
	i.Mutate(c.Mutation())
	return c
}

// CreateForecastInput represents a mutation input for creating forecasts.
type CreateForecastInput struct {
	Title       string
	Description *string
	Created     *time.Time
	Resolves    time.Time
	Closes      *time.Time
	Resolution  *forecast.Resolution
	EstimateIDs []int
}

// Mutate applies the CreateForecastInput on the ForecastMutation builder.
func (i *CreateForecastInput) Mutate(m *ForecastMutation) {
	m.SetTitle(i.Title)
	if v := i.Description; v != nil {
		m.SetDescription(*v)
	}
	if v := i.Created; v != nil {
		m.SetCreated(*v)
	}
	m.SetResolves(i.Resolves)
	if v := i.Closes; v != nil {
		m.SetCloses(*v)
	}
	if v := i.Resolution; v != nil {
		m.SetResolution(*v)
	}
	if v := i.EstimateIDs; len(v) > 0 {
		m.AddEstimateIDs(v...)
	}
}

// SetInput applies the change-set in the CreateForecastInput on the ForecastCreate builder.
func (c *ForecastCreate) SetInput(i CreateForecastInput) *ForecastCreate {
	i.Mutate(c.Mutation())
	return c
}

// CreateOutcomeInput represents a mutation input for creating outcomes.
type CreateOutcomeInput struct {
	Text           string
	Correct        *bool
	ProbabilityIDs []int
}

// Mutate applies the CreateOutcomeInput on the OutcomeMutation builder.
func (i *CreateOutcomeInput) Mutate(m *OutcomeMutation) {
	m.SetText(i.Text)
	if v := i.Correct; v != nil {
		m.SetCorrect(*v)
	}
	if v := i.ProbabilityIDs; len(v) > 0 {
		m.AddProbabilityIDs(v...)
	}
}

// SetInput applies the change-set in the CreateOutcomeInput on the OutcomeCreate builder.
func (c *OutcomeCreate) SetInput(i CreateOutcomeInput) *OutcomeCreate {
	i.Mutate(c.Mutation())
	return c
}

// CreateProbabilityInput represents a mutation input for creating probabilities.
type CreateProbabilityInput struct {
	Value      int
	EstimateID *int
	OutcomeID  *int
}

// Mutate applies the CreateProbabilityInput on the ProbabilityMutation builder.
func (i *CreateProbabilityInput) Mutate(m *ProbabilityMutation) {
	m.SetValue(i.Value)
	if v := i.EstimateID; v != nil {
		m.SetEstimateID(*v)
	}
	if v := i.OutcomeID; v != nil {
		m.SetOutcomeID(*v)
	}
}

// SetInput applies the change-set in the CreateProbabilityInput on the ProbabilityCreate builder.
func (c *ProbabilityCreate) SetInput(i CreateProbabilityInput) *ProbabilityCreate {
	i.Mutate(c.Mutation())
	return c
}
