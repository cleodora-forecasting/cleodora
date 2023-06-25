package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"
	"errors"
	"fmt"
	"html"
	"strconv"
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

// createEstimate validates the input data and creates a new Estimate.
func (r *Resolver) createEstimate(
	ctx context.Context,
	forecastID string,
	estimate model.NewEstimate,
) (*model.Estimate, error) {
	var e *model.Estimate
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error
		e, err = createEstimate(tx, ctx, forecastID, estimate)
		return err
	})
	return e, err
}

// createEstimate is a helper function that assumes that it is already wrapped
// in a DB transaction.
func createEstimate(
	tx *gorm.DB,
	ctx context.Context,
	forecastID string,
	estimate model.NewEstimate,
) (*model.Estimate, error) {
	if err := validateNewEstimate(&estimate, false); err != nil {
		return nil, fmt.Errorf("error validating NewEstimate: %w", err)
	}
	forecast := dbmodel.Forecast{}
	ret := tx.Where("id = ?", forecastID).First(&forecast)
	if ret.Error != nil {
		return nil, fmt.Errorf("error getting Forecast with ID %v: %w", forecastID, ret.Error)
	}

	var validOutcomeIds []string

	ret = tx.Model(&dbmodel.Outcome{}).Joins(
		"INNER JOIN probabilities ON probabilities.outcome_id == outcomes.id",
	).Joins(
		"INNER JOIN estimates ON estimates.id == probabilities.estimate_id",
	).Where("estimates.forecast_id = ?", forecastID).
		Distinct().Pluck("outcomes.id", &validOutcomeIds)
	if ret.Error != nil {
		return nil, fmt.Errorf("error getting outcome IDs: %w", ret.Error)
	}

	var submittedOutcomeIds []string
	for _, p := range estimate.Probabilities {
		submittedOutcomeIds = append(submittedOutcomeIds, *p.OutcomeID)
	}

	invalidOutcomeIds := stringSetDiff(submittedOutcomeIds, validOutcomeIds)
	if len(invalidOutcomeIds) > 0 {
		return nil, fmt.Errorf("invalid Outcome IDs: %v", invalidOutcomeIds)
	}

	missingOutcomeIds := stringSetDiff(validOutcomeIds, submittedOutcomeIds)
	if len(missingOutcomeIds) > 0 {
		return nil, fmt.Errorf("missing Outcome IDs: %v", missingOutcomeIds)
	}

	var probabilities []dbmodel.Probability

	for _, p := range estimate.Probabilities {
		outcomeID, err := strconv.ParseUint(*p.OutcomeID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse %v as uint: %w", outcomeID, err)
		}
		probabilities = append(
			probabilities,
			dbmodel.Probability{
				Value:     p.Value,
				OutcomeID: uint(outcomeID),
			},
		)
	}

	dbEstimate := dbmodel.Estimate{
		Created:       *estimate.Created,
		Reason:        estimate.Reason,
		Probabilities: probabilities,
	}

	err := tx.Model(&forecast).Association("Estimates").Append(&dbEstimate)
	if err != nil {
		return nil, fmt.Errorf("error creating Estimate: %w", err)
	}

	return convertEstimateDBToGQL(dbEstimate), nil
}

// stringSetDiff returns the strings in the slice a that are not in slice b.
// The slices are treated as sets so duplicate values will disappear.
func stringSetDiff(a, b []string) []string {
	setA := map[string]bool{}
	for _, e := range a {
		setA[e] = true
	}
	setB := map[string]bool{}
	for _, e := range b {
		setB[e] = true
	}
	var result []string
	for e := range setA {
		if !setB[e] {
			result = append(result, e)
		}
	}
	return result
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
func validateNewEstimate(estimate *model.NewEstimate, duringCreateForecast bool) error {
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
		if duringCreateForecast {
			if p.Outcome == nil {
				validationErr = multierror.Append(
					validationErr,
					errors.New("NewOutcome must be set when creating a new forecast"),
				)
			} else {
				if p.OutcomeID != nil {
					validationErr = multierror.Append(
						validationErr,
						errors.New("outcomeId must be unset when creating a new forecast"),
					)
				}
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
			}
		} else { // !duringCreateForecast
			if p.Outcome != nil {
				validationErr = multierror.Append(
					validationErr,
					errors.New("NewOutcome must be unset when adding an estimate"),
				)
			}
			if p.OutcomeID == nil {
				validationErr = multierror.Append(
					validationErr,
					errors.New("outcomeId must be set when adding an estimate"),
				)
			} else if *p.OutcomeID == "" {
				validationErr = multierror.Append(
					validationErr,
					errors.New("outcomeId must not be empty when adding an estimate"),
				)
			}
		}
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
		estimate.Created = timeToUTCPtr(*estimate.Created) // convert to UTC
		if estimate.Created.After(now) {
			validationErr = multierror.Append(
				validationErr,
				fmt.Errorf(
					"'created' can't be in the future: %v",
					estimate.Created,
				),
			)
		}
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
					Text:    html.EscapeString(p.Outcome.Text),
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
