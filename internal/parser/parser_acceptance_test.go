package parser

import (
	"strings"
	"testing"
)

// TestParseAcceptance 根据 spec 5.1 验收标准测试 URL 解析
func TestParseAcceptance(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantType    URLType
		wantOwner   string
		wantRepo    string
		wantNumber  int
		wantErr     bool
		errContains string
	}{
		{
			name:       "有效的 Issue URL",
			input:      "https://github.com/owner/repo/issues/123",
			wantType:   URLTypeIssue,
			wantOwner:  "owner",
			wantRepo:   "repo",
			wantNumber: 123,
			wantErr:    false,
		},
		{
			name:       "有效的 PR URL",
			input:      "https://github.com/owner/repo/pull/456",
			wantType:   URLTypePullRequest,
			wantOwner:  "owner",
			wantRepo:   "repo",
			wantNumber: 456,
			wantErr:    false,
		},
		{
			name:       "有效的 Discussion URL",
			input:      "https://github.com/owner/repo/discussions/789",
			wantType:   URLTypeDiscussion,
			wantOwner:  "owner",
			wantRepo:   "repo",
			wantNumber: 789,
			wantErr:    false,
		},
		{
			name:        "Issuecomment 子链接 - 应报错",
			input:       "github.com/owner/repo/issues/123#issuecomment-456",
			wantErr:     true,
			errContains: "unsupported URL type",
		},
		{
			name:        "无效 URL - 非 GitHub",
			input:       "https://google.com/owner/repo/issues/123",
			wantErr:     true,
			errContains: "not a GitHub URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse(%q) error = nil, want error", tt.input)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Parse(%q) error = %q, want error containing %q", tt.input, err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
				return
			}
			if got == nil {
				t.Errorf("Parse(%q) got = nil", tt.input)
				return
			}
			if got.Type != tt.wantType {
				t.Errorf("Parse(%q) Type = %v, want %v", tt.input, got.Type, tt.wantType)
			}
			if got.Owner != tt.wantOwner {
				t.Errorf("Parse(%q) Owner = %v, want %v", tt.input, got.Owner, tt.wantOwner)
			}
			if got.Repo != tt.wantRepo {
				t.Errorf("Parse(%q) Repo = %v, want %v", tt.input, got.Repo, tt.wantRepo)
			}
			if got.Number != tt.wantNumber {
				t.Errorf("Parse(%q) Number = %v, want %v", tt.input, got.Number, tt.wantNumber)
			}
		})
	}
}
