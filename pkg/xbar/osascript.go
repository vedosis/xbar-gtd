package xbar

import "fmt"

func OpenInTerminal(cmd string) []string {
	return []string{
		"osascript",
		"-e", "tell app \"Terminal\"",
		"-e", "activate",
		"-e", fmt.Sprintf("do script \"%s\"", cmd),
		"-e", "end",
	}
}
