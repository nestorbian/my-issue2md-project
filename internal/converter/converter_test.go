package converter

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/bigwhite/issue2md/internal/fetcher"
)

// TestConverter_Convert_MethodSignature 测试 Convert 方法签名
func TestConverter_Convert_MethodSignature(t *testing.T) {
	tests := []struct {
		name     string
		resource *fetcher.GitHubResource
		userLink bool
		wantErr  bool
	}{
		{
			name:     "nil resource should return error",
			resource: nil,
			userLink: false,
			wantErr:  true,
		},
		{
			name: "issue with minimal data",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Test Issue",
				Author:    "testuser",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				State:     "open",
				Body:      "Test body content",
				Labels:    []string{},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			wantErr:  false,
		},
		{
			name: "pull request",
			resource: &fetcher.GitHubResource{
				Type:      "pull_request",
				Title:     "Test PR",
				Author:    "testuser",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				State:     "closed",
				Body:      "PR body",
				Labels:    []string{"bug"},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/pull/1",
			},
			userLink: false,
			wantErr:  false,
		},
		{
			name: "discussion",
			resource: &fetcher.GitHubResource{
				Type:      "discussion",
				Title:     "Test Discussion",
				Author:    "testuser",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				State:     "open",
				Body:      "Discussion body",
				Labels:    []string{},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/discussions/1",
			},
			userLink: false,
			wantErr:  false,
		},
		{
			name: "issue with userLink true",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Test Issue",
				Author:    "testuser",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				State:     "open",
				Body:      "Test body",
				Labels:    []string{"enhancement"},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: true,
			wantErr:  false,
		},
		{
			name: "issue with labels",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Test Issue",
				Author:    "testuser",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				State:     "open",
				Body:      "Test body",
				Labels:    []string{"bug", "priority-high"},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			wantErr:  false,
		},
		{
			name: "issue with comments",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Test Issue",
				Author:    "testuser",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				State:     "open",
				Body:      "Test body",
				Labels:    []string{},
				Comments: []fetcher.Comment{
					{
						Author:    "commenter",
						CreatedAt: time.Now(),
						Body:      "This is a comment",
					},
				},
				URL: "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &converter{}
			ctx := context.Background()
			_, err := c.Convert(ctx, tt.resource, tt.userLink)

			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// =============================================================================
// 4.5.1 测试 Convert 完整流程
// =============================================================================

func TestConverter_Convert_CompleteFlow(t *testing.T) {
	tests := []struct {
		name      string
		resource  *fetcher.GitHubResource
		userLink  bool
		checkFunc func(t *testing.T, output *Output)
	}{
		{
			name: "issue with all fields",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Bug: Cannot login",
				Author:    "developer123",
				CreatedAt: time.Date(2026, 3, 20, 10, 30, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
				State:     "open",
				Body:      "Error message here",
				Labels:    []string{"bug", "oauth"},
				Comments: []fetcher.Comment{
					{
						Author:    "helper",
						CreatedAt: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
						Body:      "Can you show your config?",
					},
				},
				URL: "https://github.com/owner/repo/issues/123",
			},
			userLink: false,
			checkFunc: func(t *testing.T, output *Output) {
				// 验证 Frontmatter
				if output.Frontmatter == "" {
					t.Error("Frontmatter should not be empty")
				}
				if !strings.Contains(output.Frontmatter, "title:") {
					t.Error("Frontmatter should contain title")
				}
				if !strings.Contains(output.Frontmatter, "author:") {
					t.Error("Frontmatter should contain author")
				}
				if !strings.Contains(output.Frontmatter, "created_at:") {
					t.Error("Frontmatter should contain created_at")
				}
				if !strings.Contains(output.Frontmatter, "type:") {
					t.Error("Frontmatter should contain type")
				}
				if !strings.Contains(output.Frontmatter, "state:") {
					t.Error("Frontmatter should contain state")
				}
				if !strings.Contains(output.Frontmatter, "url:") {
					t.Error("Frontmatter should contain url")
				}

				// 验证 Body
				if output.Body == "" {
					t.Error("Body should not be empty")
				}
				if !strings.Contains(output.Body, "# Bug: Cannot login") {
					t.Error("Body should contain formatted title")
				}
				if !strings.Contains(output.Body, "**Author:**") {
					t.Error("Body should contain Author")
				}
				if !strings.Contains(output.Body, "**State:**") {
					t.Error("Body should contain State")
				}
				if !strings.Contains(output.Body, "**Labels:**") {
					t.Error("Body should contain Labels")
				}
				if !strings.Contains(output.Body, "## Description") {
					t.Error("Body should contain Description section")
				}
				if !strings.Contains(output.Body, "## Comments") {
					t.Error("Body should contain Comments section when there are comments")
				}

				// 验证 FullContent
				if output.FullContent == "" {
					t.Error("FullContent should not be empty")
				}
			},
		},
		{
			name: "issue without comments",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Simple Issue",
				Author:    "user1",
				CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Body:      "Just a simple issue",
				Labels:    []string{},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			checkFunc: func(t *testing.T, output *Output) {
				// Body should not be empty
				if output.Body == "" {
					t.Error("Body should not be empty even without comments")
				}
				// 没有评论时，Comments 部分不应该出现
				if strings.Contains(output.Body, "## Comments") {
					t.Error("Body should not contain Comments section when there are no comments")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &converter{}
			ctx := context.Background()
			output, err := c.Convert(ctx, tt.resource, tt.userLink)

			if err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			tt.checkFunc(t, output)
		})
	}
}

// =============================================================================
// 4.5.2 测试 userLink 参数
// =============================================================================

func TestConverter_Convert_UserLink(t *testing.T) {
	resource := &fetcher.GitHubResource{
		Type:      "issue",
		Title:     "Test Issue",
		Author:    "testuser",
		CreatedAt: time.Date(2026, 3, 20, 10, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 3, 20, 10, 30, 0, 0, time.UTC),
		State:     "open",
		Body:      "Test body",
		Labels:    []string{},
		Comments:  []fetcher.Comment{},
		URL:       "https://github.com/owner/repo/issues/1",
	}

	t.Run("userLink false - author as plain text", func(t *testing.T) {
		c := &converter{}
		ctx := context.Background()
		output, err := c.Convert(ctx, resource, false)

		if err != nil {
			t.Fatalf("Convert() error = %v", err)
		}

		if !strings.Contains(output.Body, "**Author:** testuser") {
			t.Errorf("Body should contain '**Author:** testuser', got %q", output.Body)
		}
		if strings.Contains(output.Body, "[testuser](https://github.com/testuser)") {
			t.Error("Body should not contain author link when userLink=false")
		}
	})

	t.Run("userLink true - author as link", func(t *testing.T) {
		c := &converter{}
		ctx := context.Background()
		output, err := c.Convert(ctx, resource, true)

		if err != nil {
			t.Fatalf("Convert() error = %v", err)
		}

		expected := "**Author:** [testuser](https://github.com/testuser)"
		if !strings.Contains(output.Body, expected) {
			t.Errorf("Body should contain %q, got %q", expected, output.Body)
		}
	})
}
