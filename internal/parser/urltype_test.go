package parser

import "testing"

func TestURLType_Values(t *testing.T) {
	tests := []struct {
		name  string
		got   URLType
		want  URLType
	}{
		{
			name:  "URLTypeIssue should be 0",
			got:   URLTypeIssue,
			want:  URLType(0),
		},
		{
			name:  "URLTypePullRequest should be 1",
			got:   URLTypePullRequest,
			want:  URLType(1),
		},
		{
			name:  "URLTypeDiscussion should be 2",
			got:   URLTypeDiscussion,
			want:  URLType(2),
		},
		{
			name:  "URLTypeUnknown should be 3",
			got:   URLTypeUnknown,
			want:  URLType(3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}
