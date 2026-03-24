package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// T6.1.1: Test environment variable reading
func TestGetTokenFromEnv(t *testing.T) {
	tests := []struct {
		name      string
		envValue  string
		setEnv    bool
		wantToken string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "GITHUB_TOKEN is set",
			envValue:  "ghp_test1234567890abcdef",
			setEnv:    true,
			wantToken: "ghp_test1234567890abcdef",
			wantErr:   false,
		},
		{
			name:      "GITHUB_TOKEN is empty string",
			envValue:  "",
			setEnv:    true,
			wantToken: "",
			wantErr:   true,
			errMsg:    "GITHUB_TOKEN environment variable is not set",
		},
		{
			name:      "GITHUB_TOKEN is not set",
			envValue:  "",
			setEnv:    false,
			wantToken: "",
			wantErr:   true,
			errMsg:    "GITHUB_TOKEN environment variable is not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment before each test
			if tt.setEnv {
				os.Setenv("GITHUB_TOKEN", tt.envValue)
				defer os.Unsetenv("GITHUB_TOKEN")
			} else {
				os.Unsetenv("GITHUB_TOKEN")
			}

			got, err := getTokenFromEnv()

			if tt.wantErr {
				if err == nil {
					t.Errorf("getTokenFromEnv() error = nil, want error")
					return
				}
				if tt.errMsg != "" && !containsString(err.Error(), tt.errMsg) {
					t.Errorf("getTokenFromEnv() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("getTokenFromEnv() unexpected error: %v", err)
				return
			}
			if got != tt.wantToken {
				t.Errorf("getTokenFromEnv() = %q, want %q", got, tt.wantToken)
			}
		})
	}
}

// T6.1.2: Test flags parsing
func TestParseFlags(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantCfg    *Config
		wantErr    bool
		errMsg     string
	}{
		{
			name: "no flags, no output file",
			args: []string{"https://github.com/owner/repo/issues/123"},
			wantCfg: &Config{
				Token:     "",
				UserLink:  false,
				OutputFile: "",
				URL:       "https://github.com/owner/repo/issues/123",
			},
			wantErr: false,
		},
		{
			name: "with --user-link flag",
			args: []string{"--user-link", "https://github.com/owner/repo/issues/123"},
			wantCfg: &Config{
				Token:     "",
				UserLink:  true,
				OutputFile: "",
				URL:       "https://github.com/owner/repo/issues/123",
			},
			wantErr: false,
		},
		{
			name: "with -u short flag",
			args: []string{"-u", "https://github.com/owner/repo/pull/456"},
			wantCfg: &Config{
				Token:     "",
				UserLink:  true,
				OutputFile: "",
				URL:       "https://github.com/owner/repo/pull/456",
			},
			wantErr: false,
		},
		{
			name: "with output file",
			args: []string{"https://github.com/owner/repo/issues/123", "output.md"},
			wantCfg: &Config{
				Token:     "",
				UserLink:  false,
				OutputFile: "output.md",
				URL:       "https://github.com/owner/repo/issues/123",
			},
			wantErr: false,
		},
		{
			name: "with --user-link and output file",
			args: []string{"--user-link", "https://github.com/owner/repo/discussions/789", "discussion.md"},
			wantCfg: &Config{
				Token:     "",
				UserLink:  true,
				OutputFile: "discussion.md",
				URL:       "https://github.com/owner/repo/discussions/789",
			},
			wantErr: false,
		},
		{
			name: "with -u short flag and output file",
			args: []string{"-u", "https://github.com/owner/repo/issues/123", "issue-123.md"},
			wantCfg: &Config{
				Token:     "",
				UserLink:  true,
				OutputFile: "issue-123.md",
				URL:       "https://github.com/owner/repo/issues/123",
			},
			wantErr: false,
		},
		{
			name:    "no arguments",
			args:    []string{},
			wantCfg: nil,
			wantErr: true,
			errMsg:  "requires a URL argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFlags(tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseFlags(%v) error = nil, want error", tt.args)
					return
				}
				if tt.errMsg != "" && !containsString(err.Error(), tt.errMsg) {
					t.Errorf("parseFlags(%v) error = %q, want error containing %q", tt.args, err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("parseFlags(%v) unexpected error: %v", tt.args, err)
				return
			}
			if got == nil {
				t.Errorf("parseFlags(%v) got = nil, want %+v", tt.args, tt.wantCfg)
				return
			}
			if got.UserLink != tt.wantCfg.UserLink {
				t.Errorf("parseFlags(%v) UserLink = %v, want %v", tt.args, got.UserLink, tt.wantCfg.UserLink)
			}
			if got.OutputFile != tt.wantCfg.OutputFile {
				t.Errorf("parseFlags(%v) OutputFile = %v, want %v", tt.args, got.OutputFile, tt.wantCfg.OutputFile)
			}
			if got.URL != tt.wantCfg.URL {
				t.Errorf("parseFlags(%v) URL = %v, want %v", tt.args, got.URL, tt.wantCfg.URL)
			}
		})
	}
}

