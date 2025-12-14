package git

import (
	"errors"
	"strings"
)

// GetRepoRoot returns the repository root path
func GetRepoRoot() (string, error) {
	out, err := Git("", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	r := strings.TrimSpace(out)
	if r == "" {
		return "", errors.New("gitリポジトリではありません")
	}
	return r, nil
}

// GetOriginHead returns the default branch (origin/HEAD)
// Falls back to "origin/main" if origin/HEAD is not set
func GetOriginHead(repoRoot string) (string, error) {
	ref, err := Git(repoRoot, "symbolic-ref", "-q", "--short", "refs/remotes/origin/HEAD")
	if err == nil {
		return strings.TrimSpace(ref), nil
	}
	// Fallback to origin/main
	return "origin/main", nil
}
