package gitcode

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Repository struct {
	ID            int64     `json:"id"`
	FullName      string    `json:"full_name"`
	Name          string    `json:"name"`
	Owner         *User     `json:"owner"`
	Description   string    `json:"description"`
	CloneURL      string    `json:"clone_url"`
	SSHURL        string    `json:"ssh_url"`
	HTMLURL       string    `json:"html_url"`
	DefaultBranch string    `json:"default_branch"`
	Private       bool      `json:"private"`
	Fork          bool      `json:"fork"`
	StarsCount    int       `json:"stars_count"`
	ForksCount    int       `json:"forks_count"`
	WatchersCount int       `json:"watchers_count"`
	OpenIssuesCount int     `json:"open_issues_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateRepositoryOptions struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     *bool  `json:"private,omitempty"`
	AutoInit    *bool  `json:"auto_init,omitempty"`
}

type UpdateRepositoryOptions struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	DefaultBranch string `json:"default_branch,omitempty"`
	Private       *bool  `json:"private,omitempty"`
}

type ListRepositoriesOptions struct {
	ListOptions
	Owner string `json:"owner,omitempty"`
	Type  string `json:"type,omitempty"`
	Sort  string `json:"sort,omitempty"`
}

func (c *Client) ListRepositories(ctx context.Context, opts ListRepositoriesOptions) ([]*Repository, error) {
	var repos []*Repository
	path := fmt.Sprintf("/user/repos?%s", opts.toQuery())
	if opts.Owner != "" {
		path = fmt.Sprintf("/users/%s/repos?%s", opts.Owner, opts.toQuery())
	}
	err := c.doRequest(ctx, http.MethodGet, path, nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*Repository, error) {
	var r Repository
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s", owner, repo), nil, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Client) CreateRepository(ctx context.Context, opts CreateRepositoryOptions) (*Repository, error) {
	var r Repository
	err := c.doRequest(ctx, http.MethodPost, "/user/repos", opts, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Client) UpdateRepository(ctx context.Context, owner, repo string, opts UpdateRepositoryOptions) (*Repository, error) {
	var r Repository
	err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("/repos/%s/%s", owner, repo), opts, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Client) DeleteRepository(ctx context.Context, owner, repo string) error {
	return c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/repos/%s/%s", owner, repo), nil, nil)
}

func (c *Client) ForkRepository(ctx context.Context, owner, repo string, opts *CreateRepositoryOptions) (*Repository, error) {
	var r Repository
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/forks", owner, repo), opts, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

type RepositoryContent struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Encoding string `json:"encoding,omitempty"`
	SHA     string `json:"sha"`
	Links   struct {
		Self string `json:"self"`
		Git  string `json:"git"`
	} `json:"_links"`
}

func (c *Client) GetRepositoryContent(ctx context.Context, owner, repo, path, ref string) (*RepositoryContent, error) {
	var content RepositoryContent
	url := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)
	if ref != "" {
		url += "?ref=" + ref
	}
	err := c.doRequest(ctx, http.MethodGet, url, nil, &content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (c *Client) ListRepositoryContents(ctx context.Context, owner, repo, path, ref string) ([]*RepositoryContent, error) {
	var contents []*RepositoryContent
	url := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)
	if ref != "" {
		url += "?ref=" + ref
	}
	err := c.doRequest(ctx, http.MethodGet, url, nil, &contents)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

type CreateFileOptions struct {
	Message string `json:"message"`
	Content string `json:"content"`
	Branch  string `json:"branch,omitempty"`
}

type UpdateFileOptions struct {
	Message string `json:"message"`
	Content string `json:"content"`
	SHA     string `json:"sha"`
	Branch  string `json:"branch,omitempty"`
}

type DeleteFileOptions struct {
	Message string `json:"message"`
	SHA     string `json:"sha"`
	Branch  string `json:"branch,omitempty"`
}

type FileResult struct {
	Content *RepositoryContent `json:"content"`
	Commit  *Commit            `json:"commit"`
}

func (c *Client) CreateFile(ctx context.Context, owner, repo, path string, opts CreateFileOptions) (*FileResult, error) {
	var result FileResult
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path), opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateFile(ctx context.Context, owner, repo, path string, opts UpdateFileOptions) (*FileResult, error) {
	var result FileResult
	err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path), opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteFile(ctx context.Context, owner, repo, path string, opts DeleteFileOptions) (*FileResult, error) {
	var result FileResult
	err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path), opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type Tag struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
	} `json:"commit"`
}

func (c *Client) ListTags(ctx context.Context, owner, repo string) ([]*Tag, error) {
	var tags []*Tag
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/tags?per_page=100", owner, repo), nil, &tags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

type Release struct {
	ID          int64     `json:"id"`
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	HTMLURL     string    `json:"html_url"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
}

type CreateReleaseOptions struct {
	TagName    string `json:"tag_name"`
	Target     string `json:"target_commitish,omitempty"`
	Title      string `json:"name"`
	Body       string `json:"body,omitempty"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

func (c *Client) ListReleases(ctx context.Context, owner, repo string) ([]*Release, error) {
	var releases []*Release
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/releases?per_page=100", owner, repo), nil, &releases)
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func (c *Client) CreateRelease(ctx context.Context, owner, repo string, opts CreateReleaseOptions) (*Release, error) {
	var r Release
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/releases", owner, repo), opts, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Client) DeleteRelease(ctx context.Context, owner, repo string, releaseID int64) error {
	return c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/repos/%s/%s/releases/%d", owner, repo, releaseID), nil, nil)
}

type Contributor struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Contributions int `json:"contributions"`
}

func (c *Client) ListContributors(ctx context.Context, owner, repo string) ([]*Contributor, error) {
	var contributors []*Contributor
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/contributors?per_page=100", owner, repo), nil, &contributors)
	if err != nil {
		return nil, err
	}
	return contributors, nil
}

type Language map[string]int

func (c *Client) GetLanguages(ctx context.Context, owner, repo string) (Language, error) {
	var lang Language
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/languages", owner, repo), nil, &lang)
	if err != nil {
		return nil, err
	}
	return lang, nil
}
