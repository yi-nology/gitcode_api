package gitcode

import (
	"context"
	"net/http"
	"testing"
)

func TestGetCommitDiff(t *testing.T) {
	diffContent := []byte("diff --git a/main.go b/main.go\n--- a/main.go\n+++ b/main.go\n@@ -1 +1 @@\n-old\n+new")
	client, server := newTestServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(200)
			w.Write(diffContent)
		}
	}())
	defer server.Close()

	diff, err := client.GetCommitDiff(context.Background(), "owner", "repo", "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(diff) != string(diffContent) {
		t.Error("diff content mismatch")
	}
}

func TestGetCommitPatch(t *testing.T) {
	patchContent := []byte("From abc123 Mon Sep 17 00:00:00 2001\nFrom: test <test@example.com>\nSubject: [PATCH] test")
	client, server := newTestServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(200)
			w.Write(patchContent)
		}
	}())
	defer server.Close()

	patch, err := client.GetCommitPatch(context.Background(), "owner", "repo", "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(patch) != string(patchContent) {
		t.Error("patch content mismatch")
	}
}
