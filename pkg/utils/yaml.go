package utils

import (
	"gopkg.in/yaml.v3"
	"os"
)

func LoadYAML(filename string, cfg interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}
	return nil
}
