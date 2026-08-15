package gitcode

import (
	"context"
	"net/http"
	"testing"
)

func TestLinkPullRequestIssue(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, nil))
	defer server.Close()

	err := client.LinkPullRequestIssue(context.Background(), "owner", "repo", 1, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnlinkPullRequestIssue(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.UnlinkPullRequestIssue(context.Background(), "owner", "repo", 1, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnassignPullRequestTesters(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.UnassignPullRequestTesters(context.Background(), "owner", "repo", 1, "qa1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListPullRequestAvailableTesters(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*User{{Login: "qa1"}}))
	defer server.Close()

	users, err := client.ListPullRequestAvailableTesters(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestAssignPullRequestApprovalReviewers(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(201)
	}))
	defer server.Close()

	err := client.AssignPullRequestApprovalReviewers(context.Background(), "owner", "repo", 1, "reviewer1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnassignPullRequestApprovalReviewers(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.UnassignPullRequestApprovalReviewers(context.Background(), "owner", "repo", 1, "reviewer1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListPullRequestAvailableReviewers(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*User{{Login: "reviewer1"}}))
	defer server.Close()

	users, err := client.ListPullRequestAvailableReviewers(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestResolvePullRequestDiscussion(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, nil))
	defer server.Close()

	err := client.ResolvePullRequestDiscussion(context.Background(), "owner", "repo", 1, "disc-123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListPullRequestModifyHistory(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*ModifyHistoryEntry{
		{ID: 1, Action: "update", Field: "title"},
	}))
	defer server.Close()

	history, err := client.ListPullRequestModifyHistory(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(history) != 1 {
		t.Errorf("expected 1 history entry, got %d", len(history))
	}
}

func TestListPullRequestUserReactions(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*IssueUserReaction{
		{ID: 1, Content: "heart"},
	}))
	defer server.Close()

	reactions, err := client.ListPullRequestUserReactions(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reactions) != 1 {
		t.Errorf("expected 1 reaction, got %d", len(reactions))
	}
}

func TestRefreshPullRequestCommentPosition(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, nil))
	defer server.Close()

	err := client.RefreshPullRequestCommentPosition(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListPullRequestFilesJSON(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*PullRequestFileChange{
		{Filename: "main.go", Status: "modified", Additions: 10, Deletions: 5},
	}))
	defer server.Close()

	files, err := client.ListPullRequestFilesJSON(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
}
