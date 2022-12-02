package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"time"

	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	forecasts []*model.Forecast
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
			Outcomes: []*model.Outcome{
				{
					ID:      "1",
					Text:    "Yes",
					Correct: false,
				},
				{
					ID:      "2",
					Text:    "No",
					Correct: false,
				},
			},
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
			Title: "Will I get an A in my upcoming exam?",
			Description: "The forecast resolves as true if and only if I get" +
				" the highest marks.",
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
			Outcomes: []*model.Outcome{
				{
					ID:      "3",
					Text:    "Yes",
					Correct: false,
				},
				{
					ID:      "4",
					Text:    "No",
					Correct: false,
				},
			},
			Estimates: []*model.Estimate{
				{
					ID:      "2",
					Created: time.Now(),
					Reason: "I'm well prepared and performed well on test" +
						" exams.",
					Probabilities: []*model.Probability{
						{
							ID:    "3",
							Value: 90,
							Outcome: &model.Outcome{
								ID:      "3",
								Text:    "Yes",
								Correct: false,
							},
						},
						{
							ID:    "4",
							Value: 10,
							Outcome: &model.Outcome{
								ID:      "4",
								Text:    "No",
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
			Outcomes: []*model.Outcome{
				{
					ID:      "5",
					Text:    "Yes",
					Correct: false,
				},
				{
					ID:      "6",
					Text:    "No",
					Correct: false,
				},
			},
			Estimates: []*model.Estimate{
				{
					ID:      "3",
					Created: time.Now(),
					Reason:  "It's a new project and people are usually busy.",
					Probabilities: []*model.Probability{
						{
							ID:    "5",
							Value: 15,
							Outcome: &model.Outcome{
								ID:      "5",
								Text:    "Yes",
								Correct: false,
							},
						},
						{
							ID:    "6",
							Value: 85,
							Outcome: &model.Outcome{
								ID:      "6",
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
