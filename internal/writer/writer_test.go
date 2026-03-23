package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWriterInterfaceSignature 测试 Writer 接口签名
func TestWriterInterfaceSignature(t *testing.T) {
	var w Writer = New()
	// 验证 Write 方法签名: Write(content string, outputFile string) error
	var _ func(string, string) error = w.Write
}

// TestWriteToStdout 测试 outputFile="" 时输出到标准输出
func TestWriteToStdout(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "simple content",
			content: "hello world",
		},
		{
			name:    "multiline content",
			content: "line1\nline2\nline3",
		},
		{
			name:    "empty content",
			content: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()

			// 捕获标准输出
			oldStdout := os.Stdout
			r, w2, err := os.Pipe()
			if err != nil {
				t.Fatalf("failed to create pipe: %v", err)
			}
			os.Stdout = w2

			err = w.Write(tt.content, "")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			w2.Close()
			os.Stdout = oldStdout

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if output != tt.content {
				t.Errorf("expected stdout %q, got %q", tt.content, output)
			}
		})
	}
}

// TestWriteToFile 表格驱动测试：写入新文件
func TestWriteToFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		filename    string
		wantContent string
		wantErr     bool
	}{
		{
			name:        "simple text file",
			content:     "hello world",
			filename:    "simple.txt",
			wantContent: "hello world",
			wantErr:     false,
		},
		{
			name:        "markdown file with frontmatter",
			content:     "---\ntitle: test\n---\n# Hello",
			filename:    "test.md",
			wantContent: "---\ntitle: test\n---\n# Hello",
			wantErr:     false,
		},
		{
			name:        "file in subdirectory",
			content:     "content",
			filename:    "subdir/test.txt",
			wantContent: "content",
			wantErr:     false,
		},
		{
			name:        "nested subdirectories",
			content:     "deep content",
			filename:    "a/b/c/deep.txt",
			wantContent: "deep content",
			wantErr:     false,
		},
		{
			name:        "empty content",
			content:     "",
			filename:    "empty.txt",
			wantContent: "",
			wantErr:     false,
		},
		{
			name:        "chinese content",
			content:     "你好世界",
			filename:    "chinese.txt",
			wantContent: "你好世界",
			wantErr:     false,
		},
		{
			name:        "special characters",
			content:     "# Title\n**bold** and *italic*",
			filename:    "special.md",
			wantContent: "# Title\n**bold** and *italic*",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, tt.filename)

			err := w.Write(tt.content, outputFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				data, err := os.ReadFile(outputFile)
				if err != nil {
					t.Errorf("failed to read file: %v", err)
					return
				}
				if string(data) != tt.wantContent {
					t.Errorf("file content = %q, want %q", string(data), tt.wantContent)
				}
			}
		})
	}
}

// TestFileAlreadyExists 表格驱动测试：文件已存在时返回错误
func TestFileAlreadyExists(t *testing.T) {
	tests := []struct {
		name       string
		setupFile  bool
		filename   string
		wantErr    bool
		errContain string
	}{
		{
			name:       "file does not exist - should succeed",
			setupFile:  false,
			filename:   "newfile.md",
			wantErr:    false,
			errContain: "",
		},
		{
			name:       "file already exists - should fail",
			setupFile:  true,
			filename:   "existing.md",
			wantErr:    true,
			errContain: "file already exists",
		},
		{
			name:       "file already exists with path",
			setupFile:  true,
			filename:   "dir/existing.md",
			wantErr:    true,
			errContain: "file already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, tt.filename)

			if tt.setupFile {
				// 确保父目录存在
				dir := filepath.Dir(outputFile)
				if dir != "" && dir != "." {
					if err := os.MkdirAll(dir, 0755); err != nil {
						t.Fatalf("failed to create directory: %v", err)
					}
				}
				if err := os.WriteFile(outputFile, []byte("existing content"), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			}

			err := w.Write("new content", outputFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContain != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContain) {
					t.Errorf("expected error containing %q, got %v", tt.errContain, err)
				}
			}
		})
	}
}

// TestDirectoryNotExists 测试目录不存在时创建目录
func TestDirectoryNotExists(t *testing.T) {
	tests := []struct {
		name      string
		dirPath   string
		wantErr   bool
		wantFile  bool
	}{
		{
			name:     "single level directory",
			dirPath:  "subdir",
			wantErr:  false,
			wantFile: true,
		},
		{
			name:     "double level directory",
			dirPath:  "a/b",
			wantErr:  false,
			wantFile: true,
		},
		{
			name:     "deep nested directory",
			dirPath:  "x/y/z/w",
			wantErr:  false,
			wantFile: true,
		},
		{
			name:     "current directory",
			dirPath:  ".",
			wantErr:  false,
			wantFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, tt.dirPath, "test.txt")
			content := "test content"

			err := w.Write(content, outputFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantFile && !tt.wantErr {
				data, err := os.ReadFile(outputFile)
				if err != nil {
					t.Errorf("failed to read file: %v", err)
					return
				}
				if string(data) != content {
					t.Errorf("file content = %q, want %q", string(data), content)
				}
			}
		})
	}
}
