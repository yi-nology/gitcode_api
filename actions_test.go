package gitcode

import (
	"context"
	"testing"
)

func TestListActionWorkflows(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*ActionWorkflow{
		{ID: 1, Name: "CI", State: "active"},
	}))
	defer server.Close()

	workflows, err := client.ListActionWorkflows(context.Background(), "owner", "repo", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow, got %d", len(workflows))
	}
}

func TestListActionRuns(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*ActionRun{
		{ID: 1, Name: "CI #1", Status: "completed", Conclusion: "success"},
	}))
	defer server.Close()

	runs, err := client.ListActionRuns(context.Background(), "owner", "repo", ListActionRunsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runs) != 1 {
		t.Errorf("expected 1 run, got %d", len(runs))
	}
}

func TestGetActionRun(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &ActionRun{
		ID: 1, Name: "CI #1", Status: "completed",
	}))
	defer server.Close()

	run, err := client.GetActionRun(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if run.Name != "CI #1" {
		t.Errorf("expected name 'CI #1', got '%s'", run.Name)
	}
}

func TestListActionJobs(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*ActionJob{
		{ID: 1, Name: "build", Status: "completed", Conclusion: "success"},
	}))
	defer server.Close()

	jobs, err := client.ListActionJobs(context.Background(), "owner", "repo", 1, ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(jobs))
	}
}

func TestGetActionJob(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &ActionJob{
		ID: 1, Name: "build", Status: "completed",
	}))
	defer server.Close()

	job, err := client.GetActionJob(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.Name != "build" {
		t.Errorf("expected name 'build', got '%s'", job.Name)
	}
}

func TestListActionArtifacts(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*ActionArtifact{
		{ID: 1, Name: "build-output", Size: 1024},
	}))
	defer server.Close()

	artifacts, err := client.ListActionArtifacts(context.Background(), "owner", "repo", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(artifacts) != 1 {
		t.Errorf("expected 1 artifact, got %d", len(artifacts))
	}
}

func TestGetActionArtifact(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &ActionArtifact{
		ID: 1, Name: "build-output", Size: 1024,
	}))
	defer server.Close()

	artifact, err := client.GetActionArtifact(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if artifact.Name != "build-output" {
		t.Errorf("expected name 'build-output', got '%s'", artifact.Name)
	}
}

func TestDeleteActionArtifact(t *testing.T) {
	client, server := newTestServer(jsonResponse(204, nil))
	defer server.Close()

	err := client.DeleteActionArtifact(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListRepoHostRunners(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*ActionRunner{
		{ID: 1, Name: "runner-1", Status: "online"},
	}))
	defer server.Close()

	runners, err := client.ListRepoHostRunners(context.Background(), "owner", "repo", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runners) != 1 {
		t.Errorf("expected 1 runner, got %d", len(runners))
	}
}

func TestListOrgRunnerGroups(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, []*RunnerGroup{
		{ID: 1, Name: "default"},
	}))
	defer server.Close()

	groups, err := client.ListOrgRunnerGroups(context.Background(), "my-org", ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestGetRunnerGroup(t *testing.T) {
	client, server := newTestServer(jsonResponse(200, &RunnerGroup{ID: 1, Name: "default"}))
	defer server.Close()

	group, err := client.GetRunnerGroup(context.Background(), "my-org", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "default" {
		t.Errorf("expected name 'default', got '%s'", group.Name)
	}
}
