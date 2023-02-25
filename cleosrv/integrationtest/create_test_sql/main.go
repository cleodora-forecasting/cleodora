// create_test_sql generates a .sql file with the complete DB state of a cleosrv instance.
// This is used for upgrade tests to ensure the database is migrated correctly.
// The tool will start 'cleosrv' (provided as a path),
// execute certain GraphQL operations and then dump the resulting SQL into a file.
// This SQL file can later be used in a test to import it into a DB,
// run the migrations on that DB and verify that the resulting DB state is as expected.
// This script will rarely be needed (probably only after making a release) and all it does could
// be replicated manually.
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/cleodora-forecasting/cleodora/cleosrv/integrationtest"
)

func main() {
	if err := do(); err != nil {
		_, _ = fmt.Fprintf(
			os.Stderr,
			"%v\n",
			err,
		)
		os.Exit(1)
	}
}

func do() error {
	args := os.Args
	usage := `Usage:
	create_test_sql CLEOSRV_BIN RESULT_DB

Where CLEOSRV_BIN is the path to a cleosrv binary and RESULT_DB is the path
where the final DB should be stored.

Example:
	go run cleosrv/integrationtest/create_test_sql/main.go \
		~/software/cleosrv \
		cleosrv/integrationtest/testdata/TestUpdate_From_0_1_1/test.db
`
	if len(args) != 3 {
		return errors.New(usage)
	}
	cleosrv := os.Args[1]
	resultDb := os.Args[2]
	tDir, err := os.MkdirTemp("", "cleosrv")
	if err != nil {
		return err
	}

	dbPath := path.Join(tDir, "test.db")
	cleosrvCmd := exec.Command(cleosrv, "--database", dbPath)
	var cleosrvStdout, cleosrvStderr bytes.Buffer
	cleosrvCmd.Stdout = &cleosrvStdout
	cleosrvCmd.Stderr = &cleosrvStderr
	err = cleosrvCmd.Start()
	if err != nil {
		return err
	}
	defer func() {
		if err := cleosrvCmd.Process.Kill(); err != nil {
			fmt.Printf("error stopping cleosrv: %v\n", err)
		}
		fmt.Println(strings.Repeat("=", 80))
		fmt.Println("cleosrv stdout:")
		fmt.Println(cleosrvStdout.String())
		fmt.Println(strings.Repeat("=", 80))
		fmt.Println("cleosrv stderr:")
		fmt.Println(cleosrvStderr.String())
		fmt.Println(strings.Repeat("=", 80))
	}()

	time.Sleep(10 * time.Second) // give time to start

	c := graphql.NewClient("http://localhost:8080/query", nil)

	if err = executeQueries(c); err != nil {
		return err
	}

	version, err := getVersion(c)
	if err != nil {
		return fmt.Errorf("getting version: %w", err)
	}
	fmt.Println("cleosrv version:", version)

	err = cleosrvCmd.Process.Kill()
	if err != nil {
		return fmt.Errorf("stopping cleosrvCmd: %w", err)
	}

	err = integrationtest.CopyFile(dbPath, resultDb)
	if err != nil {
		return err
	}
	return nil
}

func getVersion(c graphql.Client) (string, error) {
	resp, err := integrationtest.GetMetadata(context.Background(), c)
	if err != nil {
		return "", err
	}
	return resp.Metadata.Version, nil
}

func executeQueries(c graphql.Client) error {
	variables := map[string]interface{}{
		"forecast": map[string]interface{}{
			"description": "",
			"resolves":    time.Now().UTC().Add(30 * 24 * time.Hour).Format(time.RFC3339),
			"title":       "Just a regular forecast (0.1.1)",
		},
		"estimate": map[string]interface{}{
			"probabilities": []map[string]interface{}{
				{
					"outcome": map[string]interface{}{
						"text": "Yes",
					},
					"value": 20,
				},
				{
					"outcome": map[string]interface{}{
						"text": "No",
					},
					"value": 80,
				},
			},
			"reason": "Just a hunch.",
		},
	}

	err := createForecast(c, variables)
	if err != nil {
		return err
	}

	// Create another forecast with 'resolves' in the past and 'closes' after 'resolves'.
	// Also use some strange timezone for them.

	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		return fmt.Errorf("can't get TZ loc: %w", err)
	}

	variables = map[string]interface{}{
		"forecast": map[string]interface{}{
			"description": "",
			"resolves":    time.Now().In(newYork).Add(-30 * 24 * time.Hour).Format(time.RFC3339),
			"closes":      time.Now().In(newYork).Add(-15 * 24 * time.Hour).Format(time.RFC3339),
			"title":       "Forecast with illogical created/resolves/closes (0.1.1)",
		},
		"estimate": map[string]interface{}{
			"probabilities": []map[string]interface{}{
				{
					"outcome": map[string]interface{}{
						"text": "Yes",
					},
					"value": 20,
				},
				{
					"outcome": map[string]interface{}{
						"text": "No",
					},
					"value": 80,
				},
			},
			"reason": "Just a hunch.",
		},
	}

	err = createForecast(c, variables)
	if err != nil {
		return err
	}

	variables = map[string]interface{}{
		"forecast": map[string]interface{}{
			"description": "",
			"resolves":    time.Now().UTC().Add(30 * 24 * time.Hour).Format(time.RFC3339),
			"closes":      time.Now().UTC().Add(40 * 24 * time.Hour).Format(time.RFC3339),
			"title":       "Forecast with closes after resolves (0.1.1)",
		},
		"estimate": map[string]interface{}{
			"probabilities": []map[string]interface{}{
				{
					"outcome": map[string]interface{}{
						"text": "Yes",
					},
					"value": 20,
				},
				{
					"outcome": map[string]interface{}{
						"text": "No",
					},
					"value": 80,
				},
			},
			"reason": "Just a hunch.",
		},
	}

	err = createForecast(c, variables)
	if err != nil {
		return err
	}

	variables = map[string]interface{}{
		"forecast": map[string]interface{}{
			"description": "",
			"resolves":    time.Now().UTC().Add(30 * 24 * time.Hour).Format(time.RFC3339),
			"closes":      time.Time{}.Format(time.RFC3339), // null value
			"title":       "Forecast with closes set to Go time null value and 3 outcomes (0.1.1)",
		},
		"estimate": map[string]interface{}{
			"probabilities": []map[string]interface{}{
				{
					"outcome": map[string]interface{}{
						"text": "Yes",
					},
					"value": 20,
				},
				{
					"outcome": map[string]interface{}{
						"text": "No",
					},
					"value": 30,
				},
				{
					"outcome": map[string]interface{}{
						"text": "Maybe",
					},
					"value": 50,
				},
			},
			"reason": "Just a hunch.",
		},
	}

	err = createForecast(c, variables)
	if err != nil {
		return err
	}

	return nil
}

func createForecast(c graphql.Client, variables map[string]interface{}) error {
	query := `
mutation CreateForecast ($forecast: NewForecast!, $estimate: NewEstimate!) {
	createForecast(forecast: $forecast, estimate: $estimate) {
		id
		title
	}
}`

	var responseData struct {
		CreateForecast struct {
			Id    string
			Title string
		}
	}

	req := graphql.Request{
		Query:     query,
		Variables: variables,
	}
	response := graphql.Response{Data: &responseData}

	err := c.MakeRequest(
		context.Background(),
		&req,
		&response,
	)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}

	if responseData.CreateForecast.Id == "" {
		return fmt.Errorf("unexpected response: %v", responseData)
	}
	return nil
}
