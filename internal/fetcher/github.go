package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// githubClient 实现 GitHubClient 接口
type githubClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewGitHubClient 创建一个 GitHub API 客户端
func NewGitHubClient(baseURL, token string) GitHubClient {
	return &githubClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// issueResponse 是 GitHub Issue API 响应格式
type issueResponse struct {
	Title    string `json:"title"`
	User     struct {
		Login string `json:"login"`
	} `json:"user"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	State     string `json:"state"`
	Body     string `json:"body"`
	Labels   []struct {
		Name string `json:"name"`
	} `json:"labels"`
	Comments int    `json:"comments"`
	HTMLURL  string `json:"html_url"`
}

// prResponse 是 GitHub Pull Request API 响应格式
type prResponse struct {
	Title    string `json:"title"`
	User     struct {
		Login string `json:"login"`
	} `json:"user"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	State     string `json:"state"`
	Body     string `json:"body"`
	Labels   []struct {
		Name string `json:"name"`
	} `json:"labels"`
	Comments int    `json:"comments"`
	HTMLURL  string `json:"html_url"`
}

// discussionResponse 是 GitHub Discussion API 响应格式
type discussionResponse struct {
	Title    string `json:"title"`
	Author   struct {
		Login string `json:"login"`
	} `json:"author"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	URL      string `json:"url"`
	Body     string `json:"body"`
	Category struct {
		Name string `json:"name"`
	} `json:"category"`
	Answer   *discussionAnswer `json:"answer"`
	Comments struct {
		Nodes []discussionComment `json:"nodes"`
	} `json:"comments"`
}

type discussionAnswer struct {
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	CreatedAt string `json:"created_at"`
	Body     string `json:"body"`
}

type discussionComment struct {
	Author    struct {
		Login string `json:"login"`
	} `json:"author"`
	CreatedAt string `json:"created_at"`
	Body     string `json:"body"`
}

// FetchIssue 获取 Issue
func (c *githubClient) FetchIssue(ctx context.Context, owner, repo string, number int) (*GitHubResource, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", c.baseURL, owner, repo, number)
	return c.fetchIssue(ctx, url, "issue")
}

// FetchPullRequest 获取 Pull Request
func (c *githubClient) FetchPullRequest(ctx context.Context, owner, repo string, number int) (*GitHubResource, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", c.baseURL, owner, repo, number)
	return c.fetchIssue(ctx, url, "pull_request")
}

// FetchDiscussion 获取 Discussion
func (c *githubClient) FetchDiscussion(ctx context.Context, owner, repo string, number int) (*GitHubResource, error) {
	// Discussion 使用 GraphQL API 风格，但这里用 REST 模拟
	url := fmt.Sprintf("%s/repos/%s/%s/discussions/%d", c.baseURL, owner, repo, number)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("resource not found: %w", fmt.Errorf("discussion not found"))
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("GitHub API rate limit exceeded: %w", fmt.Errorf("rate limit"))
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var discResp discussionResponse
	if err := json.NewDecoder(resp.Body).Decode(&discResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	resource := &GitHubResource{
		Type:      "discussion",
		Title:     discResp.Title,
		Author:    discResp.Author.Login,
		State:     "open",
		Body:      discResp.Body,
		URL:       discResp.URL,
		Labels:    []string{},
		Comments:  []Comment{},
	}

	if discResp.CreatedAt != "" {
		resource.CreatedAt, _ = time.Parse(time.RFC3339, discResp.CreatedAt)
	}
	if discResp.UpdatedAt != "" {
		resource.UpdatedAt, _ = time.Parse(time.RFC3339, discResp.UpdatedAt)
	}

	for _, node := range discResp.Comments.Nodes {
		comment := Comment{
			Author: node.Author.Login,
			Body:   node.Body,
		}
		if node.CreatedAt != "" {
			comment.CreatedAt, _ = time.Parse(time.RFC3339, node.CreatedAt)
		}
		resource.Comments = append(resource.Comments, comment)
	}

	return resource, nil
}

func (c *githubClient) fetchIssue(ctx context.Context, urlStr, resourceType string) (*GitHubResource, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("resource not found: %w", fmt.Errorf("issue not found"))
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("GitHub API rate limit exceeded: %w", fmt.Errorf("rate limit"))
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var issueResp issueResponse
	if err := json.NewDecoder(resp.Body).Decode(&issueResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	resource := &GitHubResource{
		Type:   resourceType,
		Title:  issueResp.Title,
		Author: issueResp.User.Login,
		State:  issueResp.State,
		Body:   issueResp.Body,
		URL:    issueResp.HTMLURL,
		Labels: []string{},
	}

	if issueResp.CreatedAt != "" {
		resource.CreatedAt, _ = time.Parse(time.RFC3339, issueResp.CreatedAt)
	}
	if issueResp.UpdatedAt != "" {
		resource.UpdatedAt, _ = time.Parse(time.RFC3339, issueResp.UpdatedAt)
	}

	for _, label := range issueResp.Labels {
		resource.Labels = append(resource.Labels, label.Name)
	}

	return resource, nil
}

func (c *githubClient) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}

// ValidateURL 是用于验证 URL 的辅助函数
func ValidateURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if u.Host != "github.com" {
		return nil, fmt.Errorf("not a GitHub URL")
	}
	return u, nil
}
