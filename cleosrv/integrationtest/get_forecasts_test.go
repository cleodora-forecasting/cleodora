package integrationtest

import (
	"context"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetForecasts verifies that the forecasts are returned and
// uses the gqlgen.client for it.
func TestGetForecasts(t *testing.T) {
	c := initServerAndGetClient2(t)

	response, err := GetForecasts(context.Background(), c)
	require.NoError(t, err)

	require.Len(t, response.Forecasts, 3)
	assert.Equal(
		t,
		"Will the number of contributors to \"Cleodora\" be more than 3 at"+
			" the end of 2022?",
		response.Forecasts[2].Title,
	)
}

// TestGetForecasts_SomeFields verifies that the query can contain only a few
// fields.
func TestGetForecasts_OnlySomeFields(t *testing.T) {
	c := initServerAndGetClient2(t)

	query := `
		query GetForecasts {
			forecasts {
				id
				title
			}
		}`

	var data struct {
		Forecasts []struct {
			Id    string
			Title string
		}
	}
	request := graphql.Request{Query: query}
	response := graphql.Response{Data: &data}

	err := c.MakeRequest(context.Background(), &request, &response)
	require.NoError(t, err)

	require.Len(t, data.Forecasts, 3)
	assert.Equal(
		t,
		"Will the number of contributors to \"Cleodora\" be more than 3 at"+
			" the end of 2022?",
		data.Forecasts[2].Title,
	)
}

// TestGetForecasts_InvalidField verifies that an error is returned when
// querying for a field that does not exist.
func TestGetForecasts_InvalidField(t *testing.T) {
	c := initServerAndGetClient2(t)

	query := `
		query GetForecasts {
			forecasts {
				id
				title
				does_not_exist
			}
		}`

	req := graphql.Request{
		Query: query,
	}
	response := graphql.Response{}

	err := c.MakeRequest(
		context.Background(),
		&req,
		&response,
	)
	assert.Contains(t, err.Error(), "422")
	assert.Contains(t, err.Error(), "Cannot query field \\\"does_not_exist\\\"")
}
