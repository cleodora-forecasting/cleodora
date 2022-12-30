package cleoc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/cleodora-forecasting/cleodora/cleoc/gqclient"
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

func (a *App) AddForecast(title string, resolves string, description string) error {
	resolvesT, err := time.Parse(time.RFC3339, resolves)
	if err != nil {
		return err // todo wrap
	}
	ctx := context.Background()
	client := graphql.NewClient(
		fmt.Sprintf("%s/query", a.Config.URL),
		http.DefaultClient,
	)
	forecast := gqclient.NewForecast{
		Title:       title,
		Description: description,
		Resolves:    resolvesT,
		Closes:      resolvesT, // should be optional
	}
	estimate := gqclient.NewEstimate{
		Reason: "TODO cleoc",
		Probabilities: []gqclient.NewProbability{
			{
				Value:   50,
				Outcome: gqclient.NewOutcome{Text: "TODO cleoc"},
			},
			{
				Value:   50,
				Outcome: gqclient.NewOutcome{Text: "TODO cleoc"},
			},
		},
	}
	resp, err := gqclient.CreateForecast(ctx, client, forecast, estimate)
	if err != nil {
		return err // todo wrap
	}
	_, err = fmt.Fprint(a.Out, resp.CreateForecast.Id)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) Version() error {
	if _, err := fmt.Fprint(a.Out, cleoutils.Version); err != nil {
		return err
	}
	return nil
}
