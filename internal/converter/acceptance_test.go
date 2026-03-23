package converter

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/bigwhite/issue2md/internal/fetcher"
)

// =============================================================================
// 4.6.1 验收测试：Markdown 生成
// =============================================================================

// TestAcceptance_MarkdownGeneration 根据 spec 5.2 验收标准测试 Markdown 生成
func TestAcceptance_MarkdownGeneration(t *testing.T) {
	tests := []struct {
		name      string
		resource  *fetcher.GitHubResource
		userLink  bool
		repoPath  string
		validator func(t *testing.T, output *Output)
	}{
		{
			name: "spec 5.2: 带标签的 Issue - 标签出现在 frontmatter 和正文中",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Issue with Labels",
				Author:    "user1",
				CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Body:      "Test body",
				Labels:    []string{"bug", "enhancement"},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			repoPath: "owner/repo",
			validator: func(t *testing.T, output *Output) {
				// 验证 frontmatter 中有 labels
				if !strings.Contains(output.Frontmatter, `labels:`) {
					t.Error("Frontmatter should contain labels section")
				}
				if !strings.Contains(output.Frontmatter, `"bug"`) {
					t.Error("Frontmatter should contain bug label")
				}
				if !strings.Contains(output.Frontmatter, `"enhancement"`) {
					t.Error("Frontmatter should contain enhancement label")
				}
				// 验证 body 中有 Labels
				if !strings.Contains(output.Body, "**Labels:**") {
					t.Error("Body should contain Labels")
				}
				if !strings.Contains(output.Body, "bug") || !strings.Contains(output.Body, "enhancement") {
					t.Error("Body Labels should contain bug and enhancement")
				}
			},
		},
		{
			name: "spec 5.2: 带 @mention 的 Issue - @mention 转换为链接",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Issue with Mention",
				Author:    "user1",
				CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Body:      "Hello @developer, please review this.",
				Labels:    []string{},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			repoPath: "owner/repo",
			validator: func(t *testing.T, output *Output) {
				// @mention 应该被转换为链接
				if !strings.Contains(output.Body, "[developer](https://github.com/developer)") {
					t.Error("Body should contain @developer as a link")
				}
				// 不应该还有原始的 @developer
				if strings.Contains(output.Body, "@developer") {
					t.Error("Body should not contain raw @developer")
				}
			},
		},
		{
			name: "spec 5.2: 带 #reference 的 Issue - #reference 转换为链接",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Issue with Reference",
				Author:    "user1",
				CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Body:      "See #123 for related issue.",
				Labels:    []string{},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			repoPath: "owner/repo",
			validator: func(t *testing.T, output *Output) {
				// #123 应该被转换为链接
				if !strings.Contains(output.Body, "[#123](https://github.com/owner/repo/issues/123)") {
					t.Error("Body should contain #123 as a link")
				}
			},
		},
		{
			name: "spec 5.2: 带图片的 Issue - 图片 URL 保留",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Issue with Image",
				Author:    "user1",
				CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Body:      "![screenshot](https://user-images.githubusercontent.com/123/456.png)",
				Labels:    []string{},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			repoPath: "owner/repo",
			validator: func(t *testing.T, output *Output) {
				// 图片 URL 应该保留
				if !strings.Contains(output.Body, "https://user-images.githubusercontent.com/123/456.png") {
					t.Error("Body should preserve image URL")
				}
			},
		},
		{
			name: "spec 5.2: 带代码块的 Issue - 语言标识保留",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Issue with Code",
				Author:    "user1",
				CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Body:      "```go\nfunc main() {}\n```",
				Labels:    []string{},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			repoPath: "owner/repo",
			validator: func(t *testing.T, output *Output) {
				// 代码块和语言标识应该保留
				if !strings.Contains(output.Body, "```go") {
					t.Error("Body should preserve ```go code block")
				}
				if !strings.Contains(output.Body, "func main()") {
					t.Error("Body should preserve code content")
				}
			},
		},
		{
			name: "spec 5.2: 带任务列表的 Issue - 任务语法保留",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Issue with Tasks",
				Author:    "user1",
				CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Body:      "- [ ] TODO item\n- [x] Done item",
				Labels:    []string{},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			repoPath: "owner/repo",
			validator: func(t *testing.T, output *Output) {
				// 任务列表语法应该保留
				if !strings.Contains(output.Body, "- [ ]") {
					t.Error("Body should preserve unchecked task")
				}
				if !strings.Contains(output.Body, "- [x]") {
					t.Error("Body should preserve checked task")
				}
			},
		},
		{
			name: "spec 5.2: --user-link=false - 用户名为纯文本",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Issue",
				Author:    "plainuser",
				CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Body:      "Test",
				Labels:    []string{},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			repoPath: "owner/repo",
			validator: func(t *testing.T, output *Output) {
				// Author 应该是纯文本
				if !strings.Contains(output.Body, "**Author:** plainuser") {
					t.Error("Body should contain plain author name")
				}
				if strings.Contains(output.Body, "[plainuser](https://github.com/plainuser)") {
					t.Error("Body should not contain author link when userLink=false")
				}
			},
		},
		{
			name: "spec 5.2: --user-link=true - 用户名为 GitHub 个人主页链接",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Issue",
				Author:    "linkeduser",
				CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Body:      "Test",
				Labels:    []string{},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: true,
			repoPath: "owner/repo",
			validator: func(t *testing.T, output *Output) {
				// Author 应该是链接
				expected := "**Author:** [linkeduser](https://github.com/linkeduser)"
				if !strings.Contains(output.Body, expected) {
					t.Errorf("Body should contain author link, got %q", output.Body)
				}
			},
		},
		{
			name: "spec 5.2: 已关闭的 Issue - 状态显示 '🔴 Closed'",
			resource: &fetcher.GitHubResource{
				Type:      "issue",
				Title:     "Closed Issue",
				Author:    "user1",
				CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
				State:     "closed",
				Body:      "This issue is closed",
				Labels:    []string{},
				Comments:  []fetcher.Comment{},
				URL:       "https://github.com/owner/repo/issues/1",
			},
			userLink: false,
			repoPath: "owner/repo",
			validator: func(t *testing.T, output *Output) {
				// State 应该是 🔴 Closed
				if !strings.Contains(output.Body, "🔴 Closed") {
					t.Error("Body should contain 🔴 Closed for closed state")
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

			tt.validator(t, output)
		})
	}
}
