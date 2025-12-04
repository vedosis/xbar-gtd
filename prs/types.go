package prs

import "xbar/clients"

type GithubClientInterface interface {
	GetMyPRs() (*clients.PRList, error)
	SetHostname(hostname string)
	SetOldDays(days int)
}

type GithubCLIInterface interface {
	BinPath() string
	GetUserAndToken(hostname string) (*clients.GithubAuth, error)
}

type ConfigManager interface {
	GetConfig() (*PRConfigs, error)
	WriteConfig(config *PRConfigs) error
}

type FontFinder interface {
	FindNerdFont() (string, error)
}

type Renderer interface {
	SetTitle(title string) Renderer
	SetIcon(icon string) Renderer
	SetFont(font string) Renderer
	Output(lines ...interface{})
}

type Environment struct {
	GithubClient GithubClientInterface
	FontName     string
	AddIconSpace bool
}

type FindGithubCLIFunction func() (GithubCLIInterface, error)
type GithubCLIClientFactory func(string) GithubCLIInterface
type GithubClientFactory func(auth *clients.GithubAuth) GithubClientInterface
type ConfigManagerFactory func() ConfigManager
type FontFinderFactory func() FontFinder
type RendererFactory func() Renderer
