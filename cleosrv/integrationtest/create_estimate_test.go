package integrationtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateEstimate tests the happy case of adding a new Estimate to an existing Forecast.
func TestCreateEstimate(t *testing.T) {
	c := initServerAndGetClient(t, "")

	createForecastResponse, err := CreateForecast(
		context.Background(),
		c,
		NewForecast{
			Title:       "Test forecast",
			Description: "",
			Resolves:    time.Now().Add(24 * time.Hour),
			Closes:      nil,
			Created:     nil,
		},
		NewEstimate{
			Reason: "Just a hunch.",
			Probabilities: []NewProbability{
				{
					Value: 20,
					Outcome: &NewOutcome{
						Text: "Yes",
					},
				},
				{
					Value: 80,
					Outcome: &NewOutcome{
						Text: "No",
					},
				},
			},
			Created: nil,
		},
	)
	require.NoError(t, err)

	forecastId := createForecastResponse.CreateForecast.Id
	outcomeYesId := ""
	outcomeNoId := ""
	for _, p := range createForecastResponse.CreateForecast.Estimates[0].Probabilities {
		if p.Outcome.Text == "Yes" {
			outcomeYesId = p.Outcome.Id
		}
		if p.Outcome.Text == "No" {
			outcomeNoId = p.Outcome.Id
		}
	}
	require.NotEmpty(t, forecastId)
	require.NotEmpty(t, outcomeYesId)
	require.NotEmpty(t, outcomeNoId)

	newEstimate := NewEstimate{
		Reason: "I got some new information.",
		Probabilities: []NewProbability{
			{
				Value:     30,
				OutcomeId: &outcomeYesId,
			},
			{
				Value:     70,
				OutcomeId: &outcomeNoId,
			},
		},
		Created: nil,
	}

	response, err := CreateEstimate(
		context.Background(),
		c,
		forecastId,
		newEstimate,
	)
	require.NoError(t, err)

	assert.NotEmpty(t, response.CreateEstimate.Id)
}

