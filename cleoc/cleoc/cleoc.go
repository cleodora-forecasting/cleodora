package cleoc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func (a *App) AddForecast(
	title string,
	resolves string,
	description string,
	reason string,
	probabilities []string,
) error {
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

	reqProbabilities, err := validateAndParseProbabilities(probabilities)
	if err != nil {
		return fmt.Errorf("error parsing probabilities: %w", err)
	}

	estimate := gqclient.NewEstimate{
		Reason:        reason,
		Probabilities: reqProbabilities,
	}
	resp, err := gqclient.CreateForecast(ctx, client, forecast, estimate)
	if err != nil {
		return err // TODO wrap
	}
	_, err = fmt.Fprint(a.Out, resp.CreateForecast.Id)
	if err != nil {
		return err
	}
	return nil
}

func validateAndParseProbabilities(probabilities []string) ([]gqclient.NewProbability, error) {
	if len(probabilities) == 0 {
		return nil, errors.New("no probabilities")
	}
	var reqProbabilities []gqclient.NewProbability
	for _, p := range probabilities {
		if !strings.Contains(p, ":") {
			return nil, fmt.Errorf("'%v' must contain ':'", p)
		}
		firstSegment := p[:strings.LastIndex(p, ":")]
		lastSegment := p[strings.LastIndex(p, ":")+1:]
		if firstSegment == "" {
			return nil, fmt.Errorf("'%v' the outcome can't be empty. "+
				"Use OUTCOME:PROBABILITY", p)
		}
		if lastSegment == "" {
			return nil, fmt.Errorf("'%v' the probability can't be empty. "+
				"Use OUTCOME:PROBABILITY", p)
		}
		value, err := strconv.Atoi(lastSegment)
		if err != nil {
			return nil, fmt.Errorf("'%v' the probability is not a "+
				"valid number. Use OUTCOME:PROBABILITY", p)
		}
		outcome := firstSegment

		reqProbabilities = append(
			reqProbabilities,
			gqclient.NewProbability{
				Value: value,
				Outcome: gqclient.NewOutcome{
					Text: outcome,
				},
			},
		)
	}
	return reqProbabilities, nil
}

func (a *App) Version() error {
	if _, err := fmt.Fprint(a.Out, cleoutils.Version); err != nil {
		return err
	}
	return nil
}
