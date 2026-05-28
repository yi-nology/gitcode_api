# GitCode API Client

Go 语言的 GitCode API 客户端库，提供对 GitCode 平台的完整访问。

## 功能特性

- **仓库管理**: 创建、更新、删除、Fork 仓库
- **Issue 管理**: 创建、更新、关闭 Issue，管理标签和里程碑
- **Pull Request**: 创建、合并、审查 Pull Request
- **分支管理**: 创建、删除、保护分支
- **Webhook**: 创建、更新、删除 Webhook
- **用户和组织**: 获取用户信息，管理组织成员

## 安装

```bash
go get github.com/yi-nology/gitcode_api
```

## 使用示例

```go
package main

import (
    "context"
    "fmt"
    "log"

    gitcode "github.com/yi-nology/gitcode_api"
)

func main() {
    client := gitcode.NewClient("your-gitcode-token")
    ctx := context.Background()

    // 获取当前用户
    user, err := client.GetCurrentUser(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("当前用户: %s\n", user.Login)

    // 列出仓库
    repos, err := client.ListRepositories(ctx, gitcode.ListRepositoriesOptions{})
    if err != nil {
        log.Fatal(err)
    }
    for _, repo := range repos {
        fmt.Printf("- %s\n", repo.FullName)
    }
}
```

## API 参考

### 客户端

```go
// 创建客户端
client := gitcode.NewClient("token")

// 自定义 baseURL
client := gitcode.NewClientWithBaseURL("https://your-gitcode-instance.com/api/v5", "token")

// 设置认证方式
client.SetAuthStyle(gitcode.AuthStyleBearer)       // Bearer token (默认)
client.SetAuthStyle(gitcode.AuthStylePrivateToken) // PRIVATE-TOKEN header
client.SetAuthStyle(gitcode.AuthStyleAccessToken)  // access_token query parameter
```

### 仓库

```go
// 列出仓库
repos, err := client.ListRepositories(ctx, gitcode.ListRepositoriesOptions{})

// 获取仓库
repo, err := client.GetRepository(ctx, "owner", "repo")

// 创建仓库
repo, err := client.CreateRepository(ctx, gitcode.CreateRepositoryOptions{
    Name: "my-repo",
    Private: boolPtr(true),
})

// 更新仓库
repo, err := client.UpdateRepository(ctx, "owner", "repo", gitcode.UpdateRepositoryOptions{
    Description: "New description",
})

// 删除仓库
err := client.DeleteRepository(ctx, "owner", "repo")
```

### Issue

```go
// 列出 Issue
issues, err := client.ListIssues(ctx, "owner", "repo", gitcode.ListIssuesOptions{
    State: gitcode.IssueStateOpen,
})

// 创建 Issue
issue, err := client.CreateIssue(ctx, "owner", "repo", gitcode.CreateIssueOptions{
    Title: "Bug report",
    Body:  "Description of the bug",
    Labels: []string{"bug", "urgent"},
})

// 关闭 Issue
issue, err := client.CloseIssue(ctx, "owner", "repo", 123)
```

### Pull Request

```go
// 列出 Pull Request
prs, err := client.ListPullRequests(ctx, "owner", "repo", gitcode.ListPullRequestsOptions{
    State: gitcode.PullRequestStateOpen,
})

// 创建 Pull Request
pr, err := client.CreatePullRequest(ctx, "owner", "repo", gitcode.CreatePullRequestOptions{
    Title: "New feature",
    Body:  "Description of the feature",
    Head:  "feature-branch",
    Base:  "main",
})

// 合并 Pull Request
err := client.MergePullRequest(ctx, "owner", "repo", 123, nil)
```

### 分支

```go
// 列出分支
branches, err := client.ListBranches(ctx, "owner", "repo")

// 创建分支
branch, err := client.CreateBranch(ctx, "owner", "repo", gitcode.CreateBranchOptions{
    BranchName: "new-branch",
    Ref:         "main",
})

// 删除分支
err := client.DeleteBranch(ctx, "owner", "repo", "branch-name")
```

### Webhook

```go
// 列出 Webhook
hooks, err := client.ListWebhooks(ctx, "owner", "repo")

// 创建 Webhook
hook, err := client.CreateWebhook(ctx, "owner", "repo", gitcode.CreateWebhookOptions{
    URL:    "https://your-webhook-url.com",
    Events: []string{"push", "pull_request"},
})

// 删除 Webhook
err := client.DeleteWebhook(ctx, "owner", "repo", hookID)
```

## 认证

GitCode API 支持多种认证方式：

1. **Bearer Token** (推荐)
2. **PRIVATE-TOKEN Header**
3. **access_token Query Parameter**

Token 获取：[GitCode 个人设置 → 访问令牌](https://gitcode.com/profile/personal_access_tokens)

## 错误处理

所有 API 调用都会返回错误。错误信息包含 HTTP 状态码和响应体：

```go
repos, err := client.ListRepositories(ctx, gitcode.ListRepositoriesOptions{})
if err != nil {
    log.Printf("API 错误: %v", err)
}
```

## 许可证

MIT
