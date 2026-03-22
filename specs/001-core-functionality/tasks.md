# issue2md 任务列表

**版本：** 1.0
**日期：** 2026-03-20
**状态：** 规划中

---

## 任务列表说明

- **[T]** = 测试先行任务（TDD 强制要求先完成）
- **[P]** = 可并行执行的任务（无依赖关系）
- **依赖关系** 通过 `blockedBy` 字段标注

---

## Phase 1: Foundation (基础设施)

### 1.1 项目初始化

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 1.1.1 | 初始化 Go 模块 | `go.mod` | 创建 `go.mod`，定义模块名 `github.com/issue2md/issue2md`，Go 版本 >= 1.24 | - | completed |
| 1.1.2 | 创建目录结构 | - | 创建 `cmd/cli/`、`internal/{parser,fetcher,converter,writer,errors}/` 目录 | 1.1.1 | completed |
| 1.1.3 | 创建 Makefile | `Makefile` | 定义 `make test`、`make build`、`make web` 等标准目标 | 1.1.2 | completed |

### 1.2 错误类型定义

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 1.2.1 | 定义错误类型枚举 | `internal/errors/errors.go` | 创建 `ErrorType` 枚举：InvalidURL、UnsupportedURL、NotFound、RateLimit、AuthRequired、Network、FileExists | 1.1.2 | completed |
| 1.2.2 | 定义 AppError 结构 | `internal/errors/errors.go` | 创建 `AppError` 结构体，包含 Type、Message、Err 字段，实现 `Error()` 方法 | 1.2.1 | completed |
| 1.2.3 | [T] 编写 errors 包测试 | `internal/errors/errors_test.go` | 表格驱动测试：验证各错误类型的 Error() 输出格式 | 1.2.2 | completed |

---

## Phase 2: Parser 包 (URL 解析，TDD)

### 2.1 Parser 数据结构

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 2.1.1 | [T] 定义 URLType 测试 | `internal/parser/urltype_test.go` | 测试 URLType 枚举值：Issue=0、PullRequest=1、Discussion=2、Unknown=3 | 1.1.2 | completed |
| 2.1.2 | 定义 URLType | `internal/parser/urltype.go` | 创建 URLType 类型和常量定义 | 2.1.1 | completed |
| 2.1.3 | [T] 定义 ParsedURL 测试 | `internal/parser/parser_test.go` | 测试 ParsedURL 结构体字段：Type、Owner、Repo、Number、RawURL | 2.1.2 | completed |
| 2.1.4 | 定义 ParsedURL | `internal/parser/parser.go` | 创建 ParsedURL 结构体 | 2.1.3 | completed |

### 2.2 Parser 接口与实现

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 2.2.1 | [T] 测试 Parse Issue URL | `internal/parser/parser_test.go` | 表格驱动：测试 `https://github.com/owner/repo/issues/123` 解析为 Issue | 2.1.4 | completed |
| 2.2.2 | [T] 测试 Parse PR URL | `internal/parser/parser_test.go` | 表格驱动：测试 `https://github.com/owner/repo/pull/456` 解析为 Pull Request | 2.2.1 [P] | completed |
| 2.2.3 | [T] 测试 Parse Discussion URL | `internal/parser/parser_test.go` | 表格驱动：测试 `https://github.com/owner/repo/discussions/789` 解析为 Discussion | 2.2.1 [P] | completed |
| 2.2.4 | [T] 测试子链接拒绝 | `internal/parser/parser_test.go` | 表格驱动：测试 `#issuecomment-456` 报错 "unsupported URL type" | 2.2.2 | completed |
| 2.2.5 | [T] 测试无效 GitHub URL | `internal/parser/parser_test.go` | 表格驱动：测试 `https://google.com/owner/repo/issues/123` 报错 "not a GitHub URL" | 2.2.3 [P] | completed |
| 2.2.6 | [T] 测试 http vs https | `internal/parser/parser_test.go` | 表格驱动：测试 `http://github.com/...` 应支持 | 2.2.5 [P] | completed |
| 2.2.7 | 实现 Parse 方法 | `internal/parser/parser.go` | 实现 `Parser.Parse()` 方法，支持 Issue/PR/Discussion 解析 | 2.2.4, 2.2.5, 2.2.6 | completed |

