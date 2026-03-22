package parser

import (
	"testing"
)

func TestParsedURL_Fields(t *testing.T) {
	tests := []struct {
		name    string
		got     *ParsedURL
		wantType URLType
		wantOwner string
		wantRepo  string
		wantNumber int
		wantRawURL string
	}{
		{
			name: "ParsedURL should have all fields",
			got: &ParsedURL{
				Type:   URLTypeIssue,
				Owner:  "owner",
				Repo:   "repo",
				Number: 123,
				RawURL: "https://github.com/owner/repo/issues/123",
			},
			wantType:   URLTypeIssue,
			wantOwner:  "owner",
			wantRepo:    "repo",
			wantNumber:  123,
			wantRawURL: "https://github.com/owner/repo/issues/123",
		},
		{
			name: "ParsedURL for PullRequest",
			got: &ParsedURL{
				Type:   URLTypePullRequest,
				Owner:  "myorg",
				Repo:   "myrepo",
				Number: 456,
				RawURL: "https://github.com/myorg/myrepo/pull/456",
			},
			wantType:   URLTypePullRequest,
			wantOwner:  "myorg",
			wantRepo:   "myrepo",
			wantNumber: 456,
			wantRawURL: "https://github.com/myorg/myrepo/pull/456",
		},
		{
			name: "ParsedURL for Discussion",
			got: &ParsedURL{
				Type:   URLTypeDiscussion,
				Owner:  "testowner",
				Repo:   "testrepo",
				Number: 789,
				RawURL: "https://github.com/testowner/testrepo/discussions/789",
			},
			wantType:   URLTypeDiscussion,
			wantOwner:  "testowner",
			wantRepo:   "testrepo",
			wantNumber: 789,
			wantRawURL: "https://github.com/testowner/testrepo/discussions/789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", tt.got.Type, tt.wantType)
			}
			if tt.got.Owner != tt.wantOwner {
				t.Errorf("Owner = %v, want %v", tt.got.Owner, tt.wantOwner)
			}
			if tt.got.Repo != tt.wantRepo {
				t.Errorf("Repo = %v, want %v", tt.got.Repo, tt.wantRepo)
			}
			if tt.got.Number != tt.wantNumber {
				t.Errorf("Number = %v, want %v", tt.got.Number, tt.wantNumber)
			}
			if tt.got.RawURL != tt.wantRawURL {
				t.Errorf("RawURL = %v, want %v", tt.got.RawURL, tt.wantRawURL)
			}
		})
	}
}

func TestParsedURL_String(t *testing.T) {
	tests := []struct {
		name string
		p    *ParsedURL
	}{
		{
			name: "Issue URL representation",
			p: &ParsedURL{
				Type:   URLTypeIssue,
				Owner:  "owner",
				Repo:   "repo",
				Number: 123,
				RawURL: "https://github.com/owner/repo/issues/123",
			},
		},
		{
			name: "PullRequest URL representation",
			p: &ParsedURL{
				Type:   URLTypePullRequest,
				Owner:  "owner",
				Repo:   "repo",
				Number: 456,
				RawURL: "https://github.com/owner/repo/pull/456",
			},
		},
		{
			name: "Discussion URL representation",
			p: &ParsedURL{
				Type:   URLTypeDiscussion,
				Owner:  "owner",
				Repo:   "repo",
				Number: 789,
				RawURL: "https://github.com/owner/repo/discussions/789",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.p.Type < URLTypeIssue || tt.p.Type > URLTypeUnknown {
				t.Errorf("Type is out of valid range: %v", tt.p.Type)
			}
			if tt.p.Owner == "" {
				t.Error("Owner should not be empty")
			}
			if tt.p.Repo == "" {
				t.Error("Repo should not be empty")
			}
			if tt.p.Number <= 0 {
				t.Error("Number should be positive")
			}
			if tt.p.RawURL == "" {
				t.Error("RawURL should not be empty")
			}
		})
	}
}

