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
		util.Verbose("[verbose] base ref が未指定のため自動検出中...")
		if ref, err := git.Git(repoRoot, "symbolic-ref", "-q", "--short", "refs/remotes/origin/HEAD"); err == nil {
			base = strings.TrimSpace(ref)
			util.Verbose("[verbose] origin/HEAD から base ref を検出: %s", base)
		} else {
			base = "origin/main"
			util.Verbose("[verbose] origin/HEAD が見つからないため、デフォルトの base ref を使用: %s", base)
		}
	} else {
		util.Verbose("[verbose] base ref として %s を使用", base)
	}

	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("作成先ディレクトリが既に存在します: %s", target)
	}

	util.Verbose("[verbose] ブランチの存在を確認中...")
	existsLocal := git.GitOk(repoRoot, "show-ref", "--verify", "--quiet", "refs/heads/"+opts.Branch)
	existsRemote := git.GitOk(repoRoot, "show-ref", "--verify", "--quiet", "refs/remotes/origin/"+opts.Branch)
	util.Verbose("[verbose] ローカルブランチ %s: %v", opts.Branch, existsLocal)
	util.Verbose("[verbose] リモートブランチ origin/%s: %v", opts.Branch, existsRemote)

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
		util.Verbose("[verbose] 実行中: %s", util.ShellJoin(a))
		if err := git.Run(a[0], a[1:]...); err != nil {
			return err
		}
		util.Verbose("[verbose] 完了: %s", util.ShellJoin(a))
	}

	// TypeScriptプロジェクトの場合、node_modulesをコピー
	if err := copyNodeModulesIfExists(repoRoot, target); err != nil {
		fmt.Printf("警告: node_modules のコピーに失敗しました: %v\n", err)
	}

	if opts.OpenCmd != "" {
		util.Verbose("[verbose] エディタを開きます: %s %s", opts.OpenCmd, target)
		_ = git.Run(opts.OpenCmd, target)
	}

	return nil
}

// copyNodeModulesIfExists copies node_modules from source to target if it exists
func copyNodeModulesIfExists(repoRoot, target string) error {
	packageJSON := filepath.Join(repoRoot, "package.json")
	util.Verbose("[verbose] package.json の存在を確認中: %s", packageJSON)
	if _, err := os.Stat(packageJSON); err != nil {
		util.Verbose("[verbose] package.json が見つかりません。node_modules のコピーをスキップします")
		// package.jsonがなければスキップ
		return nil
	}

	nodeModules := filepath.Join(repoRoot, "node_modules")
	util.Verbose("[verbose] node_modules の存在を確認中: %s", nodeModules)
	if _, err := os.Stat(nodeModules); err != nil {
		util.Verbose("[verbose] node_modules が見つかりません。コピーをスキップします")
		// node_modulesがなければスキップ
		return nil
	}

	fmt.Printf("\nnode_modules をコピー中...\n")
	targetNodeModules := filepath.Join(target, "node_modules")
	util.Verbose("[verbose] コピー元: %s", nodeModules)
	util.Verbose("[verbose] コピー先: %s", targetNodeModules)

	// cp -a でシンボリックリンクや権限を保持してコピー
	util.Verbose("[verbose] 実行中: cp -a %s %s", nodeModules, targetNodeModules)
	if err := git.Run("cp", "-a", nodeModules, targetNodeModules); err != nil {
		return err
	}

	fmt.Printf("node_modules のコピーが完了しました\n")
	util.Verbose("[verbose] node_modules のコピーが正常に完了しました")
	return nil
}

// List shows worktree list
func List(repoRoot string, opts ListOptions) error {
	util.Verbose("[verbose] worktree の一覧を表示します")
	if opts.Porcelain {
		util.Verbose("[verbose] porcelain モードで出力します")
		return git.Run("git", "-C", repoRoot, "worktree", "list", "--porcelain")
	}
	return git.Run("git", "-C", repoRoot, "worktree", "list")
}

// Prune removes stale worktree references
func Prune(repoRoot string, opts PruneOptions) error {
	util.Verbose("[verbose] 削除済み worktree の参照をクリーンアップします")
	cmd := []string{"git", "-C", repoRoot, "worktree", "prune"}
	if opts.DryRun {
		cmd = append(cmd, "--dry-run")
		util.Verbose("[verbose] dry-run モードが有効です")
	}
	if opts.Verbose {
		cmd = append(cmd, "--verbose")
	}
	util.Verbose("[verbose] 実行中: %s", util.ShellJoin(cmd))
	if err := git.Run(cmd[0], cmd[1:]...); err != nil {
		return err
	}
	util.Verbose("[verbose] クリーンアップが完了しました")
	return nil
}

// Remove deletes a worktree
func Remove(repoRoot string, opts RemoveOptions) error {
	util.Verbose("[verbose] worktree の削除を開始: %s", opts.PathOrBranch)
	targetPath := opts.PathOrBranch

	util.Verbose("[verbose] パスの存在を確認中: %s", targetPath)
	if _, err := os.Stat(targetPath); err != nil {
		util.Verbose("[verbose] パスが見つかりません。ブランチ名として検索します")
		p, err2 := FindPathByBranch(repoRoot, opts.PathOrBranch)
		if err2 != nil {
			return fmt.Errorf("パスでもブランチでも見つかりませんでした: %s", opts.PathOrBranch)
		}
		targetPath = p
		util.Verbose("[verbose] ブランチ %s に対応するパスを発見: %s", opts.PathOrBranch, targetPath)
	} else {
		util.Verbose("[verbose] パスが存在します: %s", targetPath)
	}

	cmd := []string{"git", "-C", repoRoot, "worktree", "remove", targetPath}
	if opts.Force {
		cmd = append(cmd, "--force")
		util.Verbose("[verbose] 強制削除モードが有効です")
	}

	if opts.DryRun {
		fmt.Println("(dry-run) " + util.ShellJoin(cmd))
		return nil
	}

	util.Verbose("[verbose] 実行中: %s", util.ShellJoin(cmd))
	if err := git.Run(cmd[0], cmd[1:]...); err != nil {
		return err
	}
	util.Verbose("[verbose] worktree の削除が完了しました: %s", targetPath)

	return nil
}