// TestCreateEstimate_InvalidForecastID validates an error is returned when creating a new
// Estimate for an invalid Forecast ID.
func TestCreateEstimate_InvalidForecastID(t *testing.T) {
	tests := []struct {
		name        string
		forecastId  string
		expectedErr string
	}{
		{
			name:        "nonexistent ID",
			forecastId:  "3333",
			expectedErr: "error getting Forecast with ID 3333",
		},
		{
			name:        "non-numeric ID",
			forecastId:  "asdf",
			expectedErr: "error getting Forecast with ID asdf",
		},
		{
			name:        "empty ID",
			forecastId:  "",
			expectedErr: "error getting Forecast with ID :",
		},
	}
	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {

			c := initServerAndGetClient(t, "")

			createForecastResponse := simpleCreateForecastHelper(
				t,
				c,
				"Test Forecast",
				map[string]int{"Yes": 20, "No": 80},
			)

			yesId := ""
			noId := ""
			for _, p := range createForecastResponse.CreateForecast.Estimates[0].Probabilities {
				if p.Outcome.Text == "Yes" {
					yesId = p.Outcome.Id
				}
				if p.Outcome.Text == "No" {
					noId = p.Outcome.Id
				}
			}
			require.NotEmpty(t, yesId)
			require.NotEmpty(t, noId)

			newEstimate := NewEstimate{
				Reason: "I got some new information.",
				Probabilities: []NewProbability{
					{
						Value:     30,
						OutcomeId: &yesId,
					},
					{
						Value:     70,
						OutcomeId: &noId,
					},
				},
				Created: nil,
			}

			_, err := CreateEstimate(
				context.Background(),
				c,
				tt.forecastId,
				newEstimate,
			)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

// TestCreateEstimate_InvalidOutcomeID ensures that invalid Outcome IDs like non-existing or
// empty IDs are not accepted.
func TestCreateEstimate_InvalidOutcomeID(t *testing.T) {
	tests := []struct {
		name          string
		testOutcomeId *string
		expectedErr   string
	}{
		{
			name:          "non existing outcome ID",
			testOutcomeId: strPointer("4444"),
			expectedErr:   "invalid Outcome IDs: [4444]",
		},
		{
			name:          "nil outcome ID",
			testOutcomeId: nil,
			expectedErr:   "outcomeId must be set when adding an estimate",
		},
		{
			name:          "empty outcome ID",
			testOutcomeId: strPointer(""),
			expectedErr:   "outcomeId must not be empty when adding an estimate",
		},
	}
	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {
			c := initServerAndGetClient(t, "")

			createForecastResponse := simpleCreateForecastHelper(
				t,
				c,
				"Test Forecast",
				map[string]int{"Yes": 20, "No": 80},
			)

			forecastId := createForecastResponse.CreateForecast.Id
			yesId := ""
			noId := ""
			for _, p := range createForecastResponse.CreateForecast.Estimates[0].Probabilities {
				if p.Outcome.Text == "Yes" {
					yesId = p.Outcome.Id
				}
				if p.Outcome.Text == "No" {
					noId = p.Outcome.Id
				}
			}
			require.NotEmpty(t, forecastId)
			require.NotEmpty(t, yesId)
			require.NotEmpty(t, noId)

			newEstimate := NewEstimate{
				Reason: "I got some new information.",
				Probabilities: []NewProbability{
					{
						Value:     30,
						OutcomeId: tt.testOutcomeId, // NOT passing in yesId
					},
					{
						Value:     70,
						OutcomeId: &noId,
					},
				},
				Created: nil,
			}

			_, err := CreateEstimate(
				context.Background(),
				c,
				forecastId,
				newEstimate,
			)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

// TestCreateEstimate_MissingOutcomeID ensures that when creating a new Estimate all Outcomes
// that are part of the Forecast must be specified.
func TestCreateEstimate_MissingOutcomeID(t *testing.T) {
	c := initServerAndGetClient(t, "")

	createForecastResponse := simpleCreateForecastHelper(
		t,
		c,
		"Test Forecast",
		map[string]int{"Yes": 20, "No": 80},
	)

	forecastId := createForecastResponse.CreateForecast.Id
	yesId := ""
	noId := ""
	for _, p := range createForecastResponse.CreateForecast.Estimates[0].Probabilities {
		if p.Outcome.Text == "Yes" {
			yesId = p.Outcome.Id
		}
		if p.Outcome.Text == "No" {
			noId = p.Outcome.Id
		}
	}
	require.NotEmpty(t, forecastId)
	require.NotEmpty(t, yesId)
	require.NotEmpty(t, noId)

	newEstimate := NewEstimate{
		Reason: "I got some new information.",
		Probabilities: []NewProbability{
			{
				Value:     100,
				OutcomeId: &noId,
			},
		},
		Created: nil,
	}

	_, err := CreateEstimate(
		context.Background(),
		c,
		forecastId,
		newEstimate,
	)
	assert.ErrorContains(t, err, fmt.Sprintf("missing Outcome IDs: [%v]", yesId))
}

// TestCreateEstimate_WithNewOutcome ensures that when creating a new Estimate no
// NewOutcome can be passed in (only outcome IDs from existing Outcomes).
func TestCreateEstimate_WithNewOutcome(t *testing.T) {
	tests := []struct {
		name                  string
		newOutcome            *NewOutcome
		includeValidOutcomeId bool
		expectedErr           string
	}{
		{
			name:                  "no NewOutcome - happy case",
			newOutcome:            nil,
			includeValidOutcomeId: true,
			expectedErr:           "", // no error
		},
		{
			name:                  "NewOutcome and OutcomeId are set",
			newOutcome:            &NewOutcome{Text: "Maybe"},
			includeValidOutcomeId: true,
			expectedErr:           "NewOutcome must be unset when adding an estimate",
		},
		{
			name:                  "Only NewOutcome is set",
			newOutcome:            &NewOutcome{Text: "Maybe"},
			includeValidOutcomeId: false,
			expectedErr:           "NewOutcome must be unset when adding an estimate",
		},
	}
	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {
			c := initServerAndGetClient(t, "")

			createForecastResponse := simpleCreateForecastHelper(
				t,
				c,
				"Test Forecast",
				map[string]int{"Yes": 20, "No": 80},
			)

			forecastId := createForecastResponse.CreateForecast.Id
			yesId := ""
			noId := ""
			for _, p := range createForecastResponse.CreateForecast.Estimates[0].Probabilities {
				if p.Outcome.Text == "Yes" {
					yesId = p.Outcome.Id
				}
				if p.Outcome.Text == "No" {
					noId = p.Outcome.Id
				}
			}
			require.NotEmpty(t, forecastId)
			require.NotEmpty(t, yesId)
			require.NotEmpty(t, noId)

			newEstimate := NewEstimate{
				Reason: "I got some new information.",
				Probabilities: []NewProbability{
					{
						Value: 30,
						// OutcomeId is set below because it's conditional
						Outcome: tt.newOutcome,
					},
					{
						Value:     70,
						OutcomeId: &noId,
					},
				},
				Created: nil,
			}
			if tt.includeValidOutcomeId {
				newEstimate.Probabilities[0].OutcomeId = &yesId
			}

			_, err := CreateEstimate(
				context.Background(),
				c,
				forecastId,
				newEstimate,
			)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.expectedErr)
			}
		})
	}
}

