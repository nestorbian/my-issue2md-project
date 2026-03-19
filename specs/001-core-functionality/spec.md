# issue2md 规格说明书

## 1. 概述

**项目名称：** issue2md

**核心功能：** 一个命令行工具，将 GitHub Issue/PR/Discussion URL 转换为 Markdown 文件。

**目标用户：** 需要存档或导出 GitHub 讨论的开发者、技术文档作者等。

---

## 2. 用户故事

### 2.1 CLI 模式

```
作为一名用户，
我想要运行 `issue2md https://github.com/owner/repo/issues/123`，
以便将 issue 保存为 Markdown 文件。
```

**验收条件：**
- 用户提供 GitHub URL 和可选的输出文件路径
- 工具通过 GitHub API 获取数据，输出 Markdown 到标准输出或文件

### 2.2 未来计划：Web 模式

```
作为一名用户，
我想要通过 Web 界面访问 issue2md，
以便无需安装 CLI 工具即可转换 issue。
```

**注意：** Web 模式不在本期范围内，但架构设计必须便于未来扩展。

---

## 3. 功能性需求

### 3.1 URL 识别

| URL 类型 | 支持 | 示例 |
|----------|------|------|
| Issue | ✅ 是 | `github.com/{owner}/{repo}/issues/{number}` |
| Pull Request | ✅ 是 | `github.com/{owner}/{repo}/pull/{number}` |
| Discussion | ✅ 是 | `github.com/{owner}/{repo}/discussions/{number}` |
| 子链接 | ❌ 否 | `#issuecomment-123` → 报错 |

**不支持 URL 的错误信息：**
```
Error: unsupported URL type. issue2md only supports Issue, PR, and Discussion URLs.
```

### 3.2 CLI 接口

**用法：**
```
issue2md [flags] [output_file]
```

**Flags：**

| Flag | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--user-link` | `-u` | `false` | 将用户名渲染为 GitHub 个人主页链接 |

**参数：**

| 参数 | 必填 | 说明 |
|------|------|------|
| `output_file` | 否 | 输出文件路径。如不指定，则输出到标准输出 |

**行为规则：**
- 如果指定了 `output_file` 且文件已存在 → 报错
- 如果未指定 `output_file` → 输出到标准输出

### 3.3 认证

- Token 从环境变量 `GITHUB_TOKEN` 读取
- Token 为 GitHub Personal Access Token (PAT)
- **遇到限流：** 退出并报错
- **访问私有仓库但无 Token：** 退出并报错
- **需要 Token 但未配置：** 退出并报错

**错误信息：**
```
Error: GitHub API rate limit exceeded. Consider setting GITHUB_TOKEN.
Error: This repository is private. Please set GITHUB_TOKEN environment variable.
Error: GITHUB_TOKEN environment variable is not set.
```

### 3.4 Markdown 输出结构

#### 3.4.1 Frontmatter

```yaml
---
title: "Issue 标题"
author: "username"
created_at: "2026-03-20 10:30 UTC"
type: "issue"  # 或 "pull_request", "discussion"
state: "open"  # 或 "closed"
labels:
  - "bug"
  - "priority-high"
url: "https://github.com/owner/repo/issues/123"
---
```

#### 3.4.2 正文内容

```
# Issue 标题