### 2.3 Parser 验收测试

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 2.3.1 | [T] 验收测试：URL 解析 | `internal/parser/parser_acceptance_test.go` | 根据 spec 5.1 验收标准：有效 Issue/PR/Discussion URL、子链接、无效 URL | 2.2.7 | completed |

---

## Phase 3: Fetcher 包 (GitHub API 交互，TDD)

### 3.1 Fetcher 数据结构

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 3.1.1 | [T] 测试 GitHubResource 结构 | `internal/fetcher/resource_test.go` | 测试 GitHubResource 字段：Type、Title、Author、CreatedAt、State、Body、Labels、Comments、URL | 1.1.2 | completed |
| 3.1.2 | [T] 测试 Comment 结构 | `internal/fetcher/resource_test.go` | 测试 Comment 字段：Author、CreatedAt、Body | 3.1.1 [P] | completed |
| 3.1.3 | 定义 GitHubResource | `internal/fetcher/resource.go` | 创建 GitHubResource 结构体 | 3.1.2 | completed |
| 3.1.4 | 定义 Comment | `internal/fetcher/resource.go` | 创建 Comment 结构体 | 3.1.3 [P] | completed |

### 3.2 Fetcher 接口定义

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 3.2.1 | [T] 测试 Fetcher 接口签名 | `internal/fetcher/fetcher_test.go` | 测试 `Fetch(ctx, *parser.ParsedURL) (*GitHubResource, error)` 方法签名 | 3.1.4 | completed |
| 3.2.2 | 定义 Fetcher 接口 | `internal/fetcher/fetcher.go` | 创建 `Fetcher` 接口定义 | 3.2.1 | completed |
| 3.2.3 | [T] 测试 GitHubClient 接口 | `internal/fetcher/fetcher_test.go` | 测试 `FetchIssue`、`FetchPullRequest`、`FetchDiscussion` 方法签名 | 3.2.2 [P] | completed |
| 3.2.4 | 定义 GitHubClient 接口 | `internal/fetcher/fetcher.go` | 创建 `GitHubClient` 接口定义 | 3.2.3 | completed |

### 3.3 GitHub API 实现

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 3.3.1 | [T] 测试 Issue 获取 | `internal/fetcher/github_test.go` | 表格驱动：测试 FetchIssue 调用 go-github 并转换结果 | 3.2.4 | pending |
| 3.3.2 | [T] 测试 PR 获取 | `internal/fetcher/github_test.go` | 表格驱动：测试 FetchPullRequest 调用 go-github 并转换结果 | 3.3.1 [P] | pending |
| 3.3.3 | [T] 测试 Discussion 获取 | `internal/fetcher/github_test.go` | 表格驱动：测试 FetchDiscussion 调用 go-github 并转换结果 | 3.3.2 [P] | pending |
| 3.3.4 | 实现 GitHubClient | `internal/fetcher/github.go` | 实现 `githubClient` 结构体，封装 go-github v4 客户端 | 3.3.3 | pending |
| 3.3.5 | [T] 测试限流错误处理 | `internal/fetcher/github_test.go` | 测试 API 返回 403 时包装为 "rate limit exceeded" 错误 | 3.3.4 | pending |
| 3.3.6 | [T] 测试私有仓库错误 | `internal/fetcher/github_test.go` | 测试 404 或 403 私有仓库错误包装 | 3.3.5 [P] | pending |
| 3.3.7 | [T] 测试 Token 缺失错误 | `internal/fetcher/github_test.go` | 测试环境变量 GITHUB_TOKEN 未设置时的错误 | 3.3.6 [P] | pending |
| 3.3.8 | 实现限流错误处理 | `internal/fetcher/github.go` | 在 github.go 中实现错误包装逻辑 | 3.3.7 | pending |
| 3.3.9 | [T] 测试 404 Not Found | `internal/fetcher/github_test.go` | 测试资源不存在时的错误处理 | 3.3.8 [P] | pending |
| 3.3.10 | 实现 Fetch 方法 | `internal/fetcher/fetcher.go` | 实现 `Fetcher.Fetch()` 方法，根据 URLType 分发到不同获取方法 | 3.3.9 | pending |

### 3.4 Fetcher 集成测试

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 3.4.1 | [T] 集成测试：完整 Fetch 流程 | `internal/fetcher/integration_test.go` | 使用 httptest 模拟 GitHub API，测试完整获取流程 | 3.3.10 | pending |

---

## Phase 4: Converter 包 (Markdown 转换，TDD)

