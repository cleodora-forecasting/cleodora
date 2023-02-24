// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package integrationtest

import (
	"context"
	"time"

	"github.com/Khan/genqlient/graphql"
)

// CreateForecastCreateForecast includes the requested fields of the GraphQL type Forecast.
// The GraphQL type's documentation follows.
//
// A prediction about the future.
type CreateForecastCreateForecast struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	// The point in time at which you predict you will be able to resolve whether
	// how the forecast resolved.
	Resolves time.Time `json:"resolves"`
	// The point in time at which you no longer want to update your probability
	// estimates for the forecast. In most cases you won't need this. One example
	// where you might is when you want to predict the outcome of an exam. You may
	// want to set 'closes' to the time right before the exam starts, even though
	// 'resolves' is several weeks later (when the exam results are published). This
	// way your prediction history will only reflect your estimations before you
	// took the exam, which is something you may want (or not, in which case you
	// could simply not set 'closes').
	Closes    *time.Time                                       `json:"closes"`
	Estimates []*CreateForecastCreateForecastEstimatesEstimate `json:"estimates"`
}

// GetId returns CreateForecastCreateForecast.Id, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecast) GetId() string { return v.Id }

// GetTitle returns CreateForecastCreateForecast.Title, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecast) GetTitle() string { return v.Title }

// GetDescription returns CreateForecastCreateForecast.Description, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecast) GetDescription() string { return v.Description }

// GetCreated returns CreateForecastCreateForecast.Created, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecast) GetCreated() time.Time { return v.Created }

// GetResolves returns CreateForecastCreateForecast.Resolves, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecast) GetResolves() time.Time { return v.Resolves }

// GetCloses returns CreateForecastCreateForecast.Closes, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecast) GetCloses() *time.Time { return v.Closes }

// GetEstimates returns CreateForecastCreateForecast.Estimates, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecast) GetEstimates() []*CreateForecastCreateForecastEstimatesEstimate {
	return v.Estimates
}

// CreateForecastCreateForecastEstimatesEstimate includes the requested fields of the GraphQL type Estimate.
// The GraphQL type's documentation follows.
//
// A list of probabilities (one for each outcome) together with a timestamp and
// an explanation why you made this estimate. Every time you change your mind
// about a forecast you will create a new Estimate.
// All probabilities always add up to 100.
type CreateForecastCreateForecastEstimatesEstimate struct {
	Id            string                                                                   `json:"id"`
	Created       time.Time                                                                `json:"created"`
	Reason        string                                                                   `json:"reason"`
	Probabilities []*CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbability `json:"probabilities"`
}

// GetId returns CreateForecastCreateForecastEstimatesEstimate.Id, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecastEstimatesEstimate) GetId() string { return v.Id }

// GetCreated returns CreateForecastCreateForecastEstimatesEstimate.Created, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecastEstimatesEstimate) GetCreated() time.Time { return v.Created }

// GetReason returns CreateForecastCreateForecastEstimatesEstimate.Reason, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecastEstimatesEstimate) GetReason() string { return v.Reason }

// GetProbabilities returns CreateForecastCreateForecastEstimatesEstimate.Probabilities, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecastEstimatesEstimate) GetProbabilities() []*CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbability {
	return v.Probabilities
}

// CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbability includes the requested fields of the GraphQL type Probability.
// The GraphQL type's documentation follows.
//
// A number between 0 and 100 tied to a specific Outcome. It is always part of
// an Estimate.
type CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbability struct {
	Id      string                                                                       `json:"id"`
	Value   int                                                                          `json:"value"`
	Outcome CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbabilityOutcome `json:"outcome"`
}

// GetId returns CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbability.Id, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbability) GetId() string {
	return v.Id
}

// GetValue returns CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbability.Value, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbability) GetValue() int {
	return v.Value
}

// GetOutcome returns CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbability.Outcome, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbability) GetOutcome() CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbabilityOutcome {
	return v.Outcome
}

// CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbabilityOutcome includes the requested fields of the GraphQL type Outcome.
// The GraphQL type's documentation follows.
//
// The possible results of a forecast. In the simplest case you will only have
// two outcomes: Yes and No.
type CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbabilityOutcome struct {
	Id      string `json:"id"`
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
}

// GetId returns CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbabilityOutcome.Id, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbabilityOutcome) GetId() string {
	return v.Id
}

// GetText returns CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbabilityOutcome.Text, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbabilityOutcome) GetText() string {
	return v.Text
}

// GetCorrect returns CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbabilityOutcome.Correct, and is useful for accessing the field via an interface.
func (v *CreateForecastCreateForecastEstimatesEstimateProbabilitiesProbabilityOutcome) GetCorrect() bool {
	return v.Correct
}

