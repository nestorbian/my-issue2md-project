package writer

import (
	"fmt"
	"os"
	"path/filepath"
)

// Writer 写入 Markdown 输出
type Writer interface {
	// Write 将 content 写入指定目标
	// 如果 outputFile 为空，写入 stdout
	// 如果文件已存在，返回错误
	Write(content string, outputFile string) error
}

// fileWriter 实现 Writer 接口
type fileWriter struct{}

// New 创建一个新的 fileWriter
func New() Writer {
	return &fileWriter{}
}

// Write 将 content 写入指定目标
func (w *fileWriter) Write(content string, outputFile string) error {
	if outputFile == "" {
		_, err := os.Stdout.WriteString(content)
		return err
	}

	// 检查文件是否已存在
	_, err := os.Stat(outputFile)
	if err == nil {
		return fmt.Errorf("file already exists: %s", outputFile)
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check file: %w", err)
	}

	// 确保目录存在
	dir := filepath.Dir(outputFile)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	return os.WriteFile(outputFile, []byte(content), 0644)
}

