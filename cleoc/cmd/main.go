package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cleodora-forecasting/cleodora/cleoc/cleoc"
	"github.com/cleodora-forecasting/cleodora/cleoutils"
)

func main() {
	app := cleoc.NewApp()

	cmd := buildRootCommand(app)
	if err := cmd.Execute(); err != nil {
		_, err = fmt.Fprint(os.Stderr, err.Error())
		fmt.Printf("Error printing error: %v\n", err)
		os.Exit(1)
	}
}

func buildRootCommand(app *cleoc.App) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "cleoc",
		Short: "Command line tool to interact with a Cleodora server",
		Long: fmt.Sprintf(`Create and modify forecasts in Cleodora server.

cleoc version: %s
`, cleoutils.Version),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := app.Config.LoadWithViper(); err != nil {
				return err
			}
			app.Out = cmd.OutOrStdout() // this is in the example code, but not in the slides, why?
			app.Err = cmd.OutOrStderr()
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(
		&app.Config.ConfigFile,
		"config",
		"",
		"config file (default is $HOME/.cleoc.yml)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&app.Config.URL,
		"url",
		"u",
		"http://localhost:8080",
		"base URL for the API",
	)

	rootCmd.AddCommand(buildAddCommand(app))
	rootCmd.AddCommand(buildVersionCommand(app))

	return rootCmd
}
