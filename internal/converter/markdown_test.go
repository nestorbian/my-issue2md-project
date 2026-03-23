package converter

import (
	"strings"
	"testing"
	"time"

	"github.com/bigwhite/issue2md/internal/fetcher"
)

// =============================================================================
// 4.3.1 测试 Frontmatter 格式化
// =============================================================================

func TestFormatFrontmatter(t *testing.T) {
	tests := []struct {
		name      string
		resource  *fetcher.GitHubResource
		wantPrefix string
	}{
		{
			name: "issue frontmatter",
			resource: &fetcher.GitHubResource{
				Title:     "Bug: Cannot login",
				Author:    "developer123",
				CreatedAt: time.Date(2026, 3, 20, 10, 30, 0, 0, time.UTC),
				Type:      "issue",
				State:     "open",
				Labels:    []string{"bug", "oauth"},
				URL:       "https://github.com/owner/repo/issues/123",
			},
			wantPrefix: "---\ntitle:",
		},
		{
			name: "pull request frontmatter",
			resource: &fetcher.GitHubResource{
				Title:     "feat: Add dark mode",
				Author:    "designer456",
				CreatedAt: time.Date(2026, 3, 19, 8, 0, 0, 0, time.UTC),
				Type:      "pull_request",
				State:     "closed",
				Labels:    []string{"enhancement"},
				URL:       "https://github.com/owner/repo/pull/456",
			},
			wantPrefix: "---\ntitle:",
		},
		{
			name: "discussion frontmatter",
			resource: &fetcher.GitHubResource{
				Title:     "Question about API",
				Author:    "user789",
				CreatedAt: time.Date(2026, 3, 18, 14, 0, 0, 0, time.UTC),
				Type:      "discussion",
				State:     "open",
				Labels:    []string{},
				URL:       "https://github.com/owner/repo/discussions/789",
			},
			wantPrefix: "---\ntitle:",
		},
		{
			name: "frontmatter contains all required fields",
			resource: &fetcher.GitHubResource{
				Title:     "Test Issue",
				Author:    "testuser",
				CreatedAt: time.Date(2026, 3, 20, 10, 30, 0, 0, time.UTC),
				Type:      "issue",
				State:     "open",
				Labels:    []string{"bug", "priority-high"},
				URL:       "https://github.com/owner/repo/issues/123",
			},
			wantPrefix: "---\ntitle:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatFrontmatter(tt.resource)
			if !strings.HasPrefix(got, tt.wantPrefix) {
				t.Errorf("FormatFrontmatter() prefix = %q, want prefix %q", got, tt.wantPrefix)
			}
			if !strings.Contains(got, "title:") {
				t.Errorf("FormatFrontmatter() should contain title")
			}
			if !strings.Contains(got, "author:") {
				t.Errorf("FormatFrontmatter() should contain author")
			}
			if !strings.Contains(got, "created_at:") {
				t.Errorf("FormatFrontmatter() should contain created_at")
			}
			if !strings.Contains(got, "type:") {
				t.Errorf("FormatFrontmatter() should contain type")
			}
			if !strings.Contains(got, "state:") {
				t.Errorf("FormatFrontmatter() should contain state")
			}
			if !strings.Contains(got, "url:") {
				t.Errorf("FormatFrontmatter() should contain url")
			}
		})
	}
}

// =============================================================================
// 4.3.2 测试标题渲染
// =============================================================================

func TestFormatTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{
			name:  "simple title",
			title: "Bug Report",
			want:  "# Bug Report",
		},
		{
			name:  "title with colon",
			title: "Bug: Cannot login",
			want:  "# Bug: Cannot login",
		},
		{
			name:  "title with special chars",
			title: "Fix #123 - memory leak",
			want:  "# Fix #123 - memory leak",
		},
		{
			name:  "empty title",
			title: "",
			want:  "# ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatTitle(tt.title)
			if got != tt.want {
				t.Errorf("FormatTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// 4.3.3 测试 Author 渲染 (userLink=false)
// =============================================================================

func TestFormatAuthor_NoLink(t *testing.T) {
	tests := []struct {
		name     string
		author   string
		userLink bool
		want     string
	}{
		{
			name:     "author without link",
			author:   "developer123",
			userLink: false,
			want:     "**Author:** developer123",
		},
		{
			name:     "author with underscore",
			author:   "test_user",
			userLink: false,
			want:     "**Author:** test_user",
		},
		{
			name:     "author with hyphen",
			author:   "john-doe",
			userLink: false,
			want:     "**Author:** john-doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatAuthor(tt.author, tt.userLink)
			if got != tt.want {
				t.Errorf("FormatAuthor() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// 4.3.4 测试 Author 渲染 (userLink=true)
// =============================================================================

func TestFormatAuthor_WithLink(t *testing.T) {
	tests := []struct {
		name     string
		author   string
		userLink bool
		want     string
	}{
		{
			name:     "author with link",
			author:   "developer123",
			userLink: true,
			want:     "**Author:** [developer123](https://github.com/developer123)",
		},
		{
			name:     "author with link and underscore",
			author:   "test_user",
			userLink: true,
			want:     "**Author:** [test_user](https://github.com/test_user)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatAuthor(tt.author, tt.userLink)
			if got != tt.want {
				t.Errorf("FormatAuthor() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// 4.3.5 测试 State 渲染
// =============================================================================

func TestFormatState(t *testing.T) {
	tests := []struct {
		name  string
		state string
		want  string
	}{
		{
			name:  "open state",
			state: "open",
			want:  "🟢 Open",
		},
		{
			name:  "closed state",
			state: "closed",
			want:  "🔴 Closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatState(tt.state)
			if got != tt.want {
				t.Errorf("FormatState() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// 4.3.6 测试 Labels 渲染
// =============================================================================

func TestFormatLabels(t *testing.T) {
	tests := []struct {
		name   string
		labels []string
		want   string
	}{
		{
			name:   "single label",
			labels: []string{"bug"},
			want:   "**Labels:** bug",
		},
		{
			name:   "multiple labels",
			labels: []string{"bug", "priority-high"},
			want:   "**Labels:** bug, priority-high",
		},
		{
			name:   "no labels",
			labels: []string{},
			want:   "**Labels:** ",
		},
		{
			name:   "three labels",
			labels: []string{"enhancement", "ui", "dark-mode"},
			want:   "**Labels:** enhancement, ui, dark-mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatLabels(tt.labels)
			if got != tt.want {
				t.Errorf("FormatLabels() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// 4.3.7 测试 @mention 转换
// =============================================================================

func TestConvertMentions(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single mention",
			input: "Hello @username, please review",
			want:  "Hello [username](https://github.com/username), please review",
		},
		{
			name:  "multiple mentions",
			input: "@user1 and @user2 should review",
			want:  "[user1](https://github.com/user1) and [user2](https://github.com/user2) should review",
		},
		{
			name:  "mention at start",
			input: "@admin please fix this",
			want:  "[admin](https://github.com/admin) please fix this",
		},
		{
			name:  "mention with underscore",
			input: "Thanks @test_user!",
			want:  "Thanks [test_user](https://github.com/test_user)!",
		},
		{
			name:  "no mentions",
			input: "No mentions here",
			want:  "No mentions here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertMentions(tt.input)
			if got != tt.want {
				t.Errorf("ConvertMentions() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// 4.3.8 测试 #reference 转换
// =============================================================================

func TestConvertReferences(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single reference",
			input: "See #123 for details",
			want:  "See [#123](https://github.com/owner/repo/issues/123) for details",
		},
		{
			name:  "multiple references",
			input: "Related to #100 and #200",
			want:  "Related to [#100](https://github.com/owner/repo/issues/100) and [#200](https://github.com/owner/repo/issues/200)",
		},
		{
			name:  "pr reference",
			input: "Supersedes #456",
			want:  "Supersedes [#456](https://github.com/owner/repo/issues/456)",
		},
		{
			name:  "no references",
			input: "No references here",
			want:  "No references here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertReferences(tt.input, "owner/repo")
			if got != tt.want {
				t.Errorf("ConvertReferences() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// 4.3.9 测试图片保留
// =============================================================================

func TestPreserveImages(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "image from githubusercontent",
			input: "![image](https://user-images.githubusercontent.com/123/456.png)",
			want:  "![image](https://user-images.githubusercontent.com/123/456.png)",
		},
		{
			name:  "multiple images",
			input: "![img1](https://user-images.githubusercontent.com/1/1.png) and ![img2](https://user-images.githubusercontent.com/2/2.png)",
			want:  "![img1](https://user-images.githubusercontent.com/1/1.png) and ![img2](https://user-images.githubusercontent.com/2/2.png)",
		},
		{
			name:  "image with special chars in alt",
			input: "![screenshot](https://user-images.githubusercontent.com/123/789.png)",
			want:  "![screenshot](https://user-images.githubusercontent.com/123/789.png)",
		},
		{
			name:  "no images",
			input: "Just text",
			want:  "Just text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PreserveImages(tt.input)
			if got != tt.want {
				t.Errorf("PreserveImages() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// 4.3.10 测试代码块保留
// =============================================================================

func TestPreserveCodeBlocks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "go code block",
			input: "```go\nfunc main() {}\n```",
			want:  "```go\nfunc main() {}\n```",
		},
		{
			name:  "javascript code block",
			input: "```javascript\nconsole.log('hello');\n```",
			want:  "```javascript\nconsole.log('hello');\n```",
		},
		{
			name:  "python code block",
			input: "```python\nprint('hello')\n```",
			want:  "```python\nprint('hello')\n```",
		},
		{
			name:  "code block without language",
			input: "```\nsome code\n```",
			want:  "```\nsome code\n```",
		},
		{
			name:  "multiple code blocks",
			input: "```go\nfmt.Println()\n```\ntext\n```python\nprint()\n```",
			want:  "```go\nfmt.Println()\n```\ntext\n```python\nprint()\n```",
		},
		{
			name:  "no code blocks",
			input: "Just plain text",
			want:  "Just plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PreserveCodeBlocks(tt.input)
			if got != tt.want {
				t.Errorf("PreserveCodeBlocks() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// 4.3.11 测试任务列表保留
// =============================================================================

func TestPreserveTaskLists(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "unchecked task",
			input: "- [ ] TODO item",
			want:  "- [ ] TODO item",
		},
		{
			name:  "checked task",
			input: "- [x] Done item",
			want:  "- [x] Done item",
		},
		{
			name:  "mixed tasks",
			input: "- [ ] Item 1\n- [x] Item 2\n- [ ] Item 3",
			want:  "- [ ] Item 1\n- [x] Item 2\n- [ ] Item 3",
		},
		{
			name:  "tasks with content",
			input: "- [ ] Write code\n- [x] Write tests",
			want:  "- [ ] Write code\n- [x] Write tests",
		},
		{
			name:  "no tasks",
			input: "Just plain text",
			want:  "Just plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PreserveTaskLists(tt.input)
			if got != tt.want {
				t.Errorf("PreserveTaskLists() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// 4.4.1 测试 Comments 渲染
// =============================================================================

func TestFormatComments(t *testing.T) {
	tests := []struct {
		name     string
		comments []fetcher.Comment
		want     string
	}{
		{
			name:     "single comment",
			comments: []fetcher.Comment{
				{
					Author:    "username",
					CreatedAt: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
					Body:      "评论内容。",
				},
			},
			want: "### Comment by username @ 2026-03-20 12:00 UTC\n\n评论内容。",
		},
		{
			name:     "multiple comments",
			comments: []fetcher.Comment{
				{
					Author:    "username",
					CreatedAt: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
					Body:      "第一条评论。",
				},
				{
					Author:    "anotheruser",
					CreatedAt: time.Date(2026, 3, 20, 14, 30, 0, 0, time.UTC),
					Body:      "另一条评论。",
				},
			},
			want: "### Comment by username @ 2026-03-20 12:00 UTC\n\n第一条评论。\n\n### Comment by anotheruser @ 2026-03-20 14:30 UTC\n\n另一条评论。",
		},
		{
			name: "comment with underscore in author",
			comments: []fetcher.Comment{
				{
					Author:    "test_user",
					CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
					Body:      "用户名的下划线。",
				},
			},
			want: "### Comment by test_user @ 2026-03-20 10:00 UTC\n\n用户名的下划线。",
		},
		{
			name: "comment with multiline body",
			comments: []fetcher.Comment{
				{
					Author:    "user1",
					CreatedAt: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
					Body:      "第一行\n第二行\n第三行",
				},
			},
			want: "### Comment by user1 @ 2026-03-20 12:00 UTC\n\n第一行\n第二行\n第三行",
		},
		{
			name:     "comment with special chars in body",
			comments: []fetcher.Comment{
				{
					Author:    "dev",
					CreatedAt: time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
					Body:      "测试 @mention 和 #123 引用",
				},
			},
			want: "### Comment by dev @ 2026-03-20 12:00 UTC\n\n测试 @mention 和 #123 引用",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatComments(tt.comments)
			if got != tt.want {
				t.Errorf("FormatComments() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// 4.4.2 测试空 Comments
// =============================================================================

func TestFormatComments_Empty(t *testing.T) {
	tests := []struct {
		name     string
		comments []fetcher.Comment
		want     string
	}{
		{
			name:     "empty comments list",
			comments: []fetcher.Comment{},
			want:     "",
		},
		{
			name:     "nil comments",
			comments: nil,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatComments(tt.comments)
			if got != tt.want {
				t.Errorf("FormatComments() = %q, want %q", got, tt.want)
			}
		})
	}
}
