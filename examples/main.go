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
	fmt.Printf("当前用户: %s (%s)\n", user.Name, user.Login)

	// 列出仓库
	repos, err := client.ListRepositories(ctx, gitcode.ListRepositoriesOptions{
		ListOptions: gitcode.ListOptions{Page: 1, PerPage: 10},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n仓库列表:\n")
	for _, repo := range repos {
		fmt.Printf("- %s (%s)\n", repo.FullName, repo.Description)
	}

	// 获取单个仓库
	repo, err := client.GetRepository(ctx, "owner", "repo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n仓库详情: %s\n", repo.FullName)
	fmt.Printf("  描述: %s\n", repo.Description)
	fmt.Printf("  Stars: %d\n", repo.StarsCount)

	// 列出分支
	branches, err := client.ListBranches(ctx, "owner", "repo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n分支列表:\n")
	for _, branch := range branches {
		fmt.Printf("- %s\n", branch.Name)
	}

	// 列出 Issue
	issues, err := client.ListIssues(ctx, "owner", "repo", gitcode.ListIssuesOptions{
		State: gitcode.IssueStateOpen,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nIssue 列表:\n")
	for _, issue := range issues {
		fmt.Printf("#%d: %s\n", int(issue.Number), issue.Title)
	}

	// 列出 Pull Request
	prs, err := client.ListPullRequests(ctx, "owner", "repo", gitcode.ListPullRequestsOptions{
		State: gitcode.PullRequestStateOpen,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nPull Request 列表:\n")
	for _, pr := range prs {
		fmt.Printf("#%d: %s\n", pr.Number, pr.Title)
	}

	// 列出 Webhook
	hooks, err := client.ListWebhooks(ctx, "owner", "repo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nWebhook 列表:\n")
	for _, hook := range hooks {
		fmt.Printf("- %s (active: %v)\n", hook.URL, hook.Active)
	}
}
