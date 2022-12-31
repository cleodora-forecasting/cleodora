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

func (r *Resolver) AddDummyData() error {
	var forecastCount int64
	res := r.db.Model(&dbmodel.Forecast{}).Count(&forecastCount)
	if res.Error != nil {
		return res.Error
	}
	if forecastCount > 0 {
		// Since the DB is not empty don't (re-)create dummy data
		return nil
	}

	err := createDummyForecast_TheFabelmans(r.db)
	if err != nil {
		return err
	}

	err = createDummyForecast_CPEExam(r.db)
	if err != nil {
		return err
	}

	err = createDummyForecast_Contributors(r.db)
	if err != nil {
		return err
	}

	return nil
}

func createDummyForecast_TheFabelmans(db *gorm.DB) error {
	outcomeYes := dbmodel.Outcome{
		Text:    "Yes",
		Correct: false,
	}
	ret := db.Create(&outcomeYes)
	if ret.Error != nil {
		return ret.Error
	}

	outcomeNo := dbmodel.Outcome{
		Text:    "No",
		Correct: false,
	}
	ret = db.Create(&outcomeNo)
	if ret.Error != nil {
		return ret.Error
	}

	forecast := dbmodel.Forecast{
		Title:       "Will \"The Fabelmans\" win \"Best Picture\" at the Oscars 2023?",
		Description: "",
		Created: timeParseOrPanic(
			time.RFC3339,
			"2022-10-30T17:05:00+01:00",
		),
		Closes: nil,
		Resolves: timeParseOrPanic(
			time.RFC3339,
			"2023-03-11T23:59:00+01:00",
		),
		Resolution: dbmodel.ResolutionUnresolved,
		Estimates: []dbmodel.Estimate{
			{
				Created: timeParseOrPanic(
					time.RFC3339,
					"2022-10-30T17:05:00+01:00",
				),
				Reason: "It's a great film and it's of the type that the" +
					" Academy loves!",
				Probabilities: []dbmodel.Probability{
					{
						Value:     30,
						OutcomeID: outcomeYes.ID,
					},
					{
						Value:     70,
						OutcomeID: outcomeNo.ID,
					},
				},
			},
		},
	}

	ret = db.Create(&forecast)

	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

func createDummyForecast_CPEExam(db *gorm.DB) error {
	outcomeC2A := dbmodel.Outcome{
		Text:    "C2 - Grade A",
		Correct: true, // This forecast has been resolved
	}
	ret := db.Create(&outcomeC2A)
	if ret.Error != nil {
		return ret.Error
	}

	outcomeC2B := dbmodel.Outcome{
		Text:    "C2 - Grade B",
		Correct: false,
	}
	ret = db.Create(&outcomeC2B)
	if ret.Error != nil {
		return ret.Error
	}

	outcomeC2C := dbmodel.Outcome{
		Text:    "C2 - Grade C",
		Correct: false,
	}
	ret = db.Create(&outcomeC2C)
	if ret.Error != nil {
		return ret.Error
	}

	outcomeC1 := dbmodel.Outcome{
		Text:    "C1",
		Correct: false,
	}
	ret = db.Create(&outcomeC1)
	if ret.Error != nil {
		return ret.Error
	}

	outcomeFail := dbmodel.Outcome{
		Text:    "Fail",
		Correct: false,
	}
	ret = db.Create(&outcomeFail)
	if ret.Error != nil {
		return ret.Error
	}

	forecast := dbmodel.Forecast{
		Title: "What grade will I get in my upcoming CPE exam?",
		Description: "CPE C2 exam. Grade C1 is the worst passing grade. " +
			"It's a language exam using the Common European Framework of" +
			" Reference for Languages.",
		Created: timeParseOrPanic(
			time.RFC3339,
			"2022-10-15T13:10:00+02:00",
		),
		Closes: timeParseOrPanicPtr(
			time.RFC3339,
			"2022-11-11T23:59:00+01:00",
		),
		Resolves: timeParseOrPanic(
			time.RFC3339,
			"2022-12-01T09:00:00+01:00",
		),
		Resolution: dbmodel.ResolutionResolved,
		Estimates: []dbmodel.Estimate{
			{
				Created: timeParseOrPanic(
					time.RFC3339,
					"2022-10-15T13:10:00+02:00",
				),
				Reason: "I'm well prepared and performed well on test" +
					" exams.",
				Probabilities: []dbmodel.Probability{
					{
						Value:     40,
						OutcomeID: outcomeC2A.ID,
					},
					{
						Value:     30,
						OutcomeID: outcomeC2B.ID,
					},
					{
						Value:     20,
						OutcomeID: outcomeC2C.ID,
					},
					{
						Value:     8,
						OutcomeID: outcomeC1.ID,
					},
					{
						Value:     2,
						OutcomeID: outcomeFail.ID,
					},
				},
			},
		},
	}

	ret = db.Create(&forecast)

	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

func createDummyForecast_Contributors(db *gorm.DB) error {
	outcomeYes := dbmodel.Outcome{
		Text:    "Yes",
		Correct: false,
	}
	ret := db.Create(&outcomeYes)
	if ret.Error != nil {
		return ret.Error
	}

	outcomeNo := dbmodel.Outcome{
		Text:    "No",
		Correct: false,
	}
	ret = db.Create(&outcomeNo)
	if ret.Error != nil {
		return ret.Error
	}

	forecast := dbmodel.Forecast{
		Title: "Will the number of contributors to \"Cleodora\" be more " +
			"than 3 at the end of 2022?",
		Description: "A contributor is any person who has made a commit" +
			" in any Git repository of the cleodora-forecasting GitHub" +
			" organization.",
		Created: timeParseOrPanic(
			time.RFC3339,
			"2022-10-01T11:00:00+01:00",
		),
		Closes: nil,
		Resolves: timeParseOrPanic(
			time.RFC3339,
			"2022-12-31T23:59:00+01:00",
		),
		Resolution: dbmodel.ResolutionUnresolved,
		Estimates: []dbmodel.Estimate{
			{
				Created: timeParseOrPanic(
					time.RFC3339,
					"2022-10-01T11:00:00+01:00",
				),
				Reason: "It's a new project and people are usually busy.",
				Probabilities: []dbmodel.Probability{
					{
						Value:     15,
						OutcomeID: outcomeYes.ID,
					},
					{
						Value:     85,
						OutcomeID: outcomeNo.ID,
					},
				},
			},
			{
				Created: timeParseOrPanic(
					time.RFC3339,
					"2022-12-24T23:33:04+01:00",
				),
				Reason: "Despite multiple people expressing interest nobody " +
					"has contributed so far. The year is almost over.",
				Probabilities: []dbmodel.Probability{
					{
						Value:     1,
						OutcomeID: outcomeYes.ID,
					},
					{
						Value:     99,
						OutcomeID: outcomeNo.ID,
					},
				},
			},
		},
	}

	ret = db.Create(&forecast)

	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

func timeParseOrPanic(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return t
}

func timeParseOrPanicPtr(layout string, value string) *time.Time {
	t := timeParseOrPanic(layout, value)

	return &t
}

func validateNewForecast(forecast model.NewForecast) error {
	var validationErr *multierror.Error
	if forecast.Title == "" {
		validationErr = multierror.Append(
			validationErr,
			errors.New("title can't be empty"),
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

	return []dbmodel.Estimate{
		{
			Created:       time.Now(),
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
