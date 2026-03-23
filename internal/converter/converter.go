package converter

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/bigwhite/issue2md/internal/fetcher"
)

// Converter 转换 GitHub 资源为 Markdown
type Converter interface {
	// Convert 将 GitHubResource 转换为 Output
	Convert(ctx context.Context, resource *fetcher.GitHubResource, userLink bool) (*Output, error)
}

// converter 是 Converter 接口的简单实现
type converter struct{}

// NewConverter 创建一个新的 Converter
func NewConverter() Converter {
	return &converter{}
}

// extractRepoPath 从 URL 中提取 owner/repo
func extractRepoPath(url string) string {
	re := regexp.MustCompile(`github\.com/([^/]+/[^/]+)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// Convert 将 GitHubResource 转换为 Output
func (c *converter) Convert(ctx context.Context, resource *fetcher.GitHubResource, userLink bool) (*Output, error) {
	if resource == nil {
		return nil, fmt.Errorf("resource is nil")
	}

	// 提取 repoPath 用于 reference 转换
	repoPath := extractRepoPath(resource.URL)

	// 生成 Frontmatter
	frontmatter := FormatFrontmatter(resource)

	// 对 Body 内容进行转换
	bodyContent := resource.Body
	bodyContent = ConvertMentions(bodyContent)
	bodyContent = ConvertReferences(bodyContent, repoPath)

	// 构建 Body
	var bodyBuilder strings.Builder

	// 标题
	bodyBuilder.WriteString(FormatTitle(resource.Title))
	bodyBuilder.WriteString("\n\n")

	// Author
	bodyBuilder.WriteString(FormatAuthor(resource.Author, userLink))
	bodyBuilder.WriteString("\n\n")

	// Created
	bodyBuilder.WriteString(fmt.Sprintf("**Created:** %s\n\n", resource.CreatedAt.Format("2006-01-02 15:04 MST")))

	// State
	bodyBuilder.WriteString(fmt.Sprintf("**State:** %s\n\n", FormatState(resource.State)))

	// Labels
	bodyBuilder.WriteString(FormatLabels(resource.Labels))
	bodyBuilder.WriteString("\n\n")

	// 分隔符
	bodyBuilder.WriteString("---\n\n")

	// Description
	bodyBuilder.WriteString("## Description\n\n")
	bodyBuilder.WriteString(bodyContent)
	bodyBuilder.WriteString("\n\n")

	// Comments (如果有)
	if len(resource.Comments) > 0 {
		bodyBuilder.WriteString("## Comments\n\n")
		bodyBuilder.WriteString("---\n\n")
		bodyBuilder.WriteString(FormatComments(resource.Comments))
	}

	body := bodyBuilder.String()

	return &Output{
		Frontmatter: frontmatter,
		Body:        body,
		FullContent: frontmatter + "\n" + body,
	}, nil
}

// 确保 converter 实现了 Converter 接口
var _ Converter = (*converter)(nil)
