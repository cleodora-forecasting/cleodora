package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cleodora-forecasting/cleodora/cleoutils"
)

func buildVersionCommand() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(cleoutils.Version)
			return nil
		},
	}
	return versionCmd
}

func init() {
	rootCmd.AddCommand(buildVersionCommand())
}
