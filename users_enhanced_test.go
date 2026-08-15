package gitcode

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListUserWatchedRepositories(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Repository{
		{FullName: "owner/repo1"}, {FullName: "owner/repo2"},
	}))
	defer server.Close()

	repos, err := client.ListUserWatchedRepositories(context.Background(), "username", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 2 {
		t.Errorf("expected 2 repos, got %d", len(repos))
	}
}

func TestListCurrentUserWatchedRepositories(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Repository{{FullName: "owner/repo1"}}))
	defer server.Close()

	repos, err := client.ListCurrentUserWatchedRepositories(context.Background(), ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("expected 1 repo, got %d", len(repos))
	}
}

func TestUpdateCurrentUser(t *testing.T) {
	var receivedBody UpdateCurrentUserOptions
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		json.NewEncoder(w).Encode(&User{Login: "test", Name: "New Name"})
	}))
	defer server.Close()

	user, err := client.UpdateCurrentUser(context.Background(), UpdateCurrentUserOptions{
		Name: "New Name",
		Bio:  "Go developer",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "New Name" {
		t.Errorf("expected name 'New Name', got '%s'", user.Name)
	}
}

func TestListUserPullRequests(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*PullRequest{
		{Number: 1, Title: "PR 1"},
	}))
	defer server.Close()

	prs, err := client.ListUserPullRequests(context.Background(), ListPullRequestsOptions{
		State: PullRequestStateOpen,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prs) != 1 {
		t.Errorf("expected 1 PR, got %d", len(prs))
	}
}
