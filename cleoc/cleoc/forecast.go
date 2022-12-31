package cleoc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/cleodora-forecasting/cleodora/cleoc/gqclient"
)

// AddForecast creates a new forecast.
// The incoming options opts are assumed to already have been validated.
func (a *App) AddForecast(opts AddForecastOptions) error {
	resolvesT, err := time.Parse(time.RFC3339, opts.Resolves)
	if err != nil {
		return err // TODO wrap
	}
	ctx := context.Background()
	client := graphql.NewClient(
		fmt.Sprintf("%s/query", a.Config.URL),
		http.DefaultClient,
	)
	forecast := gqclient.NewForecast{
		Title:       opts.Title,
		Description: opts.Description,
		Resolves:    resolvesT,
		Closes:      resolvesT, // TODO should be optional See:
		// https://github.com/Khan/genqlient/blob/main/docs/FAQ.md#-nullable-fields
	}

	reqProbabilities, err := parseProbabilities(opts.Probabilities)
	if err != nil {
		return fmt.Errorf("error parsing probabilities: %w", err)
	}

	estimate := gqclient.NewEstimate{
		Reason:        opts.Reason,
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

func parseProbabilities(probabilities map[string]int) ([]gqclient.NewProbability, error) {
	var reqProbabilities []gqclient.NewProbability
	for outcome, value := range probabilities {
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

type AddForecastOptions struct {
	Title         string
	Description   string
	Resolves      string
	Reason        string
	Probabilities map[string]int
	// TODO Closes
}

func (opts *AddForecastOptions) Validate() error {
	// TODO consider joining as many errors as possible into one to be more
	// user friendly
	if opts.Title == "" {
		return errors.New("--title can't be empty")
	}
	if _, err := time.Parse(time.RFC3339, opts.Resolves); err != nil {
		return errors.New("--resolves must be in RFC 3339 format (2022-11-13T19:30:00+01:00)")
	}
	if opts.Reason == "" {
		return errors.New("--reason can't be empty")
	}
	if len(opts.Probabilities) == 0 {
		return errors.New("--probability is required")
	}
	sumProbabilities := 0
	for o, p := range opts.Probabilities {
		if o == "" {
			return errors.New("--probability has wrong format. Use '-p Yes=30'")
		}
		if p < 0 || p > 100 {
			return errors.New("probabilities must be between 0 and 100")
		}
		sumProbabilities += p
	}
	if sumProbabilities != 100 {
		return fmt.Errorf(
			"all probabilities must add up to 100 (here only %v)",
			sumProbabilities,
		)
	}
	return nil
}
