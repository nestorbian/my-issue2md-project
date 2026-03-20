# issue2md 技术实现方案

**版本：** 1.0
**日期：** 2026-03-20
**状态：** 已批准

---

## 1. 技术上下文总结

### 1.1 技术选型

| 组件 | 选型 | 理由 |
|------|------|------|
| **语言** | Go >= 1.24 | 满足项目要求 |
| **Web 框架** | `net/http` (标准库) | 遵循"简单性原则"，不引入 Gin/Echo |
| **GitHub API** | `google/go-github` v4 + GraphQL | 支持 Issue/PR/Discussion 获取 |
| **Markdown 处理** | 标准库 + 少量字符串处理 | 最小化外部依赖 |
| **数据存储** | 无 (实时 API 获取) | 符合本期范围 |

### 1.2 项目约束

- **无配置文件**: 所有设置通过 CLI flags 和环境变量
- **无数据库**: 数据通过 GitHub API 实时获取
- **无外部框架**: 核心逻辑使用标准库

---

## 2. "合宪性"审查

本方案严格遵循 `constitution.md` 的三条核心原则：

### 2.1 第一条：简单性原则 ✅

| 条款 | 实施方案 |
|------|----------|
| 1.1 (YAGNI) | 仅实现 spec.md 明确要求的功能：Issue/PR/Discussion 转换 |
| 1.2 (标准库优先) | Web 服务使用 `net/http`；Markdown 处理使用字符串操作 |
| 1.3 (反过度工程) | 采用简单的函数组合而非复杂接口层次 |

### 2.2 第二条：测试先行铁律 ✅

| 条款 | 实施方案 |
|------|----------|
| 2.1 (TDD 循环) | 每个包先编写失败测试，再实现功能 |
| 2.2 (表格驱动) | 所有单元测试采用 Table-Driven Tests 风格 |
| 2.3 (拒绝Mocks) | 优先使用真实 HTTP 客户端（可配置 endpoint）进行集成测试 |

### 2.3 第三条：明确性原则 ✅

| 条款 | 实施方案 |
|------|----------|
| 3.1 (错误处理) | 所有错误使用 `fmt.Errorf("...: %w", err)` 包装，区分错误类型 |
| 3.2 (无全局变量) | 依赖通过结构体成员或函数参数显式注入 |

---

## 3. 项目结构细化

```
issue2md/
├── cmd/
│   └── cli/
│       └── main.go           # CLI 入口，解析 flags 和环境变量
├── internal/
│   ├── fetcher/              # GitHub API 客户端 (使用 go-github)
│   │   ├── fetcher.go        # 核心 fetcher 接口和实现
│   │   ├── github.go         # GitHub REST/GraphQL API 封装
│   │   └── fetcher_test.go   # 表格驱动测试
│   ├── parser/               # URL 解析和类型检测
│   │   ├── parser.go         # URL 解析逻辑
│   │   ├── parser_test.go    # 表格驱动测试
│   │   └── urltype.go        # URL 类型定义
│   ├── converter/            # API 响应 → Markdown 转换
│   │   ├── converter.go      # 核心转换逻辑
│   │   ├── converter_test.go # 表格驱动测试
│   │   └── markdown.go       # Markdown 格式化辅助函数
│   └── writer/               # 输出到文件或标准输出
│       ├── writer.go         # 写入逻辑
│       └── writer_test.go    # 表格驱动测试
├── specs/
│   └── 001-core-functionality/
│       ├── spec.md            # 功能规格说明
│       └── plan.md            # 本文档
├── Makefile                   # 构建和测试入口
└── go.mod                     # Go 模块定义
```

### 3.1 包职责与依赖关系

```
┌─────────────────────────────────────────────────────────┐
│                      cmd/cli (main.go)                   │
│  - 解析 flags: --user-link, output_file                 │
│  - 读取环境变量: GITHUB_TOKEN                           │
└────────────────────────┬────────────────────────────────┘
                         │ 依赖注入
         ┌───────────────┼───────────────┐
         ▼               ▼               ▼
   ┌──────────┐   ┌──────────┐    ┌──────────┐
   │ parser   │   │ fetcher  │    │ writer   │
   └────┬─────┘   └────┬─────┘    └────┬─────┘
        │              │               │
        │              │               │
        └──────────────┼───────────────┘
                       ▼
              ┌──────────────┐
              │  converter   │
              └──────────────┘
```

