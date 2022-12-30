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

// TODO update the example
	cleoc add forecast "Will it rain tomorrow?" 2022-11-14T00:00:00+01:00 "If \
		during the day it rains for more than 2 minutes at a time the \
		forecast resolves as true."
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.AddForecast(
				opts.Title,
				opts.Resolves,
				opts.Description,
				"TODO cleoc",
				[]string{"TODO cleoc:50", "TODO cleoc:50"},
			)
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

	return forecastCmd
}
