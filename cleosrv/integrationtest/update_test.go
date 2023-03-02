package integrationtest

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate_From_0_1_1(t *testing.T) {
	tDir, err := os.MkdirTemp("", t.Name()+"_")
	require.NoError(t, err)
	dbSrc := filepath.Join("testdata", "test_0.1.1.db")
	dbPath := filepath.Join(tDir, "test.db")
	t.Log("executing with DB", dbPath)
	err = CopyFile(dbSrc, dbPath)
	require.NoError(t, err)

	c := initServerAndGetClient(t, dbPath)
	respGetForecasts, err := GetForecasts(context.Background(), c)
	require.NoError(t, err)

	expectedTitles := []string{
		"Forecast with closes set to Go time null value and 3 outcomes (0.1.1)",
		"Forecast with illogical created/resolves/closes (0.1.1)",
		"Forecast with closes after resolves (0.1.1)",
		"Just a regular forecast (0.1.1)",
	}
	var foundTitles []string
	for _, f := range respGetForecasts.Forecasts {
		foundTitles = append(foundTitles, f.Title)
		t.Log(f)
		assertUTC(t, f.Created)
		assertUTC(t, f.Resolves)
		if f.Closes != nil {
			assertUTC(t, *f.Closes)
		}
		for _, e := range f.Estimates {
			assertUTC(t, e.Created)
		}

		if f.Title == "Forecast with closes set to Go time null value and 3 outcomes (0.1.1)" {
			// In the 0.1.1 DB closes is set to 0001-01-01 00:00:00+00:00 i.e.
			// the Go time.Time null value
			assert.Nil(t, f.Closes)
		}
		if f.Title == "Forecast with illogical created/resolves/closes (0.1.1)" {
			assert.Equal(t, f.Created, f.Resolves)
			assert.Equal(t, f.Created, *f.Closes)
		}
		if f.Title == "Forecast with closes after resolves (0.1.1)" {
			assert.Equal(t, f.Resolves, *f.Closes)
		}
		if f.Title == "Just a regular forecast (0.1.1)" {
			// Verify some values that we know are contained in the DB to
			// ensure nothing got lost. If the DB is re-generated, this part of
			// the test probably has to be updated.
			expectedCreated, err := time.Parse(
				time.RFC3339,
				"2023-02-25T09:49:59.276264803Z",
			)
			assert.Nil(t, err)
			assert.Equal(t, expectedCreated, f.Created)
			expectedResolves, err := time.Parse(
				time.RFC3339,
				"2023-03-27T09:49:59Z",
			)
			assert.Nil(t, err)
			assert.Equal(t, expectedResolves, f.Resolves)
			assert.Nil(t, f.Closes)
		}
	}
	assert.ElementsMatch(t, expectedTitles, foundTitles)

	// Create a new forecast to verify it's possible after an update

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

	respCreate, err := CreateForecast(
		context.Background(),
		c,
		newForecast,
		newEstimate,
	)
	require.NoError(t, err)
	assert.Equal(t, "Will it rain tomorrow?", respCreate.CreateForecast.Title)

	respGetForecasts, err = GetForecasts(context.Background(), c)
	require.NoError(t, err)
	assert.Len(t, respGetForecasts.Forecasts, 5)
}

// TestUpdate_From_0_1_1_sqlite3_dump verifies that the sqlite3 dump (if the
// tool is installed) only contains time zones of the form +00:00 and does not
// contain the Go null time.
func TestUpdate_From_0_1_1_sqlite3_dump(t *testing.T) {
	_, err := exec.LookPath("sqlite3")
	if err != nil {
		t.Skip("'sqlite3' is not installed")
	}

	tDir, err := os.MkdirTemp("", t.Name()+"_")
	require.NoError(t, err)
	dbSrc := filepath.Join("testdata", "test_0.1.1.db")
	dbPath := filepath.Join(tDir, "test.db")
	t.Log("executing with DB", dbPath)
	err = CopyFile(dbSrc, dbPath)
	require.NoError(t, err)

	c := initServerAndGetClient(t, dbPath)
	resp, err := GetForecasts(context.Background(), c)
	require.NoError(t, err)
	require.Len(t, resp.Forecasts, 4)

	cmd := exec.Command("sqlite3", dbPath, ".dump")
	var cleosrvStdout, cleosrvStderr bytes.Buffer
	cmd.Stdout = &cleosrvStdout
	cmd.Stderr = &cleosrvStderr
	err = cmd.Run()
	require.NoError(t, err, cleosrvStderr.String())
	tzRegex, _ := regexp.Compile(`[+-]\d\d:\d\d`)
	for _, tz := range tzRegex.FindAllString(cleosrvStdout.String(), -1) {
		assert.Equal(t, "+00:00", tz)
	}
	assert.NotContains(t, cleosrvStdout.String(), "0001-01-01 00:00:00+00:00")
}