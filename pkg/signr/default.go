package signr

import (
	"fmt"
	"github.com/spf13/viper"
)

func (s *Signr) SetDefault(name string) (err error) {
	grid, _, err := s.GetList(nil)
	if err != nil {

		s.Fatal("ERROR: '%s'\n\n", err)
	}

	if s.DefaultKey == name {
		return fmt.Errorf("key '%s' was already the default", s.DefaultKey)
	}

	for _, row := range grid {

		for j := range row {

			if name == row[j] {

				s.DefaultKey = row[0]

				viper.Set("default", s.DefaultKey)

				if err = viper.WriteConfig(); err != nil {

					s.Fatal("failed to update config: '%v'\n", err)
				}

				s.PrintErr("key %s %s now default\n", row[0], row[1])
				return
			}
		}
	}

	return
}
