package git

import (
	"strings"
	"testing"
)

func TestGit(t *testing.T) {
	// Test that git command can be executed (requires git to be installed)
	out, err := Git("", "version")
	if err != nil {
		t.Fatalf("Git version failed: %v", err)
	}
	if !strings.Contains(out, "git version") {
		t.Errorf("Git version output doesn't contain 'git version': %s", out)
	}
}

func TestGitOk(t *testing.T) {
	// Test with a command that should succeed
	if !GitOk("", "version") {
		t.Error("GitOk(version) should return true")
	}

	// Test with a command that should fail
	if GitOk("", "invalid-command-that-does-not-exist") {
		t.Error("GitOk(invalid-command) should return false")
	}
}

func TestGetRepoRoot(t *testing.T) {
	// This test requires running in a git repository
	root, err := GetRepoRoot()
	if err != nil {
		t.Skipf("Not in a git repository, skipping: %v", err)
	}
	if root == "" {
		t.Error("GetRepoRoot returned empty string")
	}
	if !strings.HasPrefix(root, "/") {
		t.Errorf("GetRepoRoot should return absolute path, got: %s", root)
	}
}

func TestGetOriginHead(t *testing.T) {
	// This test requires running in a git repository
	root, err := GetRepoRoot()
	if err != nil {
		t.Skipf("Not in a git repository, skipping: %v", err)
	}

	head, err := GetOriginHead(root)
	if err != nil {
		t.Fatalf("GetOriginHead failed: %v", err)
	}

	// Should either return origin/HEAD or fallback to origin/main
	if head == "" {
		t.Error("GetOriginHead returned empty string")
	}
	if !strings.HasPrefix(head, "origin/") {
		t.Errorf("GetOriginHead should return origin/* ref, got: %s", head)
	}
}
