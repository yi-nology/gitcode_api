package gitcode

import (
	"context"
	"testing"
)

func TestInviteEnterpriseMember(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, nil))
	defer server.Close()

	err := client.InviteEnterpriseMember(context.Background(), "my-enterprise", "new-user", InviteMemberOptions{Permission: "read"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteEnterpriseMember(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteEnterpriseMember(context.Background(), "my-enterprise", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetOrgEnterprise(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &OrgEnterprise{
		ID: 1, Name: "My Enterprise", Path: "my-enterprise",
	}))
	defer server.Close()

	ent, err := client.GetOrgEnterprise(context.Background(), "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ent.Name != "My Enterprise" {
		t.Errorf("expected name 'My Enterprise', got '%s'", ent.Name)
	}
}

func TestListEnterpriseCustomizedRoles(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*EnterpriseCustomizedRole{
		{ID: 1, Name: "admin"}, {ID: 2, Name: "developer"},
	}))
	defer server.Close()

	roles, err := client.ListEnterpriseCustomizedRoles(context.Background(), "my-enterprise")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(roles))
	}
}

func TestEnterpriseMilestoneCRUD(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*EnterpriseMilestone{
		{ID: 1, Title: "v1.0", State: "active"},
	}))
	defer server.Close()

	milestones, err := client.ListEnterpriseMilestones(context.Background(), "my-enterprise", ListMilestonesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(milestones) != 1 {
		t.Errorf("expected 1 milestone, got %d", len(milestones))
	}
}

func TestListEnterpriseIssueCustomFields(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*IssueCustomField{
		{ID: 1, Name: "Priority", Type: "select"},
	}))
	defer server.Close()

	fields, err := client.ListEnterpriseIssueCustomFields(context.Background(), "my-enterprise")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 1 {
		t.Errorf("expected 1 field, got %d", len(fields))
	}
}

func TestCreateEnterpriseLabel(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &EnterpriseLabel{ID: 1, Name: "bug", Color: "#ff0000"}))
	defer server.Close()

	label, err := client.CreateEnterpriseLabel(context.Background(), "my-enterprise", CreateEnterpriseLabelOptions{
		Name: "bug", Color: "#ff0000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label.Name != "bug" {
		t.Errorf("expected name 'bug', got '%s'", label.Name)
	}
}

func TestUpdateEnterpriseLabel(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &EnterpriseLabel{ID: 1, Name: "feature", Color: "#00ff00"}))
	defer server.Close()

	label, err := client.UpdateEnterpriseLabel(context.Background(), "my-enterprise", 1, UpdateEnterpriseLabelOptions{
		Name: "feature", Color: "#00ff00",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label.Name != "feature" {
		t.Errorf("expected name 'feature', got '%s'", label.Name)
	}
}

func TestDeleteEnterpriseLabel(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteEnterpriseLabel(context.Background(), "my-enterprise", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
