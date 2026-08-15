package gitcode

import (
	"context"
	"testing"
)

func TestGetLatestRelease(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &Release{
		ID: 1, TagName: "v1.0.0", Name: "Latest Release",
	}))
	defer server.Close()

	release, err := client.GetLatestRelease(context.Background(), "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if release.TagName != "v1.0.0" {
		t.Errorf("expected tag 'v1.0.0', got '%s'", release.TagName)
	}
}

func TestGetReleaseUploadURL(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &ReleaseUploadURL{
		UploadURL: "https://upload.example.com/assets",
	}))
	defer server.Close()

	url, err := client.GetReleaseUploadURL(context.Background(), "owner", "repo", 1, "binary.tar.gz", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url.UploadURL == "" {
		t.Error("expected non-empty upload URL")
	}
}
