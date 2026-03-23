package fetcher

import (
	"context"
	"fmt"

	"github.com/bigwhite/issue2md/internal/parser"
)

// Fetcher 获取 GitHub 资源
type Fetcher interface {
	// Fetch 根据 ParsedURL 获取 GitHubResource
	Fetch(ctx context.Context, url *parser.ParsedURL) (*GitHubResource, error)
}

// GitHubClient GitHub API 客户端接口 (便于测试)
type GitHubClient interface {
	FetchIssue(ctx context.Context, owner, repo string, number int) (*GitHubResource, error)
	FetchPullRequest(ctx context.Context, owner, repo string, number int) (*GitHubResource, error)
	FetchDiscussion(ctx context.Context, owner, repo string, number int) (*GitHubResource, error)
}

// githubFetcher 实现 Fetcher 接口
type githubFetcher struct {
	client GitHubClient
}

// NewGitHubFetcher 创建一个新的 Fetcher
func NewGitHubFetcher(client GitHubClient) Fetcher {
	return &githubFetcher{client: client}
}

// Fetch 根据 URLType 分发到不同的获取方法
func (f *githubFetcher) Fetch(ctx context.Context, url *parser.ParsedURL) (*GitHubResource, error) {
	switch url.Type {
	case parser.URLTypeIssue:
		return f.client.FetchIssue(ctx, url.Owner, url.Repo, url.Number)
	case parser.URLTypePullRequest:
		return f.client.FetchPullRequest(ctx, url.Owner, url.Repo, url.Number)
	case parser.URLTypeDiscussion:
		return f.client.FetchDiscussion(ctx, url.Owner, url.Repo, url.Number)
	default:
		return nil, fmt.Errorf("unsupported URL type: %v", url.Type)
	}
}
