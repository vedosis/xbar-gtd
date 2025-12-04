package clients

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"xbar/utils"
)

type Command interface {
	Output() ([]byte, error)
}

type GithubCLI struct {
	binPath       string
	execCommandFn func(string, ...string) Command
}

func (c *GithubCLI) BinPath() string {
	return c.binPath
}

type GithubAuthStatus struct {
	Hosts map[string][]GithubAuthStatusHost `json:"hosts"`
}

type GithubAuthStatusHost struct {
	State       string `json:"state"`
	Active      bool   `json:"active"`
	Host        string `json:"host"`
	Login       string `json:"login"`
	TokenSource string `json:"tokenSource"`
	Token       string `json:"token"`
	Scopes      string `json:"scopes"`
	GitProtocol string `json:"gitProtocol"`
}

type GithubAuth struct {
	Username string
	Token    string
}

var findBin = utils.FindBin

func defaultExecCommand(command string, args ...string) Command {
	return exec.Command(command, args...)
}

func NewGithubCLIClient(binPath string) *GithubCLI {
	client := &GithubCLI{
		binPath:       binPath,
		execCommandFn: defaultExecCommand,
	}
	return client
}

func FindGithubCLI() (*GithubCLI, error) {
	binPath := findBin("gh")
	if binPath != "" {
		return NewGithubCLIClient(binPath), nil
	}
	return &GithubCLI{binPath: binPath}, fmt.Errorf("could not find gh cli in any known locations")
}

func (c *GithubCLI) GetUserAndToken(hostname string) (*GithubAuth, error) {
	var err error
	var auth = &GithubAuth{}
	cmd := c.execCommandFn(
		c.BinPath(),
		"auth",
		"status",
		"--hostname", "github.com",
		"--json", "hosts",
		"--show-token",
	)
	out, err := cmd.Output()
	if err != nil {
		return auth, fmt.Errorf("error using gh cli to check auth status: %v", err)
	}

	status := GithubAuthStatus{}
	if err = json.Unmarshal(out, &status); err != nil {
		return auth, fmt.Errorf("cannot decode JSON from gh cli: %v", err)
	}

	hostStatus, ok := status.Hosts[hostname]
	if !ok {
		return auth, fmt.Errorf("host %s not found in gh cli output", hostname)
	}

	var foundHost GithubAuthStatusHost
	for _, host := range hostStatus {
		if host.Active && host.GitProtocol == "https" {
			foundHost = host
			break
		}
	}

	if foundHost.Token == "" {
		return auth, fmt.Errorf("no active https host found for %s", hostname)
	}

	auth.Token = foundHost.Token
	auth.Username = foundHost.Login

	return auth, nil
}