// TestCreateEstimate_ErrorOutcomeIDFromAnotherForecast ensures that valid Outcome IDs that
// belong to a different Forecast are not accepted.
func TestCreateEstimate_ErrorOutcomeIDFromAnotherForecast(t *testing.T) {
	c := initServerAndGetClient(t, "")

	createForecastResponse := simpleCreateForecastHelper(
		t,
		c,
		"Test Forecast",
		map[string]int{"Yes": 20, "No": 80},
	)

	forecastId := createForecastResponse.CreateForecast.Id
	yesId := ""
	noId := ""
	for _, p := range createForecastResponse.CreateForecast.Estimates[0].Probabilities {
		if p.Outcome.Text == "Yes" {
			yesId = p.Outcome.Id
		}
		if p.Outcome.Text == "No" {
			noId = p.Outcome.Id
		}
	}
	require.NotEmpty(t, forecastId)
	require.NotEmpty(t, yesId)
	require.NotEmpty(t, noId)

	createForecastResponse2 := simpleCreateForecastHelper(
		t,
		c,
		"Test Forecast 2",
		map[string]int{"Yes": 20, "No": 80},
	)
	yesId2 := ""
	noId2 := ""
	for _, p := range createForecastResponse2.CreateForecast.Estimates[0].Probabilities {
		if p.Outcome.Text == "Yes" {
			yesId2 = p.Outcome.Id
		}
		if p.Outcome.Text == "No" {
			noId2 = p.Outcome.Id
		}
	}
	require.NotEmpty(t, yesId2)
	require.NotEmpty(t, noId2)

	newEstimate := NewEstimate{
		Reason: "I got some new information.",
		Probabilities: []NewProbability{
			{
				Value:     30,
				OutcomeId: &yesId,
			},
			{
				Value:     30,
				OutcomeId: &noId,
			},
			{
				Value:     40,
				OutcomeId: &noId2, // Outcome ID from another forecast
			},
		},
		Created: nil,
	}

	_, err := CreateEstimate(
		context.Background(),
		c,
		forecastId,
		newEstimate,
	)
	assert.ErrorContains(t, err, fmt.Sprintf("invalid Outcome IDs: [%v]", noId2))
}

// TestCreateEstimate_ProbabilitiesMustAdd100 ensures that the probabilities in a new Estimate
// add up to 100%.
func TestCreateEstimate_ProbabilitiesMustAdd100(t *testing.T) {
	c := initServerAndGetClient(t, "")

	createForecastResponse := simpleCreateForecastHelper(
		t,
		c,
		"Test Forecast",
		map[string]int{"Yes": 20, "No": 80},
	)

	forecastId := createForecastResponse.CreateForecast.Id
	yesId := ""
	noId := ""
	for _, p := range createForecastResponse.CreateForecast.Estimates[0].Probabilities {
		if p.Outcome.Text == "Yes" {
			yesId = p.Outcome.Id
		}
		if p.Outcome.Text == "No" {
			noId = p.Outcome.Id
		}
	}
	require.NotEmpty(t, forecastId)
	require.NotEmpty(t, yesId)
	require.NotEmpty(t, noId)

	newEstimate := NewEstimate{
		Reason: "I got some new information.",
		Probabilities: []NewProbability{
			{
				Value:     10,
				OutcomeId: &yesId,
			},
			{
				Value:     10,
				OutcomeId: &noId,
			},
		},
		Created: nil,
	}

	_, err := CreateEstimate(
		context.Background(),
		c,
		forecastId,
		newEstimate,
	)
	assert.ErrorContains(t, err, "probabilities must add up to 100")
}

