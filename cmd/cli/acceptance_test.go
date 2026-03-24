package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bigwhite/issue2md/internal/converter"
	"github.com/bigwhite/issue2md/internal/fetcher"
	"github.com/bigwhite/issue2md/internal/parser"
	"github.com/bigwhite/issue2md/internal/writer"
)

// T6.3.1: Acceptance test - Error handling according to spec 5.4
func TestAcceptanceErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		token      string
		handler    func(w http.ResponseWriter, r *http.Request)
		wantErr    bool
		errSubstr  string
	}{
		// spec 5.4: 限流 - API 返回 403 限流
		{
			name:  "rate limit exceeded",
			url:   "https://github.com/owner/repo/issues/123",
			token: "fake-token",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"message": "rate limit exceeded"})
			},
			wantErr:   true,
			errSubstr: "rate limit",
		},
		// spec 5.4: 私有仓库 - 无 Token 访问私有仓库
		{
			name:  "private repository",
			url:   "https://github.com/private/repo/issues/123",
			token: "fake-token",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
			},
			wantErr:   true,
			errSubstr: "not found",
		},
		// spec 5.4: 无 Token - 未设置 GITHUB_TOKEN
		{
			name:       "no token set",
			url:        "https://github.com/owner/repo/issues/123",
			token:      "",
			handler:    nil,
			wantErr:    true,
			errSubstr:  "GITHUB_TOKEN",
		},
		// spec 5.4: 未找到 - Issue/PR/Discussion 不存在
		{
			name:  "not found - issue",
			url:   "https://github.com/owner/repo/issues/999",
			token: "fake-token",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
			},
			wantErr:   true,
			errSubstr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts *httptest.Server
			if tt.handler != nil {
				ts = httptest.NewServer(http.HandlerFunc(tt.handler))
				defer ts.Close()
			}

			baseURL := "https://api.github.com"
			if ts != nil {
				baseURL = ts.URL
			}

			cfg := &Config{
				Token:     tt.token,
				UserLink:  false,
				OutputFile: "",
				URL:       tt.url,
			}

			err := run(context.Background(), cfg, baseURL)

			if tt.wantErr {
				if err == nil {
					t.Errorf("run() error = nil, want error containing %q", tt.errSubstr)
					return
				}
				if !containsString(err.Error(), tt.errSubstr) {
					t.Errorf("run() error = %q, want error containing %q", err.Error(), tt.errSubstr)
				}
				return
			}

			if err != nil {
				t.Errorf("run() unexpected error: %v", err)
			}
		})
	}
}

// Verify all packages are properly integrated
func TestPackageIntegration(t *testing.T) {
	// Verify parser.Parse works
	_, err := parser.Parse("https://github.com/owner/repo/issues/123")
	if err != nil {
		t.Errorf("parser.Parse failed: %v", err)
	}

	// Verify converter.NewConverter works
	conv := converter.NewConverter()
	if conv == nil {
		t.Error("converter.NewConverter returned nil")
	}

	// Verify writer.New works
	w := writer.New()
	if w == nil {
		t.Error("writer.New returned nil")
	}

	// Verify fetcher.NewGitHubClient works
	client := fetcher.NewGitHubClient("https://api.github.com", "test-token")
	if client == nil {
		t.Error("fetcher.NewGitHubClient returned nil")
	}
}
