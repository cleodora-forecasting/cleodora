package integrationtest

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResolveForecast tests the happy case and some correct or incorrect
// input combinations.
func TestResolveForecast(t *testing.T) {
	tests := []struct {
		name                       string
		inputResolution            *Resolution
		includeInputOutcomeId      bool
		expectedErr                string
		expectedResolution         Resolution
		expectedOutcomeCorrectness bool
	}{
		{
			name:                       "happy case - set RESOLVED",
			inputResolution:            resolutionPointer(ResolutionResolved),
			includeInputOutcomeId:      true,
			expectedErr:                "",
			expectedResolution:         ResolutionResolved,
			expectedOutcomeCorrectness: true,
		},
		{
			name:                       "with resolution RESOLVED the resolution can be omitted",
			inputResolution:            nil,
			includeInputOutcomeId:      true,
			expectedErr:                "",
			expectedResolution:         ResolutionResolved,
			expectedOutcomeCorrectness: true,
		},
		{
			name:                       "resolution RESOLVED is only allowed when passing the OutcomeId",
			inputResolution:            resolutionPointer(ResolutionResolved),
			includeInputOutcomeId:      false,
			expectedErr:                "Outcome must be specified",
			expectedResolution:         ResolutionUnresolved,
			expectedOutcomeCorrectness: false,
		},
		{
			name:                       "can apply NOT_APPLICABLE resolution",
			inputResolution:            resolutionPointer(ResolutionNotApplicable),
			includeInputOutcomeId:      false,
			expectedErr:                "",
			expectedResolution:         ResolutionNotApplicable,
			expectedOutcomeCorrectness: false,
		},
		{
			name:                       "can apply NOT_APPLICABLE resolution even with OutcomeId",
			inputResolution:            resolutionPointer(ResolutionNotApplicable),
			includeInputOutcomeId:      true,
			expectedErr:                "",
			expectedResolution:         ResolutionNotApplicable,
			expectedOutcomeCorrectness: false, // note the Outcome must not change
		},
		{
			name:                       "applying UNRESOLVED resolution is not allowed",
			inputResolution:            resolutionPointer(ResolutionUnresolved),
			includeInputOutcomeId:      true,
			expectedErr:                "not allowed",
			expectedResolution:         ResolutionUnresolved,
			expectedOutcomeCorrectness: false,
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {
			c := initServerAndGetClient2(t)

			queryResponse, err := GetForecasts(context.Background(), c)
			require.NoError(t, err)

			fmt.Println(queryResponse)

			var fabelmansForecastId string
			var yesOutcomeId string

			for _, f := range queryResponse.Forecasts {
				if strings.Contains(f.Title, "\"The Fabelmans\"") {
					fabelmansForecastId = f.Id
					require.Equal(t, ResolutionUnresolved, f.Resolution)
					require.NotEmpty(t, f.Estimates)
					for _, p := range f.Estimates[0].Probabilities {
						require.False(t, p.Outcome.Correct)
						if p.Outcome.Text == "Yes" {
							yesOutcomeId = p.Outcome.Id
							break
						}
					}
					break
				}
			}
			require.NotEmpty(t, fabelmansForecastId)
			require.NotEmpty(t, yesOutcomeId)

			var inputOutcomeId *string
			if tt.includeInputOutcomeId {
				inputOutcomeId = &yesOutcomeId
			}

			_, err = ResolveForecast(
				context.Background(),
				c,
				fabelmansForecastId,
				inputOutcomeId,
				tt.inputResolution,
			)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.expectedErr)
			}

			// Verify what has changed
			queryResponse, err = GetForecasts(context.Background(), c)
			require.NoError(t, err)

			for _, f := range queryResponse.Forecasts {
				if strings.Contains(f.Title, "\"The Fabelmans\"") {
					assert.Equal(t, tt.expectedResolution, f.Resolution)
					require.NotEmpty(t, f.Estimates)
					for _, p := range f.Estimates[0].Probabilities {
						if p.Outcome.Text == "Yes" {
							assert.Equal(t, tt.expectedOutcomeCorrectness, p.Outcome.Correct)
						} else {
							assert.False(
								t,
								p.Outcome.Correct,
								"all other outcomes must be false",
							)
						}
					}
					break
				}
			}
		})
	}
}