// Ensure ParsedURL has no unexported fields that would prevent initialization
func TestParsedURL_ZeroValue(t *testing.T) {
	var p ParsedURL
	if p.Type != 0 {
		t.Errorf("Zero value Type = %v, want 0", p.Type)
	}
	if p.Owner != "" {
		t.Errorf("Zero value Owner = %v, want empty string", p.Owner)
	}
	if p.Repo != "" {
		t.Errorf("Zero value Repo = %v, want empty string", p.Repo)
	}
	if p.Number != 0 {
		t.Errorf("Zero value Number = %v, want 0", p.Number)
	}
	if p.RawURL != "" {
		t.Errorf("Zero value RawURL = %v, want empty string", p.RawURL)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ParsedURL
		wantErr bool
		errMsg  string
	}{
		// T2.2.1: Issue URL
		{
			name:  "valid issue URL",
			input: "https://github.com/owner/repo/issues/123",
			want: &ParsedURL{
				Type:   URLTypeIssue,
				Owner:  "owner",
				Repo:   "repo",
				Number: 123,
				RawURL: "https://github.com/owner/repo/issues/123",
			},
			wantErr: false,
		},
		// T2.2.2: PullRequest URL
		{
			name:  "valid pull request URL",
			input: "https://github.com/owner/repo/pull/456",
			want: &ParsedURL{
				Type:   URLTypePullRequest,
				Owner:  "owner",
				Repo:   "repo",
				Number: 456,
				RawURL: "https://github.com/owner/repo/pull/456",
			},
			wantErr: false,
		},
		// T2.2.3: Discussion URL
		{
			name:  "valid discussion URL",
			input: "https://github.com/owner/repo/discussions/789",
			want: &ParsedURL{
				Type:   URLTypeDiscussion,
				Owner:  "owner",
				Repo:   "repo",
				Number: 789,
				RawURL: "https://github.com/owner/repo/discussions/789",
			},
			wantErr: false,
		},
		// T2.2.4: 子链接拒绝
		{
			name:     "sublink should be rejected",
			input:    "github.com/owner/repo/issues/123#issuecomment-456",
			want:     nil,
			wantErr:  true,
			errMsg:   "unsupported URL type",
		},
		// T2.2.5: 无效 GitHub URL
		{
			name:     "non-github URL should be rejected",
			input:    "https://google.com/owner/repo/issues/123",
			want:     nil,
			wantErr:  true,
			errMsg:   "not a GitHub URL",
		},
		// T2.2.6: http vs https
		{
			name:  "http URL should be supported",
			input: "http://github.com/owner/repo/issues/123",
			want: &ParsedURL{
				Type:   URLTypeIssue,
				Owner:  "owner",
				Repo:   "repo",
				Number: 123,
				RawURL: "http://github.com/owner/repo/issues/123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse(%q) error = nil, want error containing %q", tt.input, tt.errMsg)
					return
				}
				// Check error message contains expected text
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Parse(%q) error = %q, want error containing %q", tt.input, err.Error(), tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
				return
			}
			if got == nil {
				t.Errorf("Parse(%q) got = nil, want %+v", tt.input, tt.want)
				return
			}
			if got.Type != tt.want.Type {
				t.Errorf("Parse(%q) Type = %v, want %v", tt.input, got.Type, tt.want.Type)
			}
			if got.Owner != tt.want.Owner {
				t.Errorf("Parse(%q) Owner = %v, want %v", tt.input, got.Owner, tt.want.Owner)
			}
			if got.Repo != tt.want.Repo {
				t.Errorf("Parse(%q) Repo = %v, want %v", tt.input, got.Repo, tt.want.Repo)
			}
			if got.Number != tt.want.Number {
				t.Errorf("Parse(%q) Number = %v, want %v", tt.input, got.Number, tt.want.Number)
			}
			if got.RawURL != tt.want.RawURL {
				t.Errorf("Parse(%q) RawURL = %v, want %v", tt.input, got.RawURL, tt.want.RawURL)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
