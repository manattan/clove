package cmd

import (
	"fmt"

	"github.com/manattan/clove/internal/git"
	"github.com/manattan/clove/internal/worktree"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [オプション]",
	Short: "worktree の一覧を表示します",
	Args:  cobra.NoArgs,
	RunE:  runList,
}

var (
	listRepo      string
	listPorcelain bool
)

func init() {
	listCmd.Flags().StringVar(&listRepo, "repo", "", "対象リポジトリのパス（省略時: カレントから判定）")
	listCmd.Flags().BoolVar(&listPorcelain, "porcelain", false, "機械処理しやすい形式（--porcelain）で表示します")
}

func runList(cmd *cobra.Command, args []string) error {
	repoRoot := listRepo
	if repoRoot == "" {
		var err error
		repoRoot, err = git.GetRepoRoot()
		if err != nil {
			return fmt.Errorf("clove: %w", err)
		}
	}

	opts := worktree.ListOptions{
		Porcelain: listPorcelain,
	}

	return worktree.List(repoRoot, opts)
}
