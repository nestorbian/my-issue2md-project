package converter

import "testing"

func TestOutput_Fields(t *testing.T) {
	tests := []struct {
		name   string
		output Output
		want   *Output
	}{
		{
			name:   "empty output",
			output: Output{},
			want: &Output{
				Frontmatter: "",
				Body:        "",
				FullContent: "",
			},
		},
		{
			name: "output with frontmatter only",
			output: Output{
				Frontmatter: "---\ntitle: Test\n---",
				Body:        "",
				FullContent: "---\ntitle: Test\n---",
			},
			want: &Output{
				Frontmatter: "---\ntitle: Test\n---",
				Body:        "",
				FullContent: "---\ntitle: Test\n---",
			},
		},
		{
			name: "output with body only",
			output: Output{
				Frontmatter: "",
				Body:        "# Title\n\nContent",
				FullContent: "# Title\n\nContent",
			},
			want: &Output{
				Frontmatter: "",
				Body:        "# Title\n\nContent",
				FullContent: "# Title\n\nContent",
			},
		},
		{
			name: "full output",
			output: Output{
				Frontmatter: "---\ntitle: Test Issue\n---",
				Body:        "# Test Issue\n\nContent here",
				FullContent: "---\ntitle: Test Issue\n---\n# Test Issue\n\nContent here",
			},
			want: &Output{
				Frontmatter: "---\ntitle: Test Issue\n---",
				Body:        "# Test Issue\n\nContent here",
				FullContent: "---\ntitle: Test Issue\n---\n# Test Issue\n\nContent here",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证 Frontmatter 字段
			if tt.output.Frontmatter != tt.want.Frontmatter {
				t.Errorf("Frontmatter = %q, want %q", tt.output.Frontmatter, tt.want.Frontmatter)
			}
			// 验证 Body 字段
			if tt.output.Body != tt.want.Body {
				t.Errorf("Body = %q, want %q", tt.output.Body, tt.want.Body)
			}
			// 验证 FullContent 字段
			if tt.output.FullContent != tt.want.FullContent {
				t.Errorf("FullContent = %q, want %q", tt.output.FullContent, tt.want.FullContent)
			}
		})
	}
}
