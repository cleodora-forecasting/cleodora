package integrationtest

import (
	"context"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateForecast(t *testing.T) {
	c := initServerAndGetClient2(t)

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

	response, err := CreateForecast(
		context.Background(),
		c,
		newForecast,
		newEstimate,
	)
	require.NoError(t, err)

	assert.NotEmpty(t, response.CreateForecast.Id)
	assert.Equal(
		t,
		"Will it rain tomorrow?",
		response.CreateForecast.Title,
	)

	require.Len(t, response.CreateForecast.Estimates, 1)
	assert.NotEmpty(t, response.CreateForecast.Estimates[0].Id)
	assert.Equal(
		t,
		"My weather app says it will rain.",
		response.CreateForecast.Estimates[0].Reason,
	)
	require.Len(t, response.CreateForecast.Estimates[0].Probabilities, 2)
	assert.NotEmpty(t, response.CreateForecast.Estimates[0].Probabilities[0].Id)
	assert.NotEmpty(t, response.CreateForecast.Estimates[0].Probabilities[1].Id)
	assert.False(t, response.CreateForecast.Estimates[0].Probabilities[0].Outcome.Correct)
	assert.False(t, response.CreateForecast.Estimates[0].Probabilities[1].Outcome.Correct)

	// If the order is Yes, No ...
	if response.CreateForecast.Estimates[0].Probabilities[0].Outcome.Text == "Yes" {
		assert.Equal(t, 70, response.CreateForecast.Estimates[0].Probabilities[0].Value)
		assert.Equal(t, "Yes", response.CreateForecast.Estimates[0].Probabilities[0].Outcome.Text)
		assert.Equal(t, 30, response.CreateForecast.Estimates[0].Probabilities[1].Value)
		assert.Equal(t, "No", response.CreateForecast.Estimates[0].Probabilities[1].Outcome.Text)
	} else { // ... or if it is No, Yes
		assert.Equal(t, 30, response.CreateForecast.Estimates[0].Probabilities[0].Value)
		assert.Equal(t, "No", response.CreateForecast.Estimates[0].Probabilities[0].Outcome.Text)
		assert.Equal(t, 70, response.CreateForecast.Estimates[0].Probabilities[1].Value)
		assert.Equal(t, "Yes", response.CreateForecast.Estimates[0].Probabilities[1].Outcome.Text)
	}
}

// TestCreateForecast_XSS verifies that HTML is correctly escaped (to
// prevent XSS attacks).
func TestCreateForecast_XSS(t *testing.T) {
	c := initServerAndGetClient(t)

	attack := "<script>alert(document.cookie)</script>"

	newForecast := map[string]interface{}{
		"title": "Will it rain tomorrow?" + attack,
		"description": "It counts as rain if between 9am and 9pm there are " +
			"30 min or more of uninterrupted precipitation." + attack,
		"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	newEstimate := map[string]interface{}{
		"reason": "My weather app says it will rain." + attack,
		"probabilities": []map[string]interface{}{
			{
				"value": 70,
				"outcome": map[string]interface{}{
					"text": "Yes" + attack,
				},
			},
			{
				"value": 30,
				"outcome": map[string]interface{}{
					"text": "No" + attack,
				},
			},
		},
	}

	query := `
		mutation createForecast($forecast: NewForecast!, $estimate: NewEstimate!) {
			createForecast(forecast: $forecast, estimate: $estimate) {
				id
				title
                description
				estimates {
					id
					created
					reason
					probabilities {
						id
						value
						outcome {
							id
							text
							correct
						}
					}
				}
			}
		}`

	var response struct {
		CreateForecast struct {
			Id          string
			Title       string
			Description string
			Estimates   []struct {
				Id            string
				Created       string
				Reason        string
				Probabilities []struct {
					Id      string
					Value   int
					Outcome struct {
						Id      string
						Text    string
						Correct bool
					}
				}
			}
		}
	}

	err := c.Post(
		query,
		&response,
		client.Var("forecast", newForecast),
		client.Var("estimate", newEstimate),
	)
	require.NoError(t, err)

	assert.NotContains(t, response.CreateForecast.Title, attack)
	assert.Equal(
		t,
		"Will it rain tomorrow?&lt;script&gt;alert(document.cookie)&lt;/script&gt;",
		response.CreateForecast.Title,
	)
	assert.NotContains(t, response.CreateForecast.Description, attack)
	for _, e := range response.CreateForecast.Estimates {
		assert.NotContains(t, e.Reason, attack)
		for _, p := range e.Probabilities {
			assert.NotContains(t, p.Outcome.Text, attack)
		}
	}
}

// TestCreateForecast_ValidateNewEstimate verifies the expected error or no
// error with different NewEstimate values.
func TestCreateForecast_ValidateNewEstimate(t *testing.T) {
	tests := []struct {
		name        string
		newEstimate map[string]interface{}
		expectedErr string
	}{
		{
			name: "success",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "Yes",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "",
		},
		{
			name: "single probability",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 100,
						"outcome": map[string]interface{}{
							"text": "Any outcome",
						},
					},
				},
			},
			expectedErr: "",
		},
		{
			name: "success with more probabilities",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "Yes, but less than 2 hours",
						},
					},
					{
						"value": 20,
						"outcome": map[string]interface{}{
							"text": "Yes, between 2 and 5 hours",
						},
					},
					{
						"value": 20,
						"outcome": map[string]interface{}{
							"text": "Yes, more than 5 hours",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "",
		},
		{
			name: "reason cant be empty",
			newEstimate: map[string]interface{}{
				"reason": "",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "Yes",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "'reason' can't be empty",
		},
		{
			name: "probabilities must add up to 100",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain.",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "Yes",
						},
					},
					{
						"value": 20,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "probabilities must add up to 100",
		},
		{
			name: "probabilities must be between 0 and 100",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain.",
				"probabilities": []map[string]interface{}{
					{
						"value": -10,
						"outcome": map[string]interface{}{
							"text": "Yes",
						},
					},
					{
						"value": 110,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "probabilities must be between 0 and 100",
		},
		{
			name: "probabilities cant be empty",
			newEstimate: map[string]interface{}{
				"reason":        "My weather app says it will rain.",
				"probabilities": []map[string]interface{}{},
			},
			expectedErr: "probabilities can't be empty",
		},
		{
			name: "outcome must be set",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value":   70,
						"outcome": nil,
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "cannot be null",
		},
		{
			name: "outcome cant be empty",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value":   70,
						"outcome": map[string]interface{}{},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "must be defined",
		},
		{
			name: "outcome text cant be empty",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "outcome text can't be empty",
		},
		{
			name: "outcomes cant be duplicates",
			newEstimate: map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			},
			expectedErr: "outcome 'No' is a duplicate",
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {

			c := initServerAndGetClient(t)

			newForecast := map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			}

			query := `
		mutation createForecast($forecast: NewForecast!, $estimate: NewEstimate!) {
			createForecast(forecast: $forecast, estimate: $estimate) {
				id
				title
			}
		}`

			var response struct {
				CreateForecast struct {
					Id    string
					Title string
				}
			}

			err := c.Post(
				query,
				&response,
				client.Var("forecast", newForecast),
				client.Var("estimate", tt.newEstimate),
			)
			if tt.expectedErr == "" { // success
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.expectedErr)
			}
		})
	}
}

