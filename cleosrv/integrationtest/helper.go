package integrationtest

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cleodora-forecasting/cleodora/cleosrv/cleosrv"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/generated"
)

// initServerAndGetClient returns a graphql.Client generated by genqlient
func initServerAndGetClient(t *testing.T, dbPath string) graphql.Client {
	t.Helper()
	app := cleosrv.NewApp()
	if dbPath == "" {
		app.Config.Database = fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	} else {
		app.Config.Database = dbPath
	}
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

// strPointer returns the pointer of a string
func strPointer(s string) *string {
	return &s
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

func assertUTC(t *testing.T, dt time.Time) {
	t.Helper()
	_, offset := dt.Zone()
	assert.Equal(
		t,
		0,
		offset,
		"time %v should be in UTC with offset zero but has offset: %v",
		dt,
		offset,
	)
}

func CopyFile(src string, dst string) error {
	// Open original file
	srcF, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer func() { _ = srcF.Close() }()

	// Create new file
	dstF, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("open dst: %w", err)
	}
	defer func() { _ = dstF.Close() }()

	_, err = io.Copy(dstF, srcF)
	if err != nil {
		return fmt.Errorf("copying file: %w", err)
	}
	return nil
}

// simpleCreateForecastHelper creates a new Forecast with little configurability by setting
// default values (e.g. it resolves in 24 hours). It's useful for simplifying test code.
func simpleCreateForecastHelper(
	t *testing.T,
	client graphql.Client,
	title string,
	probabilities map[string]int,
) *CreateForecastResponse {
	t.Helper()
	newForecast := NewForecast{
		Title:       title,
		Description: "",
		Resolves:    time.Now().Add(24 * time.Hour),
	}
	var newProbabilities []NewProbability
	for outcome, p := range probabilities {
		newProbabilities = append(
			newProbabilities,
			NewProbability{
				Value:   p,
				Outcome: &NewOutcome{Text: outcome},
			})
	}

	newEstimate := NewEstimate{
		Reason:        "Just a hunch.",
		Probabilities: newProbabilities,
	}
	resp, err := CreateForecast(context.Background(), client, newForecast, newEstimate)
	require.NoError(t, err)
	return resp
}
