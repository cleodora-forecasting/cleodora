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
)

type App struct {
	Out    io.Writer
	Config *Config
}

func NewApp() *App {
	c := &Config{}
	return &App{
		Out:    os.Stdout,
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
	fmt.Println(resp.CreateForecast.Id)
	return nil
}
