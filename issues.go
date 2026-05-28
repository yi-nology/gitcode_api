package gitcode

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type IssueState string

const (
	IssueStateOpen   IssueState = "open"
	IssueStateClosed IssueState = "closed"
)

type Issue struct {
	ID        int64     `json:"id"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	State     IssueState `json:"state"`
	Author    *User     `json:"author"`
	Assignees []*User   `json:"assignees"`
	Labels    []*Label  `json:"labels"`
	Milestone *Milestone `json:"milestone"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ClosedAt  *time.Time `json:"closed_at"`
}

type Label struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Milestone struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	State       string     `json:"state"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type IssueComment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	Author    *User     `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListIssuesOptions struct {
	ListOptions
	State     IssueState `json:"state,omitempty"`
	Assignee  string     `json:"assignee,omitempty"`
	Creator   string     `json:"creator,omitempty"`
	Milestone string     `json:"milestone,omitempty"`
	Labels    string     `json:"labels,omitempty"`
	Sort      string     `json:"sort,omitempty"`
	Direction string     `json:"direction,omitempty"`
	Since     string     `json:"since,omitempty"`
}

type CreateIssueOptions struct {
	Title     string   `json:"title"`
	Body      string   `json:"body,omitempty"`
	Assignee  string   `json:"assignee,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Milestone int64    `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
}

type UpdateIssueOptions struct {
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	State     IssueState `json:"state,omitempty"`
	Assignee  string   `json:"assignee,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Milestone int64    `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
}

func (c *Client) ListIssues(ctx context.Context, owner, repo string, opts ListIssuesOptions) ([]*Issue, error) {
	var issues []*Issue
	query := opts.toQuery()
	if opts.State != "" {
		query += "&state=" + string(opts.State)
	}
	if opts.Assignee != "" {
		query += "&assignee=" + opts.Assignee
	}
	if opts.Creator != "" {
		query += "&creator=" + opts.Creator
	}
	if opts.Milestone != "" {
		query += "&milestone=" + opts.Milestone
	}
	if opts.Labels != "" {
		query += "&labels=" + opts.Labels
	}
	if opts.Sort != "" {
		query += "&sort=" + opts.Sort
	}
	if opts.Direction != "" {
		query += "&direction=" + opts.Direction
	}
	if opts.Since != "" {
		query += "&since=" + opts.Since
	}
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/issues?%s", owner, repo, query), nil, &issues)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (c *Client) GetIssue(ctx context.Context, owner, repo string, number int) (*Issue, error) {
	var issue Issue
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, number), nil, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

func (c *Client) CreateIssue(ctx context.Context, owner, repo string, opts CreateIssueOptions) (*Issue, error) {
	var issue Issue
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/issues", owner, repo), opts, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

func (c *Client) UpdateIssue(ctx context.Context, owner, repo string, number int, opts UpdateIssueOptions) (*Issue, error) {
	var issue Issue
	err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, number), opts, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

func (c *Client) CloseIssue(ctx context.Context, owner, repo string, number int) (*Issue, error) {
	return c.UpdateIssue(ctx, owner, repo, number, UpdateIssueOptions{State: IssueStateClosed})
}

func (c *Client) ReopenIssue(ctx context.Context, owner, repo string, number int) (*Issue, error) {
	return c.UpdateIssue(ctx, owner, repo, number, UpdateIssueOptions{State: IssueStateOpen})
}

func (c *Client) ListIssueComments(ctx context.Context, owner, repo string, number int) ([]*IssueComment, error) {
	var comments []*IssueComment
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/issues/%d/comments", owner, repo, number), nil, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *Client) CreateIssueComment(ctx context.Context, owner, repo string, number int, body string) (*IssueComment, error) {
	var comment IssueComment
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/issues/%d/comments", owner, repo, number), map[string]string{"body": body}, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *Client) UpdateIssueComment(ctx context.Context, owner, repo string, commentID int64, body string) (*IssueComment, error) {
	var comment IssueComment
	err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("/repos/%s/%s/issues/comments/%d", owner, repo, commentID), map[string]string{"body": body}, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *Client) DeleteIssueComment(ctx context.Context, owner, repo string, commentID int64) error {
	return c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/repos/%s/%s/issues/comments/%d", owner, repo, commentID), nil, nil)
}

func (c *Client) ListIssueLabels(ctx context.Context, owner, repo string) ([]*Label, error) {
	var labels []*Label
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/labels", owner, repo), nil, &labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func (c *Client) CreateIssueLabel(ctx context.Context, owner, repo string, name, color string) (*Label, error) {
	var label Label
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/labels", owner, repo), map[string]string{"name": name, "color": color}, &label)
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (c *Client) DeleteIssueLabel(ctx context.Context, owner, repo string, name string) error {
	return c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/repos/%s/%s/labels/%s", owner, repo, name), nil, nil)
}

func (c *Client) AddIssueLabels(ctx context.Context, owner, repo string, number int, labels []string) error {
	return c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/issues/%d/labels", owner, repo, number), map[string][]string{"labels": labels}, nil)
}

func (c *Client) RemoveIssueLabel(ctx context.Context, owner, repo string, number int, name string) error {
	return c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/repos/%s/%s/issues/%d/labels/%s", owner, repo, number, name), nil, nil)
}

func (c *Client) ListMilestones(ctx context.Context, owner, repo string) ([]*Milestone, error) {
	var milestones []*Milestone
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/milestones", owner, repo), nil, &milestones)
	if err != nil {
		return nil, err
	}
	return milestones, nil
}

func (c *Client) CreateMilestone(ctx context.Context, owner, repo string, title, description string) (*Milestone, error) {
	var milestone Milestone
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/milestones", owner, repo), map[string]string{"title": title, "description": description}, &milestone)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

func (c *Client) DeleteMilestone(ctx context.Context, owner, repo string, number int) error {
	return c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/repos/%s/%s/milestones/%d", owner, repo, number), nil, nil)
}
