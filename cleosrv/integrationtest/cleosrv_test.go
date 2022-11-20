package integrationtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cleodora-forecasting/cleodora/cleosrv/graph"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/generated"
)

func TestGetForecasts(t *testing.T) {
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
