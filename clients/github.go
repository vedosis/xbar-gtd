package clients

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"
)

//go:embed github.prs.graphql
var graphQLPRsTemplate string

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GithubClient struct {
	token      string
	user       string
	hostname   string
	oldDays    int
	httpClient HTTPClient
}

func NewGithubClient(auth *GithubAuth) *GithubClient {
	return &GithubClient{
		token:      auth.Token,
		user:       auth.Username,
		hostname:   "github.com",
		oldDays:    7,
		httpClient: &http.Client{Timeout: 3 * time.Second},
	}
}

func (c *GithubClient) SetHostname(hostname string) {
	c.hostname = hostname
}

func (c *GithubClient) SetOldDays(days int) {
	c.oldDays = days
}

type PullRequestMeta struct {
	Number                int
	Title                 string
	URL                   string
	RepoOwner             string
	RepoName              string
	RepoURL               string
	State                 string
	IsDraft               bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Author                string
	ReviewDecision        string
	HasUnresolvedComments bool
	Assignees             []string
	RequestedReviewers    []string
	IsOld                 bool
	DaysSinceUpdate       int
}

type PRList struct {
	MyPRs          []*PullRequestMeta
	ReviewRequests []*PullRequestMeta
	Assigned       []*PullRequestMeta
}

type GraphQLResponsePayload[T any] struct {
	Data   map[string]GraphQLData[T] `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}
type GraphQLData[T any] struct {
	Edges []struct {
		Node T `json:"node"`
	} `json:"edges"`
}

type User struct {
	Login string `json:"login"`
}
type Repository struct {
	Owner User   `json:"owner"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Comment struct {
	TotalCount int `json:"totalCount"`
}
type Review struct {
	TotalCount int `json:"totalCount"`
}

type PullRequest struct {
	Number         int        `json:"number"`
	Title          string     `json:"title"`
	URL            string     `json:"url"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	IsDraft        bool       `json:"isDraft"`
	Author         User       `json:"author"`
	Repository     Repository `json:"repository"`
	ReviewDecision string     `json:"reviewDecision"`
	Comments       Comment    `json:"comments"`
	Reviews        Review     `json:"reviews"`
}

func (c *GithubClient) GetMyPRs() (*PRList, error) {
	query := c.buildGraphQLQuery(graphQLPRsTemplate)
	resp, err := c.executeGraphQLQuery(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	prList, err := c.parseGraphQLResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return prList, nil
}

func (c *GithubClient) buildGraphQLQuery(templateString string) string {
	tmpl, _ := template.New("graphql").Parse(templateString)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf,
		map[string]string{
			"User": c.user,
		},
	); err != nil {
		return ""
	}

	return fmt.Sprintf(`{ "query": "%s" }`, strings.Join(strings.Fields(strings.ReplaceAll(buf.String(), "\"", "\\\"")), " "))
}

func (c *GithubClient) executeGraphQLQuery(query string) ([]byte, error) {
	apiURL := fmt.Sprintf("https://api.%s/graphql", c.hostname)

	req, _ := http.NewRequest("POST", apiURL, bytes.NewBufferString(query))
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

func (c *GithubClient) parseGraphQLResponse(data []byte) (*PRList, error) {
	var response GraphQLResponsePayload[PullRequest]
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Errors) > 0 {
		errMsgs := make([]string, len(response.Errors))
		for i, e := range response.Errors {
			errMsgs[i] = e.Message
		}
		return nil, fmt.Errorf("GraphQL errors: %s", strings.Join(errMsgs, "; "))
	}

	prList := &PRList{
		MyPRs:          c.convertNodeToPRMeta(response.Data["myPRs"]),
		ReviewRequests: c.convertNodeToPRMeta(response.Data["reviewRequests"]),
		Assigned:       c.convertNodeToPRMeta(response.Data["assigned"]),
	}

	return prList, nil
}

func (c *GithubClient) convertNodeToPRMeta(nodes GraphQLData[PullRequest]) []*PullRequestMeta {
	prs := make([]*PullRequestMeta, len(nodes.Edges))
	for idx, n := range nodes.Edges {
		daysSinceUpdate := int(time.Since(n.Node.UpdatedAt).Hours() / 24)
		isOld := daysSinceUpdate >= c.oldDays
		prs[idx] = &PullRequestMeta{
			Number:                n.Node.Number,
			Title:                 n.Node.Title,
			URL:                   n.Node.URL,
			RepoOwner:             n.Node.Repository.Owner.Login,
			RepoName:              n.Node.Repository.Name,
			RepoURL:               n.Node.Repository.URL,
			State:                 "OPEN",
			IsDraft:               n.Node.IsDraft,
			CreatedAt:             n.Node.CreatedAt,
			UpdatedAt:             n.Node.UpdatedAt,
			Author:                n.Node.Author.Login,
			ReviewDecision:        n.Node.ReviewDecision,
			HasUnresolvedComments: n.Node.Comments.TotalCount > 0,
			Assignees:             []string{},
			RequestedReviewers:    []string{},
			IsOld:                 isOld,
			DaysSinceUpdate:       daysSinceUpdate,
		}
	}
	return prs
}
