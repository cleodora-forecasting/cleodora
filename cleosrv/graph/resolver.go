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
