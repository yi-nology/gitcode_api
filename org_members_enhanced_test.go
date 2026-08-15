package gitcode

import (
	"context"
	"net/http"
	"testing"
)

func TestListOrgPublicMembers(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*User{{Login: "public-user"}}))
	defer server.Close()

	users, err := client.ListOrgPublicMembers(context.Background(), "my-org", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestIsOrgPublicMember(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v5/orgs/my-org/public_members/user1" {
			w.WriteHeader(204)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()

	is, err := client.IsOrgPublicMember(context.Background(), "my-org", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !is {
		t.Error("expected user1 to be public member")
	}
}

func TestPublicizeOrgMembership(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	err := client.PublicizeOrgMembership(context.Background(), "my-org", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConcealOrgMembership(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.ConcealOrgMembership(context.Background(), "my-org", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListOrgBlockedUsers(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*User{{Login: "blocked-user"}}))
	defer server.Close()

	users, err := client.ListOrgBlockedUsers(context.Background(), "my-org", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestIsOrgBlockedUser(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v5/orgs/my-org/blocks/spammer" {
			w.WriteHeader(204)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()

	is, err := client.IsOrgBlockedUser(context.Background(), "my-org", "spammer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !is {
		t.Error("expected spammer to be blocked")
	}
}

func TestBlockOrgUser(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	err := client.BlockOrgUser(context.Background(), "my-org", "spammer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnblockOrgUser(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.UnblockOrgUser(context.Background(), "my-org", "spammer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateOrganization(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &Organization{ID: 1, Login: "new-org"}))
	defer server.Close()

	org, err := client.CreateOrganization(context.Background(), CreateOrgOptions{
		Username: "new-org", Name: "New Organization",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if org.Login != "new-org" {
		t.Errorf("expected login 'new-org', got '%s'", org.Login)
	}
}

func TestDeleteOrganization(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteOrganization(context.Background(), "old-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
