package cmd

import (
	"fmt"

	"github.com/cleodora-forecasting/cleodora/utils"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the client version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(utils.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
