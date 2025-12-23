package worktree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manattan/clove/internal/git"
	"github.com/manattan/clove/internal/util"
)

// AddOptions contains options for Add operation
type AddOptions struct {
	Branch    string
	BaseRef   string
	Prefix    string
	Suffix    string
	ForceName string
	OpenCmd   string
	DryRun    bool
	NoFetch   bool
}

// RemoveOptions contains options for Remove operation
type RemoveOptions struct {
	PathOrBranch string
	Force        bool
	DryRun       bool
}

// PruneOptions contains options for Prune operation
type PruneOptions struct {
	DryRun  bool
	Verbose bool
}

// ListOptions contains options for List operation
type ListOptions struct {
	Porcelain bool
}

// Add creates a new worktree
func Add(repoRoot string, opts AddOptions) error {
	parent := filepath.Dir(repoRoot)
	repoName := filepath.Base(repoRoot)

	dirName := opts.ForceName
	if dirName == "" {
		p := opts.Prefix
		if p == "" {
			p = repoName
		}
		dirName = fmt.Sprintf("%s-%s%s", p, util.Sanitize(opts.Branch), opts.Suffix)
	}
	target := filepath.Join(parent, dirName)

	base := opts.BaseRef
	if base == "" {
		if ref, err := git.Git(repoRoot, "symbolic-ref", "-q", "--short", "refs/remotes/origin/HEAD"); err == nil {
			base = strings.TrimSpace(ref)
		} else {
			base = "origin/main"
		}
	}

	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("作成先ディレクトリが既に存在します: %s", target)
	}

	existsLocal := git.GitOk(repoRoot, "show-ref", "--verify", "--quiet", "refs/heads/"+opts.Branch)
	existsRemote := git.GitOk(repoRoot, "show-ref", "--verify", "--quiet", "refs/remotes/origin/"+opts.Branch)

	var actions [][]string
	if !opts.NoFetch {
		actions = append(actions, []string{"git", "-C", repoRoot, "fetch", "--prune", "origin"})
	}

	var wtCmd []string
	switch {
	case existsLocal:
		wtCmd = []string{"git", "-C", repoRoot, "worktree", "add", target, opts.Branch}
	case existsRemote:
		wtCmd = []string{"git", "-C", repoRoot, "worktree", "add", target, "-b", opts.Branch, "origin/" + opts.Branch}
	default:
		wtCmd = []string{"git", "-C", repoRoot, "worktree", "add", target, "-b", opts.Branch, base}
	}
	actions = append(actions, wtCmd)

	fmt.Printf("repo:   %s\n", repoRoot)
	fmt.Printf("base:   %s\n", base)
	fmt.Printf("branch: %s\n", opts.Branch)
	fmt.Printf("dir:    %s\n", target)

	if opts.DryRun {
		fmt.Println("\n(dry-run) 実行予定コマンド:")
		for _, a := range actions {
			fmt.Println("  " + util.ShellJoin(a))
		}
		return nil
	}

	for _, a := range actions {
		if err := git.Run(a[0], a[1:]...); err != nil {
			return err
		}
	}

	// TypeScriptプロジェクトの場合、node_modulesをコピー
	if err := copyNodeModulesIfExists(repoRoot, target); err != nil {
		fmt.Printf("警告: node_modules のコピーに失敗しました: %v\n", err)
	}

	if opts.OpenCmd != "" {
		_ = git.Run(opts.OpenCmd, target)
	}

	return nil
}

// copyNodeModulesIfExists copies node_modules from source to target if it exists
func copyNodeModulesIfExists(repoRoot, target string) error {
	packageJSON := filepath.Join(repoRoot, "package.json")
	if _, err := os.Stat(packageJSON); err != nil {
		// package.jsonがなければスキップ
		return nil
	}

	nodeModules := filepath.Join(repoRoot, "node_modules")
	if _, err := os.Stat(nodeModules); err != nil {
		// node_modulesがなければスキップ
		return nil
	}

	fmt.Printf("\nnode_modules をコピー中...\n")
	targetNodeModules := filepath.Join(target, "node_modules")

	// cp -a でシンボリックリンクや権限を保持してコピー
	if err := git.Run("cp", "-a", nodeModules, targetNodeModules); err != nil {
		return err
	}

	fmt.Printf("node_modules のコピーが完了しました\n")
	return nil
}

// List shows worktree list
func List(repoRoot string, opts ListOptions) error {
	if opts.Porcelain {
		return git.Run("git", "-C", repoRoot, "worktree", "list", "--porcelain")
	}
	return git.Run("git", "-C", repoRoot, "worktree", "list")
}

// Prune removes stale worktree references
func Prune(repoRoot string, opts PruneOptions) error {
	cmd := []string{"git", "-C", repoRoot, "worktree", "prune"}
	if opts.DryRun {
		cmd = append(cmd, "--dry-run")
	}
	if opts.Verbose {
		cmd = append(cmd, "--verbose")
	}
	return git.Run(cmd[0], cmd[1:]...)
}

// Remove deletes a worktree
func Remove(repoRoot string, opts RemoveOptions) error {
	targetPath := opts.PathOrBranch
	if _, err := os.Stat(targetPath); err != nil {
		p, err2 := FindPathByBranch(repoRoot, opts.PathOrBranch)
		if err2 != nil {
			return fmt.Errorf("パスでもブランチでも見つかりませんでした: %s", opts.PathOrBranch)
		}
		targetPath = p
	}

	cmd := []string{"git", "-C", repoRoot, "worktree", "remove", targetPath}
	if opts.Force {
		cmd = append(cmd, "--force")
	}

	if opts.DryRun {
		fmt.Println("(dry-run) " + util.ShellJoin(cmd))
		return nil
	}

	return git.Run(cmd[0], cmd[1:]...)
}
