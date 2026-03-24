package fetcher

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bigwhite/issue2md/internal/parser"
)

// TestIntegration_FetchIssueWorkflow 测试完整的 Issue 获取流程
func TestIntegration_FetchIssueWorkflow(t *testing.T) {
	// 准备 Mock GitHub API Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/repos/owner/repo/issues/123" {
			// 返回模拟的 Issue JSON 响应
			resp := map[string]interface{}{
				"title":    "Integration Test Issue",
				"user":     map[string]string{"login": "testuser"},
				"created_at": "2026-03-20T10:30:00Z",
				"updated_at": "2026-03-20T11:00:00Z",
				"state":    "open",
				"body":     "This is an integration test issue",
				"labels": []map[string]string{
					{"name": "test"},
					{"name": "integration"},
				},
				"comments": 1,
				"html_url": "https://github.com/owner/repo/issues/123",
			}
			json.NewEncoder(w).Encode(resp)
		} else if r.URL.Path == "/repos/owner/repo/issues/123/comments" {
			// 返回模拟的评论 JSON 响应
			resp := []map[string]interface{}{
				{
					"body":      "Test comment",
					"user":      map[string]string{"login": "commenter1"},
					"created_at": "2026-03-20T12:00:00Z",
					"updated_at": "2026-03-20T12:00:00Z",
				},
			}
			json.NewEncoder(w).Encode(resp)
		} else {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	// 创建 Client 和 Fetcher
	client := NewGitHubClient(ts.URL, "fake-token")
	fetcher := NewGitHubFetcher(client)

	// 解析 URL
	parsedURL, err := parser.Parse("https://github.com/owner/repo/issues/123")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// 执行 Fetch
	resource, err := fetcher.Fetch(context.Background(), parsedURL)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	// 验证结果
	if resource == nil {
		t.Fatal("resource is nil")
	}
	if resource.Type != "issue" {
		t.Errorf("Type = %s, want issue", resource.Type)
	}
	if resource.Title != "Integration Test Issue" {
		t.Errorf("Title = %s, want Integration Test Issue", resource.Title)
	}
	if resource.Author != "testuser" {
		t.Errorf("Author = %s, want testuser", resource.Author)
	}
	if resource.State != "open" {
		t.Errorf("State = %s, want open", resource.State)
	}
	if len(resource.Labels) != 2 {
		t.Errorf("len(Labels) = %d, want 2", len(resource.Labels))
	}
}

// TestIntegration_FetchPullRequestWorkflow 测试完整的 PR 获取流程
func TestIntegration_FetchPullRequestWorkflow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/pulls/456" {
			t.Errorf("unexpected path: got %s, want /repos/owner/repo/pulls/456", r.URL.Path)
		}

		resp := map[string]interface{}{
			"title": "Integration Test PR",
			"user":  map[string]string{"login": "prauthor"},
			"created_at": "2026-03-19T08:00:00Z",
			"updated_at": "2026-03-19T09:00:00Z",
			"state": "closed",
			"body":  "This is a test pull request",
			"labels": []map[string]string{
				{"name": "enhancement"},
			},
			"comments": 0,
			"html_url": "https://github.com/owner/repo/pull/456",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := NewGitHubClient(ts.URL, "fake-token")
	fetcher := NewGitHubFetcher(client)

	parsedURL, err := parser.Parse("https://github.com/owner/repo/pull/456")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	resource, err := fetcher.Fetch(context.Background(), parsedURL)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	if resource == nil {
		t.Fatal("resource is nil")
	}
	if resource.Type != "pull_request" {
		t.Errorf("Type = %s, want pull_request", resource.Type)
	}
	if resource.Title != "Integration Test PR" {
		t.Errorf("Title = %s, want Integration Test PR", resource.Title)
	}
	if resource.Author != "prauthor" {
		t.Errorf("Author = %s, want prauthor", resource.Author)
	}
	if resource.State != "closed" {
		t.Errorf("State = %s, want closed", resource.State)
	}
}

// TestIntegration_FetchDiscussionWorkflow 测试完整的 Discussion 获取流程
func TestIntegration_FetchDiscussionWorkflow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/discussions/789" {
			t.Errorf("unexpected path: got %s, want /repos/owner/repo/discussions/789", r.URL.Path)
		}

		resp := map[string]interface{}{
			"title": "Integration Test Discussion",
			"author": map[string]string{"login": "discauthor"},
			"created_at": "2026-03-18T14:00:00Z",
			"updated_at": "2026-03-18T14:00:00Z",
			"url":  "https://github.com/owner/repo/discussions/789",
			"body": "This is a test discussion",
			"category": map[string]string{"name": "General"},
			"answer": nil,
			"comments": map[string]interface{}{
				"nodes": []map[string]interface{}{
					{
						"author":    map[string]string{"login": "commenter1"},
						"created_at": "2026-03-18T15:00:00Z",
						"body":      "First comment",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := NewGitHubClient(ts.URL, "fake-token")
	fetcher := NewGitHubFetcher(client)

	parsedURL, err := parser.Parse("https://github.com/owner/repo/discussions/789")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	resource, err := fetcher.Fetch(context.Background(), parsedURL)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	if resource == nil {
		t.Fatal("resource is nil")
	}
	if resource.Type != "discussion" {
		t.Errorf("Type = %s, want discussion", resource.Type)
	}
	if resource.Title != "Integration Test Discussion" {
		t.Errorf("Title = %s, want Integration Test Discussion", resource.Title)
	}
	if resource.Author != "discauthor" {
		t.Errorf("Author = %s, want discauthor", resource.Author)
	}
	if len(resource.Comments) != 1 {
		t.Errorf("len(Comments) = %d, want 1", len(resource.Comments))
	}
}

// TestIntegration_FetchNotFound 测试资源不存在的情况
func TestIntegration_FetchNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Not Found"}`))
	}))
	defer ts.Close()

	client := NewGitHubClient(ts.URL, "fake-token")
	fetcher := NewGitHubFetcher(client)

	parsedURL, err := parser.Parse("https://github.com/owner/repo/issues/999")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, err = fetcher.Fetch(context.Background(), parsedURL)
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

// TestIntegration_FetchRateLimit 测试限流情况
func TestIntegration_FetchRateLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "rate limit exceeded"}`))
	}))
	defer ts.Close()

	client := NewGitHubClient(ts.URL, "fake-token")
	fetcher := NewGitHubFetcher(client)

	parsedURL, err := parser.Parse("https://github.com/owner/repo/issues/123")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, err = fetcher.Fetch(context.Background(), parsedURL)
	if err == nil {
		t.Fatal("expected error for rate limit, got nil")
	}
}
