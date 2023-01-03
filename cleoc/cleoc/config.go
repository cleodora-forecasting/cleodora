package cleoc

import (
	"fmt"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

const DefaultConfigFileName = "cleoc"
const DefaultConfigFileType = "yml"
const EnvPrefix = "CLEOC"

type Config struct {
	URL        string
	ConfigFile string
}

// LoadWithViper initializes or overwrites the Config by using the 'viper'
// library (thereby reading in config files, ENV variables etc.). You should
// probably not call it.
func (c *Config) LoadWithViper(v *viper.Viper) error {
	if c.ConfigFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(c.ConfigFile)
	} else {
		v.AddConfigPath(xdg.ConfigHome)
		v.SetConfigType(DefaultConfigFileType)
		v.SetConfigName(DefaultConfigFileName)
	}

	err := v.ReadInConfig()
	if err != nil {
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				// Note that ConfigFileNotFoundError is not returned when
				// explicitly specifying a config path (--config), which should
				// (and does) cause an error if it doesn't exist.
				return fmt.Errorf("error reading config file: %w", err)
			}
		}
	}

	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	err = v.Unmarshal(c)
	return err
}
