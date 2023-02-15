package integrationtest

import (
	"context"
	"testing"
	"time"

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
	c := initServerAndGetClient2(t)

	attack := "<script>alert(document.cookie)</script>"

	newForecast := NewForecast{
		Title: "Will it rain tomorrow?" + attack,
		Description: "It counts as rain if between 9am and 9pm there are " +
			"30 min or more of uninterrupted precipitation." + attack,
		Closes:   timePointer(time.Now().Add(24 * time.Hour)),
		Resolves: time.Now().Add(24 * time.Hour),
	}

	newEstimate := NewEstimate{
		Reason: "My weather app says it will rain." + attack,
		Probabilities: []NewProbability{
			{
				Value: 70,
				Outcome: NewOutcome{
					Text: "Yes" + attack,
				},
			},
			{
				Value: 30,
				Outcome: NewOutcome{
					Text: "No" + attack,
				},
			},
		},
	}

	response, err := CreateForecast(context.Background(), c, newForecast, newEstimate)
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
		newEstimate NewEstimate
		expectedErr string
	}{
		{
			name: "success",
			newEstimate: NewEstimate{
				Reason: "My weather app says it will rain",
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
			},
			expectedErr: "",
		},
		{
			name: "single probability",
			newEstimate: NewEstimate{
				Reason: "My weather app says it will rain",
				Probabilities: []NewProbability{
					{
						Value: 100,
						Outcome: NewOutcome{
							Text: "Any outcome",
						},
					},
				},
			},
			expectedErr: "",
		},
		{
			name: "success with more probabilities",
			newEstimate: NewEstimate{
				Reason: "My weather app says it will rain",
				Probabilities: []NewProbability{
					{
						Value: 30,
						Outcome: NewOutcome{
							Text: "Yes, but less than 2 hours",
						},
					},
					{
						Value: 20,
						Outcome: NewOutcome{
							Text: "Yes, between 2 and 5 hours",
						},
					},
					{
						Value: 20,
						Outcome: NewOutcome{
							Text: "Yes, more than 5 hours",
						},
					},
					{
						Value: 30,
						Outcome: NewOutcome{
							Text: "No",
						},
					},
				},
			},
			expectedErr: "",
		},
		{
			name: "reason cant be empty",
			newEstimate: NewEstimate{
				Reason: "",
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
			},
			expectedErr: "'reason' can't be empty",
		},
		{
			name: "probabilities must add up to 100",
			newEstimate: NewEstimate{
				Reason: "My weather app says it will rain.",
				Probabilities: []NewProbability{
					{
						Value: 70,
						Outcome: NewOutcome{
							Text: "Yes",
						},
					},
					{
						Value: 20,
						Outcome: NewOutcome{
							Text: "No",
						},
					},
				},
			},
			expectedErr: "probabilities must add up to 100",
		},
		{
			name: "probabilities must be between 0 and 100",
			newEstimate: NewEstimate{
				Reason: "My weather app says it will rain.",
				Probabilities: []NewProbability{
					{
						Value: -10,
						Outcome: NewOutcome{
							Text: "Yes",
						},
					},
					{
						Value: 110,
						Outcome: NewOutcome{
							Text: "No",
						},
					},
				},
			},
			expectedErr: "probabilities must be between 0 and 100",
		},
		{
			name: "probabilities cant be empty",
			newEstimate: NewEstimate{
				Reason:        "My weather app says it will rain.",
				Probabilities: []NewProbability{},
			},
			expectedErr: "probabilities can't be empty",
		},
		{
			name: "outcome text cant be empty",
			newEstimate: NewEstimate{
				Reason: "My weather app says it will rain",
				Probabilities: []NewProbability{
					{
						Value: 70,
						Outcome: NewOutcome{
							Text: "",
						},
					},
					{
						Value: 30,
						Outcome: NewOutcome{
							Text: "No",
						},
					},
				},
			},
			expectedErr: "outcome text can't be empty",
		},
		{
			name: "outcomes cant be duplicates",
			newEstimate: NewEstimate{
				Reason: "My weather app says it will rain",
				Probabilities: []NewProbability{
					{
						Value: 70,
						Outcome: NewOutcome{
							Text: "No",
						},
					},
					{
						Value: 30,
						Outcome: NewOutcome{
							Text: "No",
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
			c := initServerAndGetClient2(t)

			newForecast := NewForecast{
				Title: "Will it rain tomorrow?",
				Description: "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				Closes:   timePointer(time.Now().Add(24 * time.Hour)),
				Resolves: time.Now().Add(24 * time.Hour),
			}

			_, err := CreateForecast(
				context.Background(),
				c,
				newForecast,
				tt.newEstimate,
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
		newForecast NewForecast
		expectedErr string
	}{
		{
			name: "success",
			newForecast: NewForecast{
				Title: "Will it rain tomorrow?",
				Description: "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				Closes:   timePointer(time.Now().Add(24 * time.Hour)),
				Resolves: time.Now().Add(24 * time.Hour),
			},
			expectedErr: "",
		},
		{
			name: "title cant be empty",
			newForecast: NewForecast{
				Title: "",
				Description: "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				Closes:   timePointer(time.Now().Add(24 * time.Hour)),
				Resolves: time.Now().Add(24 * time.Hour),
			},
			expectedErr: "title can't be empty",
		},
		{
			name: "description can be empty",
			newForecast: NewForecast{
				Title:       "Will it rain tomorrow?",
				Description: "",
				Closes:      timePointer(time.Now().Add(24 * time.Hour)),
				Resolves:    time.Now().Add(24 * time.Hour),
			},
			expectedErr: "",
		},
		{
			name: "closes can be empty",
			newForecast: NewForecast{
				Title: "Will it rain tomorrow?",
				Description: "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				Closes:   nil,
				Resolves: time.Now().Add(24 * time.Hour),
			},
			expectedErr: "",
		},
		{
			name: "closes can be omitted",
			newForecast: NewForecast{
				Title: "Will it rain tomorrow?",
				Description: "It counts as rain if between 9am and 9pm there are " +
					"30 min or more of uninterrupted precipitation.",
				Resolves: time.Now().Add(24 * time.Hour),
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		t.Run(tt.name, func(t *testing.T) {
			c := initServerAndGetClient2(t)

			newEstimate := NewEstimate{
				Reason: "My weather app says it will rain",
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
				tt.newForecast,
				newEstimate,
			)

			if tt.expectedErr == "" { // success
				require.NoError(t, err)
				assert.NotEmpty(t, response.CreateForecast.Id)
				assert.Equal(t, tt.newForecast.Title, response.CreateForecast.Title)
			} else {
				assert.ErrorContains(t, err, tt.expectedErr)
			}
		})
	}
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
