package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cleodora-forecasting/cleodora/cleoc/cleoc"
)

func buildAddCommand(app *cleoc.App) *cobra.Command {
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("add called")
		},
	}

	addCmd.AddCommand(buildAddForecastCommand(app))

	return addCmd
}

func buildAddForecastCommand(app *cleoc.App) *cobra.Command {
	var forecastCmd = &cobra.Command{
		Use:   "forecast TITLE RESOLUTION_DATE [DESCRIPTION]",
		Short: "Add a new forecast",
		Long: `Add a new forecast to Cleodora

TITLE is the title of the forecast.
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
			description := ""
			if len(args) == 3 {
				description = args[2]
			}
			return app.AddForecast(
				args[0],
				args[1],
				description,
				"TODO cleoc",
				[]string{"TODO cleoc:50", "TODO cleoc:50"},
			)
		},
	}
	return forecastCmd
}