// T6.1.3: Test token missing error
func TestTokenMissingError(t *testing.T) {
	tests := []struct {
		name           string
		envValue       string
		setEnv         bool
		wantErr        bool
		errMsg         string
	}{
		{
			name:           "GITHUB_TOKEN not set",
			envValue:       "",
			setEnv:         false,
			wantErr:        true,
			errMsg:         "GITHUB_TOKEN environment variable is not set",
		},
		{
			name:           "GITHUB_TOKEN is empty",
			envValue:       "",
			setEnv:         true,
			wantErr:        true,
			errMsg:         "GITHUB_TOKEN environment variable is not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment before each test
			if tt.setEnv {
				os.Setenv("GITHUB_TOKEN", tt.envValue)
				defer os.Unsetenv("GITHUB_TOKEN")
			} else {
				os.Unsetenv("GITHUB_TOKEN")
			}

			token, err := getTokenFromEnv()

			if tt.wantErr {
				if err == nil {
					t.Errorf("getTokenFromEnv() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if !containsString(err.Error(), tt.errMsg) {
					t.Errorf("getTokenFromEnv() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
				return
			}

			// If no error expected, token should be valid
			if token == "" {
				t.Errorf("getTokenFromEnv() returned empty token without error")
			}
		})
	}
}

// containsString checks if s contains substr
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && containsStringHelper(s, substr)
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// T6.2.1: Test complete flow - URL → Parser → Fetcher → Converter → Writer
func TestRunCompleteFlow(t *testing.T) {
	// 创建 mock server，根据不同路径返回不同数据
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/repos/owner/repo/issues/123" || r.URL.Path == "/repos/owner/repo/pulls/456" {
			resp := map[string]interface{}{
				"title":      "Test Issue",
				"user":       map[string]string{"login": "testuser"},
				"created_at": "2026-03-20T10:30:00Z",
				"updated_at": "2026-03-20T11:00:00Z",
				"state":      "open",
				"body":       "Test body content",
				"labels":     []map[string]string{{"name": "bug"}},
				"comments":   0,
				"html_url":   "https://github.com/owner/repo/issues/123",
			}
			json.NewEncoder(w).Encode(resp)
		} else if r.URL.Path == "/repos/owner/repo/discussions/789" {
			resp := map[string]interface{}{
				"title":      "Test Discussion",
				"author":     map[string]string{"login": "testuser"},
				"created_at": "2026-03-20T10:30:00Z",
				"updated_at": "2026-03-20T11:00:00Z",
				"url":        "https://github.com/owner/repo/discussions/789",
				"body":       "Test discussion body",
				"category":   map[string]string{"name": "General"},
				"answer":     nil,
				"comments": map[string]interface{}{
					"nodes": []map[string]interface{}{},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer ts.Close()

	tests := []struct {
		name       string
		url        string
		userLink   bool
		wantErr    bool
		errMsg     string
	}{
		{
			name:     "issue URL complete flow",
			url:      "https://github.com/owner/repo/issues/123",
			userLink: false,
			wantErr:  false,
		},
		{
			name:     "pull request URL complete flow",
			url:      "https://github.com/owner/repo/pull/456",
			userLink: true,
			wantErr:  false,
		},
		{
			name:     "discussion URL complete flow",
			url:      "https://github.com/owner/repo/discussions/789",
			userLink: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Token:     "fake-token",
				UserLink:  tt.userLink,
				OutputFile: "",
				URL:       tt.url,
			}

			err := run(context.Background(), cfg, ts.URL)

			if tt.wantErr {
				if err == nil {
					t.Errorf("run() error = nil, want error")
					return
				}
				if tt.errMsg != "" && !containsString(err.Error(), tt.errMsg) {
					t.Errorf("run() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("run() unexpected error: %v", err)
			}
		})
	}
}

// T6.2.2: Test error propagation from different layers
func TestRunErrorPropagation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		handler func(w http.ResponseWriter, r *http.Request)
		wantErr bool
		errMsg  string
	}{
		{
			name: "invalid URL - sublink rejected",
			cfg: &Config{
				Token:     "fake-token",
				UserLink:  false,
				OutputFile: "",
				URL:       "github.com/owner/repo/issues/123#issuecomment-456",
			},
			handler: nil, // 不需要 handler，URL 解析会失败
			wantErr: true,
			errMsg:  "unsupported URL type",
		},
		{
			name: "non-github URL",
			cfg: &Config{
				Token:     "fake-token",
				UserLink:  false,
				OutputFile: "",
				URL:       "https://google.com/owner/repo/issues/123",
			},
			handler: nil,
			wantErr: true,
			errMsg:  "not a GitHub URL",
		},
		{
			name: "rate limit exceeded",
			cfg: &Config{
				Token:     "fake-token",
				UserLink:  false,
				OutputFile: "",
				URL:       "https://github.com/owner/repo/issues/123",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"message": "rate limit exceeded"})
			},
			wantErr: true,
			errMsg:  "rate limit",
		},
		{
			name: "private repository",
			cfg: &Config{
				Token:     "fake-token",
				UserLink:  false,
				OutputFile: "",
				URL:       "https://github.com/private/repo/issues/123",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name: "file already exists",
			cfg: &Config{
				Token:     "fake-token",
				UserLink:  false,
				OutputFile: "/tmp/existing_test.md",
				URL:       "https://github.com/owner/repo/issues/123",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				// 返回正常的 Issue 数据
				resp := map[string]interface{}{
					"title":    "Test Issue",
					"user":     map[string]string{"login": "testuser"},
					"created_at": "2026-03-20T10:30:00Z",
					"updated_at": "2026-03-20T11:00:00Z",
					"state":    "open",
					"body":     "Test body",
					"labels":   []map[string]string{{"name": "bug"}},
					"comments": 0,
					"html_url": "https://github.com/owner/repo/issues/123",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			},
			wantErr: true,
			errMsg:  "file already exists",
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

			err := run(context.Background(), tt.cfg, baseURL)

			if tt.wantErr {
				if err == nil {
					t.Errorf("run() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !containsString(err.Error(), tt.errMsg) {
					t.Errorf("run() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("run() unexpected error: %v", err)
			}
		})
	}
}

