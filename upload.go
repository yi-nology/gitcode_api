package gitcode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// UploadFileReader uploads a file using multipart form data.
// This is the proper way to upload file contents (as opposed to UploadFile which takes a path).
//
// POST /repos/{owner}/{repo}/file/upload
func (c *Client) UploadFileReader(ctx context.Context, owner, repo, filename string, reader io.Reader) (*FileUploadResult, error) {
	if err := validateOwnerRepo(owner, repo); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}

	if _, err := io.Copy(part, reader); err != nil {
		return nil, fmt.Errorf("copy file content: %w", err)
	}

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+fmt.Sprintf("/repos/%s/%s/file/upload", owner, repo), &buf)
	if err != nil {
		return nil, err
	}

	c.setAuthHeader(req)
	req.Header.Set("Content-Type", writer.FormDataContentType())

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
		return nil, newAPIError(http.MethodPost, req.URL.Path, resp.StatusCode, string(respBody))
	}

	var result FileUploadResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadImageReader uploads an image using multipart form data.
//
// POST /repos/{owner}/{repo}/img/upload
func (c *Client) UploadImageReader(ctx context.Context, owner, repo, filename string, reader io.Reader) (*FileUploadResult, error) {
	if err := validateOwnerRepo(owner, repo); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}

	if _, err := io.Copy(part, reader); err != nil {
		return nil, fmt.Errorf("copy file content: %w", err)
	}

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+fmt.Sprintf("/repos/%s/%s/img/upload", owner, repo), &buf)
	if err != nil {
		return nil, err
	}

	c.setAuthHeader(req)
	req.Header.Set("Content-Type", writer.FormDataContentType())

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
		return nil, newAPIError(http.MethodPost, req.URL.Path, resp.StatusCode, string(respBody))
	}

	var result FileUploadResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadFileBytes uploads file bytes using multipart form data.
func (c *Client) UploadFileBytes(ctx context.Context, owner, repo, filename string, data []byte) (*FileUploadResult, error) {
	return c.UploadFileReader(ctx, owner, repo, filename, bytes.NewReader(data))
}

// UploadFileWithPath uploads a file from a local path using multipart form data.
func (c *Client) UploadFileWithPath(ctx context.Context, owner, repo, localPath string) (*FileUploadResult, error) {
	// This is a convenience wrapper - the caller should read the file themselves
	// and pass an io.Reader. This method is kept for backward compatibility.
	filename := filepath.Base(localPath)
	return c.UploadFile(ctx, owner, repo, filename)
}
