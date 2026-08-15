package gitcode

import (
	"context"
	"net/http"
	"testing"
)

func TestListOrgTeams(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Team{
		{ID: 1, Name: "backend", Permission: "write"},
	}))
	defer server.Close()

	teams, err := client.ListOrgTeams(context.Background(), "my-org", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(teams) != 1 {
		t.Errorf("expected 1 team, got %d", len(teams))
	}
}

func TestCreateTeam(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &Team{ID: 1, Name: "backend"}))
	defer server.Close()

	team, err := client.CreateTeam(context.Background(), "my-org", CreateTeamOptions{Name: "backend", Permission: "write"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.Name != "backend" {
		t.Errorf("expected name 'backend', got '%s'", team.Name)
	}
}

func TestGetTeam(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &Team{ID: 1, Name: "backend"}))
	defer server.Close()

	team, err := client.GetTeam(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.Name != "backend" {
		t.Errorf("expected name 'backend', got '%s'", team.Name)
	}
}

func TestUpdateTeam(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &Team{ID: 1, Name: "backend-team"}))
	defer server.Close()

	team, err := client.UpdateTeam(context.Background(), 1, UpdateTeamOptions{Name: "backend-team"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.Name != "backend-team" {
		t.Errorf("expected name 'backend-team', got '%s'", team.Name)
	}
}

func TestDeleteTeam(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteTeam(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListTeamMembers(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*TeamMember{
		{ID: 1, Login: "member1"},
	}))
	defer server.Close()

	members, err := client.ListTeamMembers(context.Background(), 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 1 {
		t.Errorf("expected 1 member, got %d", len(members))
	}
}

func TestAddTeamMember(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	err := client.AddTeamMember(context.Background(), 1, "new-member")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoveTeamMember(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.RemoveTeamMember(context.Background(), 1, "old-member")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