### 4.1 Converter 数据结构

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 4.1.1 | [T] 测试 Output 结构 | `internal/converter/output_test.go` | 测试 Output 字段：Frontmatter、Body、FullContent | 1.1.2 | pending |
| 4.1.2 | 定义 Output | `internal/converter/output.go` | 创建 Output 结构体 | 4.1.1 | pending |

### 4.2 Converter 接口定义

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 4.2.1 | [T] 测试 Converter 接口签名 | `internal/converter/converter_test.go` | 测试 `Convert(*fetcher.GitHubResource, bool) (*Output, error)` 方法签名 | 4.1.2 | pending |
| 4.2.2 | 定义 Converter 接口 | `internal/converter/converter.go` | 创建 `Converter` 接口 | 4.2.1 | pending |

### 4.3 Markdown 格式化

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 4.3.1 | [T] 测试 Frontmatter 格式化 | `internal/converter/markdown_test.go` | 表格驱动：测试 spec 3.4.1 格式的 YAML frontmatter 生成 | 4.2.2 | pending |
| 4.3.2 | [T] 测试标题渲染 | `internal/converter/markdown_test.go` | 表格驱动：测试 `# Title` 渲染 | 4.3.1 [P] | pending |
| 4.3.3 | [T] 测试 Author 渲染 (userLink=false) | `internal/converter/markdown_test.go` | 表格驱动：测试 `**Author:** username` | 4.3.2 [P] | pending |
| 4.3.4 | [T] 测试 Author 渲染 (userLink=true) | `internal/converter/markdown_test.go` | 表格驱动：测试 `**Author:** [username](https://github.com/username)` | 4.3.3 [P] | pending |
| 4.3.5 | [T] 测试 State 渲染 | `internal/converter/markdown_test.go` | 表格驱动：测试 `🟢 Open` 和 `🔴 Closed` | 4.3.4 [P] | pending |
| 4.3.6 | [T] 测试 Labels 渲染 | `internal/converter/markdown_test.go` | 表格驱动：测试 `**Labels:** label1, label2` | 4.3.5 [P] | pending |
| 4.3.7 | [T] 测试 @mention 转换 | `internal/converter/markdown_test.go` | 表格驱动：测试 `@username` 转换为 `[username](https://github.com/username)` | 4.3.6 | pending |
| 4.3.8 | [T] 测试 #reference 转换 | `internal/converter/markdown_test.go` | 表格驱动：测试 `#123` 转换为 issue/PR 链接 | 4.3.7 [P] | pending |
| 4.3.9 | [T] 测试图片保留 | `internal/converter/markdown_test.go` | 表格驱动：测试 `![](https://user-images.githubusercontent.com/...)` URL 保留 | 4.3.8 [P] | pending |
| 4.3.10 | [T] 测试代码块保留 | `internal/converter/markdown_test.go` | 表格驱动：测试 ` ```go ` 等语言标识保留 | 4.3.9 [P] | pending |
| 4.3.11 | [T] 测试任务列表保留 | `internal/converter/markdown_test.go` | 表格驱动：测试 `- [ ]` 和 `- [x]` 语法保留 | 4.3.10 [P] | pending |
| 4.3.12 | 实现 Frontmatter 格式化 | `internal/converter/markdown.go` | 实现 frontmatter 生成函数 | 4.3.1 | pending |
| 4.3.13 | 实现 Markdown 渲染 | `internal/converter/markdown.go` | 实现标题、Author、State、Labels 渲染函数 | 4.3.12 | pending |
| 4.3.14 | 实现内容转换 | `internal/converter/markdown.go` | 实现 @mention、#reference 转换函数 | 4.3.13 | pending |

### 4.4 Comments 渲染

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 4.4.1 | [T] 测试 Comments 渲染 | `internal/converter/markdown_test.go` | 表格驱动：测试 spec 3.4.2 的评论格式 `### Comment by username @ timestamp` | 4.3.14 | pending |
| 4.4.2 | [T] 测试空 Comments | `internal/converter/markdown_test.go` | 测试无评论时不输出 Comments 部分 | 4.4.1 [P] | pending |
| 4.4.3 | 实现 Comments 渲染 | `internal/converter/markdown.go` | 实现评论渲染函数 | 4.4.2 | pending |

