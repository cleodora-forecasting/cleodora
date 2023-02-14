// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package integrationtest

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

// GetForecastsForecastsForecast includes the requested fields of the GraphQL type Forecast.
// The GraphQL type's documentation follows.
//
// A prediction about the future.
type GetForecastsForecastsForecast struct {
	Id         string                                            `json:"id"`
	Title      string                                            `json:"title"`
	Resolution Resolution                                        `json:"resolution"`
	Estimates  []*GetForecastsForecastsForecastEstimatesEstimate `json:"estimates"`
}

// GetId returns GetForecastsForecastsForecast.Id, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetId() string { return v.Id }

// GetTitle returns GetForecastsForecastsForecast.Title, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetTitle() string { return v.Title }

// GetResolution returns GetForecastsForecastsForecast.Resolution, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetResolution() Resolution { return v.Resolution }

// GetEstimates returns GetForecastsForecastsForecast.Estimates, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetEstimates() []*GetForecastsForecastsForecastEstimatesEstimate {
	return v.Estimates
}

// GetForecastsForecastsForecastEstimatesEstimate includes the requested fields of the GraphQL type Estimate.
// The GraphQL type's documentation follows.
//
// A list of probabilities (one for each outcome) together with a timestamp and
// an explanation why you made this estimate. Every time you change your mind
// about a forecast you will create a new Estimate.
// All probabilities always add up to 100.
type GetForecastsForecastsForecastEstimatesEstimate struct {
	Id            string                                                                    `json:"id"`
	Probabilities []*GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbability `json:"probabilities"`
}

// GetId returns GetForecastsForecastsForecastEstimatesEstimate.Id, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecastEstimatesEstimate) GetId() string { return v.Id }

// GetProbabilities returns GetForecastsForecastsForecastEstimatesEstimate.Probabilities, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecastEstimatesEstimate) GetProbabilities() []*GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbability {
	return v.Probabilities
}

// GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbability includes the requested fields of the GraphQL type Probability.
// The GraphQL type's documentation follows.
//
// A number between 0 and 100 tied to a specific Outcome. It is always part of
// an Estimate.
type GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbability struct {
	Id      string                                                                        `json:"id"`
	Outcome GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbabilityOutcome `json:"outcome"`
}

// GetId returns GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbability.Id, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbability) GetId() string {
	return v.Id
}

// GetOutcome returns GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbability.Outcome, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbability) GetOutcome() GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbabilityOutcome {
	return v.Outcome
}

// GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbabilityOutcome includes the requested fields of the GraphQL type Outcome.
// The GraphQL type's documentation follows.
//
// The possible results of a forecast. In the simplest case you will only have
// two outcomes: Yes and No.
type GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbabilityOutcome struct {
	Id      string `json:"id"`
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
}

// GetId returns GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbabilityOutcome.Id, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbabilityOutcome) GetId() string {
	return v.Id
}

// GetText returns GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbabilityOutcome.Text, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbabilityOutcome) GetText() string {
	return v.Text
}

// GetCorrect returns GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbabilityOutcome.Correct, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbabilityOutcome) GetCorrect() bool {
	return v.Correct
}

// GetForecastsResponse is returned by GetForecasts on success.
type GetForecastsResponse struct {
	Forecasts []GetForecastsForecastsForecast `json:"forecasts"`
}

// GetForecasts returns GetForecastsResponse.Forecasts, and is useful for accessing the field via an interface.
func (v *GetForecastsResponse) GetForecasts() []GetForecastsForecastsForecast { return v.Forecasts }

type Resolution string

const (
	ResolutionResolved      Resolution = "RESOLVED"
	ResolutionNotApplicable Resolution = "NOT_APPLICABLE"
	ResolutionUnresolved    Resolution = "UNRESOLVED"
)

// ResolveForecastResolveForecast includes the requested fields of the GraphQL type Forecast.
// The GraphQL type's documentation follows.
//
// A prediction about the future.
type ResolveForecastResolveForecast struct {
	Id         string                                             `json:"id"`
	Title      string                                             `json:"title"`
	Resolution Resolution                                         `json:"resolution"`
	Estimates  []*ResolveForecastResolveForecastEstimatesEstimate `json:"estimates"`
}

// GetId returns ResolveForecastResolveForecast.Id, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecast) GetId() string { return v.Id }

// GetTitle returns ResolveForecastResolveForecast.Title, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecast) GetTitle() string { return v.Title }

// GetResolution returns ResolveForecastResolveForecast.Resolution, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecast) GetResolution() Resolution { return v.Resolution }

