package gitcode

import (
	"context"
	"net/http"
	"testing"
)

func TestListRepoReviewers(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Reviewer{
		{ID: 1, Login: "reviewer1"},
	}))
	defer server.Close()

	reviewers, err := client.ListRepoReviewers(context.Background(), "owner", "repo", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reviewers) != 1 {
		t.Errorf("expected 1 reviewer, got %d", len(reviewers))
	}
}

func TestListPullRequestReviewers(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Reviewer{
		{ID: 1, Login: "reviewer1"},
	}))
	defer server.Close()

	reviewers, err := client.ListPullRequestReviewers(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reviewers) != 1 {
		t.Errorf("expected 1 reviewer, got %d", len(reviewers))
	}
}

func TestRequestPullRequestReviewers(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(201)
	}))
	defer server.Close()

	err := client.RequestPullRequestReviewers(context.Background(), "owner", "repo", 1, PullRequestReviewRequest{
		Reviewers: []string{"reviewer1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDismissPullRequestReview(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, nil))
	defer server.Close()

	err := client.DismissPullRequestReview(context.Background(), "owner", "repo", 1, 100, "Dismissed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetPullRequestReview(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &PullRequestReview{
		ID: 1, State: "approved", Body: "LGTM",
	}))
	defer server.Close()

	review, err := client.GetPullRequestReview(context.Background(), "owner", "repo", 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if review.State != "approved" {
		t.Errorf("expected state 'approved', got '%s'", review.State)
	}
}

func TestGetPullRequestDiff(t *testing.T) {
	diffContent := []byte("diff --git a/file.go b/file.go")
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(diffContent)
	}))
	defer server.Close()

	diff, err := client.GetPullRequestDiff(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(diff) != string(diffContent) {
		t.Error("diff content mismatch")
	}
}

func TestGetPullRequestPatch(t *testing.T) {
	patchContent := []byte("From abc Mon Sep 17 00:00:00 2001")
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(patchContent)
	}))
	defer server.Close()

	patch, err := client.GetPullRequestPatch(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(patch) != string(patchContent) {
		t.Error("patch content mismatch")
	}
}
