package gitcode

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListRepositoryTopics(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, map[string][]string{"topics": {"golang", "api", "sdk"}}))
	defer server.Close()

	topics, err := client.ListRepositoryTopics(context.Background(), "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(topics) != 3 {
		t.Errorf("expected 3 topics, got %d", len(topics))
	}
}

func TestUpdateRepositoryTopics(t *testing.T) {
	var receivedBody map[string]interface{}
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		json.NewEncoder(w).Encode(map[string][]string{"topics": {"new-topic"}})
	}))
	defer server.Close()

	topics, err := client.UpdateRepositoryTopics(context.Background(), "owner", "repo", []string{"new-topic"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(topics) != 1 {
		t.Errorf("expected 1 topic, got %d", len(topics))
	}
}

func TestAddRepositoryTopic(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.AddRepositoryTopic(context.Background(), "owner", "repo", "new-topic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteRepositoryTopic(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteRepositoryTopic(context.Background(), "owner", "repo", "old-topic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
