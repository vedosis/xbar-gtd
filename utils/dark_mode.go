package utils

import (
	"os/exec"
	"strings"
)

var colorMap = map[string][]string{
	// colorName: [lightMode, darkMode]
	"red": {"darkred", "crimson"},
}

var IsDarkMode bool

func Color(color string) string {
	value, ok := colorMap[color]
	if !ok {
		return color
	}
	if IsDarkMode {
		return value[1]
	} else {
		return value[0]
	}
}

func init() {
	cmd := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle")
	if value, err := cmd.Output(); err == nil {
		IsDarkMode = strings.Contains(strings.ToLower(string(value)), "dark")
	}
	IsDarkMode = false
}
