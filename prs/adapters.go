package prs

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
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

type DefaultConfigManager struct{}

func (d *DefaultConfigManager) GetConfig() (*PRConfigs, error) {
	return utils.GetConfig[PRConfigs]("prs", NewPRConfigs)
}

func (d *DefaultConfigManager) WriteConfig(config *PRConfigs) error {
	return utils.WriteConfig("prs", config)
}

type DefaultFontFinder struct{}

func (d *DefaultFontFinder) FindNerdFont() (string, error) {
	fcListBinPath := utils.FindBin("fc-list")
	if fcListBinPath == "" {
		return "", fmt.Errorf("could not find fc-list command")
	}

	cmd := exec.Command(fcListBinPath, ":", "family")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}
	defer cmd.Process.Kill()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "nerd") && strings.Contains(lowerLine, "mono") {
			return line, nil
		}
	}

	return "", fmt.Errorf("no nerd font found")
}

type RendererAdapter struct {
	xbarRenderer *xbar.XBarRenderer
}

func NewRendererAdapter(o ...xbar.XBarRendererOption) *RendererAdapter {
	return &RendererAdapter{
		xbarRenderer: xbar.NewXBarRenderer(o...),
	}
}

func (r *RendererAdapter) SetTitle(title string) Renderer {
	r.xbarRenderer.SetTitle(title)
	return r
}

func (r *RendererAdapter) SetIcon(icon string) Renderer {
	r.xbarRenderer.SetIcon(icon)
	return r
}

func (r *RendererAdapter) SetFont(font string) Renderer {
	r.xbarRenderer.SetFont(font)
	return r
}

func (r *RendererAdapter) Output(lines ...interface{}) {
	r.xbarRenderer.Output(lines...)
}
