package gitcode

import (
	"context"
	"testing"
)

func TestListGitignoreTemplates(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []string{"Go", "Node", "Python"}))
	defer server.Close()

	templates, err := client.ListGitignoreTemplates(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(templates) != 3 {
		t.Errorf("expected 3 templates, got %d", len(templates))
	}
}

func TestGetGitignoreTemplate(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &GitignoreTemplate{Name: "Go", Source: "*.exe\n*.test"}))
	defer server.Close()

	template, err := client.GetGitignoreTemplate(context.Background(), "Go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template.Name != "Go" {
		t.Errorf("expected name 'Go', got '%s'", template.Name)
	}
}

func TestListLicenseTemplates(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*LicenseTemplate{
		{Key: "mit", Name: "MIT License"},
		{Key: "apache-2.0", Name: "Apache License 2.0"},
	}))
	defer server.Close()

	templates, err := client.ListLicenseTemplates(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(templates) != 2 {
		t.Errorf("expected 2 templates, got %d", len(templates))
	}
}

func TestGetLicenseTemplate(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &LicenseTemplate{Key: "mit", Name: "MIT License"}))
	defer server.Close()

	template, err := client.GetLicenseTemplate(context.Background(), "mit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template.Key != "mit" {
		t.Errorf("expected key 'mit', got '%s'", template.Key)
	}
}

func TestListLabelTemplates(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []string{"default", "bug", "feature"}))
	defer server.Close()

	templates, err := client.ListLabelTemplates(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(templates) != 3 {
		t.Errorf("expected 3 templates, got %d", len(templates))
	}
}

func TestGetLabelTemplate(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Label{
		{Name: "bug", Color: "#ff0000"},
		{Name: "feature", Color: "#00ff00"},
	}))
	defer server.Close()

	labels, err := client.GetLabelTemplate(context.Background(), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
}

func TestRenderMarkdown(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, "<strong>bold</strong>"))
	defer server.Close()

	html, err := client.RenderMarkdown(context.Background(), "**bold**", "gfm", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if html != "<strong>bold</strong>" {
		t.Errorf("expected '<strong>bold</strong>', got '%s'", html)
	}
}

func TestListCurrentUserRepositories(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Repository{
		{FullName: "owner/repo1"},
	}))
	defer server.Close()

	repos, err := client.ListCurrentUserRepositories(context.Background(), ListRepositoriesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("expected 1 repo, got %d", len(repos))
	}
}

func TestListUserRepositories(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Repository{
		{FullName: "user/repo1"}, {FullName: "user/repo2"},
	}))
	defer server.Close()

	repos, err := client.ListUserRepositories(context.Background(), "user", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 2 {
		t.Errorf("expected 2 repos, got %d", len(repos))
	}
}
