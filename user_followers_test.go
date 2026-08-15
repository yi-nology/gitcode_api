package gitcode

import (
	"context"
	"net/http"
	"testing"
)

func TestListUserFollowers(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*User{
		{Login: "follower1"}, {Login: "follower2"},
	}))
	defer server.Close()

	users, err := client.ListUserFollowers(context.Background(), "username", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 followers, got %d", len(users))
	}
}

func TestListUserFollowing(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*User{{Login: "following1"}}))
	defer server.Close()

	users, err := client.ListUserFollowing(context.Background(), "username", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 following, got %d", len(users))
	}
}

func TestFollowUser(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	err := client.FollowUser(context.Background(), "target-user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnfollowUser(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	err := client.UnfollowUser(context.Background(), "target-user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIsFollowing(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v5/user/following/target" {
			w.WriteHeader(204)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()

	is, err := client.IsFollowing(context.Background(), "target")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !is {
		t.Error("expected to be following")
	}

	is, err = client.IsFollowing(context.Background(), "nobody")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if is {
		t.Error("expected not to be following")
	}
}
