package integrationtest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cleodora-forecasting/cleodora/cleosrv/cleosrv"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/generated"
)

// Tests that use the GraphQL directly without any client libraries.
// This is mainly to serve as documentation.

// TestGetForecasts_LowLevel verifies that the forecasts are returned by
// sending a low level JSON request and just checking a string in the response
// body, without any further GraphQL processing.
func TestGetForecasts_LowLevel(t *testing.T) {
	// Set up the server
	app := cleosrv.NewApp()
	app.Config.Database = fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := app.InitDB()
	require.NoError(t, err)
	resolver := graph.NewResolver(db)
	err = resolver.AddDummyData()
	require.NoError(t, err)
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
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	res := rec.Result()

	// Evaluate the result
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Contains(t, string(b), "Fabelmans")
}
