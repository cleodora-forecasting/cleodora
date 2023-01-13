package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cleodora-forecasting/cleodora/cleosrv/cleosrv"
)

var cfgFile string
var address string
var database string

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
		app := cleosrv.NewApp()
		app.Config.Address = viper.GetString("address")
		app.Config.Database = viper.GetString("database")
		app.Config.Frontend.FooterText = viper.GetString("frontend.footer_text")
		return app.Start()
	},
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.MousetrapHelpText = ""
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
		fmt.Sprintf(
			"config file (default is %v)",
			filepath.Join(xdg.ConfigHome, "cleosrv.yaml"),
		),
	)

	rootCmd.PersistentFlags().StringVar(
		&address,
		"address",
		"localhost:8080",
		"Bind the process to a network address and port number. "+
			"To bind to all IP addresses and hostnames just specify "+
			"semicolon port e.g. :8080",
	)
	rootCmd.PersistentFlags().StringVar(
		&database,
		"database",
		filepath.Join(xdg.DataHome, "cleosrv", "cleosrv.db"),
		"Path to the SQLite database to use. Will be created if it "+
			"doesn't exist.",
	)
	err := viper.BindPFlag(
		"address",
		rootCmd.PersistentFlags().Lookup("address"),
	)
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag(
		"database",
		rootCmd.PersistentFlags().Lookup("database"),
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
		viper.AddConfigPath(xdg.ConfigHome)
		viper.SetConfigType("yaml")
		viper.SetConfigName("cleosrv")
	}

	viper.SetEnvPrefix("CLEOSRV")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
