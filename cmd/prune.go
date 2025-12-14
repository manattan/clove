package cmd

import (
	"fmt"

	"github.com/manattan/clove/internal/git"
	"github.com/manattan/clove/internal/worktree"
	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:   "prune [オプション]",
	Short: "削除済み worktree の参照等を掃除します",
	Args:  cobra.NoArgs,
	RunE:  runPrune,
}

var (
	pruneRepo    string
	pruneDryRun  bool
	pruneVerbose bool
)

func init() {
	pruneCmd.Flags().StringVar(&pruneRepo, "repo", "", "対象リポジトリのパス（省略時: カレントから判定）")
	pruneCmd.Flags().BoolVar(&pruneDryRun, "dry-run", false, "削除される予定のものを表示するだけ（--dry-run）")
	pruneCmd.Flags().BoolVarP(&pruneVerbose, "verbose", "v", false, "詳細表示（--verbose）")
}

func runPrune(cmd *cobra.Command, args []string) error {
	repoRoot := pruneRepo
	if repoRoot == "" {
		var err error
		repoRoot, err = git.GetRepoRoot()
		if err != nil {
			return fmt.Errorf("clove: %w", err)
		}
	}

	opts := worktree.PruneOptions{
		DryRun:  pruneDryRun,
		Verbose: pruneVerbose,
	}

	return worktree.Prune(repoRoot, opts)
}
