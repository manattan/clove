package worktree

import (
	"errors"
	"strings"

	"github.com/manattan/clove/internal/git"
)

// WorktreeInfo represents a worktree entry
type WorktreeInfo struct {
	Path   string
	Branch string
	Head   string
}

// ParseWorktreeList parses git worktree list --porcelain output
func ParseWorktreeList(output string) ([]WorktreeInfo, error) {
	var worktrees []WorktreeInfo
	blocks := strings.Split(output, "\nworktree ")

	for i, b := range blocks {
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}

		// 最初のブロックは "worktree " プレフィックスがあるので削除
		if i == 0 && strings.HasPrefix(b, "worktree ") {
			b = strings.TrimPrefix(b, "worktree ")
		}

		lines := strings.Split(b, "\n")
		wtPath := strings.TrimSpace(lines[0])

		var br, head string
		for _, ln := range lines[1:] {
			ln = strings.TrimSpace(ln)
			if strings.HasPrefix(ln, "branch ") {
				br = strings.TrimSpace(strings.TrimPrefix(ln, "branch "))
			} else if strings.HasPrefix(ln, "HEAD ") {
				head = strings.TrimSpace(strings.TrimPrefix(ln, "HEAD "))
			}
		}

		worktrees = append(worktrees, WorktreeInfo{
			Path:   wtPath,
			Branch: br,
			Head:   head,
		})
	}

	return worktrees, nil
}

// FindPathByBranch finds worktree path by branch name
func FindPathByBranch(repoRoot, branch string) (string, error) {
	out, err := git.Git(repoRoot, "worktree", "list", "--porcelain")
	if err != nil {
		return "", err
	}

	worktrees, err := ParseWorktreeList(out)
	if err != nil {
		return "", err
	}

	targetBranch := "refs/heads/" + branch
	for _, wt := range worktrees {
		if wt.Branch == targetBranch {
			return wt.Path, nil
		}
	}

	return "", errors.New("branch not found")
}
