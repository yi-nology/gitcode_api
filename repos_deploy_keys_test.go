package gitcode

import (
	"context"
	"testing"
)

func TestListDeployKeys(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*DeployKey{
		{ID: 1, Title: "deploy-key", ReadOnly: true},
	}))
	defer server.Close()

	keys, err := client.ListDeployKeys(context.Background(), "owner", "repo", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}
}

func TestGetDeployKey(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &DeployKey{ID: 1, Title: "my-key"}))
	defer server.Close()

	key, err := client.GetDeployKey(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key.Title != "my-key" {
		t.Errorf("expected title 'my-key', got '%s'", key.Title)
	}
}

func TestCreateDeployKey(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &DeployKey{ID: 1, Title: "new-key"}))
	defer server.Close()

	readOnly := true
	key, err := client.CreateDeployKey(context.Background(), "owner", "repo", CreateDeployKeyOptions{
		Title:    "new-key",
		Key:      "ssh-ed25519 AAAA...",
		ReadOnly: &readOnly,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key.ID != 1 {
		t.Errorf("expected ID 1, got %d", key.ID)
	}
}

func TestDeleteDeployKey(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteDeployKey(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
