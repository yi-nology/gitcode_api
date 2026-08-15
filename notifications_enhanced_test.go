package gitcode

import (
	"context"
	"net/http"
	"testing"
)

func TestListNotificationsWithOptions(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*NotificationThread{
		{ID: 1, Unread: true, Subject: &NotificationSubject{Title: "New issue"}},
	}))
	defer server.Close()

	notifs, err := client.ListNotificationsWithOptions(context.Background(), ListNotificationsOptions{
		Status: "unread",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifs) != 1 {
		t.Errorf("expected 1 notification, got %d", len(notifs))
	}
}

func TestGetNotificationThread(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &NotificationThread{
		ID: 1, Unread: true, Subject: &NotificationSubject{Title: "Test"},
	}))
	defer server.Close()

	thread, err := client.GetNotificationThread(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if thread.Subject.Title != "Test" {
		t.Errorf("expected title 'Test', got '%s'", thread.Subject.Title)
	}
}

func TestMarkNotificationThreadAsRead(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	err := client.MarkNotificationThreadAsRead(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarkNotificationsAsRead(t *testing.T) {
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	err := client.MarkNotificationsAsRead(context.Background(), MarkNotificationsOptions{All: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListRepoNotifications(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*NotificationThread{
		{ID: 1, Unread: true},
	}))
	defer server.Close()

	notifs, err := client.ListRepoNotifications(context.Background(), "owner", "repo", ListNotificationsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifs) != 1 {
		t.Errorf("expected 1 notification, got %d", len(notifs))
	}
}
