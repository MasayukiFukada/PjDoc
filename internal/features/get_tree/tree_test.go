package get_tree_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/MasayukiFukada/PjDoc/internal/features/get_tree"
)

func TestBuildTree(t *testing.T) {
	// テスト用の臨時ディレクトリ構造を美しく準備いたしますわ
	tempDir := t.TempDir()

	// 構造:
	// tempDir/
	// ├── README.md
	// ├── docs/
	// │   ├── guide.md
	// │   └── ignore.txt
	// ├── .git/
	// │   └── config.md (隠しフォルダ内なので除外対象)
	// └── empty_dir/ (Markdownを含まないディレクトリなので除外対象)

	_ = os.WriteFile(filepath.Join(tempDir, "README.md"), []byte("# Root"), 0644)
	docsDir := filepath.Join(tempDir, "docs")
	_ = os.Mkdir(docsDir, 0755)
	_ = os.WriteFile(filepath.Join(docsDir, "guide.md"), []byte("# Guide"), 0644)
	_ = os.WriteFile(filepath.Join(docsDir, "ignore.txt"), []byte("txt"), 0644)

	gitDir := filepath.Join(tempDir, ".git")
	_ = os.Mkdir(gitDir, 0755)
	_ = os.WriteFile(filepath.Join(gitDir, "config.md"), []byte("# Hidden"), 0644)

	emptyDir := filepath.Join(tempDir, "empty_dir")
	_ = os.Mkdir(emptyDir, 0755)

	tree, err := get_tree.BuildTree(tempDir)
	if err != nil {
		t.Fatalf("BuildTree failed with error: %v", err)
	}

	if tree == nil {
		t.Fatal("Expected tree node, got nil")
	}

	// ディレクトリツリー検証
	if len(tree.Children) != 2 {
		t.Fatalf("Expected 2 children (docs directory & README.md), got %d", len(tree.Children))
	}

	// ソート順序の検証: ディレクトリが先頭、その後にファイルが参ります
	if !tree.Children[0].IsDir || tree.Children[0].Name != "docs" {
		t.Errorf("Expected first child to be directory 'docs', got %s (isDir=%v)", tree.Children[0].Name, tree.Children[0].IsDir)
	}

	if tree.Children[1].IsDir || tree.Children[1].Name != "README.md" {
		t.Errorf("Expected second child to be file 'README.md', got %s (isDir=%v)", tree.Children[1].Name, tree.Children[1].IsDir)
	}

	// docs ディレクトリ内の要素検証
	docsNode := tree.Children[0]
	if len(docsNode.Children) != 1 || docsNode.Children[0].Name != "guide.md" {
		t.Errorf("Expected docs node to contain 1 child 'guide.md', got %v", docsNode.Children)
	}
}

func TestHandler(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tempDir, "index.md"), []byte("# Index"), 0644)

	handler := get_tree.NewHandler(tempDir)
	req := httptest.NewRequest(http.MethodGet, "/api/tree", nil)
	rec := httptest.RecordHeaderBuffer(httptest.NewRecorder())

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status HTTP 200, got %d", rec.Code)
	}

	var node get_tree.Node
	if err := json.Unmarshal(rec.Body.Bytes(), &node); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if len(node.Children) != 1 || node.Children[0].Name != "index.md" {
		t.Errorf("Unexpected node children: %v", node.Children)
	}
}
