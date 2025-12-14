package cmd

import (
	"fmt"

	"github.com/manattan/clove/internal/git"
	"github.com/manattan/clove/internal/worktree"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove [オプション] <パス|ブランチ名>",
	Aliases: []string{"rm"},
	Short:   "worktree を削除します（パス指定 or ブランチ名指定）",
	Long: `引数が存在するパスならその worktree を削除します。
パスとして存在しない場合はブランチ名として解釈し、worktree 一覧から紐づくパスを探して削除します。

例:
  clove rm ../hogehoge-feature-update
  clove rm feature/update`,
	Args: cobra.ExactArgs(1),
	RunE: runRemove,
}

var (
	removeRepo   string
	removeForce  bool
	removeDryRun bool
)

func init() {
	removeCmd.Flags().StringVar(&removeRepo, "repo", "", "対象リポジトリのパス（省略時: カレントから判定）")
	removeCmd.Flags().BoolVar(&removeForce, "force", false, "強制削除（git worktree remove --force）")
	removeCmd.Flags().BoolVar(&removeDryRun, "dry-run", false, "実行せず、実行内容だけ表示します")
}

func runRemove(cmd *cobra.Command, args []string) error {
	pathOrBranch := args[0]

	repoRoot := removeRepo
	if repoRoot == "" {
		var err error
		repoRoot, err = git.GetRepoRoot()
		if err != nil {
			return fmt.Errorf("clove: %w", err)
		}
	}

	opts := worktree.RemoveOptions{
		PathOrBranch: pathOrBranch,
		Force:        removeForce,
		DryRun:       removeDryRun,
	}

	return worktree.Remove(repoRoot, opts)
}
