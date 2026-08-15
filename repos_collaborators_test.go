package gitcode

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListCollaborators(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Collaborator{
		{ID: 1, Login: "user1", Permission: "admin"},
		{ID: 2, Login: "user2", Permission: "write"},
	}))
	defer server.Close()

	collabs, err := client.ListCollaborators(context.Background(), "owner", "repo", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(collabs) != 2 {
		t.Errorf("expected 2 collaborators, got %d", len(collabs))
	}
	if collabs[0].Login != "user1" {
		t.Errorf("expected login 'user1', got '%s'", collabs[0].Login)
	}
}

func TestIsCollaborator(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v5/repos/owner/repo/collaborators/user1" {
			w.WriteHeader(204)
		} else {
			w.WriteHeader(404)
			w.Write([]byte(`{"message":"404 Not Found"}`))
		}
	}))
	defer server.Close()

	is, err := client.IsCollaborator(context.Background(), "owner", "repo", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !is {
		t.Error("expected user1 to be collaborator")
	}

	is, err = client.IsCollaborator(context.Background(), "owner", "repo", "nobody")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if is {
		t.Error("expected nobody to not be collaborator")
	}
}

func TestAddCollaborator(t *testing.T) {
	var receivedBody map[string]string
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		w.WriteHeader(204)
	}))
	defer server.Close()

	err := client.AddCollaborator(context.Background(), "owner", "repo", "user1", &AddCollaboratorOptions{Permission: "push"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedBody["permission"] != "push" {
		t.Errorf("expected permission 'push', got '%s'", receivedBody["permission"])
	}
}

func TestRemoveCollaborator(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.RemoveCollaborator(context.Background(), "owner", "repo", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetCollaboratorPermission(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &CollaboratorPermission{
		Permission: "admin",
		RoleName:   "admin",
	}))
	defer server.Close()

	perm, err := client.GetCollaboratorPermission(context.Background(), "owner", "repo", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if perm.Permission != "admin" {
		t.Errorf("expected permission 'admin', got '%s'", perm.Permission)
	}
}
