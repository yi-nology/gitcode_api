package gitcode

import (
	"context"
	"testing"
)

func TestListKanbanBoards(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*KanbanBoard{
		{ID: 1, Name: "Sprint Board"},
	}))
	defer server.Close()

	boards, err := client.ListKanbanBoards(context.Background(), "my-org", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(boards) != 1 {
		t.Errorf("expected 1 board, got %d", len(boards))
	}
}

func TestGetKanbanBoard(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &KanbanBoard{ID: 1, Name: "Sprint Board"}))
	defer server.Close()

	board, err := client.GetKanbanBoard(context.Background(), "my-org", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if board.Name != "Sprint Board" {
		t.Errorf("expected name 'Sprint Board', got '%s'", board.Name)
	}
}

func TestAddToKanbanBoard(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &KanbanBoardItem{ID: 1, IssueNumber: 42}))
	defer server.Close()

	item, err := client.AddToKanbanBoard(context.Background(), "my-org", 1, AddToKanbanBoardOptions{IssueNumber: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.IssueNumber != 42 {
		t.Errorf("expected issue number 42, got %d", item.IssueNumber)
	}
}

func TestUpdateKanbanBoardItem(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &KanbanBoardItem{ID: 1, Status: "done"}))
	defer server.Close()

	item, err := client.UpdateKanbanBoardItem(context.Background(), "my-org", 1, 1, UpdateKanbanBoardItemOptions{Status: "done"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Status != "done" {
		t.Errorf("expected status 'done', got '%s'", item.Status)
	}
}

func TestRemoveFromKanbanBoard(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.RemoveFromKanbanBoard(context.Background(), "my-org", 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListKanbanBoardItems(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*KanbanBoardItem{
		{ID: 1, IssueNumber: 42, Status: "todo"},
	}))
	defer server.Close()

	items, err := client.ListKanbanBoardItems(context.Background(), "my-org", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}
