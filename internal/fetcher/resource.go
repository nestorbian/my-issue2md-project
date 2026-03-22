package fetcher

import "time"

// GitHubResource 是 Issue/PR/Discussion 的统一表示
type GitHubResource struct {
	Type      string    // "issue" | "pull_request" | "discussion"
	Title     string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
	State     string // "open" | "closed"
	Body      string // Markdown 内容
	Labels    []string
	Comments  []Comment
	URL       string // 原始 GitHub URL
}

// Comment 评论结构
type Comment struct {
	Author    string
	CreatedAt time.Time
	Body      string
}
