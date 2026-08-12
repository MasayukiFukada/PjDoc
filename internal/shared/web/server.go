package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/MasayukiFukada/PjDoc/internal/features/get_tree"
	"github.com/MasayukiFukada/PjDoc/internal/features/render_document"
	"github.com/MasayukiFukada/PjDoc/internal/features/search_documents"
	"github.com/MasayukiFukada/PjDoc/internal/features/watch_changes"
)

// EmbeddedAssets には `web/` 内の全フロントエンドファイルを埋め込みます。
//
//go:embed all:assets
var EmbeddedAssets embed.FS

// ServerConfig は Web サーバー構築の基本設定構造体です。
type ServerConfig struct {
	RootDir     string
	Port        int
	Broadcaster *watch_changes.Broadcaster
}

// NewServer は各 Vertical Slice のハンドラーをルーティングした HTTP サーバーを組み立てます。
func NewServer(cfg ServerConfig, assetsFS embed.FS) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	// 1. 各スライスのハンドラー初期化
	treeHandler := get_tree.NewHandler(cfg.RootDir)
	docHandler := render_document.NewHandler(cfg.RootDir)
	searchHandler := search_documents.NewHandler(cfg.RootDir)
	watchHandler := watch_changes.NewHandler(cfg.Broadcaster)

	// 2. API ルーティング設定
	mux.Handle("/api/tree", treeHandler)
	mux.Handle("/api/document", docHandler)
	mux.Handle("/api/search", searchHandler)
	mux.Handle("/api/events", watchHandler)

	// 3. 埋め込みアセット (HTML/CSS/JS) の配信ハンドラー
	subFS, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		return nil, fmt.Errorf("failed to create sub filesystem: %w", err)
	}

	fileServer := http.FileServer(http.FS(subFS))
	mux.Handle("/", fileServer)

	return mux, nil
}
