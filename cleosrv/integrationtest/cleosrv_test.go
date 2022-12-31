package integrationtest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cleodora-forecasting/cleodora/cleosrv/cleosrv"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/generated"
)

// TestGetForecasts_LowLevel verifies that the forecasts are returned by
// sending a low level JSON request and just checking a string in the response
// body, without any further GraphQL processing.
func TestGetForecasts_LowLevel(t *testing.T) {
	// Set up the server
	db, err := cleosrv.InitDB(":memory:")
	require.Nil(t, err)
	resolver := graph.NewResolver(db)
	err = resolver.AddDummyData()
	require.Nil(t, err)
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: resolver}),
	)

	query := `
		query GetForecasts {
			forecasts {
				id
				title
				description
				created
				closes
				resolves
				resolution
				__typename
			}
		}`
	query = strings.ReplaceAll(query, "\n", `\n`)
	query = strings.ReplaceAll(query, "\t", "    ")

	body := `{
		"operationName": "GetForecasts",
		"variables": {},
		"query": "%s"
	}`
	body = fmt.Sprintf(body, query)
	t.Log("body", body)

	req, err := http.NewRequest(
		"POST",
		"localhost:8080/query",
		strings.NewReader(body),
	)
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	res := rec.Result()

	// Evaluate the result
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	assert.Contains(t, string(b), "Fabelmans")
}

// TestGetForecasts_GQClient verifies that the forecasts are returned and
// uses the gqlgen.client for it.
func TestGetForecasts_GQClient(t *testing.T) {
	c := initServer(t)

	query := `
		query GetForecasts {
			forecasts {
				id
				title
				description
				created
				closes
				resolves
				resolution
			}
		}`

	var response struct {
		Forecasts []struct {
			Closes      string
			Created     string
			Description string
			Id          string
			Resolution  string
			Resolves    string
			Title       string
		}
	}

	err := c.Post(query, &response)
	assert.Nil(t, err)

	t.Log(response)

	assert.Len(t, response.Forecasts, 3)
	assert.Equal(
		t,
		"Will the number of contributors to \"Cleodora\" be more than 3 at"+
			" the end of 2022?",
		response.Forecasts[2].Title,
	)
}

// TestGetForecasts_SomeFields verifies that the query can contain only a few
// fields.
func TestGetForecasts_OnlySomeFields(t *testing.T) {
	c := initServer(t)

	query := `
		query GetForecasts {
			forecasts {
				id
				title
			}
		}`

	var response struct {
		Forecasts []struct {
			Id    string
			Title string
		}
	}

	err := c.Post(query, &response)
	assert.Nil(t, err)

	t.Log(response)

	assert.Len(t, response.Forecasts, 3)
	assert.Equal(
		t,
		"Will the number of contributors to \"Cleodora\" be more than 3 at"+
			" the end of 2022?",
		response.Forecasts[2].Title,
	)
}

// TestGetForecasts_InvalidField verifies that an error is returned when
// querying for a field that does not exist.
func TestGetForecasts_InvalidField(t *testing.T) {
	c := initServer(t)

	query := `
		query GetForecasts {
			forecasts {
				id
				title
				does_not_exist
			}
		}`

	var response struct {
		Forecasts []struct {
			Id    string
			Title string
		}
	}

	err := c.Post(query, &response)
	assert.Contains(t, err.Error(), "http 422")
	assert.Contains(t, err.Error(), "Cannot query field \\\"does_not_exist\\\"")
}

