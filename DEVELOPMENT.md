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
│   │   ├── get_config/         # 機能1: クライアント設定提供 (config.go, handler.go, config_test.go)
│   │   ├── get_tree/           # 機能2: ドキュメントツリー構築 (tree.go, handler.go, tree_test.go)
│   │   ├── render_document/    # 機能3: Markdown 取得 (document.go, handler.go, document_test.go)
│   │   ├── search_documents/   # 機能4: 全文・ファイル名検索 (search.go, handler.go, search_test.go)
│   │   ├── show_version/       # 機能5: バージョン情報保持・補完 (version.go, version_test.go)
│   │   └── watch_changes/      # 機能6: ホットリロード通知 (watcher.go, handler.go, watcher_test.go)
│   └── shared/                 # スライス間で共有する最小限の基盤機能
│       └── web/                # embedded Webアセットおよび HTTP ルーティング初期化
├── scripts/
│   ├── build.sh               # バージョン情報埋め込みバイナリ構築スクリプト
│   ├── install.sh             # バージョン情報埋め込みバイナリインストールスクリプト
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
- [mise](https://mise.jdx.dev/) (推奨) または Go 1.22 以上

### 🚀 開発環境のセットアップ (mise を使用する場合)

本プロジェクトでは [`.mise.toml`](.mise.toml) を用意しており、`mise` を使用することで必要なツール（Go 等）のバージョンをワンコマンドで自動インストールできます。

#### 1. mise のインストール（未導入の場合）
```bash
# macOS / Linux (その他のインストール方法は公式ドキュメントを参照)
curl https://mise.run | sh
```

#### 2. 依存ツールの自動インストール
リポジトリルートで以下を実行してください。

```bash
# リポジトリの設定を信頼して必要なツールをインストール
mise trust
mise install
```

---

### 開発サーバーの起動

`go run` コマンドでローカル開発サーバーを即座に起動できます。

```bash
# カレントディレクトリ対象で起動 (デフォルトで PlantUML / Kroki API を使用)
go run ./cmd/pjdoc

# ポート番号や対象ディレクトリを指定して起動
go run ./cmd/pjdoc -port 19000 -dir ./docs

# 自前の PlantUML / Kroki ローカルサーバーを指定して起動
go run ./cmd/pjdoc -plantuml-server http://localhost:8000
```

---

### 📦 ビルド & インストールスクリプト

Git タグ・コミットハッシュ・ビルド日時（ローカルタイムゾーン）を自動注入してバイナリを作成・配置するスクリプトを用意しています。

```bash
# 1. ローカルビルド (カレントディレクトリに ./pjdoc を生成)
./scripts/build.sh           # 自動タグ/コミットハッシュ取得
./scripts/build.sh v1.0.0     # バージョン手動指定

# 2. インストール (GOPATH/GOBIN 配下に pjdoc を配置)
./scripts/install.sh         # 自動タグ/コミットハッシュ取得
./scripts/install.sh v1.0.0   # バージョン手動指定
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
