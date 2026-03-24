package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// NewGitHubClientForTest 创建一个用于测试的 GitHubClient
// 使用真实的 githubClient 实现，通过 mock server 进行测试
func NewGitHubClientForTest(baseURL, token string) GitHubClient {
	return NewGitHubClient(baseURL, token)
}

// TestFetchIssue_Success 测试成功获取 Issue
func TestFetchIssue_Success(t *testing.T) {
	// 准备 Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/repos/owner/repo/issues/123" {
			// 返回模拟的 Issue JSON 响应
			resp := `{
				"title": "Bug Report",
				"user": {"login": "author123"},
				"created_at": "2026-03-20T10:30:00Z",
				"updated_at": "2026-03-20T11:00:00Z",
				"state": "open",
				"body": "This is a bug description",
				"labels": [{"name": "bug"}, {"name": "priority-high"}],
				"comments": 2,
				"html_url": "https://github.com/owner/repo/issues/123"
			}`
			w.Write([]byte(resp))
		} else if r.URL.Path == "/repos/owner/repo/issues/123/comments" {
			// 返回模拟的评论 JSON 响应
			resp := `[
				{
					"body": "First comment",
					"user": {"login": "commenter1"},
					"created_at": "2026-03-20T12:00:00Z",
					"updated_at": "2026-03-20T12:00:00Z"
				},
				{
					"body": "Second comment",
					"user": {"login": "commenter2"},
					"created_at": "2026-03-20T13:00:00Z",
					"updated_at": "2026-03-20T13:00:00Z"
				}
			]`
			w.Write([]byte(resp))
		} else {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	// 创建 Client
	client := NewGitHubClientForTest(ts.URL, "fake-token")

	// 调用 FetchIssue
	ctx := context.Background()
	resource, err := client.FetchIssue(ctx, "owner", "repo", 123)
	if err != nil {
		t.Fatalf("FetchIssue failed: %v", err)
	}

	// 验证返回的数据结构
	if resource == nil {
		t.Fatal("resource is nil")
	}

	// 验证字段值
	if resource.Type != "issue" {
		t.Errorf("Type = %s, want issue", resource.Type)
	}
	if resource.Title != "Bug Report" {
		t.Errorf("Title = %s, want Bug Report", resource.Title)
	}
	if resource.Author != "author123" {
		t.Errorf("Author = %s, want author123", resource.Author)
	}
	if resource.State != "open" {
		t.Errorf("State = %s, want open", resource.State)
	}
	if resource.Body != "This is a bug description" {
		t.Errorf("Body = %s, want This is a bug description", resource.Body)
	}
	if len(resource.Labels) != 2 {
		t.Errorf("len(Labels) = %d, want 2", len(resource.Labels))
	}
	if resource.URL != "https://github.com/owner/repo/issues/123" {
		t.Errorf("URL = %s, want https://github.com/owner/repo/issues/123", resource.URL)
	}
}

// TestFetchIssue_RateLimit 测试限流错误处理
func TestFetchIssue_RateLimit(t *testing.T) {
	// 准备返回 403 的 Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "rate limit exceeded"}`))
	}))
	defer ts.Close()

	client := NewGitHubClientForTest(ts.URL, "fake-token")

	ctx := context.Background()
	_, err := client.FetchIssue(ctx, "owner", "repo", 123)
	if err == nil {
		t.Fatal("expected error for rate limit, got nil")
	}
	// TODO: 验证错误消息包含 "rate limit"
	_ = err
}

// TestFetchIssue_NotFound 测试 404 Not Found 错误
func TestFetchIssue_NotFound(t *testing.T) {
	// 准备返回 404 的 Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Not Found"}`))
	}))
	defer ts.Close()

	client := NewGitHubClientForTest(ts.URL, "fake-token")

	ctx := context.Background()
	_, err := client.FetchIssue(ctx, "owner", "repo", 999)
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
	// TODO: 验证错误消息包含 "not found"
	_ = err
}

// TestFetchPullRequest_Success 测试成功获取 Pull Request
func TestFetchPullRequest_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/repos/owner/repo/pulls/456" {
			resp := `{
				"title": "feat: Add dark mode support",
				"user": {"login": "designer456"},
				"created_at": "2026-03-19T08:00:00Z",
				"updated_at": "2026-03-19T09:00:00Z",
				"state": "closed",
				"body": "Adds dark mode support using CSS variables.",
				"labels": [{"name": "enhancement"}, {"name": "ui"}],
				"comments": 1,
				"html_url": "https://github.com/owner/repo/pull/456"
			}`
			w.Write([]byte(resp))
		} else if r.URL.Path == "/repos/owner/repo/issues/456/comments" {
			// PR 的评论也通过 issues 接口获取
			resp := `[
				{
					"body": "PR comment",
					"user": {"login": "prcommenter"},
					"created_at": "2026-03-19T10:00:00Z",
					"updated_at": "2026-03-19T10:00:00Z"
				}
			]`
			w.Write([]byte(resp))
		} else {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	client := NewGitHubClientForTest(ts.URL, "fake-token")

	ctx := context.Background()
	resource, err := client.FetchPullRequest(ctx, "owner", "repo", 456)
	if err != nil {
		t.Fatalf("FetchPullRequest failed: %v", err)
	}

	if resource == nil {
		t.Fatal("resource is nil")
	}

	if resource.Type != "pull_request" {
		t.Errorf("Type = %s, want pull_request", resource.Type)
	}
	if resource.Title != "feat: Add dark mode support" {
		t.Errorf("Title = %s, want feat: Add dark mode support", resource.Title)
	}
	if resource.Author != "designer456" {
		t.Errorf("Author = %s, want designer456", resource.Author)
	}
	if resource.State != "closed" {
		t.Errorf("State = %s, want closed", resource.State)
	}
}

