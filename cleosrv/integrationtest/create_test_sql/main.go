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

	version, err := getVersion(c)
	if err != nil {
		return fmt.Errorf("getting version: %w", err)
	}
	fmt.Println("cleosrv version:", version)

	if err = executeQueries(c, version); err != nil {
		return err
	}

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

func executeQueries(c graphql.Client, version string) error {
	now := time.Now().UTC()

	f1 := integrationtest.NewForecast{
		Title:       fmt.Sprintf("Just a regular forecast (%v)", version),
		Description: "",
		Resolves:    now.Add(30 * 24 * time.Hour),
	}
	e1 := integrationtest.NewEstimate{
		Reason: "Just a hunch.",
		Probabilities: []integrationtest.NewProbability{
			{
				Value:   20,
				Outcome: integrationtest.NewOutcome{Text: "Yes"},
			},
			{
				Value:   80,
				Outcome: integrationtest.NewOutcome{Text: "No"},
			},
		},
	}
	resp, err := integrationtest.CreateForecast(context.Background(), c, f1, e1)
	if err != nil {
		return fmt.Errorf("create f1: %w", err)
	}
	if resp.CreateForecast.Id == "" {
		return fmt.Errorf("unexpected response f1: %v", resp)
	}

	// Create another forecast with 'created', 'resolves' and 'closes' in the
	// past.
	// Also use some strange timezone for them.

	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		return fmt.Errorf("can't get TZ loc: %w", err)
	}

	f2 := integrationtest.NewForecast{
		Title: fmt.Sprintf(
			"Forecast with created/resolves/closes in the past (%v)",
			version,
		),
		Description: "",
		Created:     timePointer(now.Add(-30 * 24 * time.Hour).In(newYork)),
		Closes:      timePointer(now.Add(-20 * 24 * time.Hour).In(newYork)),
		Resolves:    now.Add(-10 * 24 * time.Hour).In(newYork),
	}
	e2 := integrationtest.NewEstimate{
		Reason: "Just a hunch.",
		Probabilities: []integrationtest.NewProbability{
			{
				Value:   20,
				Outcome: integrationtest.NewOutcome{Text: "Yes"},
			},
			{
				Value:   80,
				Outcome: integrationtest.NewOutcome{Text: "No"},
			},
		},
	}
	resp, err = integrationtest.CreateForecast(context.Background(), c, f2, e2)
	if err != nil {
		return fmt.Errorf("create f2: %w", err)
	}
	if resp.CreateForecast.Id == "" {
		return fmt.Errorf("unexpected response f2: %v", resp)
	}

	f3 := integrationtest.NewForecast{
		Title: fmt.Sprintf(
			"Forecast with closes set to Go time null value and 3 outcomes (%v)",
			version,
		),
		Description: "",
		Closes:      timePointer(time.Time{}),
		Resolves:    time.Now().UTC().Add(30 * 24 * time.Hour),
	}
	e3 := integrationtest.NewEstimate{
		Reason: "Just a hunch.",
		Probabilities: []integrationtest.NewProbability{
			{
				Value:   20,
				Outcome: integrationtest.NewOutcome{Text: "Yes"},
			},
			{
				Value:   30,
				Outcome: integrationtest.NewOutcome{Text: "No"},
			},
			{
				Value:   50,
				Outcome: integrationtest.NewOutcome{Text: "Maybe"},
			},
		},
	}
	resp, err = integrationtest.CreateForecast(context.Background(), c, f3, e3)
	if err != nil {
		return fmt.Errorf("create f3: %w", err)
	}
	if resp.CreateForecast.Id == "" {
		return fmt.Errorf("unexpected response f3: %v", resp)
	}

	return nil
}

func timePointer(t time.Time) *time.Time {
	return &t
}
