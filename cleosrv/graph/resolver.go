package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"time"

	"gorm.io/gorm"

	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db        *gorm.DB
	forecasts []*model.Forecast
}

func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{
		db: db,
	}
}

func (r *Resolver) AddDummyData() {
	r.forecasts = append(
		r.forecasts,
		&model.Forecast{
			ID:          "1",
			Title:       "Will \"The Fabelmans\" win \"Best Picture\" at the Oscars 2023?",
			Description: "",
			Created:     time.Now(),
			Closes: timeParseOrPanicPtr(
				time.RFC3339,
				"2023-03-11T23:59:00+00:00",
			),
			Resolves: timeParseOrPanic(
				time.RFC3339,
				"2023-03-11T23:59:00+00:00",
			),
			Resolution: model.ResolutionUnresolved,
			Estimates: []*model.Estimate{
				{
					ID:      "1",
					Created: time.Now(),
					Reason: "It's a great film and it's of the type that the" +
						" Academy loves!",
					Probabilities: []*model.Probability{
						{
							ID:    "1",
							Value: 30,
							Outcome: &model.Outcome{
								ID:      "1",
								Text:    "Yes",
								Correct: false,
							},
						},
						{
							ID:    "2",
							Value: 70,
							Outcome: &model.Outcome{
								ID:      "2",
								Text:    "No",
								Correct: false,
							},
						},
					},
				},
			},
		},
		&model.Forecast{
			ID:    "2",
			Title: "What grade will I get in my upcoming exam?",
			Description: "CPE C2 exam. Grade C1 is the worst passing grade. " +
				"It's a language exam using the Common European Framework of" +
				" Reference for Languages.",
			Created: time.Now(),
			Closes: timeParseOrPanicPtr(
				time.RFC3339,
				"2022-11-11T23:59:00+00:00",
			),
			Resolves: timeParseOrPanic(
				time.RFC3339,
				"2022-12-01T09:00:00+00:00",
			),
			Resolution: model.ResolutionUnresolved,
			Estimates: []*model.Estimate{
				{
					ID:      "2",
					Created: time.Now(),
					Reason: "I'm well prepared and performed well on test" +
						" exams.",
					Probabilities: []*model.Probability{
						{
							ID:    "3",
							Value: 40,
							Outcome: &model.Outcome{
								ID:      "3",
								Text:    "C2 Grade A",
								Correct: false,
							},
						},
						{
							ID:    "4",
							Value: 30,
							Outcome: &model.Outcome{
								ID:      "4",
								Text:    "C2 Grade B",
								Correct: false,
							},
						},
						{
							ID:    "5",
							Value: 20,
							Outcome: &model.Outcome{
								ID:      "5",
								Text:    "C2 Grade C",
								Correct: false,
							},
						},
						{
							ID:    "6",
							Value: 8,
							Outcome: &model.Outcome{
								ID:      "6",
								Text:    "C1",
								Correct: false,
							},
						},
						{
							ID:    "7",
							Value: 2,
							Outcome: &model.Outcome{
								ID:      "7",
								Text:    "Fail",
								Correct: false,
							},
						},
					},
				},
			},
		},
		&model.Forecast{
			ID:    "3",
			Title: "Will the number of contributors for \"Cleodora\" be more than 3 at the end of 2022?",
			Description: "A contributor is any person who has made a commit" +
				" in any Git repository of the cleodora-forecasting GitHub" +
				" organization.",
			Created: time.Now(),
			Closes: timeParseOrPanicPtr(
				time.RFC3339,
				"2022-12-31T23:59:00+00:00",
			),
			Resolves: timeParseOrPanic(
				time.RFC3339,
				"2022-12-31T23:59:00+00:00",
			),
			Resolution: model.ResolutionUnresolved,
			Estimates: []*model.Estimate{
				{
					ID:      "3",
					Created: time.Now(),
					Reason:  "It's a new project and people are usually busy.",
					Probabilities: []*model.Probability{
						{
							ID:    "8",
							Value: 15,
							Outcome: &model.Outcome{
								ID:      "8",
								Text:    "Yes",
								Correct: false,
							},
						},
						{
							ID:    "9",
							Value: 85,
							Outcome: &model.Outcome{
								ID:      "9",
								Text:    "No",
								Correct: false,
							},
						},
					},
				},
			},
		},
	)
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
