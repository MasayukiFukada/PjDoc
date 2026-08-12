package search_documents

import (
	"encoding/json"
	"net/http"
)

// Handler は全文検索リクエストを処理する HTTP ハンドラーでございます。
type Handler struct {
	RootDir string
}

// NewHandler は Handler インスタンスを返却いたしますわ。
func NewHandler(rootDir string) *Handler {
	return &Handler{RootDir: rootDir}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	results, err := SearchDocuments(h.RootDir, query)
	if err != nil {
		http.Error(w, "Search error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(results)
}
