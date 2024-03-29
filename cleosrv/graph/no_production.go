//go:build !production

package graph

import (
	"time"

	"gorm.io/gorm"

	"github.com/cleodora-forecasting/cleodora/cleosrv/dbmodel"
)

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

	brierScore := 0.4968

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
				BrierScore: &brierScore,
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