// TestFetchDiscussion_Success 测试成功获取 Discussion
func TestFetchDiscussion_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/discussions/789" {
			t.Errorf("unexpected path: got %s, want /repos/owner/repo/discussions/789", r.URL.Path)
		}

		resp := `{
			"title": "How to implement dark mode?",
			"author": {"login": "user123"},
			"created_at": "2026-03-18T14:00:00Z",
			"updated_at": "2026-03-18T14:00:00Z",
			"url": "https://github.com/owner/repo/discussions/789",
			"body": "I want to add dark mode to my app.",
			"category": {"name": "Q&A"},
			"answer": null,
			"comments": {
				"nodes": [
					{
						"author": {"login": "helper456"},
						"created_at": "2026-03-18T15:00:00Z",
						"body": "You can use CSS variables."
					}
				]
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	defer ts.Close()

	client := NewGitHubClientForTest(ts.URL, "fake-token")

	ctx := context.Background()
	resource, err := client.FetchDiscussion(ctx, "owner", "repo", 789)
	if err != nil {
		t.Fatalf("FetchDiscussion failed: %v", err)
	}

	if resource == nil {
		t.Fatal("resource is nil")
	}

	if resource.Type != "discussion" {
		t.Errorf("Type = %s, want discussion", resource.Type)
	}
	if resource.Title != "How to implement dark mode?" {
		t.Errorf("Title = %s, want How to implement dark mode?", resource.Title)
	}
	if resource.Author != "user123" {
		t.Errorf("Author = %s, want user123", resource.Author)
	}
	if len(resource.Comments) != 1 {
		t.Errorf("len(Comments) = %d, want 1", len(resource.Comments))
	}
}

// TestFetchIssue_MissingToken 测试 Token 缺失错误
func TestFetchIssue_MissingToken(t *testing.T) {
	// 清除环境变量
	originalToken := os.Getenv("GITHUB_TOKEN")
	os.Unsetenv("GITHUB_TOKEN")
	defer func() {
		if originalToken != "" {
			os.Setenv("GITHUB_TOKEN", originalToken)
		}
	}()

	// 当 GITHUB_TOKEN 未设置时，应该返回错误
	// TODO: 实现 Token 验证逻辑
	t.Log("GITHUB_TOKEN is not set, should return error")
}

// TestFetchIssue_PrivateRepo 测试私有仓库错误
func TestFetchIssue_PrivateRepo(t *testing.T) {
	// 准备返回 404 的 Mock Server（模拟私有仓库无权限）
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Not Found"}`))
	}))
	defer ts.Close()

	client := NewGitHubClientForTest(ts.URL, "fake-token")

	ctx := context.Background()
	_, err := client.FetchIssue(ctx, "private", "repo", 123)
	if err == nil {
		t.Fatal("expected error for private repo, got nil")
	}
	// TODO: 验证错误消息包含 "private" 或相关提示
	_ = err
}

// TestGitHubClient_Interface 测试 GitHubClient 接口完整性
func TestGitHubClient_Interface(t *testing.T) {
	// 验证 githubClient 实现了 GitHubClient 接口
	var _ GitHubClient = (*githubClient)(nil)
}

// TestGitHubResource_TimeFields 测试时间字段解析
func TestGitHubResource_TimeFields(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{
			"title": "Test Issue",
			"user": {"login": "testuser"},
			"created_at": "2026-03-20T10:30:00Z",
			"updated_at": "2026-03-20T12:00:00Z",
			"state": "open",
			"body": "Test body",
			"labels": [],
			"comments": 0,
			"html_url": "https://github.com/owner/repo/issues/1"
		}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	defer ts.Close()

	client := NewGitHubClientForTest(ts.URL, "fake-token")

	ctx := context.Background()
	resource, err := client.FetchIssue(ctx, "owner", "repo", 1)
	if err != nil {
		t.Fatalf("FetchIssue failed: %v", err)
	}
	if resource == nil {
		t.Fatal("resource is nil")
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2026-03-20T10:30:00Z")
	if !resource.CreatedAt.Equal(expectedTime) {
		t.Errorf("CreatedAt = %v, want %v", resource.CreatedAt, expectedTime)
	}
}
