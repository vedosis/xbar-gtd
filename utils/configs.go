package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var FILENAME = ".config.%s.%s"
var getWorkingDir = os.Getwd

func GetConfig[ConfigType any](name string, defaultFn func() *ConfigType) (*ConfigType, error) {
	var config = defaultFn()

	extensions := []string{"yaml", "yml"}
	searchedFor := make([]string, len(extensions))

	var data []byte
	var err error
	var foundFile string

	wd, _ := getWorkingDir()

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
		return config, nil
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
	wd, _ := getWorkingDir()
	filename := filepath.Join(wd, fmt.Sprintf(FILENAME, name, "yaml"))
	return os.WriteFile(filename, data, 0600)
}
