package gitcode

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListWikiPages(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*WikiPage{
		{Title: "Home"}, {Title: "Getting-Started"},
	}))
	defer server.Close()

	pages, err := client.ListWikiPages(context.Background(), "owner", "repo", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}
}

func TestGetWikiPage(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &WikiPage{Title: "Home", Content: "# Welcome"}))
	defer server.Close()

	page, err := client.GetWikiPage(context.Background(), "owner", "repo", "Home")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Title != "Home" {
		t.Errorf("expected title 'Home', got '%s'", page.Title)
	}
}

func TestCreateWikiPage(t *testing.T) {
	var receivedBody CreateWikiPageOptions
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&receivedBody)
		json.NewEncoder(w).Encode(&WikiPage{Title: "New-Page"})
	}))
	defer server.Close()

	page, err := client.CreateWikiPage(context.Background(), "owner", "repo", CreateWikiPageOptions{
		Title:         "New-Page",
		ContentBase64: "SGVsbG8=",
		Message:       "Create page",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Title != "New-Page" {
		t.Errorf("expected title 'New-Page', got '%s'", page.Title)
	}
}

func TestUpdateWikiPage(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &WikiPage{Title: "Updated-Page"}))
	defer server.Close()

	page, err := client.UpdateWikiPage(context.Background(), "owner", "repo", "Old-Page", UpdateWikiPageOptions{
		ContentBase64: "VXBkYXRlZA==",
		Message:       "Update page",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Title != "Updated-Page" {
		t.Errorf("expected title 'Updated-Page', got '%s'", page.Title)
	}
}

func TestDeleteWikiPage(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteWikiPage(context.Background(), "owner", "repo", "Old-Page")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
