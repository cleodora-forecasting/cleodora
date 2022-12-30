package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cleodora-forecasting/cleodora/cleoc/cleoc"
	"github.com/cleodora-forecasting/cleodora/cleoutils"
)

const (
	// Indicates that config should not be loaded for this command.
	// This is used for commands like help and version which should never
	// fail, even if cleoc is misconfigured.
	skipConfig string = "skipConfig"
)

func main() {
	app := cleoc.NewApp()

	cmd := buildRootCommand(app)
	if err := cmd.Execute(); err != nil {
		_, err = fmt.Fprintf(os.Stderr, "%s\n", err)
		if err != nil {
			fmt.Printf("Error printing error: %v\n", err)
		}
		os.Exit(1)
	}
}

func buildRootCommand(app *cleoc.App) *cobra.Command {
	var printVersion bool

	var rootCmd = &cobra.Command{
		Use:   "cleoc",
		Short: "Command line tool to interact with a Cleodora server",
		Long: fmt.Sprintf(`Create and modify forecasts in Cleodora server.

cleoc version: %s
`,
			cleoutils.Version,
		),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			app.Out = cmd.OutOrStdout()
			app.Err = cmd.OutOrStderr()

			if shouldSkipConfig(cmd) {
				return nil
			}

			if err := app.Config.LoadWithViper(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if printVersion {
				versionCmd := buildVersionCommand(app)
				err := versionCmd.PreRunE(cmd, args)
				if err != nil {
					return err
				}
				return versionCmd.RunE(cmd, args)
			}
			return cmd.Help()
		},
		SilenceUsage:  true,
		SilenceErrors: true, // Errors are printed by main
		Annotations: map[string]string{
			"skipConfig": "",
		},
	}

	// Flags for just the root command, does not apply to sub-commands
	rootCmd.Flags().BoolVarP(
		&printVersion,
		"version",
		"v",
		false,
		"Print the application version",
	)

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

func shouldSkipConfig(cmd *cobra.Command) bool {
	if cmd.Name() == "help" {
		return true
	}

	_, skip := cmd.Annotations[skipConfig]
	return skip
}
