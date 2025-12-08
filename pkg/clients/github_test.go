package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"xbar/pkg/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewGithubClient(t *testing.T) {
	client := NewGithubClient(&GithubAuth{
		Username: "testy-mc-testy-pants",
		Token:    "totally-a-legit-token",
	})
	require.NotNil(t, client)
	assert.Equal(t, client.hostname, "github.com")
	assert.Equal(t, client.user, "testy-mc-testy-pants")
	assert.Equal(t, client.token, "totally-a-legit-token")
	assert.Equal(t, client.oldDays, 7)

	client.SetHostname("gertherb.com")
	assert.Equal(t, client.hostname, "gertherb.com")

	client.SetOldDays(14)
	assert.Equal(t, client.oldDays, 14)
}

func TestBuildGraphQLQuery(t *testing.T) {
	username := "testy-mc-testy-pants"
	token := "Iamtheverymodelofamodernmajorgentlemen"
	client := NewGithubClient(&GithubAuth{
		Username: username,
		Token:    token,
	})
	data := client.buildGraphQLQuery(graphQLPRsTemplate)
	assert.NotEmpty(t, data)

	var query struct {
		Query string `json:"query"`
	}

	err := json.Unmarshal([]byte(data), &query)
	require.NoError(t, err)
	assert.NotEmpty(t, query.Query)
	assert.Equal(t, "query", query.Query[0:5])
	assert.Contains(t, query.Query, username)
	assert.NotContains(t, query.Query, token, "token should not be in the query")
}

func TestGithubClient_GetMyPRs_HappyPath(t *testing.T) {
	token := "Iamtheverymodelofamodernmajorgentlemen"
	client := NewGithubClient(&GithubAuth{
		Username: "testy-mc-testy-pants",
		Token:    token,
	})

	tests := []struct {
		name                  string
		filename              string
		failureText           string
		myPRExpected          int
		reviewRequestExpected int
		assignedExpected      int
		statusCode            int
		err                   error
	}{
		{
			name:                  "Test with valid data",
			filename:              "pr_data.json",
			myPRExpected:          2,
			reviewRequestExpected: 3,
			assignedExpected:      4,
		},
		{
			name:     "Test with no data",
			filename: "pr_empty.json",
		},
		{
			name:        "Test with error",
			filename:    "pr_error.json",
			failureText: "failed to parse response: GraphQL errors: Could not resolve to a PullRequest with the number of '123'.",
		},
		{
			name:        "Test with invalid json",
			filename:    "pr_invalid.json",
			failureText: "failed to parse response: failed to unmarshal response: unexpected end of JSON input",
		},
		{
			name:        "Test with bad status code",
			filename:    "pr_error.json",
			failureText: "failed to execute query: GitHub API returned status 400",
			statusCode:  400,
		},
		{
			name:        "Test with http exception",
			filename:    "pr_error.json",
			failureText: "failed to execute query: failed to execute request: generic error from http client",
			err:         fmt.Errorf("generic error from http client"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseData := utils.LoadTestFileData(t, "testdata/github.com/graphql/"+tt.filename)
			if tt.statusCode == 0 {
				tt.statusCode = 200
			}

			mockHttp := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					if tt.err != nil {
						return nil, tt.err
					}
					assert.Equal(t, req.URL.String(), "https://api.github.com/graphql")
					assert.Equal(t, req.Header.Get("Authorization"), fmt.Sprintf("Bearer %s", token))

					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(bytes.NewReader(responseData)),
					}, nil
				},
			}
			client.httpClient = mockHttp
			prs, err := client.GetMyPRs()
			if tt.failureText != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.failureText)
				return
			}
			require.NoError(t, err)
			assert.Len(t, prs.MyPRs, tt.myPRExpected)
			assert.Len(t, prs.ReviewRequests, tt.reviewRequestExpected)
			assert.Len(t, prs.Assigned, tt.assignedExpected)
		})
	}
}
