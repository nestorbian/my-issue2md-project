# API Sketch

本文档描述 `internal/converter` 和 `internal/fetcher` 包对外暴露的主要接口，作为后续开发的参考。

---

## 1. internal/fetcher

GitHub API 客户端，负责从 GitHub 获取 Issue/PR/Discussion 数据。

### 1.1 类型定义

```go
// ResourceType 表示 GitHub 资源类型
type ResourceType string

const (
    ResourceTypeIssue        ResourceType = "issue"
    ResourceTypePullRequest  ResourceType = "pull_request"
    ResourceTypeDiscussion   ResourceType = "discussion"
)
```

### 1.2 接口

```go
// Fetcher 从 GitHub API 获取资源的接口
type Fetcher interface {
    // Fetch 获取指定类型的 GitHub 资源
    Fetch(ctx context.Context, owner, repo string, number int, resourceType ResourceType) (*Response, error)
}

// Response GitHub API 响应
type Response struct {
    Type       ResourceType
    Title      string
    Body       string
    Author     string
    CreatedAt  time.Time
    State      string // "open" 或 "closed"
    Labels     []string
    URL        string
    Comments   []Comment
}

// Comment 评论
type Comment struct {
    Author    string
    Body      string
    CreatedAt time.Time
}
```

### 1.3 错误类型

```go
// ErrRateLimitExceeded GitHub API 限流错误
var ErrRateLimitExceeded = errors.New("GitHub API rate limit exceeded")

// ErrPrivateRepository 私有仓库访问错误
var ErrPrivateRepository = errors.New("this repository is private")

// ErrNotFound 资源不存在错误
var ErrNotFound = errors.New("resource not found")
```

---

## 2. internal/parser

URL 解析和类型检测。

### 2.1 接口

```go
// Parser 解析 GitHub URL 的接口
type Parser interface {
    // Parse 解析 URL，返回资源类型和组件
    Parse(url string) (*ParsedURL, error)
}

// ParsedURL 解析后的 URL 组件
type ParsedURL struct {
    Type       ResourceType
    Owner      string
    Repo       string
    Number     int
}
```

### 2.2 错误类型

```go
// ErrInvalidURL URL 格式无效
var ErrInvalidURL = errors.New("invalid URL format")

// ErrUnsupportedURL 不支持的 URL 类型
var ErrUnsupportedURL = errors.New("unsupported URL type")
```

---

## 3. internal/converter

将 API 响应转换为 Markdown 格式。

### 3.1 接口

```go
// Converter 将 GitHub 资源转换为 Markdown 的接口
type Converter interface {
    // Convert 将 fetcher.Response 转换为 Markdown
    Convert(resp *fetcher.Response, opts *ConvertOptions) (string, error)
}

// ConvertOptions 转换选项
type ConvertOptions struct {
    UserLink bool // 是否将用户名渲染为链接
}
```

### 3.2 转换流程

1. **Frontmatter 生成**：包含 title、author、created_at、type、state、labels、url
2. **正文生成**：标题、作者、创建时间、状态、标签、主楼内容
3. **评论生成**：扁平展示，每条包含作者、时间、內容
4. **内容转换**：
   - `@mention` → `[username](https://github.com/username)`
   - `#number` → 对应 issue/PR 的链接
   - 图片 URL 保留原始地址
   - 代码块保留语言标识
   - 任务列表保留 `- [ ]` 和 `- [x]` 语法

---

## 4. internal/writer

输出到文件或标准输出。

### 4.1 接口

```go
// Writer 输出 Markdown 的接口
type Writer interface {
    // Write 将内容写入目标（文件或 stdout）
    Write(content string, dest string) error
}
```

### 4.2 错误类型

```go
// ErrFileExists 文件已存在
var ErrFileExists = errors.New("file already exists")
```

### 4.3 行为规则

| dest | 行为 |
|------|------|
| `""` (空) | 输出到标准输出 |
| `"filename.md"` | 写入当前目录 |
| `"./path/filename.md"` | 写入指定路径 |

---

## 5. 典型使用流程

```go
// 1. 解析 URL
parser := parser.NewParser()
parsed, err := parser.Parse("https://github.com/owner/repo/issues/123")
if err != nil {
    // 处理错误
}

// 2. 获取数据
fetcher := fetcher.NewFetcher(os.Getenv("GITHUB_TOKEN"))
resp, err := fetcher.Fetch(ctx, parsed.Owner, parsed.Repo, parsed.Number, parsed.Type)
if err != nil {
    // 处理错误
}

// 3. 转换为 Markdown
converter := converter.NewConverter()
md, err := converter.Convert(resp, &converter.ConvertOptions{UserLink: false})
if err != nil {
    // 处理错误
}

// 4. 输出
writer := writer.NewWriter()
err = writer.Write(md, outputFile)
if err != nil {
    // 处理错误
}
```
