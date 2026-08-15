package gitcode

import (
	"context"
	"net/http"
	"testing"
)

func TestListRepoAssignees(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*User{
		{Login: "user1"}, {Login: "user2"},
	}))
	defer server.Close()

	users, err := client.ListRepoAssignees(context.Background(), "owner", "repo", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestIsRepoAssignee(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v5/repos/owner/repo/assignees/user1" {
			w.WriteHeader(204)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()

	is, err := client.IsRepoAssignee(context.Background(), "owner", "repo", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !is {
		t.Error("expected user1 to be assignee")
	}

	is, err = client.IsRepoAssignee(context.Background(), "owner", "repo", "nobody")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if is {
		t.Error("expected nobody to not be assignee")
	}
}

func TestAddIssueAssignees(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &Issue{Number: 1}))
	defer server.Close()

	issue, err := client.AddIssueAssignees(context.Background(), "owner", "repo", 1, []string{"user1", "user2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if int(issue.Number) != 1 {
		t.Errorf("expected issue number 1, got %d", issue.Number)
	}
}

func TestRemoveIssueAssignees(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &Issue{Number: 1}))
	defer server.Close()

	issue, err := client.RemoveIssueAssignees(context.Background(), "owner", "repo", 1, []string{"user1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if int(issue.Number) != 1 {
		t.Errorf("expected issue number 1, got %d", issue.Number)
	}
}
