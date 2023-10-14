package signr

import (
	"fmt"
	"github.com/spf13/viper"
)

func (cfg *Config) SetDefault(name string) (err error) {
	grid, _, err := cfg.GetList(nil)
	if err != nil {

		Fatal("ERROR: '%s'\n\n", err)
	}

	if cfg.DefaultKey == name {
		return fmt.Errorf("key '%s' was already the default", cfg.DefaultKey)
	}

	for _, row := range grid {

		for j := range row {

			if name == row[j] {

				cfg.DefaultKey = row[0]

				viper.Set("default", cfg.DefaultKey)

				if err = viper.WriteConfig(); err != nil {

					Fatal("failed to update config: '%v'\n", err)
				}

				PrintErr("key %s %s now default\n", row[0], row[1])
				return
			}
		}
	}

	return
}
