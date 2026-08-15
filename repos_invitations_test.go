package gitcode

import (
	"context"
	"net/http"
	"testing"
)

func TestAcceptRepoInvitation(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	err := client.AcceptRepoInvitation(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeclineRepoInvitation(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeclineRepoInvitation(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListPendingRepoInvitations(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*RepoInvitation{
		{ID: 1, Permission: "write"},
	}))
	defer server.Close()

	invitations, err := client.ListPendingRepoInvitations(context.Background(), ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(invitations) != 1 {
		t.Errorf("expected 1 invitation, got %d", len(invitations))
	}
}