### 4.5 Converter 核心实现

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 4.5.1 | [T] 测试 Convert 完整流程 | `internal/converter/converter_test.go` | 表格驱动：传入完整 GitHubResource，验证输出符合 spec 6.1/6.2 示例 | 4.4.3 | pending |
| 4.5.2 | [T] 测试 userLink 参数 | `internal/converter/converter_test.go` | 测试 userLink=true 和 userLink=false 的不同输出 | 4.5.1 [P] | pending |
| 4.5.3 | 实现 Converter | `internal/converter/converter.go` | 实现 `converter` 结构体和 `Convert` 方法 | 4.5.2 | pending |

### 4.6 Converter 验收测试

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 4.6.1 | [T] 验收测试：Markdown 生成 | `internal/converter/acceptance_test.go` | 根据 spec 5.2 验收标准：标签、@mention、#reference、图片、代码块、任务列表、userLink、closed state | 4.5.3 | pending |

---

## Phase 5: Writer 包 (输出写入，TDD)

### 5.1 Writer 接口定义

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 5.1.1 | [T] 测试 Writer 接口签名 | `internal/writer/writer_test.go` | 测试 `Write(content, outputFile string) error` 方法签名 | 1.1.2 | pending |
| 5.1.2 | 定义 Writer 接口 | `internal/writer/writer.go` | 创建 `Writer` 接口 | 5.1.1 | pending |

### 5.2 文件写入实现

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 5.2.1 | [T] 测试写入 stdout | `internal/writer/writer_test.go` | 测试 outputFile="" 时输出到标准输出 | 5.1.2 | pending |
| 5.2.2 | [T] 测试写入文件 | `internal/writer/writer_test.go` | 表格驱动：测试写入新文件 | 5.2.1 [P] | pending |
| 5.2.3 | [T] 测试文件已存在错误 | `internal/writer/writer_test.go` | 表格驱动：测试文件存在时返回 "file already exists" 错误 | 5.2.2 | pending |
| 5.2.4 | [T] 测试目录不存在 | `internal/writer/writer_test.go` | 测试 `./dir/existing.md` 目录不存在场景 | 5.2.3 [P] | pending |
| 5.2.5 | 实现 Writer | `internal/writer/writer.go` | 实现 `fileWriter` 结构体和 `Write` 方法 | 5.2.4 | pending |

### 5.3 Writer 验收测试

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 5.3.1 | [T] 验收测试：文件输出 | `internal/writer/acceptance_test.go` | 根据 spec 5.3 验收标准：stdout、文件输出、文件已存在 | 5.2.5 | pending |

---

## Phase 6: CLI Assembly (命令行入口集成)

### 6.1 CLI 入口

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 6.1.1 | [T] 测试环境变量读取 | `cmd/cli/main_test.go` | 测试 `GITHUB_TOKEN` 环境变量读取 | 1.1.2 | pending |
| 6.1.2 | [T] 测试 flags 解析 | `cmd/cli/main_test.go` | 测试 `--user-link` / `-u` 和 `output_file` 参数解析 | 6.1.1 [P] | pending |
| 6.1.3 | [T] 测试 Token 缺失错误 | `cmd/cli/main_test.go` | 测试 GITHUB_TOKEN 未设置时的错误输出 | 6.1.2 | pending |
| 6.1.4 | 实现环境变量读取 | `cmd/cli/main.go` | 实现 GITHUB_TOKEN 读取逻辑 | 6.1.3 | pending |
| 6.1.5 | 实现 flags 解析 | `cmd/cli/main.go` | 实现 `--user-link` / `-u` 和 `output_file` 参数解析 | 6.1.4 | pending |

### 6.2 组件串联

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 6.2.1 | [T] 测试完整流程 | `cmd/cli/main_test.go` | 集成测试：URL → Parser → Fetcher → Converter → Writer 完整流程 | 6.1.5 | pending |
| 6.2.2 | [T] 测试错误传播 | `cmd/cli/main_test.go` | 测试各层错误是否正确传播到 CLI 输出 | 6.2.1 [P] | pending |
| 6.2.3 | 实现 CLI 主流程 | `cmd/cli/main.go` | 实现 main 函数，串联各组件 | 6.2.2 | pending |

### 6.3 CLI 验收测试

