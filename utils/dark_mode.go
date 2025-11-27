package utils

import (
	"os/exec"
	"strings"
)

var colorMap = map[string][]string{
	// colorName: [lightMode, darkMode]
	"red":       {"crimson", "darkred"},
	"primary":   {"#000000", "#FFFFFF"},
	"secondary": {"#FF6666", "#333030"},
}

var IsDarkMode bool

func init() {
	cmd := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle")
	if value, err := cmd.Output(); err == nil {
		IsDarkMode = strings.Contains(strings.ToLower(string(value)), "dark")
		return
	}
	IsDarkMode = false
}

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
