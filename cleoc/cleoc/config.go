package cleoc

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	URL        string
	ConfigFile string
}

func (c *Config) Load() error {
	v := viper.New()

	if c.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(c.ConfigFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cleoc" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yml")
		viper.SetConfigName(".cleoc")
	}

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	// v.AutomaticEnv() // should I do this?

	err = v.Unmarshal(c)
	return err
}
