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

// TestCreateEstimate_VerifyCreated ensures that the Estimate.
// created date is interpreted as expected.
func TestCreateEstimate_VerifyCreated(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name            string
		inputCreated    *time.Time
		expectedCreated *time.Time
		expectedErr     string
	}{
		{
			name:            "default created is now",
			inputCreated:    nil,
			expectedCreated: &now,
			expectedErr:     "",
		},
		{
			name:            "created can't be in the future",
			inputCreated:    timePointer(now.Add(24 * time.Hour)),
			expectedCreated: nil,
			expectedErr:     "'created' can't be in the future",
		},
		{
			name:            "created can be in the past",
			inputCreated:    timePointer(now.Add(-24 * time.Hour)),
			expectedCreated: timePointer(now.Add(-24 * time.Hour)),
			expectedErr:     "",
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
						Value:     10,
						OutcomeId: &yesId,
					},
					{
						Value:     90,
						OutcomeId: &noId,
					},
				},
				Created: tt.inputCreated,
			}

			resp, err := CreateEstimate(
				context.Background(),
				c,
				forecastId,
				newEstimate,
			)
			if tt.expectedErr == "" {
				require.NoError(t, err)
				assertTimeAlmostEqual(t, *tt.expectedCreated, resp.CreateEstimate.Created)
			} else {
				require.ErrorContains(t, err, tt.expectedErr)
			}
		})
	}
}
