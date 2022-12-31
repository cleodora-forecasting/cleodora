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

// TestCreateForecast_CanOmitCloses verifies that a forecast can be created
// without specifying a closing date.
func TestCreateForecast_CanOmitCloses(t *testing.T) {
	c := initServer(t)

	newForecast := map[string]interface{}{
		"title": "Will it rain tomorrow?",
		"description": "It counts as rain if between 9am and 9pm there are " +
			"30 min or more of uninterrupted precipitation.",
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
				closes
			}
		}`

	var response struct {
		CreateForecast struct {
			Id     string
			Title  string
			Closes string
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
	assert.Empty(t, response.CreateForecast.Closes)
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
