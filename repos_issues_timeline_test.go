package gitcode

import (
	"context"
	"testing"
)

func TestListIssueTimelineEvents(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*IssueTimelineEvent{
		{ID: 1, Action: "opened", EventType: "opened"},
	}))
	defer server.Close()

	events, err := client.ListIssueTimelineEvents(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestListIssueSubscribers(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*IssueSubscriber{
		{ID: 1, Login: "subscriber1"},
	}))
	defer server.Close()

	subscribers, err := client.ListIssueSubscribers(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(subscribers) != 1 {
		t.Errorf("expected 1 subscriber, got %d", len(subscribers))
	}
}

func TestSubscribeToIssue(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.SubscribeToIssue(context.Background(), "owner", "repo", 1, "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnsubscribeFromIssue(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.UnsubscribeFromIssue(context.Background(), "owner", "repo", 1, "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListIssueDependencies(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*IssueDependency{
		{ID: 1, Issue: &Issue{Number: 42}},
	}))
	defer server.Close()

	deps, err := client.ListIssueDependencies(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(deps))
	}
}

func TestCreateIssueDependency(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &IssueDependency{ID: 1}))
	defer server.Close()

	dep, err := client.CreateIssueDependency(context.Background(), "owner", "repo", 1, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dep.ID != 1 {
		t.Errorf("expected ID 1, got %d", dep.ID)
	}
}

func TestDeleteIssueDependency(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteIssueDependency(context.Background(), "owner", "repo", 1, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListIssueBlockingIssues(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*Issue{
		{Number: 10, Title: "Blocked issue"},
	}))
	defer server.Close()

	issues, err := client.ListIssueBlockingIssues(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(issues))
	}
}
