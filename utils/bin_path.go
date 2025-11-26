package utils

import (
	"os/exec"
	"path"
)

var locations = []string{
	"/opt/homebrew/bin",
	"/usr/local/bin",
	"/usr/bin",
	"/usr/local/sbin",
	"/usr/sbin",
}

func FindBin(s string) string {
	if p, err := exec.LookPath(s); err == nil {
		return p
	}
	for _, p := range locations {
		p, err := exec.LookPath(path.Join(p, s))
		if err == nil {
			return p
		}
	}
	return ""
}