func TestCreateForecast(t *testing.T) {
	c := initServer(t)

	newForecast := map[string]interface{}{
		"title": "Will it rain tomorrow?",
		"description": "It counts as rain if between 9am and 9pm there are " +
			"30 min or more of uninterrupted precipitation.",
		"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	newEstimate := map[string]interface{}{
		"reason": "My weather app says it will rain.",
		"probabilities": []map[string]interface{}{
			{
				"value": 70,
				"outcome": map[string]interface{}{
					"text": "Yes",
				},
			},
			{
				"value": 30,
				"outcome": map[string]interface{}{
					"text": "No",
				},
			},
		},
	}

	query := `
		mutation createForecast($forecast: NewForecast!, $estimate: NewEstimate!) {
			createForecast(forecast: $forecast, estimate: $estimate) {
				id
				title
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
		}`

	var response struct {
		CreateForecast struct {
			Id        string
			Title     string
			Estimates []struct {
				Id            string
				Created       string
				Reason        string
				Probabilities []struct {
					Id      string
					Value   int
					Outcome struct {
						Id      string
						Text    string
						Correct bool
					}
				}
			}
		}
	}

	err := c.Post(
		query,
		&response,
		client.Var("forecast", newForecast),
		client.Var("estimate", newEstimate),
	)
	require.Nil(t, err)

	assert.NotEmpty(t, response.CreateForecast.Id)
	assert.Equal(
		t,
		"Will it rain tomorrow?",
		response.CreateForecast.Title,
	)

	assert.Len(t, response.CreateForecast.Estimates, 1)
	assert.NotEmpty(t, response.CreateForecast.Estimates[0].Id)
	assert.Equal(
		t,
		"My weather app says it will rain.",
		response.CreateForecast.Estimates[0].Reason,
	)
	assert.Len(t, response.CreateForecast.Estimates[0].Probabilities, 2)
	assert.NotEmpty(t, response.CreateForecast.Estimates[0].Probabilities[0].Id)
	assert.NotEmpty(t, response.CreateForecast.Estimates[0].Probabilities[1].Id)
	assert.False(t, response.CreateForecast.Estimates[0].Probabilities[0].Outcome.Correct)
	assert.False(t, response.CreateForecast.Estimates[0].Probabilities[1].Outcome.Correct)

	// If the order is Yes, No ...
	if response.CreateForecast.Estimates[0].Probabilities[0].Outcome.Text == "Yes" {
		assert.Equal(t, 70, response.CreateForecast.Estimates[0].Probabilities[0].Value)
		assert.Equal(t, "Yes", response.CreateForecast.Estimates[0].Probabilities[0].Outcome.Text)
		assert.Equal(t, 30, response.CreateForecast.Estimates[0].Probabilities[1].Value)
		assert.Equal(t, "No", response.CreateForecast.Estimates[0].Probabilities[1].Outcome.Text)
	} else { // ... or if it is No, Yes
		assert.Equal(t, 30, response.CreateForecast.Estimates[0].Probabilities[0].Value)
		assert.Equal(t, "No", response.CreateForecast.Estimates[0].Probabilities[0].Outcome.Text)
		assert.Equal(t, 70, response.CreateForecast.Estimates[0].Probabilities[1].Value)
		assert.Equal(t, "Yes", response.CreateForecast.Estimates[0].Probabilities[1].Outcome.Text)
	}
}

func TestCreateForecast_ValidateNewEstimate(t *testing.T) {
	tests := []struct {
		name        string
		newEstimate map[string]interface{}
		expectedErr string
	}{
		{
			name: "success",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "Yes",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "",
		},
		{
			name: "single probability",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 100,
						"outcome": map[string]interface{}{
							"text": "Any outcome",
						},
					},
				},
			},
			expectedErr: "",
		},
		{
			name: "success with more probabilities",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "Yes, but less than 2 hours",
						},
					},
					{
						"value": 20,
						"outcome": map[string]interface{}{
							"text": "Yes, between 2 and 5 hours",
						},
					},
					{
						"value": 20,
						"outcome": map[string]interface{}{
							"text": "Yes, more than 5 hours",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "",
		},
		{
			name: "reason cant be empty",
			newEstimate: map[string]interface{}{
				"reason": "",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "Yes",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "'reason' can't be empty",
		},
		{
			name: "probabilities must add up to 100",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain.",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "Yes",
						},
					},
					{
						"value": 20,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "probabilities must add up to 100",
		},
		{
			name: "probabilities must be between 0 and 100",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain.",
				"probabilities": []map[string]interface{}{
					{
						"value": -10,
						"outcome": map[string]interface{}{
							"text": "Yes",
						},
					},
					{
						"value": 110,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "probabilities must be between 0 and 100",
		},
		{
			name: "probabilities cant be empty",
			newEstimate: map[string]interface{}{
				"reason":        "My weather app says it will rain.",
				"probabilities": []map[string]interface{}{},
			},
			expectedErr: "probabilities can't be empty",
		},
		{
			name: "outcome must be set",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value":   70,
						"outcome": nil,
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "cannot be null",
		},
		{
			name: "outcome cant be empty",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value":   70,
						"outcome": map[string]interface{}{},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "must be defined",
		},
		{
			name: "outcome text cant be empty",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "outcome text can't be empty",
		},
		{
			name: "outcomes cant be duplicates",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "outcome 'No' is a duplicate",
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {

			c := initServer(t)

			newForecast := map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			}

			query := `
		mutation createForecast($forecast: NewForecast!, $estimate: NewEstimate!) {
			createForecast(forecast: $forecast, estimate: $estimate) {
				id
				title
			}
		}`

			var response struct {
				CreateForecast struct {
					Id    string
					Title string
				}
			}

			err := c.Post(
				query,
				&response,
				client.Var("forecast", newForecast),
				client.Var("estimate", tt.newEstimate),
			)
			if tt.expectedErr == "" { // success
				assert.Nil(t, err)
			} else {
				assert.ErrorContains(t, err, tt.expectedErr)
			}
		})
	}
}

