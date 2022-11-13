// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package gqclient

import (
	"context"
	"time"

	"github.com/Khan/genqlient/graphql"
)

// CreateForecastCreateForecast includes the requested fields of the GraphQL type Forecast.
type CreateForecastCreateForecast struct {
	Id string `json:"id"`
}

// GetId returns CreateForecastCreateForecast.Id, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecast) GetId() string { return v.Id }

// CreateForecastResponse is returned by CreateForecast on success.
type CreateForecastResponse struct {
	CreateForecast CreateForecastCreateForecast `json:"createForecast"`
}

// GetCreateForecast returns CreateForecastResponse.CreateForecast, and is useful for accessing the field via an interface.
func (v *CreateForecastResponse) GetCreateForecast() CreateForecastCreateForecast {
	return v.CreateForecast
}

// GetForecastsForecastsForecast includes the requested fields of the GraphQL type Forecast.
type GetForecastsForecastsForecast struct {
	Id          string     `json:"id"`
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	Created     time.Time  `json:"created"`
	Closes      time.Time  `json:"closes"`
	Resolves    time.Time  `json:"resolves"`
	Resolution  Resolution `json:"resolution"`
}

// GetId returns GetForecastsForecastsForecast.Id, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetId() string { return v.Id }

// GetSummary returns GetForecastsForecastsForecast.Summary, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetSummary() string { return v.Summary }

// GetDescription returns GetForecastsForecastsForecast.Description, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetDescription() string { return v.Description }

// GetCreated returns GetForecastsForecastsForecast.Created, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetCreated() time.Time { return v.Created }

// GetCloses returns GetForecastsForecastsForecast.Closes, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetCloses() time.Time { return v.Closes }

// GetResolves returns GetForecastsForecastsForecast.Resolves, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetResolves() time.Time { return v.Resolves }

// GetResolution returns GetForecastsForecastsForecast.Resolution, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetResolution() Resolution { return v.Resolution }

// GetForecastsResponse is returned by GetForecasts on success.
type GetForecastsResponse struct {
	Forecasts []GetForecastsForecastsForecast `json:"forecasts"`
}

// GetForecasts returns GetForecastsResponse.Forecasts, and is useful for accessing the field via an interface.
func (v *GetForecastsResponse) GetForecasts() []GetForecastsForecastsForecast { return v.Forecasts }

type NewForecast struct {
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Resolves    time.Time `json:"resolves"`
	Closes      time.Time `json:"closes"`
}

// GetSummary returns NewForecast.Summary, and is useful for accessing the field via an interface.
func (v *NewForecast) GetSummary() string { return v.Summary }

// GetDescription returns NewForecast.Description, and is useful for accessing the field via an interface.
func (v *NewForecast) GetDescription() string { return v.Description }

// GetResolves returns NewForecast.Resolves, and is useful for accessing the field via an interface.
func (v *NewForecast) GetResolves() time.Time { return v.Resolves }

// GetCloses returns NewForecast.Closes, and is useful for accessing the field via an interface.
func (v *NewForecast) GetCloses() time.Time { return v.Closes }

type Resolution string

const (
	ResolutionTrue          Resolution = "TRUE"
	ResolutionFalse         Resolution = "FALSE"
	ResolutionNotApplicable Resolution = "NOT_APPLICABLE"
	ResolutionUnresolved    Resolution = "UNRESOLVED"
)

// __CreateForecastInput is used internally by genqlient
type __CreateForecastInput struct {
	Input NewForecast `json:"input"`
}

// GetInput returns __CreateForecastInput.Input, and is useful for accessing the field via an interface.
func (v *__CreateForecastInput) GetInput() NewForecast { return v.Input }

func CreateForecast(
	ctx context.Context,
	client graphql.Client,
	input NewForecast,
) (*CreateForecastResponse, error) {
	req := &graphql.Request{
		OpName: "CreateForecast",
		Query: `
mutation CreateForecast ($input: NewForecast!) {
	createForecast(input: $input) {
		id
	}
}
`,
		Variables: &__CreateForecastInput{
			Input: input,
		},
	}
	var err error

	var data CreateForecastResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func GetForecasts(
	ctx context.Context,
	client graphql.Client,
) (*GetForecastsResponse, error) {
	req := &graphql.Request{
		OpName: "GetForecasts",
		Query: `
query GetForecasts {
	forecasts {
		id
		summary
		description
		created
		closes
		resolves
		resolution
	}
}
`,
	}
	var err error

	var data GetForecastsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
