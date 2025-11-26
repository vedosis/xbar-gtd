package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"xbar/prs"
)

type CliFunction = func()

var cliRegistry = make(map[string]CliFunction)

func init() {
	cliRegistry["prs"] = prs.CLI
}

func main() {
	base := path.Base(os.Args[0])
	name := strings.Split(base, ".")
	fn, ok := cliRegistry[name[0]]
	if !ok {
		fmt.Printf("Unknown function: %s\nRunnin all commands (sep. ----)\n", name[0])
		runAllCLIs(cliRegistry)
		return
	}
	fn()
}

func runAllCLIs(mapping map[string]CliFunction) {
	for key, value := range mapping {
		fmt.Printf("--- %s ---\n", key)
		value()
	}
}
