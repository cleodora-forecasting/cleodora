package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cleodora-forecasting/cleodora/cleoutils"
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the client version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cleoutils.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
