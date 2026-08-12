# 開発者ガイド (DEVELOPMENT.md)

本ドキュメントは、`PjDoc` の開発参加者・メンテナ向けの開発ガイドラインおよび開発手順書です。

---

## 🏗️ アーキテクチャ概要 (Vertical Slice Architecture)

`PjDoc` は **Vertical Slice Architecture (VSA / Feature-First)** に基づいて設計されています。
レイヤー層（Controller/Service/Repository等）による水平分割ではなく、**ユーザー機能（ドメイン行動）単位**で縦にスライスを切り出しています。

### ディレクトリ構造

```text
PjDoc/
├── cmd/
│   └── pjdoc/
│       └── main.go             # CLI エントリーポイント・フラグ解析
├── internal/
│   ├── features/               # Vertical Slices（機能ごとの独立ディレクトリ）
│   │   ├── get_tree/           # 機能1: ドキュメントツリー構築 (tree.go, handler.go, tree_test.go)
│   │   ├── render_document/    # 機能2: Markdown 取得 (document.go, handler.go, document_test.go)
│   │   ├── search_documents/   # 機能3: 全文・ファイル名検索 (search.go, handler.go, search_test.go)
│   │   └── watch_changes/      # 機能4: ホットリロード通知 (watcher.go, handler.go, watcher_test.go)
│   └── shared/                 # スライス間で共有する最小限の基盤機能
│       └── web/                # embedded Webアセットおよび HTTP ルーティング初期化
├── scripts/
│   └── test_coverage.sh       # テスト実行 & カバレッジHTMLレポート自動オープン
└── web/                        # Web フロントエンドアセット (HTML/CSS/JS)
```

### VSA の原則
1. **スライス間の独立性**: `features/` 配下の異なるスライス同士を直接 `import` してはいけません。
2. **ファイルの同居 (Colocation)**: 各機能に関する実装、HTTP ハンドラー、および単体テストコード (`*_test.go`) は同じスライスフォルダ内に配置します。
3. **型の非共有**: リクエスト/レスポンス型や DTO は使い回さず、各スライスが自身の型定義を所有します。

---

## 🛠️ ローカル開発手順

### 前提条件
- Go 1.22 以上
- (任意) [mise](https://mise.jdx.dev/)

### 開発サーバーの起動

`go run` コマンドでローカル開発サーバーを即座に起動できます。

```bash
# カレントディレクトリ対象で起動
go run ./cmd/pjdoc

# ポート番号やディレクトリを指定して起動
go run ./cmd/pjdoc -port 8080 -dir ./docs
```

---

## 🧪 テストの実行とカバレッジ確認

### 🚀 ワンコマンドでテスト実行 ＆ HTMLレポート表示 (推奨)

テストの実行・カバレッジプロファイルの生成・ブラウザでのHTMLレポート表示を一発で行うスクリプトを用意しています。

```bash
./scripts/test_coverage.sh
```

スクリプトを実行すると、全単体テストが実行された後、**自動的にデフォルトブラウザが立ち上がり、カバレッジ結果（緑：テスト通過行 / 赤：未通過行）を視覚的に確認できます。**

---

### 個別コマンドでの実行

#### 1. 単体テストの実行
```bash
go test -v ./...
```

#### 2. カバレッジプロファイルの生成
```bash
go test -coverprofile=coverage.out ./...
```

#### 3. HTML形式でブラウザ表示
```bash
go tool cover -html=coverage.out
```

#### 4. コンソールで関数ごとのカバレッジ割合を出力
```bash
go tool cover -func=coverage.out
```

---

## ➕ 新機能（スライス）の追加方法

新しい機能を追加する場合は、以下の手順でスライスを作成してください。

1. `internal/features/<new_feature_name>/` ディレクトリを作成する。
2. 機能に必要なロジック (`<feature>.go`)、HTTP ハンドラー (`handler.go`)、および単体テスト (`<feature>_test.go`) を同居させて実装する。
3. `internal/shared/web/server.go` 内の `NewServer()` に新しいハンドラーを登録する。
4. `./scripts/test_coverage.sh` を実行し、テストが全件パスし、カバレッジが十分であることを確認する。
