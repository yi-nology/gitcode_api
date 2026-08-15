package gitcode

import (
	"context"
	"testing"
)

func TestListGitReferences(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*GitReference{
		{Ref: "refs/heads/main", Object: &GitObject{SHA: "abc123", Type: "commit"}},
		{Ref: "refs/heads/dev", Object: &GitObject{SHA: "def456", Type: "commit"}},
	}))
	defer server.Close()

	refs, err := client.ListGitReferences(context.Background(), "owner", "repo", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(refs) != 2 {
		t.Errorf("expected 2 refs, got %d", len(refs))
	}
}

func TestGetGitReference(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &GitReference{
		Ref:    "refs/heads/main",
		Object: &GitObject{SHA: "abc123", Type: "commit"},
	}))
	defer server.Close()

	ref, err := client.GetGitReference(context.Background(), "owner", "repo", "heads/main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.Ref != "refs/heads/main" {
		t.Errorf("expected ref 'refs/heads/main', got '%s'", ref.Ref)
	}
}

func TestCreateGitReference(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &GitReference{
		Ref:    "refs/heads/new-branch",
		Object: &GitObject{SHA: "abc123"},
	}))
	defer server.Close()

	ref, err := client.CreateGitReference(context.Background(), "owner", "repo", CreateReferenceOptions{
		Ref: "refs/heads/new-branch",
		SHA: "abc123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.Ref != "refs/heads/new-branch" {
		t.Errorf("expected ref 'refs/heads/new-branch', got '%s'", ref.Ref)
	}
}

func TestUpdateGitReference(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &GitReference{
		Ref:    "refs/heads/main",
		Object: &GitObject{SHA: "new-sha"},
	}))
	defer server.Close()

	ref, err := client.UpdateGitReference(context.Background(), "owner", "repo", "heads/main", UpdateReferenceOptions{
		SHA: "new-sha",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.Object.SHA != "new-sha" {
		t.Errorf("expected SHA 'new-sha', got '%s'", ref.Object.SHA)
	}
}

func TestDeleteGitReference(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteGitReference(context.Background(), "owner", "repo", "heads/old-branch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
