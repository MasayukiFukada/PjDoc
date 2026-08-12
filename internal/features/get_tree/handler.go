package get_tree

import (
	"encoding/json"
	"net/http"
)

// Handler はディレクトリツリー取得リクエストを処理する HTTP ハンドラーでございます。
type Handler struct {
	RootDir string
}

// NewHandler は Handler のインスタンスを生成いたしますわ。
func NewHandler(rootDir string) *Handler {
	return &Handler{RootDir: rootDir}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tree, err := BuildTree(h.RootDir)
	if err != nil {
		http.Error(w, "Failed to build directory tree: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tree)
}
