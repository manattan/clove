package util

import "testing"

func TestSanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal", "feature/update", "feature-update"},
		{"with spaces", "feature update", "feature-update"},
		{"with colon", "feature:update", "feature-update"},
		{"with at", "user@branch", "user-branch"},
		{"with slash", "feat/sub/branch", "feat-sub-branch"},
		{"special chars", "feat/update#123", "feat-update-123"},
		{"empty", "", "worktree"},
		{"only special", "///", "worktree"},
		{"trim dash", "-feature-", "feature"},
		{"mixed", "feature/update:v1.0", "feature-update-v1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sanitize(tt.input)
			if result != tt.expected {
				t.Errorf("Sanitize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestShellJoin(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{"simple", []string{"git", "status"}, "git status"},
		{"with space", []string{"git", "commit", "-m", "hello world"}, "git commit -m 'hello world'"},
		{"with quote", []string{"echo", "it's"}, "echo 'it'\\''s'"},
		{"no special", []string{"ls", "-la"}, "ls -la"},
		{"with tab", []string{"echo", "hello\tworld"}, "echo 'hello\tworld'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShellJoin(tt.input)
			if result != tt.expected {
				t.Errorf("ShellJoin(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestQuote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "hello", "'hello'"},
		{"with space", "hello world", "'hello world'"},
		{"with quote", "it's", "'it'\\''s'"},
		{"empty", "", "''"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Quote(tt.input)
			if result != tt.expected {
				t.Errorf("Quote(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
