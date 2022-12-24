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

	"github.com/cleodora-forecasting/cleodora/cleosrv/graph"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/generated"
)

// TestGetForecasts_LowLevel verifies that the forecasts are returned by
// sending a low level JSON request and just checking a string in the response
// body, without any further GraphQL processing.
func TestGetForecasts_LowLevel(t *testing.T) {
	// Set up the server
	resolver := graph.Resolver{}
	resolver.AddDummyData()
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}),
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
	c := initServer()

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
		"Will the number of contributors for \"Cleodora\" be more than 3 at"+
			" the end of 2022?",
		response.Forecasts[2].Title,
	)
}

// TestGetForecasts_SomeFields verifies that the query can contain only a few
// fields.
func TestGetForecasts_OnlySomeFields(t *testing.T) {
	c := initServer()

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
		"Will the number of contributors for \"Cleodora\" be more than 3 at"+
			" the end of 2022?",
		response.Forecasts[2].Title,
	)
}

// TestGetForecasts_InvalidField verifies that an error is returned when
// querying for a field that does not exist.
func TestGetForecasts_InvalidField(t *testing.T) {
	c := initServer()

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
	c := initServer()

	newForecast := map[string]interface{}{
		"title": "Will it rain tomorrow?",
		"description": "It counts as rain if between 9am and 9pm there are " +
			"30 min or more of uninterrupted precipitation.",
		"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	query := `
		mutation createForecast($input: NewForecast!) {
			createForecast(input: $input) {
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

	err := c.Post(query, &response, client.Var("input", newForecast))
	require.Nil(t, err)

	assert.NotEmpty(t, response.CreateForecast.Id)
	assert.Equal(
		t,
		"Will it rain tomorrow?",
		response.CreateForecast.Title,
	)
}

// TestCreateForecast_OmitCloses verifies that a forecast can be created
// without specifying a closing date.
func TestCreateForecast_OmitCloses(t *testing.T) {
	c := initServer()

	newForecast := map[string]interface{}{
		"title": "Will it rain tomorrow?",
		"description": "It counts as rain if between 9am and 9pm there are " +
			"30 min or more of uninterrupted precipitation.",
		"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	query := `
		mutation createForecast($input: NewForecast!) {
			createForecast(input: $input) {
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

	err := c.Post(query, &response, client.Var("input", newForecast))
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
	c := initServer()

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

func initServer() *client.Client {
	// Set up the server
	resolver := graph.Resolver{}
	resolver.AddDummyData()
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}),
	)

	c := client.New(srv)
	return c
}
