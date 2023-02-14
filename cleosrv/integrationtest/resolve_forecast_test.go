package integrationtest

import (
	"strings"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResolveForecast tests the resolve forecast happy case.
// It resolves the forecast and set one of the Outcomes as correct.
// It also verifies that the same forecast can't be resolved again.
func TestResolveForecast(t *testing.T) {
	c := initServerAndGetClient(t)

	queryGetForecasts := `
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
		}`

	type responseForecast struct {
		Id         string
		Title      string
		Resolution string
		Estimates  []struct {
			Id            string
			Probabilities []struct {
				Id      string
				Outcome struct {
					Id      string
					Text    string
					Correct bool
				}
			}
		}
	}

	var queryResponse struct {
		Forecasts []responseForecast
	}

	err := c.Post(queryGetForecasts, &queryResponse)
	require.NoError(t, err)

	var fabelmansForecastId string
	var yesOutcomeId string

	for _, f := range queryResponse.Forecasts {
		if strings.Contains(f.Title, "\"The Fabelmans\"") {
			fabelmansForecastId = f.Id
			require.Equal(t, "UNRESOLVED", f.Resolution)
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

	mutationResolveForecast := `
		mutation resolveForecast($forecastId: ID!, $correctOutcomeId: ID) {
			resolveForecast(forecastId: $forecastId, correctOutcomeId: $correctOutcomeId) {
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
		}`

	var resolveForecastResponse struct {
		ResolveForecast responseForecast
	}

	err = c.Post(
		mutationResolveForecast,
		&resolveForecastResponse,
		client.Var("forecastId", fabelmansForecastId),
		client.Var("correctOutcomeId", yesOutcomeId),
	)
	require.NoError(t, err)

	assert.Contains(t, resolveForecastResponse.ResolveForecast.Title, "\"The Fabelmans\"")
	assert.Equal(t, "RESOLVED", resolveForecastResponse.ResolveForecast.Resolution)

	require.NotEmpty(t, resolveForecastResponse.ResolveForecast.Estimates)
	for _, p := range resolveForecastResponse.ResolveForecast.Estimates[0].Probabilities {
		if p.Outcome.Text == "Yes" {
			assert.True(t, p.Outcome.Correct)
		} else {
			assert.False(t, p.Outcome.Correct)
		}
	}

	// Verify that the same forecast can't be resolved again
	var resolveForecastResponse2 struct {
		ResolveForecast responseForecast
	}
	err = c.Post(
		mutationResolveForecast,
		&resolveForecastResponse2,
		client.Var("forecastId", fabelmansForecastId),
		client.Var("correctOutcomeId", yesOutcomeId),
	)
	assert.ErrorContains(t, err, "forecast has already been resolved")
}

// TestResolveForecast_WrongOutcomeId verifies that it's not possible to
// specify an outcome ID of a different forecast when resolving the forecast.
func TestResolveForecast_WrongOutcomeId(t *testing.T) {
	c := initServerAndGetClient(t)

	queryGetForecasts := `
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
		}`

	type responseForecast struct {
		Id         string
		Title      string
		Resolution string
		Estimates  []struct {
			Id            string
			Probabilities []struct {
				Id      string
				Outcome struct {
					Id      string
					Text    string
					Correct bool
				}
			}
		}
	}

	var queryResponse struct {
		Forecasts []responseForecast
	}

	err := c.Post(queryGetForecasts, &queryResponse)
	require.NoError(t, err)

	var fabelmansForecastId string
	var invalidOutcomeId string

	for _, f := range queryResponse.Forecasts {
		if strings.Contains(f.Title, "\"The Fabelmans\"") {
			fabelmansForecastId = f.Id
			require.Equal(t, "UNRESOLVED", f.Resolution)
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

	mutationResolveForecast := `
		mutation resolveForecast($forecastId: ID!, $correctOutcomeId: ID) {
			resolveForecast(forecastId: $forecastId, correctOutcomeId: $correctOutcomeId) {
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
		}`

	var resolveForecastResponse struct {
		ResolveForecast responseForecast
	}

	err = c.Post(
		mutationResolveForecast,
		&resolveForecastResponse,
		client.Var("forecastId", fabelmansForecastId),
		client.Var("correctOutcomeId", invalidOutcomeId),
	)
	assert.ErrorContains(t, err, "can't match")

	// Get all forecasts again and verify nothing has been changed

	queryResponse = struct {
		Forecasts []responseForecast
	}{}

	err = c.Post(queryGetForecasts, &queryResponse)
	require.NoError(t, err)

	for _, f := range queryResponse.Forecasts {
		if f.Id == fabelmansForecastId {
			assert.Equal(t, "UNRESOLVED", f.Resolution)
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

// TestResolveForecast_WithResolutionNA_NoOutcome verifies that it's
// possible to resolve a forecast as NA (not applicable) and that in that case
// no outcome is set to correct.
func TestResolveForecast_WithResolutionNA_NoOutcome(t *testing.T) {
	c := initServerAndGetClient(t)

	queryGetForecasts := `
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
		}`

	type responseForecast struct {
		Id         string
		Title      string
		Resolution string
		Estimates  []struct {
			Id            string
			Probabilities []struct {
				Id      string
				Outcome struct {
					Id      string
					Text    string
					Correct bool
				}
			}
		}
	}

	var queryResponse struct {
		Forecasts []responseForecast
	}

	err := c.Post(queryGetForecasts, &queryResponse)
	require.NoError(t, err)

	var fabelmansForecastId string

	for _, f := range queryResponse.Forecasts {
		if strings.Contains(f.Title, "\"The Fabelmans\"") {
			fabelmansForecastId = f.Id
			require.Equal(t, "UNRESOLVED", f.Resolution)
			break
		}
	}
	require.NotEmpty(t, fabelmansForecastId)

	mutationResolveForecast := `
		mutation resolveForecast($forecastId: ID!, $resolution: Resolution) {
			resolveForecast(
				forecastId: $forecastId,
				resolution: $resolution,
			) {
				id
				title
				resolution
			}
		}`

	var resolveForecastResponse struct {
		ResolveForecast responseForecast
	}

	err = c.Post(
		mutationResolveForecast,
		&resolveForecastResponse,
		client.Var("forecastId", fabelmansForecastId),
		client.Var("resolution", "NOT_APPLICABLE"),
	)
	require.NoError(t, err)

	assert.Contains(t, resolveForecastResponse.ResolveForecast.Title, "\"The Fabelmans\"")
	assert.Equal(t, "NOT_APPLICABLE", resolveForecastResponse.ResolveForecast.Resolution)

	// Get all forecasts again and verify no outcome is correct

	queryResponse = struct {
		Forecasts []responseForecast
	}{}

	err = c.Post(queryGetForecasts, &queryResponse)
	require.NoError(t, err)

	for _, f := range queryResponse.Forecasts {
		if f.Id == fabelmansForecastId {
			assert.Equal(t, "NOT_APPLICABLE", f.Resolution)
			require.NotEmpty(t, f.Estimates)
			for _, p := range f.Estimates[0].Probabilities {
				assert.False(
					t,
					p.Outcome.Correct,
					"the outcome %v should still have Correct == false",
					p.Outcome.Text,
				)
			}
		}
	}
}

// TestResolveForecast_WithResolutionNA_AndOutcome verifies that it's possible
// to resolve a forecast as NA (not applicable) and that in that case no
// outcome is set to correct even if that outcome is passed as a parameter.
func TestResolveForecast_WithResolutionNA_AndOutcome(t *testing.T) {
	c := initServerAndGetClient(t)

	queryGetForecasts := `
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
		}`

	type responseForecast struct {
		Id         string
		Title      string
		Resolution string
		Estimates  []struct {
			Id            string
			Probabilities []struct {
				Id      string
				Outcome struct {
					Id      string
					Text    string
					Correct bool
				}
			}
		}
	}

	var queryResponse struct {
		Forecasts []responseForecast
	}

	err := c.Post(queryGetForecasts, &queryResponse)
	require.NoError(t, err)

	var fabelmansForecastId string
	var yesOutcomeId string

	for _, f := range queryResponse.Forecasts {
		if strings.Contains(f.Title, "\"The Fabelmans\"") {
			fabelmansForecastId = f.Id
			require.Equal(t, "UNRESOLVED", f.Resolution)
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

	mutationResolveForecast := `
		mutation resolveForecast(
			$forecastId: ID!,
			$resolution: Resolution,
			$correctOutcomeId: ID,
		) {
			resolveForecast(
				forecastId: $forecastId,
				resolution: $resolution,
				correctOutcomeId: $correctOutcomeId,
			) {
				id
				title
				resolution
			}
		}`

	var resolveForecastResponse struct {
		ResolveForecast responseForecast
	}

	err = c.Post(
		mutationResolveForecast,
		&resolveForecastResponse,
		client.Var("forecastId", fabelmansForecastId),
		client.Var("resolution", "NOT_APPLICABLE"),
		client.Var("correctOutcomeId", yesOutcomeId),
	)
	require.NoError(t, err)

	assert.Contains(t, resolveForecastResponse.ResolveForecast.Title, "\"The Fabelmans\"")
	assert.Equal(t, "NOT_APPLICABLE", resolveForecastResponse.ResolveForecast.Resolution)

	// Get all forecasts again and verify no outcome is correct

	queryResponse = struct {
		Forecasts []responseForecast
	}{}

	err = c.Post(queryGetForecasts, &queryResponse)
	require.NoError(t, err)

	for _, f := range queryResponse.Forecasts {
		if f.Id == fabelmansForecastId {
			assert.Equal(t, "NOT_APPLICABLE", f.Resolution)
			require.NotEmpty(t, f.Estimates)
			for _, p := range f.Estimates[0].Probabilities {
				assert.False(
					t,
					p.Outcome.Correct,
					"the outcome %v should still have Correct == false",
					p.Outcome.Text,
				)
			}
		}
	}
}

// TestResolveForecast_NoOutcome verifies that it's not possible to resolve a
// forecasts as RESOLVED if no Outcome is passed in.
func TestResolveForecast_NoOutcome(t *testing.T) {
	c := initServerAndGetClient(t)

	queryGetForecasts := `
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
		}`

	type responseForecast struct {
		Id         string
		Title      string
		Resolution string
		Estimates  []struct {
			Id            string
			Probabilities []struct {
				Id      string
				Outcome struct {
					Id      string
					Text    string
					Correct bool
				}
			}
		}
	}

	var queryResponse struct {
		Forecasts []responseForecast
	}

	err := c.Post(queryGetForecasts, &queryResponse)
	require.NoError(t, err)

	var fabelmansForecastId string
	var yesOutcomeId string

	for _, f := range queryResponse.Forecasts {
		if strings.Contains(f.Title, "\"The Fabelmans\"") {
			fabelmansForecastId = f.Id
			require.Equal(t, "UNRESOLVED", f.Resolution)
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

	mutationResolveForecast := `
		mutation resolveForecast($forecastId: ID!) {
			resolveForecast(forecastId: $forecastId) {
				id
				title
				resolution
			}
		}`

	var resolveForecastResponse struct {
		ResolveForecast responseForecast
	}

	err = c.Post(
		mutationResolveForecast,
		&resolveForecastResponse,
		client.Var("forecastId", fabelmansForecastId),
	)
	assert.ErrorContains(t, err, "to resolve as RESOLVED, an Outcome must be specified")
}

// TestResolveForecast_WithResolutionUnresolved verifies that it's not possible
// to resolve with resolution UNRESOLVED.
func TestResolveForecast_WithResolutionUnresolved(t *testing.T) {
	c := initServerAndGetClient(t)

	queryGetForecasts := `
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
		}`

	type responseForecast struct {
		Id         string
		Title      string
		Resolution string
		Estimates  []struct {
			Id            string
			Probabilities []struct {
				Id      string
				Outcome struct {
					Id      string
					Text    string
					Correct bool
				}
			}
		}
	}

	var queryResponse struct {
		Forecasts []responseForecast
	}

	err := c.Post(queryGetForecasts, &queryResponse)
	require.NoError(t, err)

	var fabelmansForecastId string
	var yesOutcomeId string

	for _, f := range queryResponse.Forecasts {
		if strings.Contains(f.Title, "\"The Fabelmans\"") {
			fabelmansForecastId = f.Id
			require.Equal(t, "UNRESOLVED", f.Resolution)
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

	mutationResolveForecast := `
		mutation resolveForecast(
			$forecastId: ID!,
			$resolution: Resolution,
			$correctOutcomeId: ID,
		) {
			resolveForecast(
				forecastId: $forecastId,
				resolution: $resolution,
				correctOutcomeId: $correctOutcomeId,
			) {
				id
				title
				resolution
			}
		}`

	var resolveForecastResponse struct {
		ResolveForecast responseForecast
	}

	err = c.Post(
		mutationResolveForecast,
		&resolveForecastResponse,
		client.Var("forecastId", fabelmansForecastId),
		client.Var("resolution", "UNRESOLVED"),
		client.Var("correctOutcomeId", yesOutcomeId),
	)
	assert.ErrorContains(t, err, "resolution UNRESOLVED is not allowed")
}
