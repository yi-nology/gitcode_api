package gitcode

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"
)

// RetryPolicy defines the retry behavior for failed requests.
type RetryPolicy struct {
	// MaxRetries is the maximum number of retry attempts. Default: 0 (no retry).
	MaxRetries int
	// InitialBackoff is the initial wait time before the first retry. Default: 1s.
	InitialBackoff time.Duration
	// MaxBackoff is the maximum wait time between retries. Default: 30s.
	MaxBackoff time.Duration
	// Multiplier is the backoff multiplier for exponential backoff. Default: 2.0.
	Multiplier float64
	// RetryableStatusCodes defines which HTTP status codes should trigger a retry.
	// Default: [429, 500, 502, 503, 504]
	RetryableStatusCodes []int
}

// DefaultRetryPolicy returns a sensible default retry policy.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxRetries:           3,
		InitialBackoff:       1 * time.Second,
		MaxBackoff:           30 * time.Second,
		Multiplier:           2.0,
		RetryableStatusCodes: []int{429, 500, 502, 503, 504},
	}
}

// SetRetryPolicy sets the retry policy for the client.
func (c *Client) SetRetryPolicy(policy RetryPolicy) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.retryPolicy = &policy
}

// shouldRetry checks if a request should be retried based on the status code.
func (c *Client) shouldRetry(statusCode int) bool {
	c.mu.RLock()
	policy := c.retryPolicy
	c.mu.RUnlock()

	if policy == nil {
		return false
	}

	for _, code := range policy.RetryableStatusCodes {
		if code == statusCode {
			return true
		}
	}
	return false
}

// getRetryWaitTime calculates the wait time for a given retry attempt.
func (c *Client) getRetryWaitTime(attempt int, resp *http.Response) time.Duration {
	c.mu.RLock()
	policy := c.retryPolicy
	c.mu.RUnlock()

	if policy == nil {
		return 0
	}

	// Check Retry-After header first (for 429 responses)
	if resp != nil {
		if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
			if seconds, err := strconv.Atoi(retryAfter); err == nil {
				return time.Duration(seconds) * time.Second
			}
		}
	}

	// Exponential backoff
	backoff := float64(policy.InitialBackoff) * math.Pow(policy.Multiplier, float64(attempt))
	if backoff > float64(policy.MaxBackoff) {
		backoff = float64(policy.MaxBackoff)
	}
	return time.Duration(backoff)
}

// doRequestWithRetry executes an HTTP request with retry logic.
func (c *Client) doRequestWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	c.mu.RLock()
	policy := c.retryPolicy
	c.mu.RUnlock()

	maxRetries := 0
	if policy != nil {
		maxRetries = policy.MaxRetries
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Clone the request for retries (body can be nil for retries)
		if attempt > 0 {
			waitTime := c.getRetryWaitTime(attempt-1, resp)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(waitTime):
			}
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		// Check if we should retry
		if c.shouldRetry(resp.StatusCode) && attempt < maxRetries {
			resp.Body.Close()
			continue
		}

		return resp, nil
	}

	return resp, nil
}
