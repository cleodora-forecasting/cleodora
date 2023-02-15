package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"errors"
	"fmt"
	"html"
	"time"

	"github.com/hashicorp/go-multierror"
	"gorm.io/gorm"

	"github.com/cleodora-forecasting/cleodora/cleosrv/dbmodel"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db *gorm.DB
}

func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{
		db: db,
	}
}

func validateNewForecast(forecast model.NewForecast) error {
	var validationErr *multierror.Error
	if forecast.Title == "" {
		validationErr = multierror.Append(
			validationErr,
			errors.New("title can't be empty"),
		)
	}
	if forecast.Created.After(time.Now().UTC()) {
		validationErr = multierror.Append(
			validationErr,
			errors.New("'created' can't be in the future"),
		)
	}
	if forecast.Closes != nil && forecast.Closes.After(forecast.Resolves) {
		validationErr = multierror.Append(
			validationErr,
			fmt.Errorf(
				"'Closes' can't be set to a later date than 'Resolves'. "+
					"Closes is '%v'. Resolves is '%v'",
				*forecast.Closes,
				forecast.Resolves,
			),
		)
	}
	if forecast.Resolves.Before(*forecast.Created) {
		validationErr = multierror.Append(
			validationErr,
			fmt.Errorf(
				"'Resolves' can't be set to an earlier date than 'Created'. "+
					"Resolves is '%v'. Created is '%v'",
				forecast.Resolves,
				*forecast.Created,
			),
		)
	}
	return validationErr.ErrorOrNil()
}

func validateNewEstimate(estimate model.NewEstimate) error {
	var validationErr *multierror.Error
	if estimate.Reason == "" {
		validationErr = multierror.Append(
			validationErr,
			errors.New("'reason' can't be empty"),
		)
	}
	if len(estimate.Probabilities) == 0 {
		validationErr = multierror.Append(
			validationErr,
			errors.New("probabilities can't be empty"),
		)
	}
	sumProbabilities := 0
	existingOutcomes := map[string]bool{}
	for _, p := range estimate.Probabilities {
		if p.Outcome.Text == "" {
			validationErr = multierror.Append(
				validationErr,
				errors.New("outcome text can't be empty"),
			)
		}
		if _, ok := existingOutcomes[p.Outcome.Text]; ok {
			validationErr = multierror.Append(
				validationErr,
				fmt.Errorf("outcome '%v' is a duplicate", p.Outcome.Text),
			)
		}
		existingOutcomes[p.Outcome.Text] = true
		if p.Value < 0 || p.Value > 100 {
			validationErr = multierror.Append(
				validationErr,
				fmt.Errorf("probabilities must be between 0 and 100, not %v", p.Value),
			)
		}
		sumProbabilities += p.Value
	}
	if sumProbabilities != 100 {
		validationErr = multierror.Append(
			validationErr,
			fmt.Errorf("probabilities must add up to 100, not %v", sumProbabilities),
		)
	}
	return validationErr.ErrorOrNil()
}

func convertNewEstimateToDBEstimate(estimate model.NewEstimate) []dbmodel.Estimate {
	var probabilities []dbmodel.Probability

	for _, p := range estimate.Probabilities {
		probabilities = append(
			probabilities,
			dbmodel.Probability{
				Value: p.Value,
				Outcome: dbmodel.Outcome{
					Text:    html.EscapeString(p.Outcome.Text),
					Correct: false,
				},
			},
		)
	}

	created := time.Now().UTC()
	if estimate.Created != nil {
		created = *estimate.Created
	}

	return []dbmodel.Estimate{
		{
			Created:       created,
			Reason:        html.EscapeString(estimate.Reason),
			Probabilities: probabilities,
		},
	}
}

func convertEstimatesDBToGQL(dbEstimates []dbmodel.Estimate) []*model.Estimate {
	var gqlEstimates []*model.Estimate
	for _, e := range dbEstimates {
		gqlEstimates = append(
			gqlEstimates,
			convertEstimateDBToGQL(e),
		)
	}
	return gqlEstimates
}

func convertEstimateDBToGQL(dbEstimate dbmodel.Estimate) *model.Estimate {
	return &model.Estimate{
		ID:            fmt.Sprint(dbEstimate.ID),
		Created:       dbEstimate.Created,
		Reason:        dbEstimate.Reason,
		Probabilities: convertProbabilitiesDBToGQL(dbEstimate.Probabilities),
	}
}

func convertProbabilitiesDBToGQL(dbProbabilities []dbmodel.Probability) []*model.Probability {
	var gqlProbabilities []*model.Probability
	for _, p := range dbProbabilities {
		gqlProbabilities = append(
			gqlProbabilities,
			convertProbabilityDBToGQL(p),
		)
	}
	return gqlProbabilities
}

func convertProbabilityDBToGQL(dbProbability dbmodel.Probability) *model.Probability {
	return &model.Probability{
		ID:      fmt.Sprint(dbProbability.ID),
		Value:   dbProbability.Value,
		Outcome: convertOutcomeDBToGQL(dbProbability.Outcome),
	}
}

func convertOutcomeDBToGQL(dbOutcome dbmodel.Outcome) *model.Outcome {
	return &model.Outcome{
		ID:      fmt.Sprint(dbOutcome.ID),
		Text:    dbOutcome.Text,
		Correct: dbOutcome.Correct,
	}
}
