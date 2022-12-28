package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cleodora-forecasting/cleodora/cleoc/cleoc"
	"github.com/cleodora-forecasting/cleodora/cleoutils"
)

func buildVersionCommand(app *cleoc.App) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the client version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cleoutils.Version)
		},
	}
	return versionCmd
}