// TestCreateEstimate_ReasonCantBeEmpty ensures that the Estimate.reason is not left empty.
func TestCreateEstimate_ReasonCantBeEmpty(t *testing.T) {
	c := initServerAndGetClient(t, "")

	createForecastResponse := simpleCreateForecastHelper(
		t,
		c,
		"Test Forecast",
		map[string]int{"Yes": 20, "No": 80},
	)

	forecastId := createForecastResponse.CreateForecast.Id
	yesId := ""
	noId := ""
	for _, p := range createForecastResponse.CreateForecast.Estimates[0].Probabilities {
		if p.Outcome.Text == "Yes" {
			yesId = p.Outcome.Id
		}
		if p.Outcome.Text == "No" {
			noId = p.Outcome.Id
		}
	}
	require.NotEmpty(t, forecastId)
	require.NotEmpty(t, yesId)
	require.NotEmpty(t, noId)

	newEstimate := NewEstimate{
		Reason: "",
		Probabilities: []NewProbability{
			{
				Value:     10,
				OutcomeId: &yesId,
			},
			{
				Value:     90,
				OutcomeId: &noId,
			},
		},
		Created: nil,
	}

	_, err := CreateEstimate(
		context.Background(),
		c,
		forecastId,
		newEstimate,
	)
	assert.ErrorContains(t, err, "'reason' can't be empty")
}

// TestCreateEstimate_VerifyCreatedNilIsNow ensures that Estimate.
// created == nil is interpreted as now().
func TestCreateEstimate_VerifyCreatedNilIsNow(t *testing.T) {
	now := time.Now().UTC()

	c := initServerAndGetClient(t, "")

	createForecastResponse := simpleCreateForecastHelper(
		t,
		c,
		"Test Forecast",
		map[string]int{"Yes": 20, "No": 80},
	)

	forecastId := createForecastResponse.CreateForecast.Id
	yesId := ""
	noId := ""
	for _, p := range createForecastResponse.CreateForecast.Estimates[0].Probabilities {
		if p.Outcome.Text == "Yes" {
			yesId = p.Outcome.Id
		}
		if p.Outcome.Text == "No" {
			noId = p.Outcome.Id
		}
	}
	require.NotEmpty(t, forecastId)
	require.NotEmpty(t, yesId)
	require.NotEmpty(t, noId)

	secondEstimate := NewEstimate{
		Reason: "I got some new information.",
		Probabilities: []NewProbability{
			{
				Value:     10,
				OutcomeId: &yesId,
			},
			{
				Value:     90,
				OutcomeId: &noId,
			},
		},
		Created: nil,
	}

	resp, err := CreateEstimate(
		context.Background(),
		c,
		forecastId,
		secondEstimate,
	)

	require.NoError(t, err)
	assertTimeAlmostEqual(t, now, resp.CreateEstimate.Created)
}

