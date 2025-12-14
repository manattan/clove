package worktree

import (
	"strings"
	"testing"
)

func TestParseWorktreeList(t *testing.T) {
	input := `worktree /path/to/repo
HEAD 1234567890abcdef
branch refs/heads/main

worktree /path/to/worktree1
HEAD abcdef1234567890
branch refs/heads/feature/update

worktree /path/to/worktree2
HEAD fedcba0987654321
detached
`

	result, err := ParseWorktreeList(input)
	if err != nil {
		t.Fatalf("ParseWorktreeList failed: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 worktrees, got %d", len(result))
	}

	// Check first worktree
	if result[0].Path != "/path/to/repo" {
		t.Errorf("expected path /path/to/repo, got %s", result[0].Path)
	}
	if result[0].Branch != "refs/heads/main" {
		t.Errorf("expected branch refs/heads/main, got %s", result[0].Branch)
	}

	// Check second worktree
	if result[1].Path != "/path/to/worktree1" {
		t.Errorf("expected path /path/to/worktree1, got %s", result[1].Path)
	}
	if result[1].Branch != "refs/heads/feature/update" {
		t.Errorf("expected branch refs/heads/feature/update, got %s", result[1].Branch)
	}

	// Check third worktree (detached)
	if result[2].Path != "/path/to/worktree2" {
		t.Errorf("expected path /path/to/worktree2, got %s", result[2].Path)
	}
	if result[2].Branch != "" {
		t.Errorf("expected empty branch for detached worktree, got %s", result[2].Branch)
	}
}

func TestParseWorktreeList_Empty(t *testing.T) {
	result, err := ParseWorktreeList("")
	if err != nil {
		t.Fatalf("ParseWorktreeList failed: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 worktrees, got %d", len(result))
	}
}

func TestParseWorktreeList_SingleWorktree(t *testing.T) {
	input := `worktree /path/to/repo
HEAD 1234567890abcdef
branch refs/heads/main
`

	result, err := ParseWorktreeList(input)
	if err != nil {
		t.Fatalf("ParseWorktreeList failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 worktree, got %d", len(result))
	}

	if result[0].Path != "/path/to/repo" {
		t.Errorf("expected path /path/to/repo, got %s", result[0].Path)
	}
}

func TestFindPathByBranch_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// This test requires running in a git repository with worktrees
	// We'll just test that the function doesn't panic
	_, err := FindPathByBranch("/nonexistent", "main")
	if err == nil {
		t.Error("FindPathByBranch should fail for nonexistent repo")
	}
	if !strings.Contains(err.Error(), "exit status") && !strings.Contains(err.Error(), "no such file") {
		t.Logf("FindPathByBranch error (expected): %v", err)
	}
}
