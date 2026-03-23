package converter

// Output 包含转换后的所有 Markdown 内容
type Output struct {
	Frontmatter string // YAML frontmatter
	Body        string // Markdown 正文
	FullContent string // 完整内容 (frontmatter + body)
}
