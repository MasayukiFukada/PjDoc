package render_document_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/MasayukiFukada/PjDoc/internal/features/render_document"
)

func TestReadDocument_Success(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sample.md")
	content := "# Hello World\nThis is a test document."
	_ = os.WriteFile(filePath, []byte(content), 0644)

	doc, err := render_document.ReadDocument(tempDir, "sample.md")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if doc.Content != content {
		t.Errorf("Expected content %q, got %q", content, doc.Content)
	}
	if doc.Path != "sample.md" {
		t.Errorf("Expected path 'sample.md', got %q", doc.Path)
	}
}

func TestReadDocument_SecurityPathTraversal(t *testing.T) {
	tempDir := t.TempDir()

	// ルート外へのアクセストラバーストライ
	_, err := render_document.ReadDocument(tempDir, "../secret.txt")
	if err == nil {
		t.Fatal("Expected error for path traversal attempt, got nil")
	}

	// 絶対パス指定のトライ
	_, err = render_document.ReadDocument(tempDir, "/etc/passwd")
	if err == nil {
		t.Fatal("Expected error for absolute path attempt, got nil")
	}
}

func TestHandler(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "doc.md")
	_ = os.WriteFile(filePath, []byte("Content"), 0644)

	handler := render_document.NewHandler(tempDir)

	// 正常系
	req := httptest.NewRequest(http.MethodGet, "/api/document?path=doc.md", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d", rec.Code)
	}

	var doc render_document.Document
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if doc.Content != "Content" {
		t.Errorf("Expected Content, got %s", doc.Content)
	}

	// 404 Not Found ケース
	reqNotFound := httptest.NewRequest(http.MethodGet, "/api/document?path=missing.md", nil)
	recNotFound := httptest.NewRecorder()
	handler.ServeHTTP(recNotFound, reqNotFound)

	if recNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected HTTP 404, got %d", recNotFound.Code)
	}
}
