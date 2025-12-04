package prs

import (
	"fmt"
	"os"
	"xbar/clients"
)

/**
 * Build the core wiring of the plugin.
 * This is crap that's hard to mock and test, but with these factories, we can return different mocks for testing.
 */

func configManagerFactory() ConfigManager {
	return &DefaultConfigManager{}
}

func githubCLIFactory(binPath string) GithubCLIInterface {
	return clients.NewGithubCLIClient(binPath)
}

func findGithubCli() (GithubCLIInterface, error) {
	return clients.FindGithubCLI()
}

func githubClientFactory(auth *clients.GithubAuth) GithubClientInterface {
	return clients.NewGithubClient(auth)
}
func fontFinderFactory() FontFinder {
	return &DefaultFontFinder{}
}

func rendererFactory() Renderer {
	return NewRendererAdapter()
}

// CLI Entrypoint for the xbar plugin, called by xbar based off the $0's name. (i.e. prs.1m.cgo)
func CLI() {
	service, err := NewPRService(
		configManagerFactory,
		findGithubCli,
		githubCLIFactory,
		githubClientFactory,
		fontFinderFactory,
		rendererFactory,
	)
	if err != nil {
		renderSetupError(err)
	}

	if err := service.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func renderSetupError(err error) {
	fmt.Printf("%s%s | color=red\n", IconPlugin, IconWarning)
	fmt.Printf("---\n")
	fmt.Printf("Setup Error: %s\n", err.Error())
}
