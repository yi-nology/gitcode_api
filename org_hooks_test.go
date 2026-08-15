package gitcode

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListOrgWebhooks(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*OrgWebhook{
		{ID: 1, URL: "https://example.com/hook", Active: true},
	}))
	defer server.Close()

	hooks, err := client.ListOrgWebhooks(context.Background(), "my-org", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hooks) != 1 {
		t.Errorf("expected 1 hook, got %d", len(hooks))
	}
}

func TestGetOrgWebhook(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &OrgWebhook{ID: 1, URL: "https://example.com"}))
	defer server.Close()

	hook, err := client.GetOrgWebhook(context.Background(), "my-org", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hook.URL != "https://example.com" {
		t.Errorf("expected URL 'https://example.com', got '%s'", hook.URL)
	}
}

func TestCreateOrgWebhook(t *testing.T) {
	var receivedBody CreateOrgWebhookOptions
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		json.NewEncoder(w).Encode(&OrgWebhook{ID: 1, URL: receivedBody.URL})
	}))
	defer server.Close()

	active := true
	hook, err := client.CreateOrgWebhook(context.Background(), "my-org", CreateOrgWebhookOptions{
		URL:    "https://example.com/hook",
		Active: &active,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hook.URL != "https://example.com/hook" {
		t.Errorf("expected URL 'https://example.com/hook', got '%s'", hook.URL)
	}
}

func TestUpdateOrgWebhook(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &OrgWebhook{ID: 1, URL: "https://new-url.com"}))
	defer server.Close()

	hook, err := client.UpdateOrgWebhook(context.Background(), "my-org", 1, UpdateOrgWebhookOptions{
		URL: "https://new-url.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hook.URL != "https://new-url.com" {
		t.Errorf("expected URL 'https://new-url.com', got '%s'", hook.URL)
	}
}

func TestDeleteOrgWebhook(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteOrgWebhook(context.Background(), "my-org", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
