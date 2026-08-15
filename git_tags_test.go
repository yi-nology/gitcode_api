package gitcode

import (
	"context"
	"net/http"
	"testing"
)

func TestGetAnnotatedTag(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &AnnotatedTag{
		SHA: "abc123", Tag: "v1.0.0", Message: "Release v1.0.0",
	}))
	defer server.Close()

	tag, err := client.GetAnnotatedTag(context.Background(), "owner", "repo", "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Tag != "v1.0.0" {
		t.Errorf("expected tag 'v1.0.0', got '%s'", tag.Tag)
	}
}

func TestCreateAnnotatedTag(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &AnnotatedTag{
		SHA: "abc123", Tag: "v2.0.0",
	}))
	defer server.Close()

	tag, err := client.CreateAnnotatedTag(context.Background(), "owner", "repo", CreateAnnotatedTagOptions{
		Tag: "v2.0.0", Message: "Release v2.0.0", Object: "abc123", Type: "commit",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Tag != "v2.0.0" {
		t.Errorf("expected tag 'v2.0.0', got '%s'", tag.Tag)
	}
}

func TestGetTag(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &Tag{Name: "v1.0.0"}))
	defer server.Close()

	tag, err := client.GetTag(context.Background(), "owner", "repo", "v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Name != "v1.0.0" {
		t.Errorf("expected name 'v1.0.0', got '%s'", tag.Name)
	}
}

func TestCreateTag(t *testing.T) {
	client, server := newTestServer(jsonResponse(201, &Tag{Name: "v2.0.0"}))
	defer server.Close()

	tag, err := client.CreateTag(context.Background(), "owner", "repo", CreateTagOptions{
		TagName: "v2.0.0", Target: "abc123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Name != "v2.0.0" {
		t.Errorf("expected name 'v2.0.0', got '%s'", tag.Name)
	}
}

func TestGetReleaseByTag(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &Release{TagName: "v1.0.0", Name: "Release 1.0"}))
	defer server.Close()

	release, err := client.GetReleaseByTag(context.Background(), "owner", "repo", "v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if release.TagName != "v1.0.0" {
		t.Errorf("expected tag 'v1.0.0', got '%s'", release.TagName)
	}
}

func TestGetRelease(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &Release{ID: 1, TagName: "v1.0.0"}))
	defer server.Close()

	release, err := client.GetRelease(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if release.TagName != "v1.0.0" {
		t.Errorf("expected tag 'v1.0.0', got '%s'", release.TagName)
	}
}

func TestUpdateRelease(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &Release{ID: 1, Body: "Updated notes"}))
	defer server.Close()

	release, err := client.UpdateRelease(context.Background(), "owner", "repo", 1, UpdateReleaseOptions{
		Body: "Updated notes",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if release.Body != "Updated notes" {
		t.Errorf("expected body 'Updated notes', got '%s'", release.Body)
	}
}

func TestDeleteReleaseByID(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteReleaseByID(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListReleaseAssets(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*ReleaseAsset{
		{ID: 1, Name: "binary.tar.gz", Size: 1024},
	}))
	defer server.Close()

	assets, err := client.ListReleaseAssets(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(assets) != 1 {
		t.Errorf("expected 1 asset, got %d", len(assets))
	}
}

func TestGetReleaseAsset(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &ReleaseAsset{ID: 1, Name: "binary.tar.gz"}))
	defer server.Close()

	asset, err := client.GetReleaseAsset(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if asset.Name != "binary.tar.gz" {
		t.Errorf("expected name 'binary.tar.gz', got '%s'", asset.Name)
	}
}

func TestDeleteReleaseAsset(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteReleaseAsset(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDownloadReleaseAsset(t *testing.T) {
	assetContent := []byte("fake-asset-content")
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(assetContent)
	}))
	defer server.Close()

	data, err := client.DownloadReleaseAsset(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != string(assetContent) {
		t.Error("asset content mismatch")
	}
}
