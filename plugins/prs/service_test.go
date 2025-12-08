package prs

import (
	"fmt"
	"testing"
	"xbar/pkg/clients"
	"xbar/pkg/xbar"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockGithubClient struct {
	PRs      *clients.PRList
	Err      error
	hostname string
	oldDays  int
}

func (m *MockGithubClient) GetMyPRs() (*clients.PRList, error) {
	return m.PRs, m.Err
}

func (m *MockGithubClient) SetHostname(hostname string) {
	m.hostname = hostname
}

func (m *MockGithubClient) SetOldDays(days int) {
	m.oldDays = days
}

type MockGithubCLI struct {
	binPath string
	Auth    *clients.GithubAuth
	Err     error
}

func (m *MockGithubCLI) BinPath() string {
	return m.binPath
}

func (m *MockGithubCLI) GetUserAndToken(hostname string) (*clients.GithubAuth, error) {
	return m.Auth, m.Err
}

type MockConfigManager struct {
	Config      *PRConfigs
	Err         error
	wroteConfig bool
}

func (m *MockConfigManager) GetConfig() (*PRConfigs, error) {
	return m.Config, m.Err
}

func (m *MockConfigManager) WriteConfig(_ *PRConfigs) error {
	m.wroteConfig = true
	return m.Err
}

type MockFontFinder struct {
	Font string
	Err  error
}

func (m *MockFontFinder) FindNerdFont() (string, error) {
	return m.Font, m.Err
}

type MockRenderer struct {
	Title       string
	Icon        string
	Font        string
	OutputLines []interface{}
}

func (m *MockRenderer) SetTitle(title string) Renderer {
	m.Title = title
	return m
}

func (m *MockRenderer) SetIcon(icon string) Renderer {
	m.Icon = icon
	return m
}

func (m *MockRenderer) SetFont(font string) Renderer {
	m.Font = font
	return m
}

func (m *MockRenderer) Output(lines ...interface{}) {
	m.OutputLines = lines
}

func TestNewPRService(t *testing.T) {
	tests := []struct {
		name           string
		configs        *PRConfigs
		configError    error
		ghCliPath      string
		ghCliAuth      *clients.GithubAuth
		ghCliError     error
		ghCliAuthError error
		errorNeedle    string
	}{
		{
			name:      "happy path",
			configs:   &PRConfigs{},
			ghCliPath: "/path/to/github/cli",
			ghCliAuth: &clients.GithubAuth{Username: "testy", Token: "totally-a-legit-token"},
		},
		{
			name:        "config error",
			configError: fmt.Errorf("config error"),
			errorNeedle: "config error",
		},
		{
			name:        "gh find error",
			configs:     &PRConfigs{},
			ghCliError:  fmt.Errorf("gh cli error"),
			errorNeedle: "could not find gh CLI",
		},
		{
			name:      "gh cli in config",
			configs:   &PRConfigs{GithubCLIBinPath: "/path/to/github/cli"},
			ghCliAuth: &clients.GithubAuth{Username: "testy", Token: "totally-a-legit-token"},
		},
		{
			name: "gh auth error",
			configs: &PRConfigs{
				GithubCLIBinPath: "/path/to/github/cli",
			},
			ghCliAuthError: fmt.Errorf("gh cli auth error"),
			errorNeedle:    "could not fetch authentication tokens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testConfigMgrFactory := func() ConfigManager {
				return &MockConfigManager{
					Config: tt.configs,
					Err:    tt.configError,
				}
			}
			testFindGhCLI := func() (GithubCLIInterface, error) {
				return &MockGithubCLI{binPath: tt.ghCliPath, Err: tt.ghCliAuthError}, tt.ghCliError
			}
			testGithubCLIClientFactory := func(binPath string) GithubCLIInterface {
				return &MockGithubCLI{
					Auth: &clients.GithubAuth{},
					Err:  tt.ghCliAuthError,
				}
			}

			service, err := NewPRService(
				testConfigMgrFactory,
				testFindGhCLI,
				testGithubCLIClientFactory,
				func(auth *clients.GithubAuth) GithubClientInterface { return &MockGithubClient{} },
				func() FontFinder { return &MockFontFinder{} },
				func() Renderer { return &MockRenderer{} },
			)

			if tt.configError != nil {
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.errorNeedle)
				return
			} else if tt.ghCliError != nil {
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.errorNeedle)
				return
			} else if tt.ghCliAuthError != nil {
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.errorNeedle)
				return
			} else {
				require.NoError(t, err)
				assert.NotNil(t, service)
			}

			assert.Equal(t, "/path/to/github/cli", service.config.GithubCLIBinPath)
			assert.NotNil(t, service.renderer)
			assert.NotNil(t, service.fontFinder)
			assert.NotNil(t, service.ghClient)
			assert.NotNil(t, service.ghCLI)
		})
	}
}