| 包 | 职责 | 对外依赖 | 被谁依赖 |
|---|------|----------|----------|
| `parser` | 解析 GitHub URL，判断类型 (Issue/PR/Discussion) | 无 | `cmd/cli` |
| `fetcher` | 通过 GitHub API 获取数据 | `google/go-github` | `cmd/cli` |
| `converter` | 将 API 响应转换为 Markdown 字符串 | 无 | `cmd/cli` |
| `writer` | 将 Markdown 写入文件或 stdout | 无 | `cmd/cli` |

---

## 4. 核心数据结构

### 4.1 URL 类型枚举

```go
// internal/parser/urltype.go
package parser

type URLType int

const (
    URLTypeIssue       URLType = iota // Issue
    URLTypePullRequest                // Pull Request
    URLTypeDiscussion                 // Discussion
    URLTypeUnknown                    // 不支持的类型
)
```

### 4.2 解析后的 URL 信息

```go
// internal/parser/parser.go
type ParsedURL struct {
    Type      URLType
    Owner     string
    Repo      string
    Number    int
    RawURL    string
}
```

### 4.3 GitHub 资源统一表示

```go
// internal/fetcher/resource.go
package fetcher

// GitHubResource 是 Issue/PR/Discussion 的统一表示
type GitHubResource struct {
    Type       string    // "issue" | "pull_request" | "discussion"
    Title      string
    Author     string
    CreatedAt  time.Time
    UpdatedAt  time.Time
    State      string    // "open" | "closed"
    Body       string    // Markdown 内容
    Labels     []string
    Comments   []Comment
    URL        string    // 原始 GitHub URL
}

// Comment 评论结构
type Comment struct {
    Author    string
    CreatedAt time.Time
    Body      string
}
```

### 4.4 Markdown 输出结构

```go
// internal/converter/output.go
package converter

// Output 包含转换后的所有 Markdown 内容
type Output struct {
    Frontmatter string // YAML frontmatter
    Body        string // Markdown 正文
    FullContent string // 完整内容 (frontmatter + body)
}
```

---

## 5. 接口设计

### 5.1 Parser 接口

```go
// internal/parser/parser.go
package parser

// Parser 解析 GitHub URL
type Parser interface {
    // Parse 解析 URL，返回 ParsedURL 或错误
    Parse(rawURL string) (*ParsedURL, error)
}

// ValidateURL 验证 URL 是否为有效的 GitHub URL
type ValidateURL func(rawURL string) error
```

**实现要求：**
- 支持 `github.com/{owner}/{repo}/issues/{number}`
- 支持 `github.com/{owner}/{repo}/pull/{number}`
- 支持 `github.com/{owner}/{repo}/discussions/{number}`
- 拒绝子链接如 `#issuecomment-123`，返回明确错误

### 5.2 Fetcher 接口

```go
// internal/fetcher/fetcher.go
package fetcher

// Fetcher 获取 GitHub 资源
type Fetcher interface {
    // Fetch 根据 ParsedURL 获取 GitHubResource
    Fetch(ctx context.Context, url *parser.ParsedURL) (*GitHubResource, error)
}

// GitHubClient GitHub API 客户端接口 (便于测试)
type GitHubClient interface {
    FetchIssue(ctx context.Context, owner, repo string, number int) (*GitHubResource, error)
    FetchPullRequest(ctx context.Context, owner, repo string, number int) (*GitHubResource, error)
    FetchDiscussion(ctx context.Context, owner, repo string, number int) (*GitHubResource, error)
}
```

**错误处理：**
- 限流 (403): `fmt.Errorf("GitHub API rate limit exceeded: %w", err)`
- 私有仓库无 Token: `fmt.Errorf("repository is private: %w", err)`
- Token 未设置: `fmt.Errorf("GITHUB_TOKEN is not set: %w", err)`
- 资源不存在 (404): `fmt.Errorf("resource not found: %w", err)`

### 5.3 Converter 接口

```go
// internal/converter/converter.go
package converter

// Converter 转换 GitHub 资源为 Markdown
type Converter interface {
    // Convert 将 GitHubResource 转换为 Output
    Convert(resource *fetcher.GitHubResource, userLink bool) (*Output, error)
}
```

