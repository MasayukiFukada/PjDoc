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
	tempDir := t.TempDir()

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

	if len(tree.Children) != 2 {
		t.Fatalf("Expected 2 children (docs directory & README.md), got %d", len(tree.Children))
	}

	if !tree.Children[0].IsDir || tree.Children[0].Name != "docs" {
		t.Errorf("Expected first child to be directory 'docs', got %s (isDir=%v)", tree.Children[0].Name, tree.Children[0].IsDir)
	}

	if tree.Children[1].IsDir || tree.Children[1].Name != "README.md" {
		t.Errorf("Expected second child to be file 'README.md', got %s (isDir=%v)", tree.Children[1].Name, tree.Children[1].IsDir)
	}

	docsNode := tree.Children[0]
	if len(docsNode.Children) != 1 || docsNode.Children[0].Name != "guide.md" {
		t.Errorf("Expected docs node to contain 1 child 'guide.md', got %v", docsNode.Children)
	}
}

func TestBuildTree_InvalidDir(t *testing.T) {
	_, err := get_tree.BuildTree("/invalid/non_existent_directory_pjdoc")
	if err == nil {
		t.Error("Expected error for non existent directory, got nil")
	}
}

func TestHandler(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tempDir, "index.md"), []byte("# Index"), 0644)

	handler := get_tree.NewHandler(tempDir)
	req := httptest.NewRequest(http.MethodGet, "/api/tree", nil)
	rec := httptest.NewRecorder()

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

func TestHandler_MethodNotAllowed(t *testing.T) {
	handler := get_tree.NewHandler(".")
	req := httptest.NewRequest(http.MethodPost, "/api/tree", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected HTTP 405 Method Not Allowed, got %d", rec.Code)
	}
}

func TestHandler_BuildTreeError(t *testing.T) {
	handler := get_tree.NewHandler("/invalid/non_existent_directory_pjdoc")
	req := httptest.NewRequest(http.MethodGet, "/api/tree", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected HTTP 500 Internal Server Error, got %d", rec.Code)
	}
}
