package cleoc_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cleodora-forecasting/cleodora/cleoc/cleoc"
)

type createForecastBody struct {
	OperationName string
	Query         string
	Variables     struct {
		Estimate struct {
			Probabilities []newProbability
			Reason        string
		}
		Forecast struct {
			Title       string
			Description string
			Closes      string
			Resolves    string
		}
	}
}

type newProbability struct {
	Value   int
	Outcome newOutcome
}

type newOutcome struct {
	Text *string
}

// TestApp_AddForecast_Simple is contained in TestApp_AddForecast_Probabilities
// but serves as documentation how to write a simple test without being table
// driven, which can be sometimes a little hard to read.
func TestApp_AddForecast_Simple(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the body contains the expected request
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		var bodyStruct createForecastBody
		err = json.Unmarshal(body, &bodyStruct)
		require.NoError(t, err)

		assert.Equal(t, "CreateForecast", bodyStruct.OperationName)
		assert.Equal(t, "Will it rain tomorrow?", bodyStruct.Variables.Forecast.Title)
		assert.Equal(t,
			"The weather prediction says so",
			bodyStruct.Variables.Estimate.Reason,
		)
		assert.Len(t, bodyStruct.Variables.Estimate.Probabilities, 2)

		expectedProbabilities := []newProbability{
			{
				Value: 20,
				Outcome: newOutcome{
					Text: strPtr("Yes"),
				},
			},
			{
				Value: 80,
				Outcome: newOutcome{
					Text: strPtr("No"),
				},
			},
		}
		assertNewProbabilitiesMatch(t, expectedProbabilities, bodyStruct.Variables.Estimate.Probabilities)

		// Send a response
		w.Header().Set("Content-Type", "application/json")
		_, err = fmt.Fprint(
			w,
			"{\"data\":{\"createForecast\":{\"id\":\"999\",\"__typename\":\"Forecast\"}}}",
		)
		require.NoError(t, err)
	}))
	defer ts.Close()

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	config := &cleoc.Config{
		URL:        ts.URL,
		ConfigFile: "",
	}
	a := &cleoc.App{
		Out:    out,
		Err:    errOut,
		Config: config,
	}

	opts := cleoc.AddForecastOptions{
		Title:       "Will it rain tomorrow?",
		Description: "",
		Resolves:    time.Now().Add(time.Hour * 24).Format(time.RFC3339),
		Reason:      "The weather prediction says so",
		Probabilities: map[string]int{
			"Yes": 20,
			"No":  80,
		},
	}
	// Of course, we know the options in this test to be valid,
	// but for documentation purposes make it clear that validation is
	// expected before calling AddForecast
	err := opts.Validate()
	require.NoError(t, err)

	err = a.AddForecast(opts)
	require.NoError(t, err)
	assert.Equal(t, "999\n", out.String())
	assert.Empty(t, errOut)
}

