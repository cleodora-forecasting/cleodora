package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cleodora-forecasting/cleodora/cleosrv/cleosrv"
)

var cfgFile string
var address string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cleosrv",
	Short: "Cleodora server to track personal forecasts",
	Long: `This server is made out of a GraphQL API and an embedded user
interface you can access via a web browser. You may also use a client (e.g. the
'cleoc' tool).

The purpose of Cleodora is tracking personal forecasts (e.g. 'Will I get a
raise within the next 6 months?') and systematically improve at making such
forecasts.

Visit https://cleodora.org for more information.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cleosrv.Start(viper.GetString("address"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file (default is $HOME/.cleosrv.yaml)",
	)

	rootCmd.PersistentFlags().StringVar(
		&address,
		"address",
		"localhost:8080",
		"Bind the process to a network address and port number. "+
			"To bind to all IP addresses and hostnames just specify "+
			"semicolon port e.g. :8080",
	)
	err := viper.BindPFlag(
		"address",
		rootCmd.PersistentFlags().Lookup("address"),
	)
	if err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cleosrv" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cleosrv")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}