// TestResolveForecast_VerifyResponseValue tests the resolve forecast happy
// case and verifies the response contains the expected results. It resolves
// the forecast and set one of the Outcomes as correct.
// It also verifies that the same forecast can't be resolved again.
func TestResolveForecast_VerifyResponseValue(t *testing.T) {
	c := initServerAndGetClient2(t)

	queryResponse, err := GetForecasts(context.Background(), c)
	require.NoError(t, err)

	fmt.Println(queryResponse)

	var fabelmansForecastId string
	var yesOutcomeId string

	for _, f := range queryResponse.Forecasts {
		if strings.Contains(f.Title, "\"The Fabelmans\"") {
			fabelmansForecastId = f.Id
			require.Equal(t, ResolutionUnresolved, f.Resolution)
			require.NotEmpty(t, f.Estimates)
			for _, p := range f.Estimates[0].Probabilities {
				require.False(t, p.Outcome.Correct)
				if p.Outcome.Text == "Yes" {
					yesOutcomeId = p.Outcome.Id
					break
				}
			}
			break
		}
	}
	require.NotEmpty(t, fabelmansForecastId)
	require.NotEmpty(t, yesOutcomeId)

	resolveForecastResponse, err := ResolveForecast(
		context.Background(),
		c,
		fabelmansForecastId,
		&yesOutcomeId,
		nil,
	)
	require.NoError(t, err)

	assert.Contains(t, resolveForecastResponse.ResolveForecast.Title, "\"The Fabelmans\"")
	assert.Equal(t, ResolutionResolved, resolveForecastResponse.ResolveForecast.Resolution)

	require.NotEmpty(t, resolveForecastResponse.ResolveForecast.Estimates)
	for _, p := range resolveForecastResponse.ResolveForecast.Estimates[0].Probabilities {
		if p.Outcome.Text == "Yes" {
			assert.True(t, p.Outcome.Correct)
		} else {
			assert.False(t, p.Outcome.Correct)
		}
	}

	// Verify that the same forecast can't be resolved again
	_, err = ResolveForecast(
		context.Background(),
		c,
		fabelmansForecastId,
		&yesOutcomeId,
		nil,
	)
	assert.ErrorContains(t, err, "forecast has already been resolved")
}

// TestResolveForecast_WrongOutcomeId verifies that it's not possible to
// specify an outcome ID of a different forecast when resolving the forecast.
func TestResolveForecast_WrongOutcomeId(t *testing.T) {
	c := initServerAndGetClient2(t)

	queryResponse, err := GetForecasts(context.Background(), c)
	require.NoError(t, err)

	var fabelmansForecastId string
	var invalidOutcomeId string

	for _, f := range queryResponse.Forecasts {
		if strings.Contains(f.Title, "\"The Fabelmans\"") {
			fabelmansForecastId = f.Id
			require.Equal(t, ResolutionUnresolved, f.Resolution)
			require.NotEmpty(t, f.Estimates)
			for _, p := range f.Estimates[0].Probabilities {
				require.False(t, p.Outcome.Correct)
			}
		} else if invalidOutcomeId == "" {
			// choose a random outcomeId from another forecast
			require.NotEmpty(t, f.Estimates)
			for _, p := range f.Estimates[0].Probabilities {
				if !p.Outcome.Correct {
					invalidOutcomeId = p.Outcome.Id
					break
				}
			}
		}
	}
	require.NotEmpty(t, fabelmansForecastId)
	require.NotEmpty(t, invalidOutcomeId)

	_, err = ResolveForecast(
		context.Background(),
		c,
		fabelmansForecastId,
		&invalidOutcomeId,
		nil,
	)
	assert.ErrorContains(t, err, "can't match")

	// Get all forecasts again and verify nothing has been changed

	queryResponse, err = GetForecasts(context.Background(), c)
	require.NoError(t, err)

	for _, f := range queryResponse.Forecasts {
		if f.Id == fabelmansForecastId {
			assert.Equal(t, ResolutionUnresolved, f.Resolution)
		}
		require.NotEmpty(t, f.Estimates)
		for _, p := range f.Estimates[0].Probabilities {
			if p.Outcome.Id == invalidOutcomeId {
				assert.False(
					t,
					p.Outcome.Correct,
					"the non-matching outcome %v should still have Correct == false",
					invalidOutcomeId,
				)
			}
		}
	}
}

