package clients

import (
	"fmt"
	"os/exec"
	"testing"
	"xbar/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func dummyFn(string, ...string) Command { return nil }

func TestNewGithubCLIClient(t *testing.T) {
	client := NewGithubCLIClient("test1")
	assert.Equal(t, "test1", client.BinPath())
	assert.Equal(t, utils.GetFQFN(defaultExecCommand), utils.GetFQFN(client.execCommandFn))

	client = NewGithubCLIClient("test2")
	client.execCommandFn = dummyFn

	assert.Equal(t, "test2", client.BinPath())
	assert.Equal(t, utils.GetFQFN(dummyFn), utils.GetFQFN(client.execCommandFn))
	assert.NotEqualf(t, utils.GetFQFN(exec.Command), utils.GetFQFN(client.execCommandFn), "expected different functions")
}

func TestFindGithubCLI(t *testing.T) {
	oldFindBin := utils.FindBin
	findBin = func(string) string { return "" }
	defer func() { findBin = oldFindBin }()

	_, err := FindGithubCLI()
	assert.Error(t, err)

	findBin = func(string) string { return "/test/path/to/bin" }
	client, err := FindGithubCLI()
	assert.NoError(t, err)

	assert.Equal(t, "/test/path/to/bin", client.BinPath())
}

type DummyExecutorCommand struct {
	output []byte
	err    error
}

func (d *DummyExecutorCommand) Output() ([]byte, error) {
	return d.output, d.err
}

func TestGithubCLI_GetUserAndToken(t *testing.T) {
	client := NewGithubCLIClient("/test/fake/path")
	test := []struct {
		name       string
		fileName   string
		token      string
		errorText  string
		throwError bool
	}{
		{
			name:     "happy path",
			fileName: "testdata/gh/auth_status.json",
			token:    "gho_totallyalegittoken",
		},
		{
			name:      "no host",
			fileName:  "testdata/gh/auth_status_no_host.json",
			errorText: "host github.com not found in gh cli output",
		},
		{
			name:      "invalid json",
			fileName:  "testdata/gh/auth_status_invalid.json",
			errorText: "cannot decode JSON from gh cli",
		},
		{
			name:       "cli_error",
			throwError: true,
			errorText:  "error using gh cli to check auth status",
		},
		{
			name:      "no token",
			fileName:  "testdata/gh/auth_status_no_token.json",
			errorText: "no active https host found for github.com",
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			client.execCommandFn = func(cmd string, args ...string) Command {
				if tt.throwError {
					return &DummyExecutorCommand{
						err: fmt.Errorf("generic cli error"),
					}
				}
				require.Equal(t, "/test/fake/path", cmd)
				assert.Contains(t, args, "auth")
				assert.Contains(t, args, "status")

				data := utils.LoadTestFileData(t, tt.fileName)

				return &DummyExecutorCommand{
					output: data,
					err:    nil,
				}
			}

			auth, err := client.GetUserAndToken("github.com")
			if tt.errorText != "" {
				require.Error(t, err)
				require.Empty(t, auth)
				assert.Contains(t, err.Error(), tt.errorText)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.token, auth.Token)
		})
	}
}
