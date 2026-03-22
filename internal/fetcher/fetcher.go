package fetcher

import (
	"context"

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
