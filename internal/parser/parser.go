package parser

import (
	"fmt"
	"strings"
)

// ParsedURL 解析后的 URL 信息
type ParsedURL struct {
	Type   URLType
	Owner  string
	Repo   string
	Number int
	RawURL string
}

// splitAndValidatePath 分割并验证 URL 路径
// 返回 owner, repo, urlType, numberStr, error
func splitAndValidatePath(rawURL string) (string, string, string, string, error) {
	if strings.Contains(rawURL, "#") {
		return "", "", "", "", fmt.Errorf("unsupported URL type")
	}

	urlWithoutScheme := strings.TrimPrefix(rawURL, "https://")
	urlWithoutScheme = strings.TrimPrefix(urlWithoutScheme, "http://")

	if !strings.HasPrefix(urlWithoutScheme, "github.com/") {
		return "", "", "", "", fmt.Errorf("not a GitHub URL")
	}

	path := strings.TrimPrefix(urlWithoutScheme, "github.com/")
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return "", "", "", "", fmt.Errorf("unsupported URL type")
	}

	return parts[0], parts[1], parts[2], parts[3], nil
}

// Parse 解析 GitHub URL
func Parse(rawURL string) (*ParsedURL, error) {
	owner, repo, urlType, numberStr, err := splitAndValidatePath(rawURL)
	if err != nil {
		return nil, err
	}

	var number int
	if _, err := fmt.Sscanf(numberStr, "%d", &number); err != nil {
		return nil, fmt.Errorf("unsupported URL type")
	}

	var typ URLType
	switch urlType {
	case "issues":
		typ = URLTypeIssue
	case "pull":
		typ = URLTypePullRequest
	case "discussions":
		typ = URLTypeDiscussion
	default:
		return nil, fmt.Errorf("unsupported URL type")
	}

	return &ParsedURL{
		Type:   typ,
		Owner:  owner,
		Repo:   repo,
		Number: number,
		RawURL: rawURL,
	}, nil
}
