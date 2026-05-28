package gitcode

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type PullRequestState string

const (
	PullRequestStateOpen   PullRequestState = "open"
	PullRequestStateClosed PullRequestState = "closed"
)

type PullRequest struct {
	ID           int64           `json:"id"`
	Number       int             `json:"number"`
	Title        string          `json:"title"`
	Body         string          `json:"body"`
	State        PullRequestState `json:"state"`
	Author       *User           `json:"author"`
	Head         *PullRequestBranch `json:"head"`
	Base         *PullRequestBranch `json:"base"`
	Merged       bool            `json:"merged"`
	Mergeable    *bool           `json:"mergeable"`
	HTMLURL      string          `json:"html_url"`
	DiffURL      string          `json:"diff_url"`
	PatchURL     string          `json:"patch_url"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	ClosedAt     *time.Time      `json:"closed_at"`
	MergedAt     *time.Time      `json:"merged_at"`
}

type PullRequestBranch struct {
	Ref string `json:"ref"`
	SHA string `json:"sha"`
	Repo *Repository `json:"repo"`
}

type PullRequestFile struct {
	Filename        string `json:"filename"`
	PreviousFilename string `json:"previous_filename"`
	Status          string `json:"status"`
	Additions       int    `json:"additions"`
	Deletions       int    `json:"deletions"`
	Changes         int    `json:"changes"`
	Patch           string `json:"patch"`
}

type PullRequestComment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	Author    *User     `json:"author"`
	Path      string    `json:"path"`
	Position  int       `json:"position"`
	CommitID  string    `json:"commit_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PullRequestReview struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	State     string    `json:"state"`
	Author    *User     `json:"author"`
	CommitID  string    `json:"commit_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListPullRequestsOptions struct {
	ListOptions
	State        PullRequestState `json:"state,omitempty"`
	Sort         string           `json:"sort,omitempty"`
	Direction    string           `json:"direction,omitempty"`
	Head         string           `json:"head,omitempty"`
	Base         string           `json:"base,omitempty"`
}

type CreatePullRequestOptions struct {
	Title     string `json:"title"`
	Body      string `json:"body,omitempty"`
	Head      string `json:"head"`
	Base      string `json:"base"`
	Draft     bool   `json:"draft,omitempty"`
}

type UpdatePullRequestOptions struct {
	Title        string          `json:"title,omitempty"`
	Body         string          `json:"body,omitempty"`
	State        PullRequestState `json:"state,omitempty"`
	Base         string          `json:"base,omitempty"`
}

func (c *Client) ListPullRequests(ctx context.Context, owner, repo string, opts ListPullRequestsOptions) ([]*PullRequest, error) {
	var prs []*PullRequest
	query := opts.toQuery()
	if opts.State != "" {
		query += "&state=" + string(opts.State)
	}
	if opts.Sort != "" {
		query += "&sort=" + opts.Sort
	}
	if opts.Direction != "" {
		query += "&direction=" + opts.Direction
	}
	if opts.Head != "" {
		query += "&head=" + opts.Head
	}
	if opts.Base != "" {
		query += "&base=" + opts.Base
	}
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/pulls?%s", owner, repo, query), nil, &prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func (c *Client) GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	var pr PullRequest
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, number), nil, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (c *Client) CreatePullRequest(ctx context.Context, owner, repo string, opts CreatePullRequestOptions) (*PullRequest, error) {
	var pr PullRequest
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/pulls", owner, repo), opts, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (c *Client) UpdatePullRequest(ctx context.Context, owner, repo string, number int, opts UpdatePullRequestOptions) (*PullRequest, error) {
	var pr PullRequest
	err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, number), opts, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (c *Client) ClosePullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	return c.UpdatePullRequest(ctx, owner, repo, number, UpdatePullRequestOptions{State: PullRequestStateClosed})
}

func (c *Client) ReopenPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	return c.UpdatePullRequest(ctx, owner, repo, number, UpdatePullRequestOptions{State: PullRequestStateOpen})
}

type MergePullRequestOptions struct {
	CommitMessage string `json:"commit_message,omitempty"`
	Squash        bool   `json:"squash,omitempty"`
}

func (c *Client) MergePullRequest(ctx context.Context, owner, repo string, number int, opts *MergePullRequestOptions) error {
	body := map[string]interface{}{}
	if opts != nil {
		if opts.CommitMessage != "" {
			body["commit_message"] = opts.CommitMessage
		}
		if opts.Squash {
			body["merge_method"] = "squash"
		}
	}
	return c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/repos/%s/%s/pulls/%d/merge", owner, repo, number), body, nil)
}

func (c *Client) IsPullRequestMerged(ctx context.Context, owner, repo string, number int) (bool, error) {
	_, err := c.doRawRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/pulls/%d/merge", owner, repo, number))
	if err != nil {
		if err.Error() == "404 Not Found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *Client) ListPullRequestFiles(ctx context.Context, owner, repo string, number int) ([]*PullRequestFile, error) {
	var files []*PullRequestFile
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/pulls/%d/files", owner, repo, number), nil, &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (c *Client) ListPullRequestComments(ctx context.Context, owner, repo string, number int) ([]*PullRequestComment, error) {
	var comments []*PullRequestComment
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/pulls/%d/comments", owner, repo, number), nil, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *Client) CreatePullRequestComment(ctx context.Context, owner, repo string, number int, body, path string, position, commitID string) (*PullRequestComment, error) {
	var comment PullRequestComment
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/pulls/%d/comments", owner, repo, number), map[string]interface{}{
		"body":      body,
		"path":      path,
		"position":  position,
		"commit_id": commitID,
	}, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *Client) UpdatePullRequestComment(ctx context.Context, owner, repo string, commentID int64, body string) (*PullRequestComment, error) {
	var comment PullRequestComment
	err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("/repos/%s/%s/pulls/comments/%d", owner, repo, commentID), map[string]string{"body": body}, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *Client) DeletePullRequestComment(ctx context.Context, owner, repo string, commentID int64) error {
	return c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/repos/%s/%s/pulls/comments/%d", owner, repo, commentID), nil, nil)
}

func (c *Client) ListPullRequestReviews(ctx context.Context, owner, repo string, number int) ([]*PullRequestReview, error) {
	var reviews []*PullRequestReview
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews", owner, repo, number), nil, &reviews)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (c *Client) CreatePullRequestReview(ctx context.Context, owner, repo string, number int, body, event string) (*PullRequestReview, error) {
	var review PullRequestReview
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews", owner, repo, number), map[string]interface{}{
		"body":  body,
		"event": event,
	}, &review)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (c *Client) ListPullRequestCommits(ctx context.Context, owner, repo string, number int) ([]*Commit, error) {
	var commits []*Commit
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/pulls/%d/commits", owner, repo, number), nil, &commits)
	if err != nil {
		return nil, err
	}
	return commits, nil
}