func TestApp_AddForecast_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the body contains the expected request
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		var bodyStruct createForecastBody
		err = json.Unmarshal(body, &bodyStruct)
		require.NoError(t, err)

		assert.Equal(t, "CreateForecast", bodyStruct.OperationName)
		assert.Equal(t, "Will it rain tomorrow?", bodyStruct.Variables.Forecast.Title)
		assert.Equal(t,
			"The weather prediction says so",
			bodyStruct.Variables.Estimate.Reason,
		)
		assert.Len(t, bodyStruct.Variables.Estimate.Probabilities, 2)

		expectedProbabilities := []newProbability{
			{
				Value: 20,
				Outcome: newOutcome{
					Text: strPtr("Yes"),
				},
			},
			{
				Value: 80,
				Outcome: newOutcome{
					Text: strPtr("No"),
				},
			},
		}
		assertNewProbabilitiesMatch(t, expectedProbabilities, bodyStruct.Variables.Estimate.Probabilities)

		// Send a response, copied from an actual response due to an erroneous
		// request where 'title' was missing
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = fmt.Fprint(
			w,
			"{\"errors\":[{\"message\":\"json request body could not be "+
				"decoded: invalid character 'd' looking for beginning of "+
				"object key string body:{\\\"operationName\\\":\\\""+
				"createForecast\\\",\\\"variables\\\":{\\\"forecast\\\":{"+
				"description\\\":\\\"asdf\\\",\\\"closes\\\":\\\""+
				"2022-12-31T11:28:52.431Z\\\",\\\"resolves\\\":\\\""+
				"2022-12-31T11:28:52.431Z\\\"},\\\"estimate\\\":{\\\""+
				"reason\\\":\\\"asdf\\\",\\\"probabilities\\\":[{\\\""+
				"value\\\":50,\\\"outcome\\\":{\\\"text\\\":\\\""+
				"Yes\\\"}},{\\\"value\\\":50,\\\"outcome\\\":{"+
				"\\\"text\\\":\\\"No\\\"}}]}},\\\"query\\\":"+
				"\\\"mutation createForecast($forecast: NewForecast!, "+
				"$estimate: NewEstimate!) {\\\\n  createForecast(forecast: "+
				"$forecast, estimate: $estimate) {\\\\n    id\\\\n    "+
				"title\\\\n    __typename\\\\n  }\\\\n}\\\"}\"}],\"data\":"+
				"null}",
		)
		require.NoError(t, err)
	}))
	defer ts.Close()

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	config := &cleoc.Config{
		URL:        ts.URL,
		ConfigFile: "",
	}
	a := &cleoc.App{
		Out:    out,
		Err:    errOut,
		Config: config,
	}

	opts := cleoc.AddForecastOptions{
		Title:       "Will it rain tomorrow?",
		Description: "",
		Resolves:    time.Now().Add(time.Hour * 24).Format(time.RFC3339),
		Reason:      "The weather prediction says so",
		Probabilities: map[string]int{
			"Yes": 20,
			"No":  80,
		},
	}
	// Of course, we know the options in this test to be valid,
	// but for documentation purposes make it clear that validation is
	// expected before calling AddForecast
	err := opts.Validate()
	require.NoError(t, err)

	err = a.AddForecast(opts)
	assert.ErrorContains(t, err, "400 Bad Request")
	assert.Empty(t, out)
	assert.Empty(t, errOut)
}

// TestApp_AddForecast_Probabilities verifies correct 'probabilities'
// parameter during forecast creation.
func TestApp_AddForecast_Probabilities(t *testing.T) {
	tests := []struct {
		name                  string
		inputProbabilities    map[string]int
		expectedProbabilities []newProbability
	}{
		{
			name:               "20-80",
			inputProbabilities: map[string]int{"Yes": 20, "No": 80},
			expectedProbabilities: []newProbability{
				{
					Value: 20,
					Outcome: newOutcome{
						Text: strPtr("Yes"),
					},
				},
				{
					Value: 80,
					Outcome: newOutcome{
						Text: strPtr("No"),
					},
				},
			},
		},
		{
			name: "multi word outcome",
			inputProbabilities: map[string]int{
				"Yes, more than 1 hour": 70,
				"Yes, less than 1 hour": 20,
				"No":                    10,
			},
			expectedProbabilities: []newProbability{
				{
					Value: 70,
					Outcome: newOutcome{
						Text: strPtr("Yes, more than 1 hour"),
					},
				},
				{
					Value: 20,
					Outcome: newOutcome{
						Text: strPtr("Yes, less than 1 hour"),
					},
				},
				{
					Value: 10,
					Outcome: newOutcome{
						Text: strPtr("No"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the body contains the expected request
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				var bodyStruct createForecastBody
				err = json.Unmarshal(body, &bodyStruct)
				require.NoError(t, err)

				assert.Equal(t, "CreateForecast", bodyStruct.OperationName)
				assert.Equal(t, "Will it rain tomorrow?", bodyStruct.Variables.Forecast.Title)
				assert.Equal(t,
					"The weather prediction says so",
					bodyStruct.Variables.Estimate.Reason,
				)
				assertNewProbabilitiesMatch(
					t,
					tt.expectedProbabilities,
					bodyStruct.Variables.Estimate.Probabilities,
				)

				// Send a response
				w.Header().Set("Content-Type", "application/json")
				_, err = fmt.Fprint(
					w,
					"{\"data\":{\"createForecast\":{\"id\":\"999\",\"__typename\":\"Forecast\"}}}",
				)
				require.NoError(t, err)
			}))
			defer server.Close()

			out := &bytes.Buffer{}
			errOut := &bytes.Buffer{}
			a := &cleoc.App{
				Out: out,
				Err: errOut,
				Config: &cleoc.Config{
					URL:        server.URL,
					ConfigFile: "",
				},
			}

			opts := cleoc.AddForecastOptions{
				Title:         "Will it rain tomorrow?",
				Description:   "",
				Resolves:      time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				Reason:        "The weather prediction says so",
				Probabilities: tt.inputProbabilities,
			}
			// Of course, we know the options in this test to be valid,
			// but for documentation purposes make it clear that validation is
			// expected before calling AddForecast
			err := opts.Validate()
			require.NoError(t, err)

			err = a.AddForecast(opts)
			require.NoError(t, err)
			assert.Equal(t, "999\n", out.String())
			assert.Empty(t, errOut)
		})
	}
}

