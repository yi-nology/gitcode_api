package gitcode

import (
	"context"
	"net/http"
	"testing"
)

func TestGetRepositoryArchive(t *testing.T) {
	archiveContent := []byte("fake-tar-gz-content")
	client, server := newTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.Write(archiveContent)
	}))
	defer server.Close()

	data, err := client.GetRepositoryArchive(context.Background(), "owner", "repo", "main.tar.gz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != string(archiveContent) {
		t.Error("archive content mismatch")
	}
}

func TestGetRepoParticipation(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &ParticipationStats{
		All:   []int{1, 2, 3, 4, 5},
		Owner: []int{1, 1, 1, 1, 1},
	}))
	defer server.Close()

	stats, err := client.GetRepoParticipation(context.Background(), "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats.All) != 5 {
		t.Errorf("expected 5 entries, got %d", len(stats.All))
	}
}

func TestGetRepoCodeFrequency(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []CodeFrequencyEntry{
		{1000, 50, -20},
	}))
	defer server.Close()

	freq, err := client.GetRepoCodeFrequency(context.Background(), "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(freq) != 1 {
		t.Errorf("expected 1 entry, got %d", len(freq))
	}
}

func TestGetRepoCommitActivity(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*CommitActivity{
		{Total: 10, Days: []int{1, 2, 3, 4, 5, 6, 7}},
	}))
	defer server.Close()

	activity, err := client.GetRepoCommitActivity(context.Background(), "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(activity) != 1 {
		t.Errorf("expected 1 entry, got %d", len(activity))
	}
}

func TestGetRepoPunchCard(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []PunchCardEntry{
		{0, 0, 5}, {0, 1, 10},
	}))
	defer server.Close()

	card, err := client.GetRepoPunchCard(context.Background(), "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(card) != 2 {
		t.Errorf("expected 2 entries, got %d", len(card))
	}
}
