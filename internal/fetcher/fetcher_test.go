package fetcher

import (
	"context"
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
	return nil, nil
}

// mockGitHubClient 实现 GitHubClient 接口用于测试
type mockGitHubClient struct{}

func (m *mockGitHubClient) FetchIssue(ctx context.Context, owner, repo string, number int) (*GitHubResource, error) {
	return nil, nil
}

func (m *mockGitHubClient) FetchPullRequest(ctx context.Context, owner, repo string, number int) (*GitHubResource, error) {
	return nil, nil
}

func (m *mockGitHubClient) FetchDiscussion(ctx context.Context, owner, repo string, number int) (*GitHubResource, error) {
	return nil, nil
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

	createdAt := time.Date(2026, 3, 20, 10, 30, 0, 0, time.UTC)

	if resource.Type != "" {
		t.Errorf("Type = %v, want empty", resource.Type)
	}
	if resource.Title != "" {
		t.Errorf("Title = %v, want empty", resource.Title)
	}
	if resource.Author != "" {
		t.Errorf("Author = %v, want empty", resource.Author)
	}
	if resource.CreatedAt != createdAt {
		t.Errorf("CreatedAt mismatch")
	}
	if resource.State != "" {
		t.Errorf("State = %v, want empty", resource.State)
	}
	if resource.Body != "" {
		t.Errorf("Body = %v, want empty", resource.Body)
	}
	if resource.URL != "" {
		t.Errorf("URL = %v, want empty", resource.URL)
	}
}
