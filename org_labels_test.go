package gitcode

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListOrgLabels(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*OrgLabel{
		{ID: 1, Name: "bug", Color: "#ff0000"},
	}))
	defer server.Close()

	labels, err := client.ListOrgLabels(context.Background(), "my-org", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 1 {
		t.Errorf("expected 1 label, got %d", len(labels))
	}
}

func TestGetOrgLabel(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &OrgLabel{ID: 1, Name: "bug"}))
	defer server.Close()

	label, err := client.GetOrgLabel(context.Background(), "my-org", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label.Name != "bug" {
		t.Errorf("expected name 'bug', got '%s'", label.Name)
	}
}

func TestCreateOrgLabel(t *testing.T) {
	var receivedBody CreateOrgLabelOptions
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		json.NewEncoder(w).Encode(&OrgLabel{ID: 1, Name: receivedBody.Name, Color: receivedBody.Color})
	}))
	defer server.Close()

	label, err := client.CreateOrgLabel(context.Background(), "my-org", CreateOrgLabelOptions{
		Name: "feature", Color: "#00ff00",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label.Name != "feature" {
		t.Errorf("expected name 'feature', got '%s'", label.Name)
	}
}

func TestUpdateOrgLabel(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &OrgLabel{ID: 1, Name: "updated"}))
	defer server.Close()

	label, err := client.UpdateOrgLabel(context.Background(), "my-org", 1, UpdateOrgLabelOptions{Name: "updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label.Name != "updated" {
		t.Errorf("expected name 'updated', got '%s'", label.Name)
	}
}

func TestDeleteOrgLabel(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteOrgLabel(context.Background(), "my-org", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