// CreateForecastResponse is returned by CreateForecast on success.
type CreateForecastResponse struct {
	CreateForecast CreateForecastCreateForecast `json:"createForecast"`
}

// GetCreateForecast returns CreateForecastResponse.CreateForecast, and is useful for accessing the field via an interface.
func (v *CreateForecastResponse) GetCreateForecast() CreateForecastCreateForecast {
	return v.CreateForecast
}

// GetForecastsForecastsForecast includes the requested fields of the GraphQL type Forecast.
// The GraphQL type's documentation follows.
//
// A prediction about the future.
type GetForecastsForecastsForecast struct {
	Id         string     `json:"id"`
	Created    time.Time  `json:"created"`
	Title      string     `json:"title"`
	Resolution Resolution `json:"resolution"`
	// The point in time at which you predict you will be able to resolve whether
	// how the forecast resolved.
	Resolves time.Time `json:"resolves"`
	// The point in time at which you no longer want to update your probability
	// estimates for the forecast. In most cases you won't need this. One example
	// where you might is when you want to predict the outcome of an exam. You may
	// want to set 'closes' to the time right before the exam starts, even though
	// 'resolves' is several weeks later (when the exam results are published). This
	// way your prediction history will only reflect your estimations before you
	// took the exam, which is something you may want (or not, in which case you
	// could simply not set 'closes').
	Closes    *time.Time                                        `json:"closes"`
	Estimates []*GetForecastsForecastsForecastEstimatesEstimate `json:"estimates"`
}

// GetId returns GetForecastsForecastsForecast.Id, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetId() string { return v.Id }

// GetCreated returns GetForecastsForecastsForecast.Created, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetCreated() time.Time { return v.Created }

// GetTitle returns GetForecastsForecastsForecast.Title, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetTitle() string { return v.Title }

// GetResolution returns GetForecastsForecastsForecast.Resolution, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetResolution() Resolution { return v.Resolution }

// GetResolves returns GetForecastsForecastsForecast.Resolves, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetResolves() time.Time { return v.Resolves }

// GetCloses returns GetForecastsForecastsForecast.Closes, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecast) GetCloses() *time.Time { return v.Closes }

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
	Created       time.Time                                                                 `json:"created"`
	Probabilities []*GetForecastsForecastsForecastEstimatesEstimateProbabilitiesProbability `json:"probabilities"`
}

// GetId returns GetForecastsForecastsForecastEstimatesEstimate.Id, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecastEstimatesEstimate) GetId() string { return v.Id }

// GetCreated returns GetForecastsForecastsForecastEstimatesEstimate.Created, and is useful for accessing the field via an interface.
func (v *GetForecastsForecastsForecastEstimatesEstimate) GetCreated() time.Time { return v.Created }

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

// GetMetadataMetadata includes the requested fields of the GraphQL type Metadata.
// The GraphQL type's documentation follows.
//
// Information about the application itself e.g. the current version.
type GetMetadataMetadata struct {
	Version string `json:"version"`
}

// GetVersion returns GetMetadataMetadata.Version, and is useful for accessing the field via an interface.
func (v *GetMetadataMetadata) GetVersion() string { return v.Version }

// GetMetadataResponse is returned by GetMetadata on success.
type GetMetadataResponse struct {
	Metadata GetMetadataMetadata `json:"metadata"`
}

// GetMetadata returns GetMetadataResponse.Metadata, and is useful for accessing the field via an interface.
func (v *GetMetadataResponse) GetMetadata() GetMetadataMetadata { return v.Metadata }

type NewEstimate struct {
	Reason        string           `json:"reason"`
	Probabilities []NewProbability `json:"probabilities"`
	// An optional date in the past when you created this estimate. This can be
	// useful for cases when you wrote it down on a piece of paper or when importing
	// from other software. When creating a new Forecast this value will be for
	// the first Estimate (which will get the same timestamp as the
	// Forecast.Created).
	Created *time.Time `json:"created"`
}

// GetReason returns NewEstimate.Reason, and is useful for accessing the field via an interface.
func (v *NewEstimate) GetReason() string { return v.Reason }

// GetProbabilities returns NewEstimate.Probabilities, and is useful for accessing the field via an interface.
func (v *NewEstimate) GetProbabilities() []NewProbability { return v.Probabilities }

// GetCreated returns NewEstimate.Created, and is useful for accessing the field via an interface.
func (v *NewEstimate) GetCreated() *time.Time { return v.Created }

type NewForecast struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Resolves    time.Time  `json:"resolves"`
	Closes      *time.Time `json:"closes"`
	// An optional date in the past when you created this forecast. This can be
	// useful for cases when you wrote it down on a piece of paper or when importing
	// from other software.
	Created *time.Time `json:"created"`
}

