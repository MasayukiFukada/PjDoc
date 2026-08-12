package search_documents_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/MasayukiFukada/PjDoc/internal/features/search_documents"
)

func TestSearchDocuments(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "architecture.md")
	_ = os.WriteFile(file1, []byte("# Architecture\nThis document describes Vertical Slice Architecture."), 0644)

	file2 := filepath.Join(tempDir, "guide.md")
	_ = os.WriteFile(file2, []byte("# User Guide\nFollow the instructions here."), 0644)

	// 'Vertical' キーワード検索
	results, err := search_documents.SearchDocuments(tempDir, "Vertical")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result for 'Vertical', got %d", len(results))
	}
	if results[0].Path != "architecture.md" || results[0].LineNumber != 2 {
		t.Errorf("Unexpected search result: %+v", results[0])
	}

	// ファイル名検索 'guide'
	fileNameResults, err := search_documents.SearchDocuments(tempDir, "guide")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(fileNameResults) < 1 {
		t.Fatalf("Expected results for filename search 'guide', got 0")
	}

	foundTitle := false
	for _, res := range fileNameResults {
		if res.IsTitle && res.Path == "guide.md" {
			foundTitle = true
			break
		}
	}
	if !foundTitle {
		t.Errorf("Expected filename title match for guide.md")
	}
}

func TestHandler(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tempDir, "test.md"), []byte("Special keyword inside"), 0644)

	handler := search_documents.NewHandler(tempDir)
	req := httptest.NewRequest(http.MethodGet, "/api/search?q=Special", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d", rec.Code)
	}

	var results []*search_documents.MatchResult
	if err := json.Unmarshal(rec.Body.Bytes(), &results); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if len(results) != 1 || results[0].Path != "test.md" {
		t.Errorf("Unexpected search handler response: %+v", results)
	}
}
