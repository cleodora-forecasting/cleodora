package cleosrv

import (
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"

	"github.com/cleodora-forecasting/cleodora/cleoutils/errors"
)

const DefaultConfigFileName = "cleosrv"
const DefaultConfigFileType = "yml"
const EnvPrefix = "CLEOSRV"

type Config struct {
	ConfigFile string
	Address    string
	Database   string
	Frontend   struct {
		FooterText string
	}
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
				return errors.Wrap(err, "error reading config file")
			}
		}
	}

	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	err = v.Unmarshal(c)
	return err
}
