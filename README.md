# GitCode API Client

Go 语言的 [GitCode](https://gitcode.com) / AtomGit API 客户端库,提供对 GitCode 平台几乎所有 REST API (`/api/v5`) 的类型安全访问。

[![Go Reference](https://pkg.go.dev/badge/github.com/yi-nology/gitcode_api.svg)](https://pkg.go.dev/github.com/yi-nology/gitcode_api)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev)

## 功能特性

- **认证** — Bearer Token / PRIVATE-TOKEN Header / `access_token` Query 三种鉴权 + OAuth 2.0 授权码流程
- **仓库** — 创建/更新/删除/Fork/归档/转让,组织仓库,文件 CRUD,Tree/Blob,Raw,文件&图片上传
- **Issue** — 仓库 Issue CRUD、评论、标签、里程碑、操作日志、关联 PR;用户/组织/企业级 Issue 查询
- **Pull Request** — CRUD、合并、文件/提交/评论/审查、标签、审查人/测试人分配、操作日志、关联 Issue
- **分支** — 列表/创建/删除,保护分支规则,提交比较,提交历史
- **Webhook** — CRUD、测试推送,Push/TagPull/Issue/PullRequest/Note 事件解析
- **用户** — 当前用户/指定用户、SSH 公钥、邮箱、动态、Star 仓库、Namespace
- **组织 / 企业** — 组织信息/成员/关注者,企业成员、企业 Issue/PR、Issue 扩展状态、企业标签
- **搜索** — 仓库 / Issue / 用户
- **其他** — 仓库设置、Push 配置、PR 设置、模块开关、审查规则、下载统计、限流、通知
- **类型友好** — `FlexInt`/`FlexString`/`NullableTime` 等 JSON 容错类型,适配 GitCode 返回值类型漂移

## 环境要求

- Go ≥ 1.26(`go.mod` 已锁定)
- 一个 GitCode 个人访问令牌(PAT)或 OAuth 应用凭据

## 安装

```bash
go get github.com/yi-nology/gitcode_api@latest
```

当前最新版本为 `v0.2.0`,完整版本历史见 [releases](https://github.com/yi-nology/gitcode_api/releases)。

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
    client := gitcode.NewClient("your-gitcode-token")
    ctx := context.Background()

    user, err := client.GetCurrentUser(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("当前用户: %s (%s)\n", user.Name, user.Login)

    repos, err := client.ListRepositories(ctx, gitcode.ListRepositoriesOptions{
        ListOptions: gitcode.ListOptions{Page: 1, PerPage: 10},
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, r := range repos {
        fmt.Printf("- %s\n", r.FullName)
    }
}
```

更完整的示例参见 [`examples/main.go`](examples/main.go)。

## 客户端与认证

### 获取 Token

访问 [GitCode → 个人设置 → 私人令牌](https://gitcode.com/profile/personal_access_tokens) 创建 Token。

### 三种鉴权方式

`Client` 默认使用 `Authorization: Bearer <token>`,可通过 `SetAuthStyle` 切换:

```go
client := gitcode.NewClient("your-token")                       // Bearer (默认)
client.SetAuthStyle(gitcode.AuthStylePrivateToken)              // PRIVATE-TOKEN Header
client.SetAuthStyle(gitcode.AuthStyleAccessToken)               // access_token Query 参数
```

### 私有部署 / 自定义 HTTP 客户端

```go
client := gitcode.NewClientWithBaseURL("https://your-gitcode.com/api/v5", "token")
client.SetHTTPClient(&http.Client{Timeout: 60 * time.Second})   // 自定义超时
```

默认 BaseURL 为 `https://api.gitcode.com/api/v5`,默认超时 30s。

### OAuth 2.0

适用于需要用户授权而非固定 Token 的场景:

```go
oauth := gitcode.NewOAuthClient("client-id", "client-secret", "https://app.example.com/callback")

// 1. 引导用户跳转授权
url := oauth.AuthorizeURL("user_info projects", "random-state")
http.Redirect(w, r, url, http.StatusFound)

// 2. 用回调拿到的 code 换 token
token, err := oauth.ExchangeToken(ctx, code)

// 3. 用 token 创建 API 客户端
client := gitcode.NewClientFromOAuthToken(token)

// 4. 后续可用 refresh token 续期
token, err = oauth.RefreshToken(ctx, token.RefreshToken)
```

## API 示例

> 下文示例中 `ctx := context.Background()`,`owner`/`repo` 为仓库所有者和路径名。

### 仓库与文件

```go
// 创建仓库(私有)
private := false
repo, _ := client.CreateRepository(ctx, gitcode.CreateRepositoryOptions{
    Name:        "my-repo",
    Description: "描述",
    Private:     &private,
    AutoInit:    &private, // true,初始化 README
})

// 更新 / 删除 / Fork
client.UpdateRepository(ctx, "owner", "my-repo", gitcode.UpdateRepositoryOptions{Description: "新描述"})
client.DeleteRepository(ctx, "owner", "my-repo")
client.ForkRepository(ctx, "owner", "repo", nil)

// 组织下创建仓库
client.CreateOrgRepository(ctx, "my-org", gitcode.CreateOrgRepoOptions{Name: "team-repo"})

// 文件 CRUD(content 需 base64 编码)
content, _ := client.GetRepositoryContent(ctx, "owner", "repo", "README.md", "main")
res, _ := client.CreateFile(ctx, "owner", "repo", "a.txt", gitcode.CreateFileOptions{
    Message: "add a",
    Content: base64.StdEncoding.EncodeToString([]byte("hi")),
    Branch:  "main",
})
client.UpdateFile(ctx, "owner", "repo", "a.txt", gitcode.UpdateFileOptions{
    Message: "upd", Content: "...", SHA: res.Content.SHA, Branch: "main",
})

// 原始内容 / Tree / Blob / 贡献者 / 语言
raw, _      := client.GetRawFile(ctx, "owner", "repo", "README.md", "main")
tree, _     := client.GetTree(ctx, "owner", "repo", "main", true)
blob, _     := client.GetBlob(ctx, "owner", "repo", sha)
contribs, _ := client.ListContributors(ctx, "owner", "repo")
langs, _    := client.GetLanguages(ctx, "owner", "repo")

// 标签 / Release
tags, _     := client.ListTags(ctx, "owner", "repo")
release, _  := client.CreateRelease(ctx, "owner", "repo", gitcode.CreateReleaseOptions{
    TagName: "v1.0.0", Title: "v1.0.0", Body: "release notes",
})
client.DeleteRelease(ctx, "owner", "repo", release.TagName)
```

### Issue / 标签 / 里程碑

```go
// Issue CRUD
issue, _ := client.CreateIssue(ctx, "owner", "repo", gitcode.CreateIssueOptions{
    Title:  "Bug 报告",
    Body:   "问题描述",
    Labels: []string{"bug"},
})
client.UpdateIssue(ctx, "owner", "repo", int(issue.Number), gitcode.UpdateIssueOptions{State: gitcode.IssueStateClosed})
client.ReopenIssue(ctx, "owner", "repo", int(issue.Number))

// 评论
client.CreateIssueComment(ctx, "owner", "repo", int(issue.Number), "+1")
client.ListIssueComments(ctx, "owner", "repo", int(issue.Number))

// 仓库标签
client.CreateIssueLabel(ctx, "owner", "repo", "bug", "#ee0701")
client.UpdateIssueLabel(ctx, "owner", "repo", "bug", gitcode.UpdateLabelOptions{Color: "#d73a4a"})
client.AddIssueLabels(ctx, "owner", "repo", int(issue.Number), []string{"bug"})
client.ReplaceIssueLabels(ctx, "owner", "repo", int(issue.Number), []string{"bug", "p0"})

// 里程碑
ms, _ := client.CreateMilestone(ctx, "owner", "repo", "v2.0", "下个版本")
client.ListMilestones(ctx, "owner", "repo")
client.UpdateMilestone(ctx, "owner", "repo", int(ms.ID), gitcode.UpdateMilestoneOptions{State: "closed"})

// 操作日志 / 关联 PR
logs, _ := client.GetIssueOperateLogs(ctx, "owner", "repo", int(issue.Number))
prs, _  := client.GetIssueLinkedPRs(ctx, "owner", "repo", int(issue.Number))
```

### Pull Request

```go
pr, _ := client.CreatePullRequest(ctx, "owner", "repo", gitcode.CreatePullRequestOptions{
    Title: "feat: xxx", Head: "feature", Base: "main",
})

// 合并 (Squash)
client.MergePullRequest(ctx, "owner", "repo", pr.Number, &gitcode.MergePullRequestOptions{
    CommitMessage: "merge feat",
    Squash:        true,
})

// 文件 / 提交 / 评论 / 审查
client.ListPullRequestFiles(ctx, "owner", "repo", pr.Number)
client.ListPullRequestCommits(ctx, "owner", "repo", pr.Number)
client.CreatePullRequestReview(ctx, "owner", "repo", pr.Number, "LGTM", "APPROVE")

// 标签 / 审查人 / 测试人
client.AddPullRequestLabels(ctx, "owner", "repo", pr.Number, []string{"review"})
client.AssignPullRequestReviewers(ctx, "owner", "repo", pr.Number, "user1,user2")
client.AssignPullRequestTesters(ctx, "owner", "repo", pr.Number, "qa1")
client.HandlePullRequestReview(ctx, "owner", "repo", pr.Number, false) // force=false

// 关联 Issue / 操作日志
issues, _ := client.GetPullRequestLinkedIssues(ctx, "owner", "repo", pr.Number, gitcode.ListOptions{})
logs, _   := client.GetPullRequestOperateLogs(ctx, "owner", "repo", pr.Number)
```

### 分支与提交

```go
client.CreateBranch(ctx, "owner", "repo", gitcode.CreateBranchOptions{
    BranchName: "dev", Refs: "main",
})
client.DeleteBranch(ctx, "owner", "repo", "dev")

// 保护分支规则
client.CreateBranchProtection(ctx, "owner", "repo", gitcode.CreateBranchProtectionOptions{
    Name:                     "main",
    RequiredApprovingReviews: 2,
    AllowForcePushes:         false,
})
client.ListBranchProtections(ctx, "owner", "repo")
client.DeleteBranchProtection(ctx, "owner", "repo", "main")

// 提交
commits, _ := client.ListCommits(ctx, "owner", "repo", gitcode.ListCommitsOptions{Branch: "main"})
commit, _  := client.GetCommit(ctx, "owner", "repo", commits[0].SHA)
cmp, _     := client.CompareCommits(ctx, "owner", "repo", "main", "dev")
fmt.Printf("ahead=%d behind=%d\n", cmp.AheadBy, cmp.BehindBy)
```

### Webhook

```go
active := true
hook, _ := client.CreateWebhook(ctx, "owner", "repo", gitcode.CreateWebhookOptions{
    URL:    "https://example.com/hook",
    Secret: "s3cret",
    Events: []string{"push", "pull_request"},
    Active: &active,
})
client.TestWebhook(ctx, "owner", "repo", hook.ID)
client.DeleteWebhook(ctx, "owner", "repo", hook.ID)
```

接收并解析事件(以 Gin 为例):

```go
payload, _ := io.ReadAll(r.Body)
switch r.Header.Get("X-Gitcode-Event") {
case "push":
    e, _ := client.ParsePushEvent(payload)
    log.Printf("push %s -> %s", e.Before[:8], e.After[:8])
case "pull_request":
    e, _ := client.ParsePullRequestEvent(payload)
    log.Printf("PR #%d %s", e.Number, e.Action)
case "issues":
    e, _ := client.ParseIssueEvent(payload)
    log.Printf("issue #%d %s", int(e.Issue.Number), e.Action)
case "note":
    e, _ := client.ParseNoteEvent(payload)
    log.Printf("note on %s", e.NoteType)
case "tag_push":
    e, _ := client.ParseTagPushEvent(payload)
    log.Printf("tag %s", e.Ref)
}
```

> 还提供了 `PushEvent`/`PullRequestWebhookEvent`/`IssueWebhookEvent`/`NoteWebhookEvent`/`TagPushEvent` 五种事件类型,见 `webhooks.go`。

### 用户 / SSH 公钥 / 邮箱

```go
me, _     := client.GetCurrentUser(ctx)
user, _   := client.GetUser(ctx, "somebody")
emails, _ := client.ListEmails(ctx)

key, _ := client.CreateSSHKey(ctx, gitcode.CreateSSHKeyOptions{
    Title: "mbp", Key: "ssh-ed25519 AAAA...",
})
client.ListSSHKeys(ctx, gitcode.ListOptions{PerPage: 50})
client.DeleteSSHKey(ctx, key.ID)

events, _ := client.GetUserEvents(ctx, me.Login, "2025", "") // 年度动态
starred, _ := client.ListStarredRepositories(ctx, gitcode.ListStarredReposOptions{})
ns, _      := client.GetNamespace(ctx, "somepath")
```

### 组织 / 企业

```go
// 组织
client.ListUserOrganizations(ctx, "username", gitcode.ListOptions{})
org, _      := client.GetOrgInfo(ctx, "my-org")
client.UpdateOrganization(ctx, "my-org", gitcode.UpdateOrgOptions{Description: "..."})
client.InviteOrgMember(ctx, "my-org", "newbie", gitcode.InviteMemberOptions{Permission: "write"})
client.RemoveOrgMember(ctx, "my-org", "newbie")
client.ListOrgMembers(ctx, "my-org", "admin", gitcode.ListOptions{})

// 企业
client.ListEnterpriseMembers(ctx, "my-ent", "", gitcode.ListOptions{})
client.UpdateEnterpriseMember(ctx, "my-ent", "user", gitcode.UpdateEnterpriseMemberOptions{Role: "admin"})

// 企业 Issue / PR
client.ListEnterpriseIssues(ctx, "my-ent", gitcode.ListUserIssuesOptions{State: "open"})
client.ListEnterprisePullRequests(ctx, "my-ent", gitcode.ListEnterprisePRsOptions{State: "open"})
client.ListEnterpriseLabels(ctx, "my-ent")
client.GetOrgIssueExtendSettings(ctx, "my-org") // 自定义状态扩展
```

### 搜索

```go
repos, _ := client.SearchRepositories(ctx, gitcode.SearchRepositoriesOptions{
    Query: "gin web", Sort: "stars_count",
})
issues, _ := client.SearchIssues(ctx, gitcode.SearchIssuesOptions{
    Query: "memory leak", Repo: "owner/repo", State: "open",
})
users, _ := client.SearchUsers(ctx, gitcode.SearchUsersOptions{Query: "octocat"})
```

### Star / Watch / 其他

```go
client.StarRepository(ctx, "owner", "repo")
starred, _ := client.IsRepositoryStarred(ctx, "owner", "repo") // bool
client.UnstarRepository(ctx, "owner", "repo")

stargazers, _ := client.ListStargazers(ctx, "owner", "repo", gitcode.ListOptions{})
watchers, _   := client.ListWatchers(ctx, "owner", "repo", gitcode.ListOptions{})

// 仓库设置类
client.UpdateRepoSettings(ctx, "owner", "repo", &gitcode.RepoSettings{HasIssues: true})
client.UpdatePushConfig(ctx, "owner", "repo", &gitcode.PushConfig{MaxFileSize: 104857600})
client.UpdatePRSettings(ctx, "owner", "repo", &gitcode.PRSettings{DefaultMergeMethod: "merge"})
client.SetModuleSetting(ctx, "owner", "repo", gitcode.ModuleSetting{Wiki: false})
client.UpdateReviewerConfig(ctx, "owner", "repo", gitcode.ReviewerConfig{MinApprovingReviews: 1})

// 归档 / 转让
client.ArchiveRepository(ctx, "owner", "repo")
client.TransferRepository(ctx, "owner", "repo", gitcode.TransferRepoOptions{NewOwner: "new-owner"})

// 限流 / 通知
rl, _ := client.GetRateLimit(ctx)
nots, _ := client.ListNotifications(ctx)
```

## 通用类型

`types.go` 中定义了若干容错类型,用于处理 GitCode API 返回值类型不稳定的情况:

| 类型 | 用途 |
|---|---|
| `FlexInt` | 同一字段可能是 `123` 或 `"123"` |
| `FlexString` | 同一字段可能是 `"abc"` 或 `123` |
| `NullableTime` | 时间字段可能为空字符串 |
| `Timestamp` | 字符串格式时间(`RFC3339`) |
| `Error` | API 错误响应结构(见下) |

## 错误处理

`Client` 在 HTTP 状态码 ≥ 400 时返回标准 `error`(`fmt.Errorf`),错误信息包含方法、路径、状态码和响应体:

```go
_, err := client.GetRepository(ctx, "owner", "missing")
if err != nil {
    log.Println(err)
    // 输出形如: GitCode API GET /repos/owner/missing returned 404: {"message":"404 Not Found"}
}
```

如需进一步解析错误体,可使用 `types.Error`:

```go
var apiErr gitcode.Error
if json.Unmarshal([]byte(extractBody(err)), &apiErr) == nil {
    log.Printf("message=%s", apiErr.Message)
}
```

## 项目结构

```
gitcode_api/
├── client.go             # Client 构造、鉴权、HTTP 请求、User、ListOptions、RateLimit
├── oauth.go              # OAuthClient 授权码流程
├── repos.go              # 仓库 / 文件 / Tree / Blob / 标签 / Release / 设置 / 归档 / 转让
├── issues.go             # Issue CRUD / 评论 / 标签 / 里程碑
├── enterprise_issues.go  # 用户 / 组织 / 企业级 Issue、操作日志、关联 PR
├── pulls.go              # Pull Request 全套 + 审查 / 测试 / 标签 / 企业 PR
├── branches.go           # 分支 / 保护分支 / 提交 / 比较
├── webhooks.go           # Webhook CRUD + 5 类事件解析
├── search.go             # 搜索仓库 / Issue / 用户
├── users.go              # SSH 公钥 / 邮箱 / 动态 / Star / Namespace
├── orgs.go               # 组织 / 企业成员管理 / Issue 扩展配置
├── milestones.go         # 里程碑(带选项版本)
├── labels.go             # 仓库标签更新 / 替换 / 企业标签
├── types.go              # FlexInt/FlexString/NullableTime/Error 等通用类型 + Star/通知 API
├── gitcode_test.go       # 单元测试 + 真实 API 集成测试
└── examples/
    └── main.go           # 可运行的使用示例
```

## 测试

测试分为两类:纯函数测试(不需要网络)和 **真实 API 集成测试**(会在你的账户下创建/删除临时仓库)。

```bash
# 默认运行(需在 gitcode_test.go 顶部填入有效 Token)
GITCODE_TOKEN="your-token" go test -v ./...

# 跳过需要联网的集成测试
go test -short ./...
```

> 集成测试会真实创建仓库(`test-api-<timestamp>`),测试结束后自动清理。请使用测试账号的 Token。

## 与 gitcode-cli 的关系

本项目参考了 [gitcode-cli](https://gitcode.com/gitcode-cli/cli) 的设计理念,提供 Go 语言的 API 客户端实现:

| gitcode-cli (`gc`) | `gitcode_api` |
|---|---|
| `gc auth login` | `NewClient(token)` / `OAuthClient` |
| `gc repo list` | `ListRepositories()` |
| `gc issue create` | `CreateIssue()` |
| `gc pr create` | `CreatePullRequest()` |

## 相关项目

- [gitcode-cli](https://gitcode.com/gitcode-cli/cli) — GitCode 官方 CLI 工具
- [git-platform-sdk](https://github.com/yi-nology/git-platform-sdk) — 多平台 Git SDK 统一接口(GitHub/GitLab/Gitea/Forgejo/GitCode)

## 许可证

MIT。
