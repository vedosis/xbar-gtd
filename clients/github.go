package clients

type GithubClient struct {
	token string
	user  string
}

func NewGithubClient(auth *GithubAuth) *GithubClient {
	return &GithubClient{token: auth.Token, user: auth.Username}
}

func (c *GithubClient) GetMyPRs() ([]*PullRequestMeta, error) {
	return make([]*PullRequestMeta, 0), nil
}

type PullRequestMeta struct {
}
