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

var lookPathCmd = exec.LookPath

func FindBin(s string) string {
	if p1, err := lookPathCmd(s); err == nil {
		return p1
	}
	for _, p2 := range locations {
		p3, err := lookPathCmd(path.Join(p2, s))
		if err == nil {
			return p3
		}
	}
	return ""
}
