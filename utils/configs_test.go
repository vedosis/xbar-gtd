package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var badTestConfigFileData = `
name: test
name: test
`

type ConfigTestType struct {
	Name string `yaml:"name"`
}

func TestGetConfig(t *testing.T) {
	dir, err := os.MkdirTemp("", "unit-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	getWorkingDir = func() (string, error) {
		return dir, nil
	}
	defer func() {
		getWorkingDir = os.Getwd
		os.RemoveAll(dir)
	}()

	// Happy found filepath and write config
	if err = WriteConfig("test1", ConfigTestType{Name: "test-name"}); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	config, err := GetConfig[ConfigTestType]("test1", func() *ConfigTestType {
		return &ConfigTestType{Name: "missing-data"}
	})

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "test-name", config.Name)

	// No config file found (first file is empty)
	if err = os.WriteFile(filepath.Join(dir, ".config.test2.yml"), []byte(""), 0600); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}
	config, err = GetConfig[ConfigTestType]("test2", func() *ConfigTestType {
		return &ConfigTestType{Name: "default"}
	})
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, config.Name, "default")

	// Bad config file
	if err = os.WriteFile(filepath.Join(dir, ".config.test3.yaml"), []byte(badTestConfigFileData), 0600); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}
	_, err = GetConfig[ConfigTestType]("test3", func() *ConfigTestType { return nil })
	assert.Error(t, err)
}
