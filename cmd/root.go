package cmd

import (
	"github.com/manattan/clove/internal/util"
	"github.com/spf13/cobra"
)

var (
	globalVerbose bool
)

var rootCmd = &cobra.Command{
	Use:   "clove",
	Short: "git worktree を並列開発向けに扱うためのコマンド",
	Long: `clove: git worktree を並列開発向けに扱うためのコマンド

使い方:
  clove <サブコマンド> [オプション]

サブコマンド:
  add <ブランチ名>     worktree を作成し、指定ブランチをチェックアウトします
  list                worktree の一覧を表示します
  prune               削除済み worktree の参照等を掃除します
  rm <パス|ブランチ>  worktree を削除します（パス指定 or ブランチ名指定）
  help                このヘルプを表示します

例:
  cd ~/manattan/hogehoge
  clove add feature/update
  # => ~/manattan/hogehoge-feature-update が作られる

  clove list
  clove prune --dry-run
  clove rm ../hogehoge-feature-update
  clove rm feature/update

各サブコマンドの詳細:
  clove add   -h
  clove list  -h
  clove prune -h
  clove rm    -h`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&globalVerbose, "verbose", "v", false, "詳細なログを出力します")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		util.SetVerbose(globalVerbose)
	}

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(pruneCmd)
	rootCmd.AddCommand(removeCmd)
}
