package gitcode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultBaseURL = "https://api.gitcode.com/api/v5"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	authStyle  AuthStyle
}

type AuthStyle int

const (
	AuthStyleBearer AuthStyle = iota
	AuthStylePrivateToken
	AuthStyleAccessToken
)

func NewClient(token string) *Client {
	return &Client{
		baseURL: DefaultBaseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authStyle: AuthStyleBearer,
	}
}

func NewClientWithBaseURL(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authStyle: AuthStyleBearer,
	}
}

func (c *Client) SetAuthStyle(style AuthStyle) {
	c.authStyle = style
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

func (c *Client) setAuthHeader(req *http.Request) {
	switch c.authStyle {
	case AuthStylePrivateToken:
		req.Header.Set("PRIVATE-TOKEN", c.token)
	case AuthStyleAccessToken:
		q := req.URL.Query()
		q.Set("access_token", c.token)
		req.URL.RawQuery = q.Encode()
	default:
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	_, err := c.doRequestWithHeaders(ctx, method, path, body, result)
	return err
}

func (c *Client) doRequestWithHeaders(ctx context.Context, method, path string, body interface{}, result interface{}) (http.Header, error) {
	var reqBody io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	c.setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("GitCode API %s %s returned %d: %s", method, path, resp.StatusCode, string(respBody))
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.Unmarshal(respBody, result); err != nil {
			return nil, err
		}
	}

	return resp.Header, nil
}

func (c *Client) doRawRequest(ctx context.Context, method, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("GitCode API %s %s returned %d: %s", method, path, resp.StatusCode, string(body[:min(len(body), 200)]))
	}

	return body, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	var user User
	err := c.doRequest(ctx, http.MethodGet, "/user", nil, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) GetUser(ctx context.Context, username string) (*User, error) {
	var user User
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/users/%s", username), nil, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type ListOptions struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func (o ListOptions) withDefaults() ListOptions {
	if o.Page == 0 {
		o.Page = 1
	}
	if o.PerPage == 0 {
		o.PerPage = 20
	}
	return o
}

func (o ListOptions) toQuery() string {
	opts := o.withDefaults()
	return fmt.Sprintf("page=%d&per_page=%d", opts.Page, opts.PerPage)
}

type RateLimit struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     time.Time `json:"reset"`
}

func (c *Client) GetRateLimit(ctx context.Context) (*RateLimit, error) {
	var rl RateLimit
	err := c.doRequest(ctx, http.MethodGet, "/rate_limit", nil, &rl)
	if err != nil {
		return nil, err
	}
	return &rl, nil
}