// TestCreateEstimate_VerifyCreatedDateBehavior verifies which Estimate.
// created values are valid and which are not in different circumstances.
// The reason why it is allowed to set Estimate.
// created in the past even for resolved forecasts is for importing data.
// It's convenient to be able to directly mark the Forecast as resolved and only afterwards
// import all the Estimates.
// It is debatable whether this is a good idea or whether it should not be possible to add
// Estimates to resolved Forecasts.
// This tests exists in any case to document the current behavior (yes,
// it is possible to add Estimates to resolved Forecasts).
func TestCreateEstimate_VerifyCreatedDateBehavior(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name                      string
		resolution                Resolution
		forecastCreatedDate       *time.Time
		forecastResolvesDate      time.Time
		forecastClosesDate        *time.Time
		secondEstimateCreatedDate *time.Time
		expectedErr               string
	}{
		{
			name:                      "(unresolved) default successful case",
			resolution:                ResolutionUnresolved,
			forecastCreatedDate:       nil,
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: nil,
			expectedErr:               "",
		},
		{
			name:                      "(unresolved) error if estimate.created is in the future",
			resolution:                ResolutionUnresolved,
			forecastCreatedDate:       nil,
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(1 * time.Hour)),
			expectedErr:               "'created' can't be in the future",
		},
		{
			name:                      "(resolved) error if estimate.created is in the future",
			resolution:                ResolutionResolved,
			forecastCreatedDate:       nil,
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(1 * time.Hour)),
			expectedErr:               "'created' can't be in the future",
		},
		{
			name:                      "(NA resolved) error if estimate.created is in the future",
			resolution:                ResolutionNotApplicable,
			forecastCreatedDate:       nil,
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(1 * time.Hour)),
			expectedErr:               "'created' can't be in the future",
		},
		{
			// This test will lead to an error because when resolving the forecast the resolves
			// date is set to now(). When adding another estimate without date, now() will be
			// used again, but it will then lead to a later timestamp than the first now().
			name:                "(resolved) error if estimate has no created date",
			resolution:          ResolutionResolved,
			forecastCreatedDate: nil,
			// note that if the resolves date is in the future it is set to now() when resolving
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: nil,
			expectedErr: "'estimate.created' is set to a later date than 'forecast." +
				"resolves'",
		},
		{
			// See the previous test for an explanation
			name:                      "(NA resolved) error if estimate.created is nil",
			resolution:                ResolutionNotApplicable,
			forecastCreatedDate:       nil,
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: nil,
			expectedErr: "'estimate.created' is set to a later date than 'forecast." +
				"resolves'",
		},
		{
			name:                      "(unresolved) success if estimate.created date is in the past",
			resolution:                ResolutionUnresolved,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(-23 * time.Hour)),
			expectedErr:               "",
		},
		{
			name:                      "(resolved) success if estimate.created date is in the past",
			resolution:                ResolutionResolved,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(-23 * time.Hour)),
			expectedErr:               "",
		},
		{
			name:                      "(NA resolved) success if estimate.created date is in the past",
			resolution:                ResolutionNotApplicable,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(-23 * time.Hour)),
			expectedErr:               "",
		},
		{
			name: "(unresolved) error if estimate.created is before forecast." +
				"created",
			resolution:                ResolutionUnresolved,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(-25 * time.Hour)),
			expectedErr: "'estimate.created' is set to an earlier date than 'forecast." +
				"created'",
		},
		{
			name: "(resolved) error if estimate.created is before forecast." +
				"created",
			resolution:                ResolutionResolved,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(-25 * time.Hour)),
			expectedErr: "'estimate.created' is set to an earlier date than 'forecast." +
				"created'",
		},
		{
			name: "(NA resolved) error if estimate.created is before " +
				"forecast.created",
			resolution:                ResolutionNotApplicable,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(24 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(-25 * time.Hour)),
			expectedErr: "'estimate.created' is set to an earlier date than 'forecast." +
				"created'",
		},
		{
			name: "(unresolved) success if estimate.created is after forecast." +
				"resolves, which is in the past",
			resolution:                ResolutionUnresolved,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(-23 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(-22 * time.Hour)),
			expectedErr:               "",
		},
		{
			name: "(resolved) error if estimate.created is after forecast." +
				"resolves, which is in the past",
			resolution:                ResolutionResolved,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(-23 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(-22 * time.Hour)),
			expectedErr: "'estimate.created' is set to a later date than 'forecast." +
				"resolves'",
		},
		{
			name: "(NA resolved) error if estimate.created is after forecast." +
				"resolves which is in the past",
			resolution:                ResolutionNotApplicable,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(-23 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(-22 * time.Hour)),
			expectedErr: "'estimate.created' is set to a later date than 'forecast." +
				"resolves'",
		},
		{
			name: "(resolved) success if estimate.created is equal to " +
				"forecast.resolves",
			resolution:                ResolutionResolved,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(-23 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(-23 * time.Hour)),
			expectedErr:               "",
		},
		{
			name: "(NA resolved) success if estimate.created is equal to " +
				"forecast.resolves",
			resolution:                ResolutionNotApplicable,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(-23 * time.Hour),
			forecastClosesDate:        nil,
			secondEstimateCreatedDate: timePointer(now.Add(-23 * time.Hour)),
			expectedErr:               "",
		},
		{
			name: "(unresolved) error if estimate.created is after " +
				"forecast.closes",
			resolution:                ResolutionUnresolved,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(-20 * time.Hour),
			forecastClosesDate:        timePointer(now.Add(-22 * time.Hour)),
			secondEstimateCreatedDate: timePointer(now.Add(-21 * time.Hour)),
			expectedErr:               "'estimate.created' is set to a later date than 'forecast.closes'",
		},
		{
			name: "(resolved) error if estimate.created is after " +
				"forecast.closes",
			resolution:                ResolutionResolved,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(-20 * time.Hour),
			forecastClosesDate:        timePointer(now.Add(-22 * time.Hour)),
			secondEstimateCreatedDate: timePointer(now.Add(-21 * time.Hour)),
			expectedErr:               "'estimate.created' is set to a later date than 'forecast.closes'",
		},
		{
			name: "(NA resolved) error if estimate.created is after " +
				"forecast.closes",
			resolution:                ResolutionNotApplicable,
			forecastCreatedDate:       timePointer(now.Add(-24 * time.Hour)),
			forecastResolvesDate:      now.Add(-20 * time.Hour),
			forecastClosesDate:        timePointer(now.Add(-22 * time.Hour)),
			secondEstimateCreatedDate: timePointer(now.Add(-21 * time.Hour)),
			expectedErr:               "'estimate.created' is set to a later date than 'forecast.closes'",
		},
	}
	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {
			c := initServerAndGetClient(t, "")

			// 1. Create a new forecast

			newForecast := NewForecast{
				Title:       "Test Forecast",
				Description: "",
				Created:     tt.forecastCreatedDate,
				Resolves:    tt.forecastResolvesDate,
				Closes:      tt.forecastClosesDate,
			}

			initialEstimate := NewEstimate{
				Reason: "Just a hunch.",
				Probabilities: []NewProbability{
					{
						Value:   20,
						Outcome: &NewOutcome{Text: "Yes"},
					},
					{
						Value:   80,
						Outcome: &NewOutcome{Text: "No"},
					},
				},
			}
			createForecastResponse, err := CreateForecast(
				context.Background(),
				c,
				newForecast,
				initialEstimate,
			)
			require.NoError(t, err)

			forecastId := createForecastResponse.CreateForecast.Id
			yesId := ""
			noId := ""
			for _, p := range createForecastResponse.CreateForecast.Estimates[0].Probabilities {
				if p.Outcome.Text == "Yes" {
					yesId = p.Outcome.Id
				}
				if p.Outcome.Text == "No" {
					noId = p.Outcome.Id
				}
			}
			require.NotEmpty(t, forecastId)
			require.NotEmpty(t, yesId)
			require.NotEmpty(t, noId)

			// 2. Resolve the forecast (if necessary)

			if tt.resolution == ResolutionResolved {
				_, err := ResolveForecast(
					context.Background(),
					c,
					forecastId,
					&yesId,
					&tt.resolution,
				)
				require.NoError(t, err)
			} else if tt.resolution == ResolutionNotApplicable {
				_, err := ResolveForecast(
					context.Background(),
					c,
					forecastId,
					nil,
					&tt.resolution,
				)
				require.NoError(t, err)
			}

			// 3. Add a new estimate to the forecast

			secondEstimate := NewEstimate{
				Reason: "I got some new information.",
				Probabilities: []NewProbability{
					{
						Value:     10,
						OutcomeId: &yesId,
					},
					{
						Value:     90,
						OutcomeId: &noId,
					},
				},
				Created: tt.secondEstimateCreatedDate,
			}

			_, err = CreateEstimate(
				context.Background(),
				c,
				forecastId,
				secondEstimate,
			)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.expectedErr)
			}

			// 4. Done
		})
	}
}
