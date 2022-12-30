package cleoc_test

import (
	"bytes"
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

func TestApp_AddForecast(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.Nil(t, err)
		assert.Contains(t, string(body), "CreateForecast")
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
	)
	require.Nil(t, err)
	assert.Equal(t, "999", out.String())
	assert.Empty(t, errOut)
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
