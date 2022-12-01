package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/cleodora-forecasting/cleodora/cleoc/gqclient"
)

// forecastCmd represents the forecast command.
var forecastCmd = &cobra.Command{
	Use:   "forecast SUMMARY RESOLUTION_DATE [DESCRIPTION]",
	Short: "Add a new forecast",
	Long: `Add a new forecast to Cleodora

SUMMARY is the title of the forecast.
RESOLUTION_DATE needs to be in the format 2022-11-13T19:30:00+01:00
DESCRIPTION is optional.

It returns the ID of the forecast that was just created.

Example:

	cleoc add forecast "Will it rain tomorrow?" 2022-11-14T00:00:00+01:00 "If \
		during the day it rains for more than 2 minutes at a time the \
		forecast resolves as true."
`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]
		resolves, err := time.Parse(time.RFC3339, args[1])
		if err != nil {
			return err // todo wrap
		}
		description := ""
		if len(args) == 3 {
			description = args[2]
		}
		ctx := context.Background()
		client := graphql.NewClient(
			fmt.Sprintf("%s/query", URL),
			http.DefaultClient,
		)
		f := gqclient.NewForecast{
			Title:       title,
			Description: description,
			Resolves:    resolves,
			Closes:      resolves, // should be optional
		}
		resp, err := gqclient.CreateForecast(ctx, client, f)
		if err != nil {
			return err // todo wrap
		}
		fmt.Println(resp.CreateForecast.Id)
		return nil
	},
}

func init() {
	addCmd.AddCommand(forecastCmd)
}