**Markdown 转换规则 (spec 3.4.3)：**
| 元素 | 转换 |
|------|------|
| @mention | `[username](https://github.com/username)` |
| #reference | 转换为指向对应 issue/PR 的链接 |
| 图片 | 保留原始 CDN URL |
| 代码块 | 保留语言标识 |
| 任务列表 | 保留 `- [ ]` 和 `- [x]` |

### 5.4 Writer 接口

```go
// internal/writer/writer.go
package writer

// Writer 写入 Markdown 输出
type Writer interface {
    // Write 将 content 写入指定目标
    // 如果 outputFile 为空，写入 stdout
    // 如果文件已存在，返回错误
    Write(content string, outputFile string) error
}
```

**错误处理：**
- 文件已存在: `fmt.Errorf("file already exists: %s", outputFile)`

---

## 6. 错误类型体系

```go
// internal/errors/errors.go
package errors

type ErrorType int

const (
    ErrTypeInvalidURL     ErrorType = iota // URL 格式无效
    ErrTypeUnsupportedURL                  // 不支持的 URL 类型
    ErrTypeNotFound                         // 资源不存在 (404)
    ErrTypeRateLimit                        // API 限流
    ErrTypeAuthRequired                     // 需要认证
    ErrTypeNetwork                          // 网络故障
    ErrTypeFileExists                       // 文件已存在
)

// AppError 包含错误类型和原始错误
type AppError struct {
    Type    ErrorType
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %w", e.Message, e.Err)
    }
    return e.Message
}
```

---

## 7. CLI 设计

### 7.1 命令行接口

```
issue2md [flags] [output_file]
```

| Flag | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--user-link` | `-u` | false | 将用户名渲染为链接 |

**位置参数：**
- `output_file`: 可选，输出文件路径

### 7.2 环境变量

| 变量 | 必填 | 说明 |
|------|------|------|
| `GITHUB_TOKEN` | 是 | GitHub Personal Access Token |

---

## 8. 测试策略

### 8.1 表格驱动测试示例

```go
// internal/parser/parser_test.go
func TestParser_Parse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *ParsedURL
        wantErr error
    }{
        {
            name:  "valid issue URL",
            input: "https://github.com/owner/repo/issues/123",
            want: &ParsedURL{
                Type:   URLTypeIssue,
                Owner:  "owner",
                Repo:   "repo",
                Number: 123,
            },
            wantErr: nil,
        },
        // ... 更多测试用例
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := NewParser()
            got, err := p.Parse(tt.input)
            // 验证逻辑
        })
    }
}
```

### 8.2 集成测试

- 使用 `net/http/httptest` 启动本地 HTTP 服务器
- 或使用 `os.Getenv("GITHUB_TOKEN")` 进行真实 API 测试

---

## 9. 实现步骤

### Phase 1: 基础设施
1. 初始化 Go 模块 (`go mod init`)
2. 创建目录结构
3. 定义错误类型
4. 编写 Makefile

### Phase 2: Parser 包
1. 定义 URLType 和 ParsedURL
2. 实现 URL 解析逻辑
3. 编写表格驱动测试
4. 验证所有 spec 5.1 测试用例

### Phase 3: Fetcher 包
1. 封装 go-github 客户端
2. 实现 Issue/PR/Discussion 获取
3. 处理认证和限流错误
4. 编写测试

### Phase 4: Converter 包
1. 实现 Markdown 生成逻辑
2. 实现 Frontmatter 格式化
3. 处理 @mention 和 #reference 转换
4. 验证所有 spec 5.2 测试用例

### Phase 5: Writer 包
1. 实现文件写入逻辑
2. 实现 stdout 输出
3. 处理文件已存在错误
4. 验证所有 spec 5.3 测试用例

### Phase 6: CLI 集成
1. 实现 main.go
2. 解析 flags 和环境变量
3. 串联所有组件
4. 端到端测试

---

## 10. 验收清单

| 功能 | 验收条件 |
|------|----------|
| URL 解析 | Issue/PR/Discussion URL 正确解析 |
| URL 校验 | 子链接和无效 URL 正确报错 |
| Markdown 生成 | Frontmatter 和正文格式正确 |
| 内容转换 | @mention、#reference、图片、代码块、任务列表正确处理 |
| 文件输出 | stdout 和文件输出正确 |
| 错误处理 | 所有错误场景正确报错 |
| 认证 | Token 缺失和限流正确处理 |
