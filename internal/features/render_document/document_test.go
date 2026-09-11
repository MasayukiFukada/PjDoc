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

func TestReadDocument_DirectoryAsFile(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "sub_dir")
	_ = os.Mkdir(subDir, 0755)

	_, err := render_document.ReadDocument(tempDir, "sub_dir")
	if err != render_document.ErrInvalidPath {
		t.Errorf("Expected ErrInvalidPath when trying to read directory, got %v", err)
	}
}

func TestReadDocument_SecurityPathTraversal(t *testing.T) {
	tempDir := t.TempDir()

	_, err := render_document.ReadDocument(tempDir, "../secret.txt")
	if err == nil {
		t.Fatal("Expected error for path traversal attempt, got nil")
	}

	_, err = render_document.ReadDocument(tempDir, "/etc/passwd")
	if err == nil {
		t.Fatal("Expected error for absolute path attempt, got nil")
	}

	// 外部を指すシンボリックリンクによるトラバーサル攻撃テスト
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "secret.txt")
	_ = os.WriteFile(outsideFile, []byte("sensitive-secret-token"), 0644)

	maliciousSymlink := filepath.Join(tempDir, "evil_symlink.md")
	if err := os.Symlink(outsideFile, maliciousSymlink); err == nil {
		_, err = render_document.ReadDocument(tempDir, "evil_symlink.md")
		if err != render_document.ErrInvalidPath {
			t.Fatalf("Expected ErrInvalidPath for external symlink attack, got %v", err)
		}
	}

	// 内部ファイルを指す正常なシンボリックリンクの動作テスト
	validFile := filepath.Join(tempDir, "valid.md")
	_ = os.WriteFile(validFile, []byte("Valid content"), 0644)

	validSymlink := filepath.Join(tempDir, "valid_symlink.md")
	if err := os.Symlink(validFile, validSymlink); err == nil {
		doc, err := render_document.ReadDocument(tempDir, "valid_symlink.md")
		if err != nil {
			t.Fatalf("Expected successful read for internal symlink, got %v", err)
		}
		if doc.Content != "Valid content" {
			t.Errorf("Expected 'Valid content', got %q", doc.Content)
		}
	}
}

func TestHandler(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "doc.md")
	_ = os.WriteFile(filePath, []byte("Content"), 0644)

	handler := render_document.NewHandler(tempDir)

	// 1. 正常系
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

	// 2. 404 Not Found ケース
	reqNotFound := httptest.NewRequest(http.MethodGet, "/api/document?path=missing.md", nil)
	recNotFound := httptest.NewRecorder()
	handler.ServeHTTP(recNotFound, reqNotFound)

	if recNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected HTTP 404, got %d", recNotFound.Code)
	}

	// 3. 400 Bad Request (クエリ欠落) ケース
	reqMissing := httptest.NewRequest(http.MethodGet, "/api/document", nil)
	recMissing := httptest.NewRecorder()
	handler.ServeHTTP(recMissing, reqMissing)

	if recMissing.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP 400 Bad Request for missing path, got %d", recMissing.Code)
	}

	// 4. 400 Bad Request (不正パス) ケース
	reqInvalid := httptest.NewRequest(http.MethodGet, "/api/document?path=../secret.txt", nil)
	recInvalid := httptest.NewRecorder()
	handler.ServeHTTP(recInvalid, reqInvalid)

	if recInvalid.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP 400 Bad Request for invalid path, got %d", recInvalid.Code)
	}

	// 5. 405 Method Not Allowed ケース
	reqPost := httptest.NewRequest(http.MethodPost, "/api/document?path=doc.md", nil)
	recPost := httptest.NewRecorder()
	handler.ServeHTTP(recPost, reqPost)

	if recPost.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected HTTP 405 Method Not Allowed, got %d", recPost.Code)
	}
}
