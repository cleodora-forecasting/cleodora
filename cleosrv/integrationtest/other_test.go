package integrationtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	c := initServerAndGetClient(t)

	query := `
		query GetMetadata {
			metadata {
				version
			}
		}`

	var resp struct {
		Metadata struct {
			Version string
		}
	}

	err := c.Post(query, &resp)
	require.NoError(t, err)

	t.Log(resp)
	assert.Equal(t, "dev", resp.Metadata.Version)
}
