package parser

type URLType int

const (
	URLTypeIssue       URLType = iota // Issue
	URLTypePullRequest                // Pull Request
	URLTypeDiscussion                 // Discussion
	URLTypeUnknown                    // 不支持的类型
)
