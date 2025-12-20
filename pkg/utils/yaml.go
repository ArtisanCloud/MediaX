package utils

import (
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

var envPlaceholderPattern = regexp.MustCompile(`\$\{([A-Za-z0-9_]+)(:-([^}]*))?\}`)

func LoadYAML(filename string, cfg interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	processed := envPlaceholderPattern.ReplaceAllFunc(data, func(match []byte) []byte {
		groups := envPlaceholderPattern.FindSubmatch(match)
		if len(groups) == 0 {
			return match
		}
		key := string(groups[1])
		defaultVal := ""
		if len(groups) >= 4 {
			defaultVal = string(groups[3])
		}
		if val := os.Getenv(key); val != "" {
			return []byte(val)
		}
		return []byte(defaultVal)
	})

	if err := yaml.Unmarshal(processed, cfg); err != nil {
		return err
	}
	return nil
}
