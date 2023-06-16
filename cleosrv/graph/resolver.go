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

// validateNewForecast validates and makes some automatic changes to the
// forecast if appropriate.
func validateNewForecast(forecast *model.NewForecast) error {
	var validationErr *multierror.Error
	now := time.Now().UTC()

	// Validate title
	if forecast.Title == "" {
		validationErr = multierror.Append(
			validationErr,
			errors.New("title can't be empty"),
		)
	}

	//  Validate created
	if forecast.Created == nil || forecast.Created.IsZero() {
		forecast.Created = &now
	}
	if forecast.Created.After(now) {
		validationErr = multierror.Append(
			validationErr,
			errors.New("'created' can't be in the future"),
		)
	}
	forecast.Created = timeToUTCPtr(*forecast.Created)

	// Validate closes
	if forecast.Closes != nil {
		forecast.Closes = timeToUTCPtr(*forecast.Closes)
		if forecast.Closes.After(forecast.Resolves) {
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
		if forecast.Closes.IsZero() {
			// for more consistent DB handling of the data
			forecast.Closes = nil
		}
	}

	// Validate resolves
	forecast.Resolves = timeToUTC(forecast.Resolves)
	if forecast.Resolves.IsZero() {
		validationErr = multierror.Append(
			validationErr,
			errors.New("'resolves' can't be the zero time"),
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

// validateNewEstimate validates and makes some automatic changes to the
// estimate if appropriate.
func validateNewEstimate(estimate model.NewEstimate) error {
	var validationErr *multierror.Error
	now := time.Now().UTC()

	// Validate reason
	if estimate.Reason == "" {
		validationErr = multierror.Append(
			validationErr,
			errors.New("'reason' can't be empty"),
		)
	}

	// Validate probabilities
	if len(estimate.Probabilities) == 0 {
		validationErr = multierror.Append(
			validationErr,
			errors.New("probabilities can't be empty"),
		)
	}
	sumProbabilities := 0
	existingOutcomes := map[string]bool{}
	for _, p := range estimate.Probabilities {
		outcomeText := *p.Outcome.Text
		if outcomeText == "" {
			validationErr = multierror.Append(
				validationErr,
				errors.New("outcome text can't be empty"),
			)
		}
		if _, ok := existingOutcomes[outcomeText]; ok {
			validationErr = multierror.Append(
				validationErr,
				fmt.Errorf("outcome '%v' is a duplicate", outcomeText),
			)
		}
		existingOutcomes[outcomeText] = true
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

	// Validate created
	if estimate.Created != nil {
		estimate.Created = timeToUTCPtr(*estimate.Created)
	}
	if estimate.Created == nil || estimate.Created.IsZero() {
		estimate.Created = &now
	}

	return validationErr.ErrorOrNil()
}

func timeToUTC(t time.Time) time.Time {
	return t.UTC()
}

func timeToUTCPtr(t time.Time) *time.Time {
	temp := t.UTC()
	return &temp
}

func convertNewEstimateToDBEstimate(estimate model.NewEstimate) []dbmodel.Estimate {
	var probabilities []dbmodel.Probability

	for _, p := range estimate.Probabilities {
		probabilities = append(
			probabilities,
			dbmodel.Probability{
				Value: p.Value,
				Outcome: dbmodel.Outcome{
					Text:    html.EscapeString(*p.Outcome.Text),
					Correct: false,
				},
			},
		)
	}

	return []dbmodel.Estimate{
		{
			Created:       *estimate.Created,
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
