package converter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bigwhite/issue2md/internal/fetcher"
)

// FormatFrontmatter 生成 YAML frontmatter
func FormatFrontmatter(resource *fetcher.GitHubResource) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %q\n", resource.Title))
	sb.WriteString(fmt.Sprintf("author: %q\n", resource.Author))
	sb.WriteString(fmt.Sprintf("created_at: %q\n", resource.CreatedAt.Format("2006-01-02 15:04 MST")))
	sb.WriteString(fmt.Sprintf("type: %q\n", resource.Type))
	sb.WriteString(fmt.Sprintf("state: %q\n", resource.State))

	if len(resource.Labels) > 0 {
		sb.WriteString("labels:\n")
		for _, label := range resource.Labels {
			sb.WriteString(fmt.Sprintf("  - %q\n", label))
		}
	}

	sb.WriteString(fmt.Sprintf("url: %q\n", resource.URL))
	sb.WriteString("---\n")
	return sb.String()
}

// FormatTitle 格式化标题
func FormatTitle(title string) string {
	return "# " + title
}

// FormatAuthor 格式化作者
func FormatAuthor(author string, userLink bool) string {
	if userLink {
		return fmt.Sprintf("**Author:** [%s](https://github.com/%s)", author, author)
	}
	return fmt.Sprintf("**Author:** %s", author)
}

// FormatState 格式化状态
func FormatState(state string) string {
	switch state {
	case "open":
		return "🟢 Open"
	case "closed":
		return "🔴 Closed"
	default:
		return state
	}
}

// FormatLabels 格式化标签
func FormatLabels(labels []string) string {
	return "**Labels:** " + strings.Join(labels, ", ")
}

// ConvertMentions 转换 @mention
func ConvertMentions(body string) string {
	re := regexp.MustCompile(`@([a-zA-Z0-9_-]+)`)
	return re.ReplaceAllString(body, "[$1](https://github.com/$1)")
}

// ConvertReferences 转换 #reference
func ConvertReferences(body string, repoPath string) string {
	re := regexp.MustCompile(`#(\d+)`)
	return re.ReplaceAllString(body, fmt.Sprintf("[#%s](https://github.com/%s/issues/%s)", "$1", repoPath, "$1"))
}

// PreserveImages 保留图片
func PreserveImages(body string) string {
	return body
}

// PreserveCodeBlocks 保留代码块
func PreserveCodeBlocks(body string) string {
	return body
}

// PreserveTaskLists 保留任务列表
func PreserveTaskLists(body string) string {
	return body
}

// FormatBody 格式化正文内容
func FormatBody(resource *fetcher.GitHubResource, userLink bool) string {
	var sb strings.Builder
	sb.WriteString("## Description\n\n")
	sb.WriteString(resource.Body)
	return sb.String()
}

// FormatComments 格式化评论
func FormatComments(comments []fetcher.Comment) string {
	if len(comments) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, comment := range comments {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(fmt.Sprintf("### Comment by %s @ %s\n\n",
			comment.Author,
			comment.CreatedAt.Format("2006-01-02 15:04 MST")))
		sb.WriteString(comment.Body)
	}
	return sb.String()
}
