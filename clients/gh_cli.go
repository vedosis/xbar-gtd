package clients

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type GithubCLI struct {
	BinPath string
}

type GithubAuthStatus struct {
	Hosts map[string]GithubAuthStatusHost `json:"hosts"`
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

func NewGithubCLIClient(binPath string) *GithubCLI {
	return &GithubCLI{BinPath: binPath}
}

func FindGithubCLI() (*GithubCLI, error) {
	binPath, err := exec.LookPath("gh")
	if err != nil {
		return nil, fmt.Errorf("could not find gh: %w", err)
	}
	return &GithubCLI{BinPath: binPath}, nil
}

func (c *GithubCLI) GetUserAndToken(hostname string) (*GithubAuth, error) {
	var err error
	var auth = &GithubAuth{}
	cmd := exec.Command(
		c.BinPath,
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

	auth.Token = hostStatus.Token
	auth.Username = hostStatus.Login

	return auth, nil
}
