package cleoc

import (
	"fmt"
	"io"
	"os"

	"github.com/cleodora-forecasting/cleodora/cleoutils"
)

type App struct {
	Out    io.Writer
	Err    io.Writer
	Config *Config
}

func NewApp() *App {
	c := &Config{}
	return &App{
		Out:    os.Stdout,
		Err:    os.Stderr,
		Config: c,
	}
}

func (a *App) Version() error {
	if _, err := fmt.Fprintf(a.Out, "%v\n", cleoutils.Version); err != nil {
		return err
	}
	return nil
}