func TestAddForecastOptions_Validate(t *testing.T) {
	type fields struct {
		Title         string
		Description   string
		Resolves      string
		Closes        string
		Reason        string
		Probabilities map[string]int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				Title:       "Will it rain?",
				Description: "Yes if it rains for more than 30 minutes.",
				Resolves:    "2020-10-01T13:13:00+02:00",
				Closes:      "2020-10-01T13:13:00+02:00",
				Reason:      "The weather forecast says so",
				Probabilities: map[string]int{
					"Yes": 20,
					"No":  80,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "title cant be empty",
			fields: fields{
				Title:       "",
				Description: "Yes if it rains for more than 30 minutes.",
				Resolves:    "2020-10-01T13:13:00+02:00",
				Closes:      "2020-10-01T13:13:00+02:00",
				Reason:      "The weather forecast says so",
				Probabilities: map[string]int{
					"Yes": 20,
					"No":  80,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "description",
			fields: fields{
				Title:       "Will it rain?",
				Description: "",
				Resolves:    "2020-10-01T13:13:00+02:00",
				Closes:      "2020-10-01T13:13:00+02:00",
				Reason:      "The weather forecast says so",
				Probabilities: map[string]int{
					"Yes": 20,
					"No":  80,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "resolves cant be empty",
			fields: fields{
				Title:       "Will it rain?",
				Description: "Yes if it rains for more than 30 minutes.",
				Resolves:    "",
				Closes:      "2020-10-01T13:13:00+02:00",
				Reason:      "The weather forecast says so",
				Probabilities: map[string]int{
					"Yes": 20,
					"No":  80,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "resolves must be formatted properly",
			fields: fields{
				Title:       "Will it rain?",
				Description: "Yes if it rains for more than 30 minutes.",
				Resolves:    "2020-10-01",
				Closes:      "2020-10-01T13:13:00+02:00",
				Reason:      "The weather forecast says so",
				Probabilities: map[string]int{
					"Yes": 20,
					"No":  80,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "reason cant be empty",
			fields: fields{
				Title:       "Will it rain?",
				Description: "Yes if it rains for more than 30 minutes.",
				Resolves:    "2020-10-01T13:13:00+02:00",
				Closes:      "2020-10-01T13:13:00+02:00",
				Reason:      "",
				Probabilities: map[string]int{
					"Yes": 20,
					"No":  80,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "probabilities cant be empty",
			fields: fields{
				Title:         "Will it rain?",
				Description:   "Yes if it rains for more than 30 minutes.",
				Resolves:      "2020-10-01T13:13:00+02:00",
				Closes:        "2020-10-01T13:13:00+02:00",
				Reason:        "The weather forecast says so",
				Probabilities: map[string]int{},
			},
			wantErr: assert.Error,
		},
		{
			name: "probabilities cant be empty 2",
			fields: fields{
				Title:         "Will it rain?",
				Description:   "Yes if it rains for more than 30 minutes.",
				Resolves:      "2020-10-01T13:13:00+02:00",
				Closes:        "2020-10-01T13:13:00+02:00",
				Reason:        "The weather forecast says so",
				Probabilities: nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "probabilities must be between 0 and 100",
			fields: fields{
				Title:       "Will it rain?",
				Description: "Yes if it rains for more than 30 minutes.",
				Resolves:    "2020-10-01T13:13:00+02:00",
				Closes:      "2020-10-01T13:13:00+02:00",
				Reason:      "The weather forecast says so",
				Probabilities: map[string]int{
					"Yes": 110,
					"No":  -10,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "probabilities must add up to 100",
			fields: fields{
				Title:       "Will it rain?",
				Description: "Yes if it rains for more than 30 minutes.",
				Resolves:    "2020-10-01T13:13:00+02:00",
				Closes:      "2020-10-01T13:13:00+02:00",
				Reason:      "The weather forecast says so",
				Probabilities: map[string]int{
					"Yes": 10,
					"No":  10,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "probabilities outcome string cant be empty",
			fields: fields{
				Title:       "Will it rain?",
				Description: "Yes if it rains for more than 30 minutes.",
				Resolves:    "2020-10-01T13:13:00+02:00",
				Closes:      "2020-10-01T13:13:00+02:00",
				Reason:      "The weather forecast says so",
				Probabilities: map[string]int{
					"":   90,
					"No": 10,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "closes can be empty",
			fields: fields{
				Title:       "Will it rain?",
				Description: "Yes if it rains for more than 30 minutes.",
				Resolves:    "2020-10-01T13:13:00+02:00",
				Closes:      "",
				Reason:      "The weather forecast says so",
				Probabilities: map[string]int{
					"Yes": 20,
					"No":  80,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {
			opts := &cleoc.AddForecastOptions{
				Title:         tt.fields.Title,
				Description:   tt.fields.Description,
				Resolves:      tt.fields.Resolves,
				Closes:        tt.fields.Closes,
				Reason:        tt.fields.Reason,
				Probabilities: tt.fields.Probabilities,
			}
			tt.wantErr(t, opts.Validate())
		})
	}
}

// strPtr returns the pointer of a string
func strPtr(s string) *string {
	return &s
}

func assertNewProbabilitiesMatch(
	t *testing.T,
	expectedProbs []newProbability,
	actualProbs []newProbability,
) {
	t.Helper()
	if len(expectedProbs) != len(actualProbs) {
		assert.Fail(
			t,
			"lengths don't match",
			"expectedProbs %v, actualProbs %v",
			len(expectedProbs),
			len(actualProbs),
		)
		return
	}
	foundOutcomes := map[string]int{}
	for _, p := range expectedProbs {
		if p.Outcome.Text != nil {
			foundOutcomes[*p.Outcome.Text] = 0
		}
	}
	if len(expectedProbs) != len(foundOutcomes) {
		assert.Fail(
			t,
			"length of found outcomes does not match",
			"expectedProbs %v, outcomes %v",
			len(expectedProbs),
			len(foundOutcomes),
		)
		return
	}

	for _, exP := range expectedProbs {
		for _, acP := range actualProbs {
			if exP.Outcome.Text != nil &&
				acP.Outcome.Text != nil &&
				*exP.Outcome.Text == *acP.Outcome.Text &&
				exP.Value == acP.Value {
				foundOutcomes[*exP.Outcome.Text] += 1
			}
		}
	}

	for outcome, count := range foundOutcomes {
		assert.Equal(t,
			1,
			count,
			"outcome %v should have been found once but was found %v times",
			outcome,
			count,
		)
	}
}
