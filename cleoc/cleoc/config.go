package cleoc

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	URL        string
	ConfigFile string
}

// LoadWithViper initializes or overwrites the Config by using the 'viper'
// library (thereby reading in config files, ENV variables etc.). You should
// probably not call it.
func (c *Config) LoadWithViper() error {
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
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return fmt.Errorf("error reading config file: %w", err)
			}
		}
	}

	// v.AutomaticEnv() // should I do this?

	err = v.Unmarshal(c)
	return err
}
