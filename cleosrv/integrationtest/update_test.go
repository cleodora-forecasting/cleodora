package integrationtest

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdate_From_0_1_1(t *testing.T) {
	tDir, err := os.MkdirTemp("", t.Name()+"_")
	require.NoError(t, err)
	dbSrc := filepath.Join("testdata", t.Name(), "test.db")
	dbPath := filepath.Join(tDir, "test.db")
	t.Log("executing with DB", dbPath)
	err = CopyFile(dbSrc, dbPath)
	require.NoError(t, err)

	c := initServerAndGetClient(t, dbPath)
	resp, err := GetForecasts(context.Background(), c)
	require.NoError(t, err)
	require.Len(t, resp.Forecasts, 4)
}
