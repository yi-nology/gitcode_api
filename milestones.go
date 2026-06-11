package gitcode

import (
	"context"
	"fmt"
	"net/http"
)

type ListMilestonesOptions struct {
	ListOptions
	State     string `json:"state,omitempty"`
	Sort      string `json:"sort,omitempty"`
	Direction string `json:"direction,omitempty"`
}

func (c *Client) ListMilestonesWithOptions(ctx context.Context, owner, repo string, opts ListMilestonesOptions) ([]*Milestone, error) {
	var milestones []*Milestone
	query := opts.toQuery()
	if opts.State != "" {
		query += "&state=" + opts.State
	}
	if opts.Sort != "" {
		query += "&sort=" + opts.Sort
	}
	if opts.Direction != "" {
		query += "&direction=" + opts.Direction
	}
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/milestones?%s", owner, repo, query), nil, &milestones)
	if err != nil {
		return nil, err
	}
	return milestones, nil
}

func (c *Client) GetMilestone(ctx context.Context, owner, repo string, number int) (*Milestone, error) {
	var milestone Milestone
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/milestones/%d", owner, repo, number), nil, &milestone)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

type UpdateMilestoneOptions struct {
	Title       string `json:"title"`
	State       string `json:"state,omitempty"`
	Description string `json:"description,omitempty"`
	DueOn       string `json:"due_on"`
}

func (c *Client) UpdateMilestone(ctx context.Context, owner, repo string, number int, opts UpdateMilestoneOptions) (*Milestone, error) {
	var milestone Milestone
	err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("/repos/%s/%s/milestones/%d", owner, repo, number), opts, &milestone)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

type CreateMilestoneOptions struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueOn       string `json:"due_on"`
}

func (c *Client) CreateMilestoneWithOptions(ctx context.Context, owner, repo string, opts CreateMilestoneOptions) (*Milestone, error) {
	var milestone Milestone
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/milestones", owner, repo), opts, &milestone)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}
