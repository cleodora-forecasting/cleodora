package cleoc_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cleodora-forecasting/cleodora/cleoc/cleoc"
)

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
	assert.NoError(t, err)
	assert.Equal(t, "dev\n", out.String())
	assert.Empty(t, errOut)
}