// TestCreateForecast_ValidateNewForecast tests forecast creation with
// different input values, some leading to errors, others not.
func TestCreateForecast_ValidateNewForecast(t *testing.T) {
	tests := []struct {
		name        string
		newForecast map[string]interface{}
		expectedErr string
	}{
		{
			name: "success",
			newForecast: map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "",
		},
		{
			name: "title cant be empty",
			newForecast: map[string]interface{}{
				"title": "",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "title can't be empty",
		},
		{
			name: "description can be empty",
			newForecast: map[string]interface{}{
				"title":       "Will it rain tomorrow?",
				"description": "",
				"closes":      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"resolves":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "",
		},
		{
			name: "closes can be empty",
			newForecast: map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes":   nil,
				"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "",
		},
		{
			name: "closes can be omitted",
			newForecast: map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "",
		},
		{
			name: "resolves cant be empty",
			newForecast: map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"resolves": nil,
			},
			expectedErr: "cannot be null",
		},
		{
			name: "resolves cant be omitted",
			newForecast: map[string]interface{}{
				"title": "Will it rain tomorrow?",
				"description": "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				"closes": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedErr: "must be defined",
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {

			c := initServerAndGetClient(t)

			newEstimate := map[string]interface{}{
				"reason": "My weather app says it will rain",
				"probabilities": []map[string]interface{}{
					{
						"value": 70,
						"outcome": map[string]interface{}{
							"text": "Yes",
						},
					},
					{
						"value": 30,
						"outcome": map[string]interface{}{
							"text": "No",
						},
					},
				},
			}

			query := `
		mutation createForecast($forecast: NewForecast!, $estimate: NewEstimate!) {
			createForecast(forecast: $forecast, estimate: $estimate) {
				id
				title
			}
		}`

			var response struct {
				CreateForecast struct {
					Id    string
					Title string
				}
			}

			err := c.Post(
				query,
				&response,
				client.Var("forecast", tt.newForecast),
				client.Var("estimate", newEstimate),
			)
			if tt.expectedErr == "" { // success
				require.NoError(t, err)
				assert.NotEmpty(t, response.CreateForecast.Id)
				assert.Equal(t, tt.newForecast["title"], response.CreateForecast.Title)
			} else {
				assert.ErrorContains(t, err, tt.expectedErr)
			}
		})
	}
}

func TestCreateForecast_FailsWithoutEstimate(t *testing.T) {
	c := initServerAndGetClient(t)

	newForecast := map[string]interface{}{
		"title": "Will it rain tomorrow?",
		"description": "It counts as rain if between 9am and 9pm there are " +
			"30 min or more of uninterrupted precipitation.",
		"closes":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"resolves": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	query := `
		mutation createForecast($forecast: NewForecast!) {
			createForecast(forecast: $forecast) {
				id
				title
			}
		}`

	var response struct {
		CreateForecast struct {
			Id    string
			Title string
		}
	}

	// Note that the client 'c' never returns the JSON errors in response but
	// stuffs it all into a string in 'err' and that's it. Otherwise, it would
	// be nicer to more precisely check the error result.
	err := c.Post(query, &response, client.Var("forecast", newForecast))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "http 422")
	assert.Contains(
		t,
		err.Error(),
		"argument \\\"estimate\\\" of type \\\"NewEstimate!\\\" is required",
	)
}

// TestCreateForecast_WithTimestamps verifies which different combinations of
// the input timestamps (Created, Resolves, Closes) are valid and which are
// not.
func TestCreateForecast_WithTimestamps(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name             string
		inputCreated     *time.Time
		inputResolves    *time.Time
		inputCloses      *time.Time
		expectedCreated  *time.Time
		expectedResolves *time.Time
		expectedCloses   *time.Time
		expectedErr      string
	}{
		{
			name:             "created in the past succeeds",
			inputCreated:     timePointer(now.Add(-1 * 3 * 30 * 24 * time.Hour)), // 3 months ago
			inputResolves:    timePointer(now.Add(24 * time.Hour)),
			inputCloses:      nil,
			expectedCreated:  timePointer(now.Add(-1 * 3 * 30 * 24 * time.Hour)),
			expectedResolves: timePointer(now.Add(24 * time.Hour)),
			expectedCloses:   nil,
			expectedErr:      "",
		},
		{
			name:             "created in the future fails",
			inputCreated:     timePointer(now.Add(3 * 30 * 24 * time.Hour)), // in 3 months
			inputResolves:    timePointer(now.Add(24 * time.Hour)),
			inputCloses:      nil,
			expectedCreated:  nil,
			expectedResolves: nil,
			expectedCloses:   nil,
			expectedErr:      "'created' can't be in the future",
		},
		{
			name:             "created is set to now if empty",
			inputCreated:     nil,
			inputResolves:    timePointer(now.Add(24 * time.Hour)),
			inputCloses:      nil,
			expectedCreated:  timePointer(now),
			expectedResolves: timePointer(now.Add(24 * time.Hour)),
			expectedCloses:   nil,
			expectedErr:      "",
		},
		{
			name:             "closes can be set",
			inputCreated:     timePointer(now),
			inputResolves:    timePointer(now.Add(24 * time.Hour)),
			inputCloses:      timePointer(now.Add(2 * time.Hour)),
			expectedCreated:  timePointer(now),
			expectedResolves: timePointer(now.Add(24 * time.Hour)),
			expectedCloses:   timePointer(now.Add(2 * time.Hour)),
			expectedErr:      "",
		},
		{
			name:             "closes can't be later than resolves",
			inputCreated:     timePointer(now),
			inputResolves:    timePointer(now.Add(24 * time.Hour)),
			inputCloses:      timePointer(now.Add(25 * time.Hour)),
			expectedCreated:  nil,
			expectedResolves: nil,
			expectedCloses:   nil,
			expectedErr:      "'Closes' can't be set to a later date than 'Resolves'",
		},
		{
			name:             "closes and resolves can be the same",
			inputCreated:     timePointer(now),
			inputResolves:    timePointer(now.Add(24 * time.Hour)),
			inputCloses:      timePointer(now.Add(24 * time.Hour)),
			expectedCreated:  timePointer(now),
			expectedResolves: timePointer(now.Add(24 * time.Hour)),
			expectedCloses:   timePointer(now.Add(24 * time.Hour)),
			expectedErr:      "",
		},
		{
			name:             "resolves can't be earlier than created",
			inputCreated:     timePointer(now),
			inputResolves:    timePointer(now.Add(-24 * time.Hour)),
			inputCloses:      nil,
			expectedCreated:  nil,
			expectedResolves: nil,
			expectedCloses:   nil,
			expectedErr:      "'Resolves' can't be set to an earlier date than 'Created'",
		},
		{
			name:             "resolves and created can be the same time",
			inputCreated:     timePointer(now),
			inputResolves:    timePointer(now),
			inputCloses:      nil,
			expectedCreated:  timePointer(now),
			expectedResolves: timePointer(now),
			expectedCloses:   nil,
			expectedErr:      "",
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {
			c := initServerAndGetClient2(t)

			newForecast := NewForecast{
				Title: "Will it rain tomorrow?",
				Description: "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				Closes:   tt.inputCloses,
				Resolves: *tt.inputResolves,
				Created:  tt.inputCreated,
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

			response, err := CreateForecast(
				context.Background(),
				c,
				newForecast,
				newEstimate,
			)
			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, response.CreateForecast.Id)
				assertTimeAlmostEqual(
					t,
					*tt.expectedCreated,
					response.CreateForecast.Created,
				)
				assertTimeAlmostEqual(
					t,
					*tt.expectedResolves,
					response.CreateForecast.Resolves,
				)
				if tt.expectedCloses != nil {
					assertTimeAlmostEqual(
						t,
						*tt.expectedCloses,
						*response.CreateForecast.Closes,
					)
				} else {
					assert.Nil(t, response.CreateForecast.Closes)
				}
				// We always expect the Estimate.Created to match Forecast.Created
				assert.Equal(
					t,
					response.CreateForecast.Created,
					response.CreateForecast.Estimates[0].Created,
				)
			}
		})
	}
}

// TestCreateForecast_NewEstimateCreatedIsIgnored verifies that the 'Created'
// value in NewEstimate is ignored when creating a new forecast because we want
// the initial estimate to match the forecast. Later estimates can have
// different Created timestamps.
func TestCreateForecast_NewEstimateCreatedIsIgnored(t *testing.T) {
	c := initServerAndGetClient2(t)
	now := time.Now().UTC()
	hoursAgo24 := now.Add(-24 * time.Hour)
	hoursAgo5 := now.Add(-5 * time.Hour)

	newForecast := NewForecast{
		Title: "Will it rain tomorrow?",
		Description: "It counts as rain if between 9am and 9pm there are " +
			"30 min or more of uninterrupted precipitation.",
		Closes:   nil,
		Resolves: now.Add(24 * time.Hour),
		Created:  timePointer(hoursAgo24),
	}

	newEstimate := NewEstimate{
		Reason:  "My weather app says it will rain.",
		Created: timePointer(hoursAgo5),
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

	response, err := CreateForecast(
		context.Background(),
		c,
		newForecast,
		newEstimate,
	)
	require.NoError(t, err)
	assertTimeAlmostEqual(
		t,
		hoursAgo24,
		response.CreateForecast.Estimates[0].Created,
	)
}