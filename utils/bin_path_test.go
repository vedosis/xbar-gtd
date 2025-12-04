package utils

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindBin(t *testing.T) {
	assert.Equal(t, GetFQFN(exec.LookPath), GetFQFN(lookPathCmd))

	sentLocations := []string{}
	lookPathCmd = func(searchPath string) (string, error) {
		sentLocations = append(sentLocations, searchPath)
		return "", fmt.Errorf("generic error")
	}
	defer func() { lookPathCmd = exec.LookPath }()

	result := FindBin("test")
	assert.Empty(t, result)
	assert.Contains(t, sentLocations, "test")
	assert.Contains(t, sentLocations, "/usr/local/bin/test")

	sentLocations = []string{}
	lookPathCmd = func(searchPath string) (string, error) {
		return "/from/path/golangunittest", nil
	}
	result = FindBin("test")
	assert.Equal(t, "/from/path/golangunittest", result)
	assert.NotContains(t, sentLocations, "/usr/local/bin/test")

	sentLocations = []string{}
	lookPathCmd = func(searchPath string) (string, error) {
		if strings.Contains(searchPath, "sbin") {
			return "/usr/sbin/golangunittest", nil
		}
		return "", fmt.Errorf("generic error")
	}
	result = FindBin("test")
	assert.Equal(t, "/usr/sbin/golangunittest", result)
}
