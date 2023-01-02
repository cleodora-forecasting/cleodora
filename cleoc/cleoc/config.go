package cleoc

import (
	"fmt"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

const DefaultConfigFileName = "cleoc"
const DefaultConfigFileType = "yml"

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
		viper.AddConfigPath(xdg.ConfigHome)
		viper.SetConfigType(DefaultConfigFileType)
		viper.SetConfigName(DefaultConfigFileName)
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
