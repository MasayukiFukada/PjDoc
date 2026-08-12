package render_document

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Handler は Markdown ファイル読み込みリクエストを安全に処理する HTTP ハンドラーです。
type Handler struct {
	RootDir string
}

// NewHandler は Handler のインスタンスを生成します。
func NewHandler(rootDir string) *Handler {
	return &Handler{RootDir: rootDir}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		http.Error(w, "Missing 'path' query parameter", http.StatusBadRequest)
		return
	}

	doc, err := ReadDocument(h.RootDir, relPath)
	if err != nil {
		if errors.Is(err, ErrInvalidPath) {
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to read document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(doc)
}
