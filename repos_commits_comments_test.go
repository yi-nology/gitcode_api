package gitcode

import (
	"context"
	"testing"
)

func TestListCommitComments(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*CommitComment{
		{ID: 1, Body: "Nice code!", CommitID: "abc123"},
	}))
	defer server.Close()

	comments, err := client.ListCommitComments(context.Background(), "owner", "repo", "abc123", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(comments))
	}
}

func TestListRepoCommitComments(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*CommitComment{
		{ID: 1, Body: "Comment 1"}, {ID: 2, Body: "Comment 2"},
	}))
	defer server.Close()

	comments, err := client.ListRepoCommitComments(context.Background(), "owner", "repo", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}
}

func TestGetCommitComment(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &CommitComment{ID: 1, Body: "Nice!"}))
	defer server.Close()

	comment, err := client.GetCommitComment(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.Body != "Nice!" {
		t.Errorf("expected body 'Nice!', got '%s'", comment.Body)
	}
}

func TestCreateCommitComment(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &CommitComment{ID: 1, Body: "New comment"}))
	defer server.Close()

	comment, err := client.CreateCommitComment(context.Background(), "owner", "repo", "abc123", CreateCommitCommentOptions{
		Body: "New comment",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.Body != "New comment" {
		t.Errorf("expected body 'New comment', got '%s'", comment.Body)
	}
}

func TestUpdateCommitComment(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &CommitComment{ID: 1, Body: "Updated"}))
	defer server.Close()

	comment, err := client.UpdateCommitComment(context.Background(), "owner", "repo", 1, UpdateCommitCommentOptions{Body: "Updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.Body != "Updated" {
		t.Errorf("expected body 'Updated', got '%s'", comment.Body)
	}
}

func TestDeleteCommitComment(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteCommitComment(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
