package integrationtest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	client := initServerAndGetClient(t, "")
	resp, err := GetMetadata(context.Background(), client)
	require.NoError(t, err)
	assert.Equal(t, "dev", resp.Metadata.Version)
}
