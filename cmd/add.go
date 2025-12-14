package cmd

import (
	"fmt"

	"github.com/manattan/clove/internal/git"
	"github.com/manattan/clove/internal/worktree"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [オプション] <ブランチ名>",
	Short: "worktree を作成し、指定ブランチをチェックアウトします",
	Long: `現在いるリポジトリの「隣」に worktree 用ディレクトリを作成します。
例: ~/hogehoge で実行 → ~/hogehoge-<branch> が作られる

例:
  clove add feature/update
  clove add -open code feature/update
  clove add -base origin/develop feature/update`,
	Args: cobra.ExactArgs(1),
	RunE: runAdd,
}

var (
	addBaseRef   string
	addPrefix    string
	addSuffix    string
	addOpenCmd   string
	addDryRun    bool
	addForceName string
	addNoFetch   bool
)

func init() {
	addCmd.Flags().StringVar(&addBaseRef, "base", "", "起点にするref（省略時: origin/HEAD を試し、ダメなら origin/main）")
	addCmd.Flags().StringVar(&addPrefix, "prefix", "", "作成するディレクトリ名の接頭辞（省略時: リポジトリ名）")
	addCmd.Flags().StringVar(&addSuffix, "suffix", "", "作成するディレクトリ名の接尾辞（任意）")
	addCmd.Flags().StringVar(&addOpenCmd, "open", "", "作成後にディレクトリを開くコマンド（例: code / cursor / open）")
	addCmd.Flags().BoolVar(&addDryRun, "dry-run", false, "実行せず、実行内容だけ表示します")
	addCmd.Flags().StringVar(&addForceName, "dir", "", "ディレクトリ名を明示します（repoの親ディレクトリ配下に作る）")
	addCmd.Flags().BoolVar(&addNoFetch, "no-fetch", false, "git fetch origin をスキップします")
}

func runAdd(cmd *cobra.Command, args []string) error {
	branch := args[0]

	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("clove: %w", err)
	}

	opts := worktree.AddOptions{
		Branch:    branch,
		BaseRef:   addBaseRef,
		Prefix:    addPrefix,
		Suffix:    addSuffix,
		ForceName: addForceName,
		OpenCmd:   addOpenCmd,
		DryRun:    addDryRun,
		NoFetch:   addNoFetch,
	}

	return worktree.Add(repoRoot, opts)
}
