package fetcher

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bigwhite/issue2md/internal/parser"
)

func TestFetcherInterface(t *testing.T) {
	// Test that Fetcher interface is implemented correctly
	var fetcher Fetcher = &mockFetcher{}
	if fetcher == nil {
		t.Error("Fetcher interface should be implemented")
	}
}

func TestGitHubClientInterface(t *testing.T) {
	// Test that GitHubClient interface is implemented correctly
	var client GitHubClient = &mockGitHubClient{}
	if client == nil {
		t.Error("GitHubClient interface should be implemented")
	}
}

func TestFetcherFetch_Signature(t *testing.T) {
	tests := []struct {
		name    string
		url     *parser.ParsedURL
		wantErr bool
	}{
		{
			name: "Issue URL",
			url: &parser.ParsedURL{
				Type:   parser.URLTypeIssue,
				Owner:  "owner",
				Repo:   "repo",
				Number: 123,
				RawURL: "https://github.com/owner/repo/issues/123",
			},
			wantErr: false,
		},
		{
			name: "PullRequest URL",
			url: &parser.ParsedURL{
				Type:   parser.URLTypePullRequest,
				Owner:  "owner",
				Repo:   "repo",
				Number: 456,
				RawURL: "https://github.com/owner/repo/pull/456",
			},
			wantErr: false,
		},
		{
			name: "Discussion URL",
			url: &parser.ParsedURL{
				Type:   parser.URLTypeDiscussion,
				Owner:  "owner",
				Repo:   "repo",
				Number: 789,
				RawURL: "https://github.com/owner/repo/discussions/789",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcher := &mockFetcher{}
			ctx := context.Background()
			got, err := fetcher.Fetch(ctx, tt.url)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Fetch() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Errorf("Fetch() unexpected error: %v", err)
				return
			}
			if got == nil {
				t.Errorf("Fetch() got = nil")
			}
		})
	}
}

func TestGitHubClientMethods(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		owner     string
		repo      string
		number    int
		wantErr   bool
	}{
		{
			name:   "FetchIssue",
			method: "FetchIssue",
			owner:  "owner",
			repo:   "repo",
			number: 123,
			wantErr: false,
		},
		{
			name:   "FetchPullRequest",
			method: "FetchPullRequest",
			owner:  "owner",
			repo:   "repo",
			number: 456,
			wantErr: false,
		},
		{
			name:   "FetchDiscussion",
			method: "FetchDiscussion",
			owner:  "owner",
			repo:   "repo",
			number: 789,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockGitHubClient{}
			ctx := context.Background()

			var got *GitHubResource
			var err error

			switch tt.method {
			case "FetchIssue":
				got, err = client.FetchIssue(ctx, tt.owner, tt.repo, tt.number)
			case "FetchPullRequest":
				got, err = client.FetchPullRequest(ctx, tt.owner, tt.repo, tt.number)
			case "FetchDiscussion":
				got, err = client.FetchDiscussion(ctx, tt.owner, tt.repo, tt.number)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("%s() error = nil, want error", tt.method)
				}
				return
			}
			if err != nil {
				t.Errorf("%s() unexpected error: %v", tt.method, err)
				return
			}
			if got == nil {
				t.Errorf("%s() got = nil", tt.method)
			}
		})
	}
}

// mockFetcher 实现 Fetcher 接口用于测试
type mockFetcher struct{}

func (m *mockFetcher) Fetch(ctx context.Context, url *parser.ParsedURL) (*GitHubResource, error) {
	// 返回一个模拟的 GitHubResource
	resource := &GitHubResource{
		Type:      "issue",
		Title:     "Mock Issue",
		Author:    "mockuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		State:     "open",
		Body:      "Mock body content",
		Labels:    []string{"mock-label"},
		Comments:  []Comment{},
		URL:       url.RawURL,
	}
	return resource, nil
}

// mockGitHubClient 实现 GitHubClient 接口用于测试
type mockGitHubClient struct{}

func (m *mockGitHubClient) FetchIssue(ctx context.Context, owner, repo string, number int) (*GitHubResource, error) {
	return &GitHubResource{
		Type:      "issue",
		Title:     "Mock Issue",
		Author:    "mockuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		State:     "open",
		Body:      "Mock body content",
		Labels:    []string{"mock-label"},
		Comments:  []Comment{},
		URL:       fmt.Sprintf("https://github.com/%s/%s/issues/%d", owner, repo, number),
	}, nil
}

func (m *mockGitHubClient) FetchPullRequest(ctx context.Context, owner, repo string, number int) (*GitHubResource, error) {
	return &GitHubResource{
		Type:      "pull_request",
		Title:     "Mock PR",
		Author:    "mockuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		State:     "closed",
		Body:      "Mock PR body",
		Labels:    []string{},
		Comments:  []Comment{},
		URL:       fmt.Sprintf("https://github.com/%s/%s/pull/%d", owner, repo, number),
	}, nil
}

func (m *mockGitHubClient) FetchDiscussion(ctx context.Context, owner, repo string, number int) (*GitHubResource, error) {
	return &GitHubResource{
		Type:      "discussion",
		Title:     "Mock Discussion",
		Author:    "mockuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		State:     "open",
		Body:      "Mock discussion body",
		Labels:    []string{},
		Comments:  []Comment{},
		URL:       fmt.Sprintf("https://github.com/%s/%s/discussions/%d", owner, repo, number),
	}, nil
}

func TestGitHubResource_NotNil(t *testing.T) {
	client := &mockGitHubClient{}
	ctx := context.Background()

	resource, err := client.FetchIssue(ctx, "owner", "repo", 123)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resource == nil {
		t.Errorf("FetchIssue() returned nil resource")
	}
}

func TestGitHubResource_FieldsComplete(t *testing.T) {
	client := &mockGitHubClient{}
	ctx := context.Background()

	resource, _ := client.FetchIssue(ctx, "testowner", "testrepo", 999)
	if resource == nil {
		t.Skip("resource is nil, skipping field checks")
	}

	// 验证 mock 返回的资源包含正确的字段值
	if resource.Type != "issue" {
		t.Errorf("Type = %v, want issue", resource.Type)
	}
	if resource.Title != "Mock Issue" {
		t.Errorf("Title = %v, want Mock Issue", resource.Title)
	}
	if resource.Author != "mockuser" {
		t.Errorf("Author = %v, want mockuser", resource.Author)
	}
	if resource.State != "open" {
		t.Errorf("State = %v, want open", resource.State)
	}
	if resource.Body != "Mock body content" {
		t.Errorf("Body = %v, want Mock body content", resource.Body)
	}
	if resource.URL != "https://github.com/testowner/testrepo/issues/999" {
		t.Errorf("URL = %v, want https://github.com/testowner/testrepo/issues/999", resource.URL)
	}
}
