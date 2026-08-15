package gitcode

import (
	"context"
	"testing"
)

func TestListOrgCustomizedRoles(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*OrgCustomizedRole{
		{ID: 1, Name: "developer", Permission: "write"},
		{ID: 2, Name: "viewer", Permission: "read"},
	}))
	defer server.Close()

	roles, err := client.ListOrgCustomizedRoles(context.Background(), "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(roles))
	}
}

func TestListOrgDiscussions(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*OrgDiscussion{
		{ID: 1, Number: 1, Title: "Discussion 1"},
	}))
	defer server.Close()

	discussions, err := client.ListOrgDiscussions(context.Background(), "my-org", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(discussions) != 1 {
		t.Errorf("expected 1 discussion, got %d", len(discussions))
	}
}

func TestGetOrgDiscussion(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &OrgDiscussion{
		ID: 1, Number: 1, Title: "Test Discussion",
	}))
	defer server.Close()

	disc, err := client.GetOrgDiscussion(context.Background(), "my-org", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if disc.Title != "Test Discussion" {
		t.Errorf("expected title 'Test Discussion', got '%s'", disc.Title)
	}
}

func TestListOrgDiscussionComments(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*OrgDiscussionComment{
		{ID: 1, Body: "Comment 1"},
	}))
	defer server.Close()

	comments, err := client.ListOrgDiscussionComments(context.Background(), "my-org", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(comments))
	}
}

func TestListOrgDiscussionCommentReplies(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*OrgDiscussionCommentReply{
		{ID: 1, Body: "Reply 1"},
	}))
	defer server.Close()

	replies, err := client.ListOrgDiscussionCommentReplies(context.Background(), "my-org", 1, 100, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(replies) != 1 {
		t.Errorf("expected 1 reply, got %d", len(replies))
	}
}
