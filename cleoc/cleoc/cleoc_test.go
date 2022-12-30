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
			Probabilities []probability
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

type probability struct {
	Value   int
	Outcome outcome
}

type outcome struct {
	Text string
}

// TestApp_AddForecast_Simple is contained in TestApp_AddForecast_Probabilities
// but serves as documentation how to write a simple test without being table
// driven, which can be sometimes a little hard to read.
func TestApp_AddForecast_Simple(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the body contains the expected request
		body, err := io.ReadAll(r.Body)
		require.Nil(t, err)
		var bodyStruct createForecastBody
		err = json.Unmarshal(body, &bodyStruct)
		require.Nil(t, err)

		assert.Equal(t, "CreateForecast", bodyStruct.OperationName)
		assert.Equal(t, "Will it rain tomorrow?", bodyStruct.Variables.Forecast.Title)
		assert.Equal(t,
			"The weather prediction says so",
			bodyStruct.Variables.Estimate.Reason,
		)
		assert.Len(t, bodyStruct.Variables.Estimate.Probabilities, 2)

		expectedProbabilities := []probability{
			{
				Value: 20,
				Outcome: outcome{
					Text: "Yes",
				},
			},
			{
				Value: 80,
				Outcome: outcome{
					Text: "No",
				},
			},
		}
		assert.ElementsMatch(t, expectedProbabilities, bodyStruct.Variables.Estimate.Probabilities)

		// Send a response
		w.Header().Set("Content-Type", "application/json")
		_, err = fmt.Fprint(
			w,
			"{\"data\":{\"createForecast\":{\"id\":\"999\",\"__typename\":\"Forecast\"}}}",
		)
		require.Nil(t, err)
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

	err := a.AddForecast(
		"Will it rain tomorrow?",
		time.Now().Add(time.Hour*24).Format(time.
			RFC3339),
		"",
		"The weather prediction says so",
		[]string{"Yes:20", "No:80"},
	)
	require.Nil(t, err)
	assert.Equal(t, "999", out.String())
	assert.Empty(t, errOut)
}

// TestApp_AddForecast_Probabilities verifies correct and incorrect
// 'probabilities' parameter during forecast creation.
func TestApp_AddForecast_Probabilities(t *testing.T) {
	tests := []struct {
		name                  string
		inputProbabilities    []string
		expectedProbabilities []probability
		expectedErr           string
	}{
		{
			name:               "20-80",
			inputProbabilities: []string{"Yes:20", "No:80"},
			expectedProbabilities: []probability{
				{
					Value: 20,
					Outcome: outcome{
						Text: "Yes",
					},
				},
				{
					Value: 80,
					Outcome: outcome{
						Text: "No",
					},
				},
			},
			expectedErr: "",
		},
		{
			name:               "80-20 order doesnt matter",
			inputProbabilities: []string{"No:80", "Yes:20"},
			expectedProbabilities: []probability{
				{
					Value: 20,
					Outcome: outcome{
						Text: "Yes",
					},
				},
				{
					Value: 80,
					Outcome: outcome{
						Text: "No",
					},
				},
			},
			expectedErr: "",
		},
		{
			name: "multi word outcome",
			inputProbabilities: []string{
				"Yes, more than 1 hour:70",
				"Yes, less than 1 hour:20",
				"No:10",
			},
			expectedProbabilities: []probability{
				{
					Value: 70,
					Outcome: outcome{
						Text: "Yes, more than 1 hour",
					},
				},
				{
					Value: 20,
					Outcome: outcome{
						Text: "Yes, less than 1 hour",
					},
				},
				{
					Value: 10,
					Outcome: outcome{
						Text: "No",
					},
				},
			},
			expectedErr: "",
		},
		{
			name: "multi word outcome with colon",
			inputProbabilities: []string{
				"Yes: More than: 1 hour:70",
				"Yes: Less than: 1 hour:20",
				"No:10",
			},
			expectedProbabilities: []probability{
				{
					Value: 70,
					Outcome: outcome{
						Text: "Yes: More than: 1 hour",
					},
				},
				{
					Value: 20,
					Outcome: outcome{
						Text: "Yes: Less than: 1 hour",
					},
				},
				{
					Value: 10,
					Outcome: outcome{
						Text: "No",
					},
				},
			},
			expectedErr: "",
		},
		{
			name: "invalid not a number",
			inputProbabilities: []string{
				"Yes:Tree",
				"No:10",
			},
			expectedProbabilities: nil,
			expectedErr: "error parsing probabilities: 'Yes:Tree'" +
				" the probability is not a valid number. Use OUTCOME:PROBABILITY",
		},
		{
			name: "invalid no colon",
			inputProbabilities: []string{
				"Yes",
				"No:10",
			},
			expectedProbabilities: nil,
			expectedErr:           "error parsing probabilities: 'Yes' must contain ':'",
		},
		{
			name: "invalid no outcome",
			inputProbabilities: []string{
				":30",
				"No:10",
			},
			expectedProbabilities: nil,
			expectedErr: "error parsing probabilities: ':30' the " +
				"outcome can't be empty. Use OUTCOME:PROBABILITY",
		},
		{
			name: "invalid no probability",
			inputProbabilities: []string{
				"Yes:",
				"No:10",
			},
			expectedProbabilities: nil,
			expectedErr: "error parsing probabilities: 'Yes:' the " +
				"probability can't be empty. Use OUTCOME:PROBABILITY",
		},
		{
			name:                  "invalid no probabilities",
			inputProbabilities:    []string{},
			expectedProbabilities: nil,
			expectedErr:           "error parsing probabilities: no probabilities",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectedErr != "" {
					t.Fatalf("If we expected an error the server should never be called")
				}
				// Verify the body contains the expected request
				body, err := io.ReadAll(r.Body)
				require.Nil(t, err)
				var bodyStruct createForecastBody
				err = json.Unmarshal(body, &bodyStruct)
				require.Nil(t, err)

				assert.Equal(t, "CreateForecast", bodyStruct.OperationName)
				assert.Equal(t, "Will it rain tomorrow?", bodyStruct.Variables.Forecast.Title)
				assert.Equal(t,
					"The weather prediction says so",
					bodyStruct.Variables.Estimate.Reason,
				)
				assert.ElementsMatch(t, tt.expectedProbabilities, bodyStruct.Variables.Estimate.Probabilities)

				// Send a response
				w.Header().Set("Content-Type", "application/json")
				_, err = fmt.Fprint(
					w,
					"{\"data\":{\"createForecast\":{\"id\":\"999\",\"__typename\":\"Forecast\"}}}",
				)
				require.Nil(t, err)
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

			err := a.AddForecast(
				"Will it rain tomorrow?",
				time.Now().Add(time.Hour*24).Format(time.
					RFC3339),
				"",
				"The weather prediction says so",
				tt.inputProbabilities,
			)
			if tt.expectedErr == "" {
				require.Nil(t, err)
				assert.Equal(t, "999", out.String())
				assert.Empty(t, errOut)
			} else {
				assert.Error(t, err)
				assert.Empty(t, out.String())
				assert.Equal(t, tt.expectedErr, err.Error())
			}
		})
	}
}

func TestApp_Version(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	config := &cleoc.Config{
		URL:        "",
		ConfigFile: "",
	}
	a := &cleoc.App{
		Out:    out,
		Err:    errOut,
		Config: config,
	}
	err := a.Version()
	assert.Nil(t, err)
	assert.Equal(t, "dev", out.String())
	assert.Empty(t, errOut)
}

// TODO Add test that verifies error server response with the forecasts creation