// TestResolveForecast_ResolvesInFuture verifies that when resolves/closes are
// set in the future they are set to now() when resolving.
func TestResolveForecast_ResolvesInFuture(t *testing.T) {
	c := initServerAndGetClient2(t)

	// Create a new forecast as set up

	newForecast := NewForecast{
		Title: "Will it rain tomorrow?",
		Description: "It counts as rain if between 9am and 9pm there are " +
			"30 min or more of uninterrupted precipitation.",
		Closes:   timePointer(time.Now().Add(24 * time.Hour)),
		Resolves: time.Now().Add(24 * time.Hour),
	}

	newEstimate := NewEstimate{
		Reason: "My weather app says it will rain.",
		Probabilities: []NewProbability{
			{
				Value: 70,
				Outcome: NewOutcome{
					Text: "Yes",
				},
			},
			{
				Value: 30,
				Outcome: NewOutcome{
					Text: "No",
				},
			},
		},
	}

	response, err := CreateForecast(
		context.Background(),
		c,
		newForecast,
		newEstimate,
	)
	require.NoError(t, err)
	require.NotEmpty(t, response.CreateForecast.Id)
	require.Len(t, response.CreateForecast.Estimates, 1)
	require.Len(t, response.CreateForecast.Estimates[0].Probabilities, 2)

	forecastId := response.CreateForecast.Id
	outcomeId := ""

	for _, p := range response.CreateForecast.Estimates[0].Probabilities {
		if p.Outcome.Text == "Yes" {
			outcomeId = p.Outcome.Id
		}
	}
	require.NotEmpty(t, forecastId)
	require.NotEmpty(t, outcomeId)

	// Set up done

	resolveResponse, err := ResolveForecast(
		context.Background(),
		c,
		forecastId,
		&outcomeId,
		nil,
	)
	require.NoError(t, err)

	// Verify the dates have been set to 'now' (+/- 2 minutes to account for
	// test execution delays).
	assert.Equal(t, ResolutionResolved, resolveResponse.ResolveForecast.Resolution)
	assertTimeAlmostEqual(
		t,
		time.Now(),
		resolveResponse.ResolveForecast.Resolves,
	)
	assertTimeAlmostEqual(
		t,
		time.Now(),
		*resolveResponse.ResolveForecast.Closes,
	)
}

// TestResolveForecast_ResolvesInPast verifies that when resolves/closes are
// set in the past they stay as they are.
func TestResolveForecast_ResolvesInPast(t *testing.T) {
	client := initServerAndGetClient2(t)
	getResponse, err := GetForecasts(context.Background(), client)
	require.NoError(t, err)

	var currentResolves time.Time
	var currentCloses *time.Time
	var forecastId string
	var outcomeNoId string

	for _, f := range getResponse.Forecasts {
		if f.Title == "Will the number of contributors to \"Cleodora\" be more "+
			"than 3 at the end of 2022?" {
			forecastId = f.Id
			require.Equal(t, ResolutionUnresolved, f.Resolution)
			currentResolves = f.Resolves
			currentCloses = f.Closes
			for _, p := range f.Estimates[0].Probabilities {
				if p.Outcome.Text == "No" {
					outcomeNoId = p.Outcome.Id
					break
				}
			}
			break
		}
	}
	require.NotEmpty(t, forecastId, "did not find the expected forecast")
	require.NotEmpty(t, outcomeNoId)

	resolveResponse, err := ResolveForecast(
		context.Background(),
		client,
		forecastId,
		&outcomeNoId,
		nil,
	)
	require.NoError(t, err)

	assert.Equal(t, currentResolves.UTC(), resolveResponse.ResolveForecast.Resolves.UTC())
	if currentCloses == nil {
		assert.Nil(t, resolveResponse.ResolveForecast.Closes)
	} else {
		assert.Equal(t, currentCloses.UTC(), resolveResponse.ResolveForecast.Closes.UTC())
	}
}
