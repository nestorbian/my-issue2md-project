package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// T6.3.2: End-to-end test with mock GitHub API
func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		userLink   bool
		wantTitle  string
		wantAuthor string
		wantState  string
	}{
		// Issue 完整流程
		{
			name:       "e2e issue flow",
			url:        "https://github.com/owner/repo/issues/123",
			userLink:   false,
			wantTitle:  "E2E Test Issue",
			wantAuthor: "e2euser",
			wantState:  "open",
		},
		// PR 完整流程 with user link
		{
			name:       "e2e pr flow with user link",
			url:        "https://github.com/owner/repo/pull/456",
			userLink:   true,
			wantTitle:  "E2E Test PR",
			wantAuthor: "e2eauthor",
			wantState:  "closed",
		},
		// Discussion 完整流程
		{
			name:       "e2e discussion flow",
			url:        "https://github.com/owner/repo/discussions/789",
			userLink:   false,
			wantTitle:  "E2E Test Discussion",
			wantAuthor: "e2ediscauthor",
			wantState:  "open",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 mock server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path == "/repos/owner/repo/issues/123" || r.URL.Path == "/repos/owner/repo/pulls/456" {
					resp := map[string]interface{}{
						"title":      tt.wantTitle,
						"user":       map[string]string{"login": tt.wantAuthor},
						"created_at": "2026-03-20T10:30:00Z",
						"updated_at": "2026-03-20T11:00:00Z",
						"state":      tt.wantState,
						"body":       "E2E test body content",
						"labels":     []map[string]string{{"name": "e2e-label"}},
						"comments":   0,
						"html_url":   "https://github.com/owner/repo/issues/123",
					}
					json.NewEncoder(w).Encode(resp)
				} else if r.URL.Path == "/repos/owner/repo/discussions/789" {
					resp := map[string]interface{}{
						"title":      tt.wantTitle,
						"author":     map[string]string{"login": tt.wantAuthor},
						"created_at": "2026-03-20T10:30:00Z",
						"updated_at": "2026-03-20T11:00:00Z",
						"url":        "https://github.com/owner/repo/discussions/789",
						"body":       "E2E discussion body",
						"category":   map[string]string{"name": "General"},
						"answer":     nil,
						"comments": map[string]interface{}{
							"nodes": []map[string]interface{}{
								{
									"author":    map[string]string{"login": "commenter1"},
									"created_at": "2026-03-20T12:00:00Z",
									"body":      "First comment",
								},
							},
						},
					}
					json.NewEncoder(w).Encode(resp)
				}
			}))
			defer ts.Close()

			cfg := &Config{
				Token:     "fake-token",
				UserLink:  tt.userLink,
				OutputFile: "",
				URL:       tt.url,
			}

			err := run(context.Background(), cfg, ts.URL)
			if err != nil {
				t.Errorf("run() error: %v", err)
			}
		})
	}
}

// TestOutputFileWrite tests file output via e2e
func TestOutputFileWrite(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"title":      "File Output Test",
			"user":       map[string]string{"login": "testuser"},
			"created_at": "2026-03-20T10:30:00Z",
			"updated_at": "2026-03-20T11:00:00Z",
			"state":      "open",
			"body":       "Test body for file output",
			"labels":     []map[string]string{},
			"comments":   0,
			"html_url":   "https://github.com/owner/repo/issues/123",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "issue2md_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	// Clean up after test
	defer os.RemoveAll(tmpDir)

	cfg := &Config{
		Token:     "fake-token",
		UserLink:  false,
		OutputFile: tmpDir + string(os.PathSeparator) + "output.md",
		URL:       "https://github.com/owner/repo/issues/123",
	}

	err = run(context.Background(), cfg, ts.URL)
	if err != nil {
		t.Errorf("run() error: %v", err)
	}
}
