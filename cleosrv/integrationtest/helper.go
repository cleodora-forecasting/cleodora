package integrationtest

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cleodora-forecasting/cleodora/cleosrv/cleosrv"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/generated"
)

// initServerAndGetClient returns a gqlgen Client that is not meant for public
// consumption by gqlgen. For that reason this function is deprecated, and it
// should be replaced where possible with initServerAndGetClient2 .
func initServerAndGetClient(t *testing.T) *client.Client {
	t.Helper()
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

	c := client.New(srv)
	return c
}

// initServerAndGetClient2 returns a graphql.Client generated by genqlient
func initServerAndGetClient2(t *testing.T) graphql.Client {
	t.Helper()
	app := cleosrv.NewApp()
	app.Config.Database = fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := app.InitDB()
	require.NoError(t, err)
	resolver := graph.NewResolver(db)
	err = resolver.AddDummyData()
	require.NoError(t, err)
	srv := httptest.NewServer(
		handler.NewDefaultServer(
			generated.NewExecutableSchema(generated.Config{Resolvers: resolver}),
		),
	)
	return graphql.NewClient(srv.URL, srv.Client())
}

func resolutionPointer(r Resolution) *Resolution {
	return &r
}

func timePointer(t time.Time) *time.Time {
	return &t
}

// assertTimeAlmostEqual asserts that the two time stamps are within 2 minutes
// of each other. The two minutes are chosen to account for some delay in
// test execution.
func assertTimeAlmostEqual(t *testing.T, expected, actual time.Time) {
	assert.InDelta(
		t,
		expected.UTC().Unix(),
		actual.UTC().Unix(),
		(2 * time.Minute).Seconds(),
	)
}