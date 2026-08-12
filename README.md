# PjDoc

リポジトリ内の Markdown ドキュメントをローカル Web ブラウザ上で高速・快適に閲覧するための CLI ツールです。

---

## 🌟 主な特徴

* **単一バイナリ動作**: Web フロントエンドアセット（HTML/CSS/JS）を Go バイナリ内に埋め込んでいるため、1つの実行ファイルで動作します。
* **ドキュメントツリー & インクリメンタル検索**: ディレクトリ階層表示と、ファイル名・本文のリアルタイム全文検索を備えています。
* **高度な Markdown & ダイヤグラムサポート**: GitHub Flavored Markdown (GFM)、Highlight.js（コードハイライト）、Mermaid.js（図表・シーケンス図）、PlantUML（Kroki / PlantUML サーバー連動）、KaTeX（数式）をサポートします。
* **多彩なカラーテーマ**: Dark, Light, Sepia, Cyberpunk, Auto (System) のテーマ切り替えと設定の自動保存 (LocalStorage) に対応します。
* **ホットリロード (SSE)**: ドキュメントの変更を検知し、ブラウザを自動更新します。
* **Vertical Slice Architecture (VSA)**: 堅牢で保守性の高い機能別スライス構造で設計されています。

---

## 🛠️ 前提条件

* Go 1.22 以上（または [mise](https://mise.jdx.dev/)）

---

## 📦 ビルド & インストール

```bash
# リポジトリのクローン
git clone https://github.com/MasayukiFukada/PjDoc.git
cd PjDoc

# バイナリのビルド
go build -o pjdoc ./cmd/pjdoc

# PATH の通った場所へインストールする場合（任意）
go install ./cmd/pjdoc
```

### 🔄 アップデート方法

最新のソースコードを取得し、再ビルド/再インストールすることで上書きアップデートが行えます。

```bash
cd PjDoc
git pull origin main
go install ./cmd/pjdoc
```

リモートリポジトリから直接最新版に更新する場合：

```bash
go install github.com/MasayukiFukada/PjDoc/cmd/pjdoc@latest
```

---

## 🚀 使い方

閲覧したい Git リポジトリのルートで `pjdoc` コマンドを実行します。

```bash
# カレントディレクトリの Markdown を閲覧（デフォルト: http://localhost:18080）
pjdoc

# ポート番号やディレクトリを明示して実行
pjdoc -port 9000 -dir /path/to/repository

# 自前の PlantUML / Kroki サーバーを指定して実行
pjdoc -plantuml-server http://localhost:8000
```

### 主な CLI オプション

| オプション | デフォルト値 | 説明 |
|---|---|---|
| `-port` | `18080` | Web サーバーのポート番号 |
| `-dir` | `.` | 閲覧対象のルートディレクトリ |
| `-open` | `true` | 起動時にブラウザを自動で開く (`-open=false` で無効化) |
| `-plantuml-server` | `https://kroki.io` | PlantUML / Kroki サーバーのエンドポイント URL |

---

## 📜 ライセンス

MIT License