| Task ID | 任务名称 | 文件 | 描述 | 依赖 | 状态 |
|---------|----------|------|------|------|------|
| 6.3.1 | [T] 验收测试：错误处理 | `cmd/cli/acceptance_test.go` | 根据 spec 5.4 验收标准：限流、私有仓库、无Token、网络错误、未找到 | 6.2.3 | pending |
| 6.3.2 | [T] 端到端测试 | `cmd/cli/e2e_test.go` | 使用真实 GitHub API Token（或 mock）进行端到端测试 | 6.3.1 | pending |

---

## 任务依赖图

```
Phase 1: Foundation
├── 1.1.1 (go.mod) ──→ 1.1.2 (目录结构) ──→ 1.1.3 (Makefile)
└── 1.2.1 ──→ 1.2.2 ──→ 1.2.3

Phase 2: Parser
├── 2.1.1 ──→ 2.1.2 ──→ 2.1.3 ──→ 2.1.4
├── 2.2.1 ──→ 2.2.2 ──→ 2.2.3 ──→ 2.2.4 ──→ 2.2.5 ──→ 2.2.6 ──→ 2.2.7
└── 2.3.1

Phase 3: Fetcher
├── 3.1.1 ──→ 3.1.2 ──→ 3.1.3 ──→ 3.1.4
├── 3.2.1 ──→ 3.2.2 ──→ 3.2.3 ──→ 3.2.4
├── 3.3.1 ──→ 3.3.2 ──→ 3.3.3 ──→ 3.3.4 ──→ 3.3.5 ──→ 3.3.6 ──→ 3.3.7 ──→ 3.3.8 ──→ 3.3.9 ──→ 3.3.10
└── 3.4.1

Phase 4: Converter
├── 4.1.1 ──→ 4.1.2
├── 4.2.1 ──→ 4.2.2
├── 4.3.1 ──→ 4.3.2 ──→ 4.3.3 ──→ 4.3.4 ──→ 4.3.5 ──→ 4.3.6 ──→ 4.3.7 ──→ 4.3.8 ──→ 4.3.9 ──→ 4.3.10 ──→ 4.3.11 ──→ 4.3.12 ──→ 4.3.13 ──→ 4.3.14
├── 4.4.1 ──→ 4.4.2 ──→ 4.4.3
├── 4.5.1 ──→ 4.5.2 ──→ 4.5.3
└── 4.6.1

Phase 5: Writer
├── 5.1.1 ──→ 5.1.2
├── 5.2.1 ──→ 5.2.2 ──→ 5.2.3 ──→ 5.2.4 ──→ 5.2.5
└── 5.3.1

Phase 6: CLI
├── 6.1.1 ──→ 6.1.2 ──→ 6.1.3 ──→ 6.1.4 ──→ 6.1.5
├── 6.2.1 ──→ 6.2.2 ──→ 6.2.3
└── 6.3.1 ──→ 6.3.2
```

---

## 并行任务组

以下任务可以并行执行（无依赖关系）：

| 并行组 | 任务 |
|--------|------|
| **Group A** | 2.2.2, 2.2.3, 2.2.5, 2.2.6 (Parser 各类 URL 测试) |
| **Group B** | 3.1.2, 3.2.3, 3.3.2, 3.3.3 (Fetcher 各类测试) |
| **Group C** | 4.3.2 ~ 4.3.11 (Markdown 格式化测试，并行) |
| **Group D** | 4.4.2, 4.5.2 (额外测试) |
| **Group E** | 5.2.1, 5.2.2, 5.2.4 (Writer 测试) |
| **Group F** | 6.1.2, 6.2.2 (CLI 测试) |

---

## 验收标准映射

| 验收标准 | 相关任务 |
|----------|----------|
| **spec 5.1** URL 解析 | 2.3.1 |
| **spec 5.2** Markdown 生成 | 4.6.1 |
| **spec 5.3** 文件输出 | 5.3.1 |
| **spec 5.4** 错误处理 | 6.3.1 |
| **端到端** | 6.3.2 |

---

## 统计信息

| 阶段 | 测试任务 | 实现任务 | 小计 |
|------|----------|----------|------|
| Phase 1: Foundation | 1 | 3 | 4 |
| Phase 2: Parser | 7 | 3 | 10 |
| Phase 3: Fetcher | 10 | 4 | 14 |
| Phase 4: Converter | 12 | 5 | 17 |
| Phase 5: Writer | 4 | 1 | 5 |
| Phase 6: CLI | 6 | 2 | 8 |
| **总计** | **40** | **18** | **58** |
