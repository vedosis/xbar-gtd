package prs

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"
	"xbar/clients"
	"xbar/utils"
	xbar "xbar/xbar_utils"
)

type PRConfigs struct {
	GithubCLIBinPath string `yaml:"github_cli"`
	FontName         string `yaml:"nerd_font"`
	OldPRDays        int    `yaml:"old_pr_days"`
	GithubHostname   string `yaml:"github_hostname"`
	AddIconSpace     bool   `yaml:"add_icon_space"`
}

func NewPRConfigs() *PRConfigs {
	return &PRConfigs{
		OldPRDays:      7,
		GithubHostname: "github.com",
		AddIconSpace:   false,
	}
}

func CLI() {
	renderer := xbar.NewXBarRenderer()

	env, messages := environmentSetup()
	if messages != nil && len(messages) > 0 {
		renderError(renderer, messages...)
		return
	}
	renderer.SetFont(env.FontName)

	prList, err := env.GithubClient.GetMyPRs()
	if err != nil {
		renderError(renderer,
			xbar.NewXBarLine("Error fetching PRs", xbar.WithColor("red")),
			xbar.NewXBarDebugError("Error fetching PRs", err),
		)
		return
	}

	// Set header title with counts
	renderer.SetTitle(fmt.Sprintf("%s(%d,%d,%d)", IconPlugin, len(prList.MyPRs), len(prList.ReviewRequests), len(prList.Assigned)))

	// Render PR lists
	prLines := RenderPRList(prList, env)
	if len(prLines) == 0 {
		renderer.SetTitle(IconPlugin)
		renderer.Output(
			xbar.NewXBarLine("No PRs to display"),
			xbar.NewXBarLine(fmt.Sprintf("Last updated: %s", time.Now().Format(time.RFC3339))),
		)
		return
	}

	// Convert to interface{} for renderer
	outputLines := make([]interface{}, len(prLines))
	for i, line := range prLines {
		outputLines[i] = line
	}
	renderer.Output(outputLines...)
}

type Environment struct {
	GithubClient *clients.GithubClient
	FontName     string
	AddIconSpace bool
}

func environmentSetup() (*Environment, []*xbar.XBarLine) {
	env := &Environment{}
	var err error
	var ghCli *clients.GithubCLI
	redLine := xbar.NewXBarLine("Red Line", xbar.WithColor("red"))

	shouldWriteConfig := false
	config, err := utils.GetConfig[PRConfigs]("prs", NewPRConfigs)

	if err != nil {
		env.AddIconSpace = config.AddIconSpace
	}

	if err != nil || config.GithubCLIBinPath == "" {
		ghCli, err = clients.FindGithubCLI()
		if err != nil {
			return nil, []*xbar.XBarLine{
				redLine.Clone("Could not find github cli"),
				xbar.NewXBarLine("'gh' is used to fetch authentication tokens'"),
				xbar.NewXBarLine("Run 'brew install gh' to install github cli",
					xbar.WithCommand(
						xbar.OpenInTerminal("brew install gh")...,
					),
				),
				xbar.NewXBarDebugError("error finding gh cli", err),
			}
		}
		config.GithubCLIBinPath = ghCli.BinPath
		shouldWriteConfig = true
	} else {
		ghCli = clients.NewGithubCLIClient(config.GithubCLIBinPath)
	}

	auth, err := ghCli.GetUserAndToken("github.com")
	if err != nil {
		return nil, []*xbar.XBarLine{
			redLine.Clone("Could not fetch authentication tokens"),
			xbar.NewXBarLine("Make sure you're logged in"),
			xbar.NewXBarLine("Run 'gh auth login' to log in", xbar.WithCommand(
				xbar.OpenInTerminal("gh auth login -w -h github.com -p https")...)),
			xbar.NewXBarDebugError("error fetching auth tokens", err),
		}
	}

	ghClient := clients.NewGithubClient(auth)
	ghClient.SetHostname(config.GithubHostname)
	ghClient.SetOldDays(config.OldPRDays)
	env.GithubClient = ghClient

	font := config.FontName
	if font == "" {
		fcListBinPath := utils.FindBin("fc-list")
		if fcListBinPath == "" {
			return env, []*xbar.XBarLine{
				redLine.Clone("Could not find command 'fc-list'"),
				xbar.NewXBarLine("Click here to execute 'brew install fontconfig", xbar.WithCommand(
					xbar.OpenInTerminal("/opt/homebrew/bin/brew install fontconfig")...,
				)),
			}
		}

		findFontCmd := exec.Command(fcListBinPath, ":", "family")
		font := findFontFromList(findFontCmd)
		if font == "" {
			return env, []*xbar.XBarLine{
				redLine.Clone("Couldn't find monospaced 'Nerd Font' to render icons."),
				xbar.NewXBarLine("Click here to open font selection.", xbar.WithHref("https://www.nerdfonts.com/font-downloads")),
			}
		}
		config.FontName = font
		shouldWriteConfig = true
	}
	env.FontName = font

	if shouldWriteConfig {
		err = utils.WriteConfig("prs", config)
		if err != nil {
			panic(fmt.Errorf("could not save configurations after successful env discovery: %w", err))
		}
	}

	return env, nil
}

func findFontFromList(cmd *exec.Cmd) string {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return ""
	}
	if err := cmd.Start(); err != nil {
		return ""
	}

	scanner := bufio.NewScanner(stdout)
	var result string
	defer cmd.Process.Kill()
	for scanner.Scan() {
		line := scanner.Text()
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "nerd") && strings.Contains(lowerLine, "mono") {
			return line
		}
	}
	return result
}

func renderError(renderer *xbar.XBarRenderer, lines ...*xbar.XBarLine) {
	rLines := make([]interface{}, len(lines))
	for i, line := range lines {
		rLines[i] = line
	}

	renderer.
		SetTitle("PRs").
		SetIcon("⚠️").
		Output(rLines...)
}
