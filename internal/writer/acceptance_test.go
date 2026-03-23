package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAcceptance_FileOutput 根据 spec 5.3 验收标准：文件输出
func TestAcceptance_FileOutput(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, tmpDir string) (content string, outputFile string)
		wantErr    bool
		errContain string
		validate   func(t *testing.T, tmpDir string, content string, outputFile string)
	}{
		{
			name: "spec 5.3 - 无输出文件参数 - 输出到标准输出",
			setup: func(t *testing.T, tmpDir string) (string, string) {
				return "test content for stdout", ""
			},
			wantErr:    false,
			errContain: "",
			validate: func(t *testing.T, tmpDir string, content string, outputFile string) {
				// 验证写入 stdout 成功（无错误即可，内容通过捕获 stdout 验证）
			},
		},
		{
			name: "spec 5.3 - 指定输出文件 - 写入 issue-123.md",
			setup: func(t *testing.T, tmpDir string) (string, string) {
				return "# Issue 123\n\nContent here", filepath.Join(tmpDir, "issue-123.md")
			},
			wantErr:    false,
			errContain: "",
			validate: func(t *testing.T, tmpDir string, content string, outputFile string) {
				data, err := os.ReadFile(outputFile)
				if err != nil {
					t.Errorf("failed to read file: %v", err)
					return
				}
				if string(data) != content {
					t.Errorf("file content mismatch, got %q, want %q", string(data), content)
				}
			},
		},
		{
			name: "spec 5.3 - 文件已存在 - 报错：文件已存在",
			setup: func(t *testing.T, tmpDir string) (string, string) {
				outputFile := filepath.Join(tmpDir, "existing.md")
				if err := os.WriteFile(outputFile, []byte("old content"), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
				return "new content", outputFile
			},
			wantErr:    true,
			errContain: "file already exists",
			validate:   nil,
		},
		{
			name: "spec 5.3 - 目录中文件已存在 - 报错：文件已存在",
			setup: func(t *testing.T, tmpDir string) (string, string) {
				dir := filepath.Join(tmpDir, "dir")
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				outputFile := filepath.Join(dir, "existing.md")
				if err := os.WriteFile(outputFile, []byte("old content"), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
				return "new content", outputFile
			},
			wantErr:    true,
			errContain: "file already exists",
			validate:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			tmpDir := t.TempDir()

			content, outputFile := tt.setup(t, tmpDir)

			err := w.Write(content, outputFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContain != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContain) {
					t.Errorf("expected error containing %q, got %v", tt.errContain, err)
				}
			}

			if tt.validate != nil && !tt.wantErr {
				tt.validate(t, tmpDir, content, outputFile)
			}
		})
	}
}

// TestAcceptance_FullMarkdownOutput 验收测试：完整 Markdown 输出
func TestAcceptance_FullMarkdownOutput(t *testing.T) {
	w := New()
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "pull_request-456.md")

	frontmatter := `---
title: "feat: Add dark mode support"
author: "designer456"
created_at: "2026-03-19 08:00 UTC"
type: "pull_request"
state: "closed"
labels:
  - "enhancement"
  - "ui"
url: "https://github.com/owner/repo/pull/456"
---`

	body := `# feat: Add dark mode support

**Author:** designer456
**Created:** 2026-03-19 08:00 UTC
**State:** 🔴 Closed
**Labels:** enhancement, ui

---

## Description

Adds dark mode support using CSS variables.`

	fullContent := frontmatter + "\n" + body

	err := w.Write(fullContent, outputFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(data) != fullContent {
		t.Errorf("content mismatch, got:\n%s\nwant:\n%s", string(data), fullContent)
	}
}
