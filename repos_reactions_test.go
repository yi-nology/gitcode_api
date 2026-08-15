package gitcode

import (
	"context"
	"testing"
)

func TestListIssueReactions(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Reaction{
		{ID: 1, Content: "+1"}, {ID: 2, Content: "heart"},
	}))
	defer server.Close()

	reactions, err := client.ListIssueReactions(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reactions) != 2 {
		t.Errorf("expected 2 reactions, got %d", len(reactions))
	}
}

func TestCreateIssueReaction(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &Reaction{ID: 1, Content: "heart"}))
	defer server.Close()

	reaction, err := client.CreateIssueReaction(context.Background(), "owner", "repo", 1, ReactionHeart)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reaction.Content != "heart" {
		t.Errorf("expected content 'heart', got '%s'", reaction.Content)
	}
}

func TestDeleteIssueReaction(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteIssueReaction(context.Background(), "owner", "repo", 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListIssueCommentReactions(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Reaction{{ID: 1, Content: "+1"}}))
	defer server.Close()

	reactions, err := client.ListIssueCommentReactions(context.Background(), "owner", "repo", 100, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reactions) != 1 {
		t.Errorf("expected 1 reaction, got %d", len(reactions))
	}
}

func TestCreatePullRequestCommentReaction(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &Reaction{ID: 1, Content: "rocket"}))
	defer server.Close()

	reaction, err := client.CreatePullRequestCommentReaction(context.Background(), "owner", "repo", 100, ReactionRocket)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reaction.Content != "rocket" {
		t.Errorf("expected content 'rocket', got '%s'", reaction.Content)
	}
}