// GetTitle returns NewForecast.Title, and is useful for accessing the field via an interface.
func (v *NewForecast) GetTitle() string { return v.Title }

// GetDescription returns NewForecast.Description, and is useful for accessing the field via an interface.
func (v *NewForecast) GetDescription() string { return v.Description }

// GetResolves returns NewForecast.Resolves, and is useful for accessing the field via an interface.
func (v *NewForecast) GetResolves() time.Time { return v.Resolves }

// GetCloses returns NewForecast.Closes, and is useful for accessing the field via an interface.
func (v *NewForecast) GetCloses() *time.Time { return v.Closes }

// GetCreated returns NewForecast.Created, and is useful for accessing the field via an interface.
func (v *NewForecast) GetCreated() *time.Time { return v.Created }

type NewOutcome struct {
	Text string `json:"text"`
}

// GetText returns NewOutcome.Text, and is useful for accessing the field via an interface.
func (v *NewOutcome) GetText() string { return v.Text }

type NewProbability struct {
	Value   int        `json:"value"`
	Outcome NewOutcome `json:"outcome"`
}

// GetValue returns NewProbability.Value, and is useful for accessing the field via an interface.
func (v *NewProbability) GetValue() int { return v.Value }

// GetOutcome returns NewProbability.Outcome, and is useful for accessing the field via an interface.
func (v *NewProbability) GetOutcome() NewOutcome { return v.Outcome }

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
	Id         string     `json:"id"`
	Title      string     `json:"title"`
	Resolution Resolution `json:"resolution"`
	// The point in time at which you predict you will be able to resolve whether
	// how the forecast resolved.
	Resolves time.Time `json:"resolves"`
	// The point in time at which you no longer want to update your probability
	// estimates for the forecast. In most cases you won't need this. One example
	// where you might is when you want to predict the outcome of an exam. You may
	// want to set 'closes' to the time right before the exam starts, even though
	// 'resolves' is several weeks later (when the exam results are published). This
	// way your prediction history will only reflect your estimations before you
	// took the exam, which is something you may want (or not, in which case you
	// could simply not set 'closes').
	Closes    *time.Time                                         `json:"closes"`
	Estimates []*ResolveForecastResolveForecastEstimatesEstimate `json:"estimates"`
}

// GetId returns ResolveForecastResolveForecast.Id, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecast) GetId() string { return v.Id }

// GetTitle returns ResolveForecastResolveForecast.Title, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecast) GetTitle() string { return v.Title }

// GetResolution returns ResolveForecastResolveForecast.Resolution, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecast) GetResolution() Resolution { return v.Resolution }

// GetResolves returns ResolveForecastResolveForecast.Resolves, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecast) GetResolves() time.Time { return v.Resolves }

// GetCloses returns ResolveForecastResolveForecast.Closes, and is useful for accessing the field via an interface.
func (v *ResolveForecastResolveForecast) GetCloses() *time.Time { return v.Closes }

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

// __CreateForecastInput is used internally by genqlient
type __CreateForecastInput struct {
	Forecast NewForecast `json:"forecast"`
	Estimate NewEstimate `json:"estimate"`
}

// GetForecast returns __CreateForecastInput.Forecast, and is useful for accessing the field via an interface.
func (v *__CreateForecastInput) GetForecast() NewForecast { return v.Forecast }

// GetEstimate returns __CreateForecastInput.Estimate, and is useful for accessing the field via an interface.
func (v *__CreateForecastInput) GetEstimate() NewEstimate { return v.Estimate }

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

func CreateForecast(
	ctx context.Context,
	client graphql.Client,
	forecast NewForecast,
	estimate NewEstimate,
) (*CreateForecastResponse, error) {
	req := &graphql.Request{
		OpName: "CreateForecast",
		Query: `
mutation CreateForecast ($forecast: NewForecast!, $estimate: NewEstimate!) {
	createForecast(forecast: $forecast, estimate: $estimate) {
		id
		title
		description
		created
		resolves
		closes
		estimates {
			id
			created
			reason
			probabilities {
				id
				value
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
		Variables: &__CreateForecastInput{
			Forecast: forecast,
			Estimate: estimate,
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
		created
		title
		resolution
		resolves
		closes
		estimates {
			id
			created
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

func GetMetadata(
	ctx context.Context,
	client graphql.Client,
) (*GetMetadataResponse, error) {
	req := &graphql.Request{
		OpName: "GetMetadata",
		Query: `
query GetMetadata {
	metadata {
		version
	}
}
`,
	}
	var err error

	var data GetMetadataResponse
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
		resolves
		closes
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
