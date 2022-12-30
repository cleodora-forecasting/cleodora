package cleoc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/cleodora-forecasting/cleodora/cleoc/gqclient"
)

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

type AddForecastOptions struct {
	Title       string
	Description string
	Resolves    string
	// TODO Closes
}

func (opts *AddForecastOptions) Validate() error {
	if opts.Title == "" {
		return errors.New("--title can't be empty")
	}
	if _, err := time.Parse(time.RFC3339, opts.Resolves); err != nil {
		return errors.New("--resolves must be in RFC 3339 format (2022-11-13T19:30:00+01:00)")
	}
	return nil
}
