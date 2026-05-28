# GitCode API Client

Go 语言的 GitCode / AtomGit API 客户端库，提供对 GitCode 平台的完整访问。

[![Go Reference](https://pkg.go.dev/badge/github.com/yi-nology/gitcode_api.svg)](https://pkg.go.dev/github.com/yi-nology/gitcode_api)

## 简介

`gitcode_api` 是一个轻量级、易用的 Go SDK，用于与 [GitCode](https://gitcode.com) / AtomGit 平台进行交互。参考了 [gitcode-cli](https://gitcode.com/gitcode-cli/cli) 的设计理念，提供类型安全的 API 调用。

## 功能特性

- 🔐 **认证管理** - 支持 Bearer Token、PRIVATE-TOKEN、access_token 多种认证方式
- 📦 **仓库管理** - 创建、更新、删除、Fork 仓库，文件操作
- 🐛 **Issue 管理** - 创建、更新、关闭 Issue，标签、里程碑管理
- 🔀 **Pull Request** - 创建、合并、审查 Pull Request
- 🌿 **分支管理** - 创建、删除、保护分支，提交比较
- 🔔 **Webhook** - 创建、更新、删除 Webhook，事件解析
- 👤 **用户与组织** - 获取用户信息，管理组织成员
- ⭐ **Star 管理** - Star/Unstar 仓库

## 安装

```bash
go get github.com/yi-nology/gitcode_api@v0.1.0
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"

    gitcode "github.com/yi-nology/gitcode_api"
)

func main() {
    // 创建客户端
    client := gitcode.NewClient("your-gitcode-token")
    ctx := context.Background()

    // 获取当前用户
    user, err := client.GetCurrentUser(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("当前用户: %s (%s)\n", user.Name, user.Login)

    // 列出仓库
    repos, err := client.ListRepositories(ctx, gitcode.ListRepositoriesOptions{})
    if err != nil {
        log.Fatal(err)
    }
    for _, repo := range repos {
        fmt.Printf("- %s: %s\n", repo.FullName, repo.Description)
    }
}
```

## 认证配置

### 获取 Token

访问 [GitCode 个人设置 → 访问令牌](https://gitcode.com/profile/personal_access_tokens) 创建 Token。

### 认证方式

```go
// 方式 1: Bearer Token (推荐)
client := gitcode.NewClient("your-token")

// 方式 2: PRIVATE-TOKEN Header
client := gitcode.NewClient("your-token")
client.SetAuthStyle(gitcode.AuthStylePrivateToken)

// 方式 3: access_token Query Parameter
client := gitcode.NewClient("your-token")
client.SetAuthStyle(gitcode.AuthStyleAccessToken)

// 方式 4: 自定义 baseURL (私有部署)
client := gitcode.NewClientWithBaseURL("https://your-gitcode.com/api/v5", "your-token")
```

## API 使用示例

### 仓库操作

```go
// 列出仓库
repos, _ := client.ListRepositories(ctx, gitcode.ListRepositoriesOptions{
    ListOptions: gitcode.ListOptions{Page: 1, PerPage: 10},
})

// 获取仓库详情
repo, _ := client.GetRepository(ctx, "owner", "repo")

// 创建仓库
repo, _ := client.CreateRepository(ctx, gitcode.CreateRepositoryOptions{
    Name:    "my-repo",
    Private: boolPtr(true),
})

// 更新仓库
repo, _ = client.UpdateRepository(ctx, "owner", "repo", gitcode.UpdateRepositoryOptions{
    Description: "新描述",
})

// Fork 仓库
fork, _ := client.ForkRepository(ctx, "owner", "repo", nil)

// 删除仓库
err := client.DeleteRepository(ctx, "owner", "repo")
```

### Issue 操作

```go
// 列出 Issue
issues, _ := client.ListIssues(ctx, "owner", "repo", gitcode.ListIssuesOptions{
    State: gitcode.IssueStateOpen,
})

// 创建 Issue
issue, _ := client.CreateIssue(ctx, "owner", "repo", gitcode.CreateIssueOptions{
    Title:  "Bug 报告",
    Body:   "问题描述...",
    Labels: []string{"bug", "urgent"},
})

// 关闭 Issue
issue, _ = client.CloseIssue(ctx, "owner", "repo", issue.Number)

// 添加评论
comment, _ := client.CreateIssueComment(ctx, "owner", "repo", issue.Number, "评论内容")
```

### Pull Request 操作

```go
// 列出 PR
prs, _ := client.ListPullRequests(ctx, "owner", "repo", gitcode.ListPullRequestsOptions{
    State: gitcode.PullRequestStateOpen,
})

// 创建 PR
pr, _ := client.CreatePullRequest(ctx, "owner", "repo", gitcode.CreatePullRequestOptions{
    Title: "新功能",
    Body:  "功能描述",
    Head:  "feature-branch",
    Base:  "main",
})

// 合并 PR
err := client.MergePullRequest(ctx, "owner", "repo", pr.Number, &gitcode.MergePullRequestOptions{
    CommitMessage: "合并新功能",
    Squash:        true,
})

// 获取 PR 文件列表
files, _ := client.ListPullRequestFiles(ctx, "owner", "repo", pr.Number)
```

### 分支操作

```go
// 列出分支
branches, _ := client.ListBranches(ctx, "owner", "repo")

// 创建分支
branch, _ := client.CreateBranch(ctx, "owner", "repo", gitcode.CreateBranchOptions{
    BranchName: "new-branch",
    Ref:         "main",
})

// 比较提交
cmp, _ := client.CompareCommits(ctx, "owner", "repo", "main", "feature-branch")
fmt.Printf("领先 %d 个提交, 落后 %d 个提交\n", cmp.AheadBy, cmp.BehindBy)

// 删除分支
err := client.DeleteBranch(ctx, "owner", "repo", "branch-name")
```

### Webhook 操作

```go
// 列出 Webhook
hooks, _ := client.ListWebhooks(ctx, "owner", "repo")

// 创建 Webhook
hook, _ := client.CreateWebhook(ctx, "owner", "repo", gitcode.CreateWebhookOptions{
    URL:    "https://your-webhook-url.com",
    Secret: "your-secret",
    Events: []string{"push", "pull_request"},
})

// 测试 Webhook
err := client.TestWebhook(ctx, "owner", "repo", hook.ID)

// 删除 Webhook
err = client.DeleteWebhook(ctx, "owner", "repo", hook.ID)
```

### 文件操作

```go
// 获取文件内容
content, _ := client.GetRepositoryContent(ctx, "owner", "repo", "README.md", "main")

// 创建文件
result, _ := client.CreateFile(ctx, "owner", "repo", "new-file.txt", gitcode.CreateFileOptions{
    Message: "添加新文件",
    Content: "文件内容",
    Branch:  "main",
})

// 更新文件
result, _ = client.UpdateFile(ctx, "owner", "repo", "new-file.txt", gitcode.UpdateFileOptions{
    Message: "更新文件",
    Content: "新内容",
    SHA:     result.Content.SHA,
    Branch:  "main",
})

// 删除文件
_, err := client.DeleteFile(ctx, "owner", "repo", "new-file.txt", gitcode.DeleteFileOptions{
    Message: "删除文件",
    SHA:     result.Content.SHA,
    Branch:  "main",
})
```

### 标签与发布

```go
// 列出标签
tags, _ := client.ListTags(ctx, "owner", "repo")

// 列出发布
releases, _ := client.ListReleases(ctx, "owner", "repo")

// 创建发布
release, _ := client.CreateRelease(ctx, "owner", "repo", gitcode.CreateReleaseOptions{
    TagName:    "v1.0.0",
    Title:      "Version 1.0.0",
    Body:       "发布说明",
    Prerelease: false,
})
```

## 项目结构

```
gitcode_api/
├── client.go      # 基础客户端、认证、用户 API
├── repos.go       # 仓库管理、文件操作、标签、发布
├── issues.go      # Issue 管理、标签、里程碑
├── pulls.go       # Pull Request 管理、审查、合并
├── branches.go    # 分支管理、保护分支、提交比较
├── webhooks.go    # Webhook 管理、事件解析
├── types.go       # 通用类型定义
└── examples/
    └── main.go    # 使用示例
```

## 与 gitcode-cli 的关系

本项目参考了 [gitcode-cli](https://gitcode.com/gitcode-cli/cli) 的设计理念，提供 Go 语言的 API 客户端实现：

| gitcode-cli (gc) | gitcode_api |
|------------------|-------------|
| `gc auth login` | `NewClient(token)` |
| `gc repo list` | `ListRepositories()` |
| `gc issue create` | `CreateIssue()` |
| `gc pr create` | `CreatePullRequest()` |

## 错误处理

```go
repos, err := client.ListRepositories(ctx, gitcode.ListRepositoriesOptions{})
if err != nil {
    // 检查是否为 API 错误
    if apiErr, ok := err.(*gitcode.Error); ok {
        fmt.Printf("API 错误: %s\n", apiErr.Message)
    }
}
```

## 相关项目

- [gitcode-cli](https://gitcode.com/gitcode-cli/cli) - GitCode 官方 CLI 工具
- [git-platform-sdk](https://github.com/yi-nology/git-platform-sdk) - 多平台 Git SDK 统一接口

## 许可证

MIT

## 贡献

欢迎提交 Issue 和 Pull Request！
