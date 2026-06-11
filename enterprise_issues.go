package gitcode

import (
	"context"
	"fmt"
	"net/http"
)

type ListUserIssuesOptions struct {
	ListOptions
	Filter    string `json:"filter,omitempty"`
	State     string `json:"state,omitempty"`
	Labels    string `json:"labels,omitempty"`
	Sort      string `json:"sort,omitempty"`
	Direction string `json:"direction,omitempty"`
	Since     string `json:"since,omitempty"`
}

func (c *Client) ListUserIssues(ctx context.Context, opts ListUserIssuesOptions) ([]*Issue, error) {
	var issues []*Issue
	query := opts.toQuery()
	if opts.Filter != "" {
		query += "&filter=" + opts.Filter
	}
	if opts.State != "" {
		query += "&state=" + opts.State
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
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/user/issues?%s", query), nil, &issues)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (c *Client) ListOrgIssues(ctx context.Context, org string, opts ListUserIssuesOptions) ([]*Issue, error) {
	var issues []*Issue
	query := opts.toQuery()
	if opts.Filter != "" {
		query += "&filter=" + opts.Filter
	}
	if opts.State != "" {
		query += "&state=" + opts.State
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
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/orgs/%s/issues?%s", org, query), nil, &issues)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (c *Client) ListRepoAllIssueComments(ctx context.Context, owner, repo string, opts ListOptions) ([]*IssueComment, error) {
	var comments []*IssueComment
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/issues/comments?%s", owner, repo, opts.toQuery()), nil, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *Client) GetIssueComment(ctx context.Context, owner, repo string, commentID int64) (*IssueComment, error) {
	var comment IssueComment
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/issues/comments/%d", owner, repo, commentID), nil, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

type IssueOperateLog struct {
	ID        int64  `json:"id"`
	Action    string `json:"action"`
	CreatedAt string `json:"created_at"`
	User      *User  `json:"user"`
}

func (c *Client) GetIssueOperateLogs(ctx context.Context, owner, repo string, number int) ([]*IssueOperateLog, error) {
	var logs []*IssueOperateLog
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/issues/%d/operate_logs", owner, repo, number), nil, &logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (c *Client) GetIssueLinkedPRs(ctx context.Context, owner, repo string, number int) ([]*PullRequest, error) {
	var prs []*PullRequest
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/issues/%d/pull_requests", owner, repo, number), nil, &prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func (c *Client) ListEnterpriseIssues(ctx context.Context, enterprise string, opts ListUserIssuesOptions) ([]*Issue, error) {
	var issues []*Issue
	query := opts.toQuery()
	if opts.Filter != "" {
		query += "&filter=" + opts.Filter
	}
	if opts.State != "" {
		query += "&state=" + opts.State
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
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/enterprises/%s/issues?%s", enterprise, query), nil, &issues)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (c *Client) GetEnterpriseIssue(ctx context.Context, enterprise string, number int) (*Issue, error) {
	var issue Issue
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/enterprises/%s/issues/%d", enterprise, number), nil, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

func (c *Client) ListEnterpriseIssueComments(ctx context.Context, enterprise string, number int, opts ListOptions) ([]*IssueComment, error) {
	var comments []*IssueComment
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/enterprises/%s/issues/%d/comments?%s", enterprise, number, opts.toQuery()), nil, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *Client) ListEnterpriseIssueLabels(ctx context.Context, enterprise string, issueID int64) ([]*Label, error) {
	var labels []*Label
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/enterprises/%s/issues/%d/labels", enterprise, issueID), nil, &labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}
