package cleoc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/go-multierror"

	"github.com/cleodora-forecasting/cleodora/cleoc/gqclient"
	"github.com/cleodora-forecasting/cleodora/cleoutils/errors"
)

// AddForecast creates a new forecast.
// The incoming options opts are assumed to already have been validated.
func (a *App) AddForecast(opts AddForecastOptions) error {
	resolvesT, err := time.Parse(time.RFC3339, opts.Resolves)
	if err != nil {
		return errors.Wrap(err, "could not parse 'resolves'")
	}
	var closesT *time.Time
	if opts.Closes != "" {
		parsedTime, err := time.Parse(time.RFC3339, opts.Closes)
		if err != nil {
			return errors.Wrap(err, "could not parse 'closes'")
		}
		closesT = &parsedTime
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
		Closes:      closesT,
	}

	reqProbabilities, err := parseProbabilities(opts.Probabilities)
	if err != nil {
		return errors.Wrap(err, "error parsing probabilities")
	}

	estimate := gqclient.NewEstimate{
		Reason:        opts.Reason,
		Probabilities: reqProbabilities,
	}
	resp, err := gqclient.CreateForecast(ctx, client, forecast, estimate)
	if err != nil {
		return errors.Wrap(err, "error calling the API")
	}
	_, err = fmt.Fprint(a.Out, resp.CreateForecast.Id+"\n")
	if err != nil {
		return err
	}
	return nil
}

func parseProbabilities(probabilities map[string]int) ([]gqclient.NewProbability, error) {
	var reqProbabilities []gqclient.NewProbability
	for outcome, value := range probabilities {
		newOutcome := gqclient.NewOutcome{
			Text: outcome,
		}
		reqProbabilities = append(
			reqProbabilities,
			gqclient.NewProbability{
				Value:   value,
				Outcome: &newOutcome,
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
	Closes        string
}

func (opts *AddForecastOptions) Validate() error {
	var validationErr *multierror.Error
	if opts.Title == "" {
		validationErr = multierror.Append(
			validationErr,
			errors.New("--title can't be empty"),
		)
	}
	if _, err := time.Parse(time.RFC3339, opts.Resolves); err != nil {
		validationErr = multierror.Append(
			validationErr,
			errors.New("--resolves must be in RFC 3339 format "+
				"(2022-11-13T19:30:00+01:00)"),
		)
	}
	if opts.Closes != "" { // it's allowed to be empty
		if _, err := time.Parse(time.RFC3339, opts.Closes); err != nil {
			validationErr = multierror.Append(
				validationErr,
				errors.New("--closes must be in RFC 3339 format "+
					"(2022-11-13T19:30:00+01:00)"),
			)
		}
	}
	if opts.Reason == "" {
		validationErr = multierror.Append(
			validationErr,
			errors.New("--reason can't be empty"),
		)
	}
	if len(opts.Probabilities) == 0 {
		validationErr = multierror.Append(
			validationErr,
			errors.New("--probability is required"),
		)
	}
	sumProbabilities := 0
	for o, p := range opts.Probabilities {
		if o == "" {
			validationErr = multierror.Append(
				validationErr,
				errors.New("--probability has wrong format. Use '-p Yes=30'"),
			)
		}
		if p < 0 || p > 100 {
			validationErr = multierror.Append(
				validationErr,
				errors.New("probabilities must be between 0 and 100"),
			)
		}
		sumProbabilities += p
	}
	if sumProbabilities != 100 {
		validationErr = multierror.Append(
			validationErr,
			errors.Newf(
				"all probabilities must add up to 100 (here only %v)",
				sumProbabilities,
			),
		)
	}
	return validationErr.ErrorOrNil()
}