func TestCreateForecast_ValidateNewForecast(t *testing.T) {
	tests := []struct {
		name        string
		newForecast map[string]interface{}
		expectedErr string
	}{
		{
			name: "success",
			newForecast: map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "",
		},
		{
			name: "title cant be empty",
			newForecast: map[string]interface{}{
				"title": "",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "title can't be empty",
		},
		{
			name: "description can be empty",
			newForecast: map[string]interface{}{
				"title":       "Will it rain tomorrow?",
				"description": "",
				"closes":      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"resolves":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "",
		},
		{
			name: "closes can be empty",
			newForecast: map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes":   nil,
				"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "",
		},
		{
			name: "closes can be omitted",
			newForecast: map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "",
		},
		{
			name: "resolves cant be empty",
			newForecast: map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"resolves": nil,
			},
			expectedErr: "cannot be null",
		},
		{
			name: "resolves cant be omitted",
			newForecast: map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "must be defined",
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {

			c := initServer(t)

			newEstimate := map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "Yes",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			}

			query := `
		mutation createForecast($forecast: NewForecast!, $estimate: NewEstimate!) {
			createForecast(forecast: $forecast, estimate: $estimate) {
				id
				title
			}
		}`

			var response struct {
				CreateForecast struct {
					Id    string
					Title string
				}
			}

			err := c.Post(
				query,
				&response,
				client.Var("forecast", tt.newForecast),
				client.Var("estimate", newEstimate),
			)
			if tt.expectedErr == "" { // success
				assert.Nil(t, err)
			} else {
				assert.ErrorContains(t, err, tt.expectedErr)
			}
		})
	}
}

func TestCreateForecast_FailsWithoutEstimate(t *testing.T) {
	c := initServer(t)

	newForecast := map[string]interface{}{
		"title": "Will it rain tomorrow?",
		"description": "It counts as rain if between 9am and 9pm there are " +
			"30 min or more of uninterrupted precipitation.",
		"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	query := `
		mutation createForecast($forecast: NewForecast!) {
			createForecast(forecast: $forecast) {
				id
				title
			}
		}`

	var response struct {
		CreateForecast struct {
			Id    string
			Title string
		}
	}

	err := c.Post(query, &response, client.Var("forecast", newForecast))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "http 422")
	assert.Contains(
		t,
		err.Error(),
		"argument \\\"estimate\\\" of type \\\"NewEstimate!\\\" is required",
	)
}

func TestGetVersion(t *testing.T) {
	c := initServer(t)

	query := `
		query GetMetadata {
			metadata {
				version
			}
		}`

	var resp struct {
		Metadata struct {
			Version string
		}
	}

	err := c.Post(query, &resp)
	assert.Nil(t, err)

	t.Log(resp)
	assert.Equal(t, "dev", resp.Metadata.Version)
}

func initServer(t *testing.T) *client.Client {
	t.Helper()
	// Set up the server
	db, err := cleosrv.InitDB(":memory:")
	require.Nil(t, err)
	resolver := graph.NewResolver(db)
	err = resolver.AddDummyData()
	require.Nil(t, err)
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: resolver}),
	)

	c := client.New(srv)
	return c
}
