package main

import (
	"github.com/spf13/cobra"

	"github.com/cleodora-forecasting/cleodora/cleoc/cleoc"
)

func buildVersionCommand(app *cleoc.App) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the client version",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// to be used for validation
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Version()
		},
		Annotations: map[string]string{
			"skipConfig": "",
		},
	}
	return versionCmd
}
