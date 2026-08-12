package web_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/MasayukiFukada/PjDoc/internal/features/watch_changes"
	sharedweb "github.com/MasayukiFukada/PjDoc/internal/shared/web"
)

func TestNewServerIntegration(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tempDir, "README.md"), []byte("# Title"), 0644)

	broadcaster := watch_changes.NewBroadcaster()
	cfg := sharedweb.ServerConfig{
		RootDir:     tempDir,
		Port:        8080,
		Broadcaster: broadcaster,
	}

	mux, err := sharedweb.NewServer(cfg, sharedweb.EmbeddedAssets)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 1. API ルーティングテスト (/api/tree)
	reqTree := httptest.NewRequest(http.MethodGet, "/api/tree", nil)
	recTree := httptest.NewRecorder()
	mux.ServeHTTP(recTree, reqTree)
	if recTree.Code != http.StatusOK {
		t.Errorf("Expected HTTP 200 for /api/tree, got %d", recTree.Code)
	}

	// 2. 埋め込みアセットテスト (/style.css)
	reqAsset := httptest.NewRequest(http.MethodGet, "/style.css", nil)
	recAsset := httptest.NewRecorder()
	mux.ServeHTTP(recAsset, reqAsset)
	if recAsset.Code != http.StatusOK {
		t.Errorf("Expected HTTP 200 for /style.css, got %d", recAsset.Code)
	}
}
