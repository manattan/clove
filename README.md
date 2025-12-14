# clove

**clove** は git worktree を並列開発向けに扱うための CLI ツールです。

Simplify parallel development with git worktree - A modern CLI tool for managing git worktrees efficiently.

## 特徴 (Features)

- 🚀 **簡単な worktree 作成** - ブランチ名を指定するだけで、リポジトリの隣に worktree ディレクトリを自動作成
- 🎯 **直感的なコマンド** - `add`, `list`, `prune`, `rm` のシンプルな操作
- 🔍 **ブランチ名での削除** - パスだけでなく、ブランチ名でも worktree を削除可能
- 🛠️ **IDE 連携** - `--open` オプションでエディタを自動起動
- ✅ **安全設計** - dry-run モードで事前確認可能

## インストール (Installation)

### Go ツールチェーン経由

```bash
go install github.com/manattan/clove@latest
```

### ソースからビルド

```bash
git clone https://github.com/manattan/clove.git
cd clove
make install
```

バイナリが `$GOBIN` (通常は `~/go/bin`) にインストールされます。

## 使い方 (Usage)

### worktree を作成

```bash
# ブランチ名を指定して worktree を作成
clove add feature/new-ui

# 作成後に VS Code で開く
clove add feature/new-ui --open code

# 特定のブランチから分岐
clove add hotfix/bug-123 --base origin/develop

# dry-run で実行内容を確認
clove add feature/test --dry-run
```

**例**: `~/projects/myapp` で実行すると、`~/projects/myapp-feature-new-ui` が作成されます。

### worktree 一覧を表示

```bash
clove list

# 機械処理しやすい形式で出力
clove list --porcelain
```

### worktree を削除

```bash
# パスを指定して削除
clove rm ../myapp-feature-new-ui

# ブランチ名で削除
clove rm feature/new-ui

# 強制削除
clove rm feature/new-ui --force
```

### 削除済み worktree の参照をクリーンアップ

```bash
clove prune

# dry-run で確認
clove prune --dry-run
```

## コマンド一覧 (Commands)

| コマンド | 説明 |
|---------|------|
| `clove add <ブランチ名>` | worktree を作成 |
| `clove list` | worktree の一覧を表示 |
| `clove prune` | 削除済み worktree の参照を掃除 |
| `clove rm <パス\|ブランチ名>` | worktree を削除 |
| `clove help` | ヘルプを表示 |

各コマンドの詳細は `clove <コマンド> -h` で確認できます。

## オプション (Options)

### `clove add` のオプション

| オプション | 説明 |
|-----------|------|
| `--base <ref>` | 起点にする ref (デフォルト: origin/HEAD) |
| `--prefix <string>` | ディレクトリ名の接頭辞 (デフォルト: リポジトリ名) |
| `--suffix <string>` | ディレクトリ名の接尾辞 |
| `--dir <string>` | ディレクトリ名を明示的に指定 |
| `--open <command>` | 作成後に実行するコマンド (例: `code`, `cursor`) |
| `--dry-run` | 実行せず、実行内容だけ表示 |
| `--no-fetch` | git fetch をスキップ |

## 開発 (Development)

### 必要環境

- Go 1.23.4 以降
- Git

### ビルド

```bash
# ローカルビルド
make build

# インストール
make install

# テスト実行
make test

# フォーマット
make fmt
```

### ディレクトリ構造

```
.
├── cmd/
├── internal/
│   ├── git/         # Git 操作
│   ├── worktree/    # Worktree ビジネスロジック
│   └── util/        # ユーティリティ
├── main.go
└── Makefile
```

## ライセンス (License)

MIT License

## 貢献 (Contributing)

Issue や Pull Request を歓迎します！

---

**Note**: このツールは git worktree の wrapper です。Git 2.15+ が必要です。
