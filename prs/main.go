package prs

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"xbar/clients"
	"xbar/utils"
	xbar "xbar/xbar_utils"
)

type PRConfigs struct {
	GithubCLIBinPath string `yaml:"github_cli"`
	NerdFontName     string `yaml:"nerd_font"`
}

func NewPRConfigs() *PRConfigs {
	return &PRConfigs{}
}

func CLI() {
	renderer := xbar.NewXBarRenderer()

	env, messages := environmentSetup()
	if messages != nil && len(messages) > 0 {
		renderer.
			SetTitle("PRs").
			SetIcon("⚠️").
			Output(linesToInterfaces(messages)...)
		return
	}
	renderer.SetFont(env.FontName)

	myPRs := env.GithubClient.getMyPRs()

	renderer.SetTitle("\uF408 (0,0)")
	renderer.Output(
		xbar.NewXBarLine("My Open PRs",
			xbar.WithChildren(
				xbar.NewXBarLine("Things!!"),
				xbar.NewXBarLine("Things!!"),
				xbar.NewXBarLine("Things!!"),
				xbar.NewXBarLine("Things!!"),
				xbar.NewXBarLine("Things!!"),
				xbar.NewXBarLine("Stuff", xbar.WithChildren(xbar.NewXBarLine("Things!!"))),
			),
		),
		xbar.NewXBarLine("Test!"),
		xbar.NewXBarLine(fmt.Sprintf("FONT: %s", fontName)),
	)
}

type Environment struct {
	GithubClient *clients.GithubClient
	FontName     string
}

func environmentSetup() (*Environment, []*xbar.XBarLine) {
	env := &Environment{}
	var err error
	var ghCli *clients.GithubCLI
	redLine := xbar.NewXBarLine("Red Line", xbar.WithColor("red"))

	shouldWriteConfig := false
	config, err := utils.GetConfig[PRConfigs]("prs", NewPRConfigs)

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
			xbar.NewXBarLine("error details", xbar.WithChildren()),
		}
	}

	ghClient := clients.NewGithubClient(auth)
	env.GithubClient = ghClient

	font := config.NerdFontName
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
		config.NerdFontName = font
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

func linesToInterfaces(l []*xbar.XBarLine) []interface{} {
	result := make([]interface{}, len(l))
	for i, line := range l {
		result[i] = line
	}
	return result
}