// GetEstimates returns ResolveForecastResolveForecast.Estimates, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecast) GetEstimates() []*ResolveForecastResolveForecastEstimatesEstimate {
	return v.Estimates
}

// ResolveForecastResolveForecastEstimatesEstimate includes the requested fields of the GraphQL type Estimate.
// The GraphQL type's documentation follows.
//
// A list of probabilities (one for each outcome) together with a timestamp and
// an explanation why you made this estimate. Every time you change your mind
// about a forecast you will create a new Estimate.
// All probabilities always add up to 100.
type ResolveForecastResolveForecastEstimatesEstimate struct {
	Id            string                                                                     `json:"id"`
	Probabilities []*ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbability `json:"probabilities"`
}

// GetId returns ResolveForecastResolveForecastEstimatesEstimate.Id, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecastEstimatesEstimate) GetId() string { return v.Id }

// GetProbabilities returns ResolveForecastResolveForecastEstimatesEstimate.Probabilities, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecastEstimatesEstimate) GetProbabilities() []*ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbability {
	return v.Probabilities
}

// ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbability includes the requested fields of the GraphQL type Probability.
// The GraphQL type's documentation follows.
//
// A number between 0 and 100 tied to a specific Outcome. It is always part of
// an Estimate.
type ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbability struct {
	Id      string                                                                         `json:"id"`
	Outcome ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbabilityOutcome `json:"outcome"`
}

// GetId returns ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbability.Id, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbability) GetId() string {
	return v.Id
}

// GetOutcome returns ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbability.Outcome, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbability) GetOutcome() ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbabilityOutcome {
	return v.Outcome
}

// ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbabilityOutcome includes the requested fields of the GraphQL type Outcome.
// The GraphQL type's documentation follows.
//
// The possible results of a forecast. In the simplest case you will only have
// two outcomes: Yes and No.
type ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbabilityOutcome struct {
	Id      string `json:"id"`
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
}

// GetId returns ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbabilityOutcome.Id, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbabilityOutcome) GetId() string {
	return v.Id
}

// GetText returns ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbabilityOutcome.Text, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbabilityOutcome) GetText() string {
	return v.Text
}

// GetCorrect returns ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbabilityOutcome.Correct, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecastEstimatesEstimateProbabilitiesProbabilityOutcome) GetCorrect() bool {
	return v.Correct
}

// ResolveForecastResponse is returned by ResolveForecast on success.
type ResolveForecastResponse struct {
	ResolveForecast *ResolveForecastResolveForecast `json:"resolveForecast"`
}

// GetResolveForecast returns ResolveForecastResponse.ResolveForecast, and is useful for accessing the field via an interface.
func (v *ResolveForecastResponse) GetResolveForecast() *ResolveForecastResolveForecast {
	return v.ResolveForecast
}

// __ResolveForecastInput is used internally by genqlient
type __ResolveForecastInput struct {
	ForecastId       string      `json:"forecastId"`
	CorrectOutcomeId *string     `json:"correctOutcomeId"`
	Resolution       *Resolution `json:"resolution"`
}

// GetForecastId returns __ResolveForecastInput.ForecastId, and is useful for accessing the field via an interface.
func (v *__ResolveForecastInput) GetForecastId() string { return v.ForecastId }

// GetCorrectOutcomeId returns __ResolveForecastInput.CorrectOutcomeId, and is useful for accessing the field via an interface.
func (v *__ResolveForecastInput) GetCorrectOutcomeId() *string { return v.CorrectOutcomeId }

// GetResolution returns __ResolveForecastInput.Resolution, and is useful for accessing the field via an interface.
func (v *__ResolveForecastInput) GetResolution() *Resolution { return v.Resolution }

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
		title
		resolution
		estimates {
			id
			probabilities {
				id
				outcome {
					id
					text
					correct
				}
			}
		}
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

func ResolveForecast(
	ctx context.Context,
	client graphql.Client,
	forecastId string,
	correctOutcomeId *string,
	resolution *Resolution,
) (*ResolveForecastResponse, error) {
	req := &graphql.Request{
		OpName: "ResolveForecast",
		Query: `
mutation ResolveForecast ($forecastId: ID!, $correctOutcomeId: ID, $resolution: Resolution) {
	resolveForecast(forecastId: $forecastId, correctOutcomeId: $correctOutcomeId, resolution: $resolution) {
		id
		title
		resolution
		estimates {
			id
			probabilities {
				id
				outcome {
					id
					text
					correct
				}
			}
		}
	}
}
`,
		Variables: &__ResolveForecastInput{
			ForecastId:       forecastId,
			CorrectOutcomeId: correctOutcomeId,
			Resolution:       resolution,
		},
	}
	var err error

	var data ResolveForecastResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
