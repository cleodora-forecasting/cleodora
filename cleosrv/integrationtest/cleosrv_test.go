package integrationtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	// Prepare the request
	body := "{\"operationName\":\"GetForecasts\",\"variables\":{}," +
		"\"query\":\"query GetForecasts {\\n  forecasts {\\n    id" +
		"\\n    summary\\n    description\\n    created\\n    closes" +
		"\\n    resolves\\n    resolution\\n    __typename\\n  }\\n}\"}"
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
	// Set up the server
	resolver := graph.Resolver{}
	resolver.AddDummyData()
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}),
	)

	c := client.New(srv)

	query := `
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
			}`

	var resp struct {
		Forecasts []struct {
			Closes      string
			Created     string
			Description string
			Id          string
			Resolution  string
			Resolves    string
			Summary     string
		}
	}

	err := c.Post(query, &resp)
	assert.Nil(t, err)

	t.Log(resp)

	assert.Len(t, resp.Forecasts, 3)
	assert.Equal(
		t,
		resp.Forecasts[2].Summary,
		"Will the number of contributors for \"Cleodora\" be more than 3 at"+
			" the end of 2022?",
	)
}