**Author:** [username](https://github.com/username) *（仅当 --user-link 开启时渲染为链接，否则为纯文本）*
**Created:** 2026-03-20 10:30 UTC
**State:** 🟢 Open / 🔴 Closed
**Labels:** bug, priority-high

---

## Description

Issue 正文内容。支持 **Markdown** 和 *GFM*。

## Comments

---

### Comment by username @ 2026-03-20 12:00 UTC

评论内容。

### Comment by anotheruser @ 2026-03-20 14:30 UTC

另一条评论。
```

#### 3.4.3 内容转换规则

| 元素 | 转换规则 |
|------|----------|
| @mention | 转换为 `[username](https://github.com/username)` |
| #reference | 转换为指向对应 issue/PR 的链接 |
| 图片 | 保留原始 CDN URL (`https://user-images.githubusercontent.com/...`) |
| 代码块 | 保留语言标识（如 ```go） |
| 任务列表 | 保留 `- [ ]` 和 `- [x]` 语法 |

### 3.5 文件命名

- 格式：`{type}-{number}.md`
- 示例：`issue-123.md`、`pull_request-456.md`、`discussion-789.md`
- 与 Markdown/文件名规范冲突的字符需转义或移除

### 3.6 错误处理

| 场景 | 行为 |
|------|------|
| URL 格式无效 | 明确报错 |
| 不支持的 URL 类型（子链接） | 明确报错 |
| 资源不存在（404） | 明确报错 |
| 网络故障（请求中途断开） | 完全失败 + 明确报错，不输出部分内容 |
| 触发限流 | 明确报错 |
| 文件已存在 | 明确报错 |

---

## 4. 非功能性需求

### 4.1 架构设计

```
issue2md/
├── cmd/            # CLI 入口
├── internal/
│   ├── fetcher/    # GitHub API 客户端
│   ├── parser/     # URL 解析和类型检测
│   ├── converter/ # 将 API 响应转换为 Markdown
│   └── writer/     # 输出到文件或标准输出
└── specs/          # 规格说明文档
```

**核心原则：** 核心逻辑（fetcher、parser、converter、writer）不得依赖 CLI，便于未来提取为库或 Web 服务。

### 4.2 依赖

- 仅使用 Go 标准库（`net/http`、`encoding/json` 等）
- 核心功能无外部依赖

### 4.3 配置

- 无配置文件
- 所有设置通过 CLI flags 和环境变量

---

## 5. 验收标准

### 5.1 URL 解析

| 测试用例 | 输入 | 期望结果 |
|----------|------|----------|
| 有效的 Issue URL | `https://github.com/owner/repo/issues/123` | 解析为 Issue |
| 有效的 PR URL | `https://github.com/owner/repo/pull/456` | 解析为 Pull Request |
| 有效的 Discussion URL | `https://github.com/owner/repo/discussions/789` | 解析为 Discussion |
| Issuecomment 子链接 | `github.com/owner/repo/issues/123#issuecomment-456` | 报错：不支持的 URL 类型 |
| 无效 URL | `https://google.com/owner/repo/issues/123` | 报错：不是 GitHub URL |

### 5.2 Markdown 生成

| 测试用例 | 期望结果 |
|----------|----------|
| 带标签的 Issue | 标签出现在 frontmatter 和正文中 |
| 带 @mention 的 Issue | @mention 转换为链接 |
| 带 #reference 的 Issue | #reference 转换为链接 |
| 带图片的 Issue | 图片 URL 保留 |
| 带代码块的 Issue | 语言标识保留 |
| 带任务列表的 Issue | 任务语法保留 |
| --user-link=false | 用户名为纯文本 |
| --user-link=true | 用户名为 GitHub 个人主页链接 |
| 已关闭的 Issue | 状态显示 "🔴 Closed" |

### 5.3 文件输出

| 测试用例 | 输入 | 期望结果 |
|----------|------|----------|
| 无输出文件参数 | `issue2md <url>` | 输出到标准输出 |
| 指定输出文件 | `issue2md <url> issue-123.md` | 写入 issue-123.md |
| 文件已存在 | `issue2md <url> existing.md` | 报错：文件已存在 |
| 目录中文件已存在 | `issue2md <url> ./dir/existing.md` | 报错：文件已存在 |

### 5.4 错误处理

| 测试用例 | 场景 | 期望结果 |
|----------|------|----------|
| 限流 | API 返回 403 限流 | 报错：限流信息 |
| 私有仓库 | 无 Token 访问私有仓库 | 报错：私有仓库信息 |
| 无 Token | 未设置 GITHUB_TOKEN | 报错：缺少 Token |
| 网络错误 | 连接超时 | 报错：网络故障 |
| 未找到 | Issue/PR/Discussion 不存在 | 报错：404 Not Found |

---

## 6. 输出格式示例

### 6.1 Issue 示例

**输入：** `issue2md https://github.com/owner/repo/issues/123`

**输出：**

```markdown
---
title: "Bug: Cannot login with OAuth"
author: "developer123"
created_at: "2026-03-20 10:30 UTC"
type: "issue"
state: "open"
labels:
  - "bug"
  - "oauth"
url: "https://github.com/owner/repo/issues/123"
---

# Bug: Cannot login with OAuth

**Author:** developer123
**Created:** 2026-03-20 10:30 UTC
**State:** 🟢 Open
**Labels:** bug, oauth

---

## Description

When using OAuth to login, the following error occurs:

```
Error: redirect_uri_mismatch
```

## Steps to Reproduce

1. Go to login page
2. Click "Login with GitHub"
3. See error

## Comments

---

### Comment by username @ 2026-03-20 12:00 UTC

Can you provide the `redirect_uri` you configured?

### Comment by developer123 @ 2026-03-20 14:30 UTC

It's set to `https://myapp.com/callback`. See #issue 120 for related discussion.
```

### 6.2 PR 示例

**输入：** `issue2md https://github.com/owner/repo/pull/456 -u`

**输出：**

```markdown
---
title: "feat: Add dark mode support"
author: "designer456"
created_at: "2026-03-19 08:00 UTC"
type: "pull_request"
state: "closed"
labels:
  - "enhancement"
  - "ui"
url: "https://github.com/owner/repo/pull/456"
---

# feat: Add dark mode support

**Author:** [designer456](https://github.com/designer456)
**Created:** 2026-03-19 08:00 UTC
**State:** 🔴 Closed
**Labels:** enhancement, ui

---

## Description

Adds dark mode support using CSS variables.

## Comments

---

### Comment by reviewer @ 2026-03-19 10:00 UTC

Looks good! Merging now.
```

---

## 7. 本期范围外（v1.0）

- Web 界面
- 批量转换（多个 URL）
- GitHub Enterprise 支持
- 图片本地下载
- Discussion 分类信息转换
- `--force` / `-f` 强制覆盖参数
