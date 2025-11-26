package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var FILENAME = ".config.%s.%s"

func GetConfig[ConfigType any](name string, defaultFn func() *ConfigType) (*ConfigType, error) {
	var config = defaultFn()

	extensions := []string{"yaml", "yml"}
	searchedFor := make([]string, len(extensions))

	var data []byte
	var err error
	var foundFile string

	wd, _ := os.Getwd()

	for idx, ext := range extensions {
		filename := fmt.Sprintf(FILENAME, name, ext)
		fqfn := filepath.Join(wd, filename)
		searchedFor[idx] = fqfn
		if data, err = os.ReadFile(fqfn); err == nil {
			if len(data) == 0 {
				continue
			}
			foundFile = fqfn
			break
		}
	}

	if data == nil || len(data) == 0 {
		return config, fmt.Errorf("config file not found or only empty files found, searched at (%s)", strings.Join(searchedFor, " | "))
	}

	if err = yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config file (%s): %w", foundFile, err)
	}

	return config, nil
}

func WriteConfig[ConfigType any](name string, config ConfigType) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config file (%s): %w", name, err)
	}
	wd, _ := os.Getwd()
	filename := filepath.Join(wd, fmt.Sprintf(FILENAME, name, "yaml"))
	return os.WriteFile(filename, data, 0600)
}
