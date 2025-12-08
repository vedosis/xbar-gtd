package prs

import (
	"fmt"
	"time"
	"xbar/pkg/clients"
	"xbar/pkg/xbar"
)

// PRService handles all PR-related operations
type PRService struct {
	ghClient   GithubClientInterface
	ghCLI      GithubCLIInterface
	configMgr  ConfigManager
	fontFinder FontFinder
	renderer   Renderer
	config     *PRConfigs
}

func NewPRService(
	configManagerFactory ConfigManagerFactory,
	findGithubCliFn FindGithubCLIFunction,
	githubCLIClientFactory GithubCLIClientFactory,
	githubClientFactory GithubClientFactory,
	fontFinderFactory FontFinderFactory,
	rendererFactory RendererFactory,
) (*PRService, error) {
	configMgr := configManagerFactory()
	config, err := configMgr.GetConfig()
	if err != nil {
		return nil, err
	}

	var ghCLI GithubCLIInterface
	if config.GithubCLIBinPath == "" {
		ghCLI, err = findGithubCliFn()
		if err != nil {
			return nil, fmt.Errorf("could not find gh CLI: %w", err)
		}
		config.GithubCLIBinPath = ghCLI.BinPath()
	} else {
		ghCLI = githubCLIClientFactory(config.GithubCLIBinPath)
	}

	auth, err := ghCLI.GetUserAndToken(config.GithubHostname)
	if err != nil {
		return nil, fmt.Errorf("could not fetch authentication tokens: %w", err)
	}

	ghClient := githubClientFactory(auth)
	ghClient.SetHostname(config.GithubHostname)
	ghClient.SetOldDays(config.OldPRDays)

	return &PRService{
		ghClient:   ghClient,
		ghCLI:      ghCLI,
		configMgr:  configMgr,
		fontFinder: fontFinderFactory(),
		renderer:   rendererFactory(),
		config:     config,
	}, nil
}

func (s *PRService) Run() error {
	if err := s.ensureFont(); err != nil {
		err = fmt.Errorf("failed to find Nerd Font: %w", err)
		s.renderError(err)
		return err
	}

	prList, err := s.ghClient.GetMyPRs()
	if err != nil {
		s.renderError(fmt.Errorf("failed to fetch PRs: %w", err))
		return err
	}

	s.renderPRList(prList)
	return nil
}

func (s *PRService) ensureFont() error {
	if s.config.FontName == "" {
		font, err := s.fontFinder.FindNerdFont()
		if err != nil {
			return err
		}

		s.config.FontName = font
		s.configMgr.WriteConfig(s.config)
	}
	s.renderer.SetFont(s.config.FontName)
	return nil
}

func (s *PRService) renderPRList(prList *clients.PRList) {
	if len(prList.MyPRs)+len(prList.Assigned)+len(prList.ReviewRequests) == 0 {
		s.renderer.SetTitle(IconPlugin)
		s.renderer.Output(
			xbar.NewXBarLine("No PRs to display"),
			xbar.NewXBarLine(fmt.Sprintf("Last update: %s", time.Now().Format(time.RFC3339))),
		)
		return
	}

	title := fmt.Sprintf("%s(%d,%d,%d)", IconPlugin, len(prList.MyPRs), len(prList.Assigned), len(prList.ReviewRequests))
	s.renderer.SetTitle(title)

	env := &Environment{
		GithubClient: s.ghClient,
		FontName:     s.config.FontName,
		AddIconSpace: s.config.AddIconSpace,
	}

	prLines := RenderPRList(prList, env)
	outputLines := make([]interface{}, len(prLines))
	for i, line := range prLines {
		outputLines[i] = line
	}
	s.renderer.Output(outputLines...)
}

func (s *PRService) renderError(err error) {
	s.renderer.SetTitle(IconWarning)
	s.renderer.SetIcon(IconWarning)
	s.renderer.Output(
		xbar.NewXBarLine("Error", xbar.WithColor("red")),
		xbar.NewXBarLine(err.Error(), xbar.WithWrapTextLength(120)),
	)
}
