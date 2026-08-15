package gitcode

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListProtectedTags(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*ProtectedTag{
		{ID: 1, NamePattern: "v*"},
	}))
	defer server.Close()

	tags, err := client.ListProtectedTags(context.Background(), "owner", "repo", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}
}

func TestGetProtectedTag(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &ProtectedTag{ID: 1, NamePattern: "v*"}))
	defer server.Close()

	tag, err := client.GetProtectedTag(context.Background(), "owner", "repo", "v*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.NamePattern != "v*" {
		t.Errorf("expected pattern 'v*', got '%s'", tag.NamePattern)
	}
}

func TestCreateProtectedTag(t *testing.T) {
	var receivedBody CreateProtectedTagOptions
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&receivedBody)
		json.NewEncoder(w).Encode(&ProtectedTag{ID: 1, NamePattern: "release-*"})
	}))
	defer server.Close()

	tag, err := client.CreateProtectedTag(context.Background(), "owner", "repo", CreateProtectedTagOptions{
		NamePattern: "release-*",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.NamePattern != "release-*" {
		t.Errorf("expected pattern 'release-*', got '%s'", tag.NamePattern)
	}
}

func TestUpdateProtectedTag(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &ProtectedTag{ID: 1, NamePattern: "v2*"}))
	defer server.Close()

	tag, err := client.UpdateProtectedTag(context.Background(), "owner", "repo", "v*", UpdateProtectedTagOptions{
		NamePattern: "v2*",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.NamePattern != "v2*" {
		t.Errorf("expected pattern 'v2*', got '%s'", tag.NamePattern)
	}
}

func TestDeleteProtectedTag(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteProtectedTag(context.Background(), "owner", "repo", "v*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
