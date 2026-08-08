package gitcode

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// KanbanBoard represents a kanban board.
type KanbanBoard struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// KanbanBoardItem represents an item on a kanban board.
type KanbanBoardItem struct {
	ID          int64  `json:"id"`
	BoardID     int64  `json:"board_id"`
	IssueNumber int    `json:"issue_number,omitempty"`
	PRNumber    int    `json:"pr_number,omitempty"`
	Status      string `json:"status"`
	Sort        int    `json:"sort,omitempty"`
}

// KanbanBoardStatus represents a kanban board status.
type KanbanBoardStatus struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color,omitempty"`
	Sort      int    `json:"sort,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
}

// ListKanbanBoards lists all kanban boards for an organization.
//
// GET /orgs/{org}/kanban/boards
func (c *Client) ListKanbanBoards(ctx context.Context, org string, opts ListOptions) ([]*KanbanBoard, error) {
	var boards []*KanbanBoard
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/orgs/%s/kanban/boards?%s", org, opts.toQuery()), nil, &boards)
	if err != nil {
		return nil, err
	}
	return boards, nil
}

// GetKanbanBoard gets a single kanban board.
//
// GET /orgs/{org}/kanban/boards/{id}
func (c *Client) GetKanbanBoard(ctx context.Context, org string, boardID int64) (*KanbanBoard, error) {
	var board KanbanBoard
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/orgs/%s/kanban/boards/%d", org, boardID), nil, &board)
	if err != nil {
		return nil, err
	}
	return &board, nil
}

// AddToKanbanBoardOptions specifies options for adding an item to a kanban board.
type AddToKanbanBoardOptions struct {
	IssueNumber int    `json:"issue_number,omitempty"`
	PRNumber    int    `json:"pr_number,omitempty"`
	Status      string `json:"status,omitempty"`
}

// AddToKanbanBoard adds an issue or pull request to a kanban board.
//
// POST /orgs/{org}/kanban/boards/{id}/items
func (c *Client) AddToKanbanBoard(ctx context.Context, org string, boardID int64, opts AddToKanbanBoardOptions) (*KanbanBoardItem, error) {
	var item KanbanBoardItem
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/orgs/%s/kanban/boards/%d/items", org, boardID), opts, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateKanbanBoardItemOptions specifies options for updating a kanban board item.
type UpdateKanbanBoardItemOptions struct {
	Status string `json:"status,omitempty"`
	Sort   int    `json:"sort,omitempty"`
}

// UpdateKanbanBoardItem updates an item on a kanban board.
//
// PUT /orgs/{org}/kanban/boards/{id}/items/{item_id}
func (c *Client) UpdateKanbanBoardItem(ctx context.Context, org string, boardID, itemID int64, opts UpdateKanbanBoardItemOptions) (*KanbanBoardItem, error) {
	var item KanbanBoardItem
	err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/orgs/%s/kanban/boards/%d/items/%d", org, boardID, itemID), opts, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// RemoveFromKanbanBoard removes an item from a kanban board.
//
// DELETE /orgs/{org}/kanban/boards/{id}/items/{item_id}
func (c *Client) RemoveFromKanbanBoard(ctx context.Context, org string, boardID, itemID int64) error {
	return c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/orgs/%s/kanban/boards/%d/items/%d", org, boardID, itemID), nil, nil)
}

// ListKanbanBoardItems lists all items on a kanban board.
//
// GET /orgs/{org}/kanban/boards/{id}/items
func (c *Client) ListKanbanBoardItems(ctx context.Context, org string, boardID int64, opts ListOptions) ([]*KanbanBoardItem, error) {
	var items []*KanbanBoardItem
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/orgs/%s/kanban/boards/%d/items?%s", org, boardID, opts.toQuery()), nil, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateKanbanBoardStatusOptions specifies options for updating a kanban board status.
type UpdateKanbanBoardStatusOptions struct {
	Statuses []KanbanBoardStatus `json:"statuses"`
}

// UpdateKanbanBoardStatus updates the statuses of a kanban board.
//
// PUT /orgs/{org}/kanban/boards/{id}/status
func (c *Client) UpdateKanbanBoardStatus(ctx context.Context, org string, boardID int64, opts UpdateKanbanBoardStatusOptions) error {
	return c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/orgs/%s/kanban/boards/%d/status", org, boardID), opts, nil)
}
