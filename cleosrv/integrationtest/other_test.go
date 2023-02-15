package integrationtest

import (
	"context"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	client := initServerAndGetClient(t)

	query := `
		query GetMetadata {
			metadata {
				version
			}
		}`

	req := graphql.Request{
		Query: query,
	}

	var data struct {
		Metadata struct {
			Version string
		}
	}
	response := graphql.Response{Data: &data}

	err := client.MakeRequest(context.Background(), &req, &response)
	require.NoError(t, err)

	assert.Equal(t, "dev", data.Metadata.Version)
}