func TestPRService_Run(t *testing.T) {
	tests := []struct {
		name          string
		config        *PRConfigs
		errorHaystack string
		fontFinder    FontFinder
		prList        *clients.PRList
		prErr         error
		renderNeedles []string
	}{
		{
			name:       "happy path",
			config:     &PRConfigs{AddIconSpace: true},
			fontFinder: &MockFontFinder{Font: "Test Nerd Font"},
			prList: &clients.PRList{
				MyPRs: []*clients.PullRequestMeta{
					{Number: 123, Title: "Test PR"},
					{Number: 456, Title: "Another Test PR"},
				},
				ReviewRequests: []*clients.PullRequestMeta{
					{Number: 789, Title: "Review Request"},
				},
				Assigned: []*clients.PullRequestMeta{
					{Number: 78910, Title: "Assigned PR"},
					{Number: 101112, Title: "Yet Another Test PR", IsOld: true, DaysSinceUpdate: 100},
				},
			},
			renderNeedles: []string{IconPlugin, "(2,2,1)", "\uf407  /", "Test PR", "Another Test PR", "Assigned PR", "Yet Another Test PR"},
		},
		{
			name:          "font error",
			config:        &PRConfigs{},
			errorHaystack: "failed to find Nerd Font",
			fontFinder:    &MockFontFinder{Err: fmt.Errorf("font error")},
		},
		{
			name:          "fetch error",
			config:        &PRConfigs{FontName: "Test Nerd Font"},
			errorHaystack: "fetch error",
			prErr:         fmt.Errorf("fetch error"),
		},
		{
			name:          "empty list",
			config:        &PRConfigs{FontName: "Test Nerd Font"},
			renderNeedles: []string{IconPlugin, "No PRs to display"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outText := ""

			shouldWriteConfig := false
			if tt.config.FontName == "" {
				shouldWriteConfig = true
			}

			if tt.prList == nil {
				tt.prList = &clients.PRList{
					MyPRs:          make([]*clients.PullRequestMeta, 0),
					ReviewRequests: make([]*clients.PullRequestMeta, 0),
					Assigned:       make([]*clients.PullRequestMeta, 0),
				}
			}

			configMgr := &MockConfigManager{Config: tt.config}
			renderer := NewRendererAdapter(
				xbar.WithCmdPrintFn(func(a ...any) (int, error) {
					outText += fmt.Sprint(a...)
					return 0, nil
				}),
				xbar.WithCmdPrintLnFn(func(a ...any) (int, error) {
					outText += fmt.Sprint(append(a, "\n")...)
					return 0, nil
				}),
			)

			service := &PRService{
				config:     tt.config,
				fontFinder: tt.fontFinder,
				renderer:   renderer,
				ghClient:   &MockGithubClient{PRs: tt.prList, Err: tt.prErr},
				ghCLI:      &MockGithubCLI{},
				configMgr:  configMgr,
			}

			err := service.Run()

			if tt.errorHaystack == "" {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.errorHaystack)
				return
			}

			require.NotNil(t, tt.config.FontName)
			assert.Equal(t, shouldWriteConfig, configMgr.wroteConfig)

			if len(tt.renderNeedles) > 0 {
				for _, needle := range tt.renderNeedles {
					assert.Contains(t, outText, needle)
				}
				return
			} else {
				assert.Empty(t, outText)
			}
		})
	}
}
