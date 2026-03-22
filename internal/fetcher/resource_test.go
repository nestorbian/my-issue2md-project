package fetcher

import (
	"testing"
	"time"
)

func TestGitHubResource_Fields(t *testing.T) {
	createdAt := time.Date(2026, 3, 20, 10, 30, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC)
	comments := []Comment{
		{
			Author:    "user1",
			CreatedAt: time.Date(2026, 3, 20, 14, 0, 0, 0, time.UTC),
			Body:      "First comment",
		},
	}

	tests := []struct {
		name      string
		got       *GitHubResource
		wantType  string
		wantTitle string
		wantAuthor string
		wantState string
		wantBody  string
		wantURL   string
	}{
		{
			name: "GitHubResource should have all fields - Issue",
			got: &GitHubResource{
				Type:      "issue",
				Title:     "Bug Report",
				Author:    "developer123",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				State:     "open",
				Body:      "Issue body content",
				Labels:    []string{"bug", "priority-high"},
				Comments:  comments,
				URL:       "https://github.com/owner/repo/issues/123",
			},
			wantType:   "issue",
			wantTitle:  "Bug Report",
			wantAuthor: "developer123",
			wantState:  "open",
			wantBody:   "Issue body content",
			wantURL:    "https://github.com/owner/repo/issues/123",
		},
		{
			name: "GitHubResource for PullRequest",
			got: &GitHubResource{
				Type:      "pull_request",
				Title:     "feat: Add new feature",
				Author:    "contributor",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				State:     "closed",
				Body:      "PR description",
				Labels:    []string{"enhancement"},
				Comments:  nil,
				URL:       "https://github.com/owner/repo/pull/456",
			},
			wantType:   "pull_request",
			wantTitle:  "feat: Add new feature",
			wantAuthor: "contributor",
			wantState:  "closed",
			wantBody:   "PR description",
			wantURL:    "https://github.com/owner/repo/pull/456",
		},
		{
			name: "GitHubResource for Discussion",
			got: &GitHubResource{
				Type:      "discussion",
				Title:     "Question about architecture",
				Author:    "questioner",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				State:     "open",
				Body:      "Discussion content",
				Labels:    []string{"question"},
				Comments:  []Comment{},
				URL:       "https://github.com/owner/repo/discussions/789",
			},
			wantType:   "discussion",
			wantTitle:  "Question about architecture",
			wantAuthor: "questioner",
			wantState:  "open",
			wantBody:   "Discussion content",
			wantURL:    "https://github.com/owner/repo/discussions/789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", tt.got.Type, tt.wantType)
			}
			if tt.got.Title != tt.wantTitle {
				t.Errorf("Title = %v, want %v", tt.got.Title, tt.wantTitle)
			}
			if tt.got.Author != tt.wantAuthor {
				t.Errorf("Author = %v, want %v", tt.got.Author, tt.wantAuthor)
			}
			if tt.got.State != tt.wantState {
				t.Errorf("State = %v, want %v", tt.got.State, tt.wantState)
			}
			if tt.got.Body != tt.wantBody {
				t.Errorf("Body = %v, want %v", tt.got.Body, tt.wantBody)
			}
			if tt.got.URL != tt.wantURL {
				t.Errorf("URL = %v, want %v", tt.got.URL, tt.wantURL)
			}
		})
	}
}

func TestGitHubResource_Labels(t *testing.T) {
	tests := []struct {
		name   string
		got    *GitHubResource
		wantLen int
	}{
		{
			name: "Labels with multiple items",
			got: &GitHubResource{
				Labels: []string{"bug", "enhancement", "documentation"},
			},
			wantLen: 3,
		},
		{
			name: "Empty labels",
			got: &GitHubResource{
				Labels: []string{},
			},
			wantLen: 0,
		},
		{
			name: "Nil labels",
			got: &GitHubResource{
				Labels: nil,
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.got.Labels) != tt.wantLen {
				t.Errorf("len(Labels) = %v, want %v", len(tt.got.Labels), tt.wantLen)
			}
		})
	}
}

func TestGitHubResource_Comments(t *testing.T) {
	comment1 := Comment{
		Author:    "user1",
		CreatedAt: time.Date(2026, 3, 20, 14, 0, 0, 0, time.UTC),
		Body:      "First comment",
	}
	comment2 := Comment{
		Author:    "user2",
		CreatedAt: time.Date(2026, 3, 20, 15, 0, 0, 0, time.UTC),
		Body:      "Second comment",
	}

	tests := []struct {
		name      string
		got       *GitHubResource
		wantLen   int
	}{
		{
			name: "Multiple comments",
			got: &GitHubResource{
				Comments: []Comment{comment1, comment2},
			},
			wantLen: 2,
		},
		{
			name: "Empty comments",
			got: &GitHubResource{
				Comments: []Comment{},
			},
			wantLen: 0,
		},
		{
			name: "Nil comments",
			got: &GitHubResource{
				Comments: nil,
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.got.Comments) != tt.wantLen {
				t.Errorf("len(Comments) = %v, want %v", len(tt.got.Comments), tt.wantLen)
			}
		})
	}
}

func TestComment_Fields(t *testing.T) {
	createdAt := time.Date(2026, 3, 20, 14, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		got         Comment
		wantAuthor  string
		wantBody    string
	}{
		{
			name:       "Comment with all fields",
			got:        Comment{
				Author:    "commenter",
				CreatedAt: createdAt,
				Body:      "This is a comment",
			},
			wantAuthor: "commenter",
			wantBody:   "This is a comment",
		},
		{
			name: "Comment with empty body",
			got: Comment{
				Author:    "user",
				CreatedAt: createdAt,
				Body:      "",
			},
			wantAuthor: "user",
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got.Author != tt.wantAuthor {
				t.Errorf("Author = %v, want %v", tt.got.Author, tt.wantAuthor)
			}
			if tt.got.Body != tt.wantBody {
				t.Errorf("Body = %v, want %v", tt.got.Body, tt.wantBody)
			}
		})
	}
}
