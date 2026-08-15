package gitcode

import (
	"context"
	"testing"
)

func TestListIssueTemplates(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*IssueTemplate{
		{Name: "bug_report", Title: "Bug Report"},
		{Name: "feature_request", Title: "Feature Request"},
	}))
	defer server.Close()

	templates, err := client.ListIssueTemplates(context.Background(), "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(templates) != 2 {
		t.Errorf("expected 2 templates, got %d", len(templates))
	}
}

func TestGetIssueTemplate(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &IssueTemplate{
		Name: "bug_report", Title: "Bug Report", Body: "Describe the bug",
	}))
	defer server.Close()

	template, err := client.GetIssueTemplate(context.Background(), "owner", "repo", "bug_report")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template.Name != "bug_report" {
		t.Errorf("expected name 'bug_report', got '%s'", template.Name)
	}
}

func TestListPullRequestMergeTemplates(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*PRMergeTemplate{
		{Name: "default", Body: "Merge template"},
	}))
	defer server.Close()

	templates, err := client.ListPullRequestMergeTemplates(context.Background(), "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(templates) != 1 {
		t.Errorf("expected 1 template, got %d", len(templates))
	}
}
