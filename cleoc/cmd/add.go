package main

import (
	"github.com/spf13/cobra"

	"github.com/cleodora-forecasting/cleodora/cleoc/cleoc"
)

func buildAddCommand(app *cleoc.App) *cobra.Command {
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add new things to Cleodora",
		Long: `Add new things to Cleodora, such as forecasts
`,
	}

	addCmd.AddCommand(buildAddForecastCommand(app))

	return addCmd
}

func buildAddForecastCommand(app *cleoc.App) *cobra.Command {
	var opts cleoc.AddForecastOptions

	var forecastCmd = &cobra.Command{
		Use:   "forecast",
		Short: "Add a new forecast",
		Long: `Add a new forecast to Cleodora

It returns the ID of the forecast that was just created.

Example:

	cleoc add forecast \
        --title "Will it rain tomorrow?" \
        --resolves 2022-11-14T00:00:00+01:00 \
        --description "If during the day it" \
            "rains for more than 2 minutes at" \
            "a time the forecast resolves as" \
            "true." \
        --probability Yes=70 \
        --probability No=30 \
        --reason "The weather forecast said" \
            "it would rain"
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.AddForecast(opts)
		},
	}

	forecastCmd.Flags().StringVarP(
		&opts.Title,
		"title",
		"t",
		"",
		"Title of the forecast",
	)
	forecastCmd.Flags().StringVarP(
		&opts.Description,
		"description",
		"d",
		"",
		"Description of the forecast",
	)
	forecastCmd.Flags().StringVarP(
		&opts.Resolves,
		"resolves",
		"r",
		"",
		"Resolution date of the forecast (format 2022-11-13T19:30:00+01:00)",
	)
	forecastCmd.Flags().StringVar(
		&opts.Closes,
		"closes",
		"",
		"Closing date of the forecast (format 2022-11-13T19:30:00+01:00)",
	)
	forecastCmd.Flags().StringVar(
		&opts.Reason,
		"reason",
		"",
		"The reason why you chose these probabilities",
	)
	forecastCmd.Flags().StringToIntVarP(
		&opts.Probabilities,
		"probability",
		"p",
		nil,
		"Outcome probability pair. Can be specified multiple times. "+
			"-p Green=30 -p Red=30 -p \"Light red=40\"",
	)

	return forecastCmd
}
