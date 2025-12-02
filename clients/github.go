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

type GithubClient struct {
	token    string
	user     string
	hostname string
	oldDays  int
}

func NewGithubClient(auth *GithubAuth) *GithubClient {
	return &GithubClient{
		token:    auth.Token,
		user:     auth.Username,
		hostname: "github.com",
		oldDays:  7,
	}
}

func (c *GithubClient) SetHostname(hostname string) {
	c.hostname = hostname
}

func (c *GithubClient) SetOldDays(days int) {
	c.oldDays = days
}

// PullRequestMeta represents metadata about a pull request
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

// PRList categorizes PRs by type
type PRList struct {
	MyPRs          []*PullRequestMeta
	ReviewRequests []*PullRequestMeta
	Assigned       []*PullRequestMeta
}

// GetMyPRs fetches all relevant PRs from GitHub
func (c *GithubClient) GetMyPRs() (*PRList, error) {
	query := c.buildGraphQLQuery()
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

func (c *GithubClient) buildGraphQLQuery() string {
	tmpl, err := template.New("graphql").Parse(graphQLPRsTemplate)
	if err != nil {
		panic(fmt.Errorf("failed to parse GraphQL query template: %w", err))
	}

	data := map[string]string{
		"User": c.user,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(fmt.Errorf("failed to execute GraphQL query template: %w", err))
	}

	return fmt.Sprintf(`{ "query": "%s" }`, strings.Join(strings.Fields(strings.ReplaceAll(buf.String(), "\"", "\\\"")), " "))
}

func (c *GithubClient) executeGraphQLQuery(query string) ([]byte, error) {
	apiURL := fmt.Sprintf("https://api.%s/graphql", c.hostname)
	if c.hostname == "github.com" {
		apiURL = "https://api.github.com/graphql"
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(query))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func (c *GithubClient) parseGraphQLResponse(data []byte) (*PRList, error) {
	var response struct {
		Data struct {
			MyPRs struct {
				Edges []struct {
					Node struct {
						Number    int       `json:"number"`
						Title     string    `json:"title"`
						URL       string    `json:"url"`
						CreatedAt time.Time `json:"createdAt"`
						UpdatedAt time.Time `json:"updatedAt"`
						IsDraft   bool      `json:"isDraft"`
						Author    struct {
							Login string `json:"login"`
						} `json:"author"`
						Repository struct {
							Owner struct {
								Login string `json:"login"`
							} `json:"owner"`
							Name string `json:"name"`
							URL  string `json:"url"`
						} `json:"repository"`
						ReviewDecision string `json:"reviewDecision"`
						Comments       struct {
							TotalCount int `json:"totalCount"`
						} `json:"comments"`
						Reviews struct {
							TotalCount int `json:"totalCount"`
						} `json:"reviews"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"myPRs"`
			ReviewRequests struct {
				Edges []struct {
					Node struct {
						Number    int       `json:"number"`
						Title     string    `json:"title"`
						URL       string    `json:"url"`
						CreatedAt time.Time `json:"createdAt"`
						UpdatedAt time.Time `json:"updatedAt"`
						IsDraft   bool      `json:"isDraft"`
						Author    struct {
							Login string `json:"login"`
						} `json:"author"`
						Repository struct {
							Owner struct {
								Login string `json:"login"`
							} `json:"owner"`
							Name string `json:"name"`
							URL  string `json:"url"`
						} `json:"repository"`
						ReviewDecision string `json:"reviewDecision"`
						Comments       struct {
							TotalCount int `json:"totalCount"`
						} `json:"comments"`
						Reviews struct {
							TotalCount int `json:"totalCount"`
						} `json:"reviews"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"reviewRequests"`
			Assigned struct {
				Edges []struct {
					Node struct {
						Number    int       `json:"number"`
						Title     string    `json:"title"`
						URL       string    `json:"url"`
						CreatedAt time.Time `json:"createdAt"`
						UpdatedAt time.Time `json:"updatedAt"`
						IsDraft   bool      `json:"isDraft"`
						Author    struct {
							Login string `json:"login"`
						} `json:"author"`
						Repository struct {
							Owner struct {
								Login string `json:"login"`
							} `json:"owner"`
							Name string `json:"name"`
							URL  string `json:"url"`
						} `json:"repository"`
						ReviewDecision string `json:"reviewDecision"`
						Comments       struct {
							TotalCount int `json:"totalCount"`
						} `json:"comments"`
						Reviews struct {
							TotalCount int `json:"totalCount"`
						} `json:"reviews"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"assigned"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

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
		MyPRs:          make([]*PullRequestMeta, 0),
		ReviewRequests: make([]*PullRequestMeta, 0),
		Assigned:       make([]*PullRequestMeta, 0),
	}

	// Parse My PRs
	for _, edge := range response.Data.MyPRs.Edges {
		pr := c.convertNodeToPRMeta(
			edge.Node.Number,
			edge.Node.Title,
			edge.Node.URL,
			edge.Node.CreatedAt,
			edge.Node.UpdatedAt,
			edge.Node.IsDraft,
			edge.Node.Author.Login,
			edge.Node.Repository.Owner.Login,
			edge.Node.Repository.Name,
			edge.Node.Repository.URL,
			edge.Node.ReviewDecision,
			edge.Node.Comments.TotalCount,
		)
		prList.MyPRs = append(prList.MyPRs, pr)
	}

	// Parse Review Requests
	for _, edge := range response.Data.ReviewRequests.Edges {
		pr := c.convertNodeToPRMeta(
			edge.Node.Number,
			edge.Node.Title,
			edge.Node.URL,
			edge.Node.CreatedAt,
			edge.Node.UpdatedAt,
			edge.Node.IsDraft,
			edge.Node.Author.Login,
			edge.Node.Repository.Owner.Login,
			edge.Node.Repository.Name,
			edge.Node.Repository.URL,
			edge.Node.ReviewDecision,
			edge.Node.Comments.TotalCount,
		)
		prList.ReviewRequests = append(prList.ReviewRequests, pr)
	}

	// Parse Assigned
	for _, edge := range response.Data.Assigned.Edges {
		pr := c.convertNodeToPRMeta(
			edge.Node.Number,
			edge.Node.Title,
			edge.Node.URL,
			edge.Node.CreatedAt,
			edge.Node.UpdatedAt,
			edge.Node.IsDraft,
			edge.Node.Author.Login,
			edge.Node.Repository.Owner.Login,
			edge.Node.Repository.Name,
			edge.Node.Repository.URL,
			edge.Node.ReviewDecision,
			edge.Node.Comments.TotalCount,
		)
		prList.Assigned = append(prList.Assigned, pr)
	}

	return prList, nil
}

func (c *GithubClient) convertNodeToPRMeta(
	number int,
	title string,
	url string,
	createdAt time.Time,
	updatedAt time.Time,
	isDraft bool,
	author string,
	repoOwner string,
	repoName string,
	repoURL string,
	reviewDecision string,
	commentsCount int,
) *PullRequestMeta {
	daysSinceUpdate := int(time.Since(updatedAt).Hours() / 24)
	isOld := daysSinceUpdate >= c.oldDays

	return &PullRequestMeta{
		Number:                number,
		Title:                 title,
		URL:                   url,
		RepoOwner:             repoOwner,
		RepoName:              repoName,
		RepoURL:               repoURL,
		State:                 "OPEN",
		IsDraft:               isDraft,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
		Author:                author,
		ReviewDecision:        reviewDecision,
		HasUnresolvedComments: commentsCount > 0,
		Assignees:             []string{},
		RequestedReviewers:    []string{},
		IsOld:                 isOld,
		DaysSinceUpdate:       daysSinceUpdate,
	}
}
