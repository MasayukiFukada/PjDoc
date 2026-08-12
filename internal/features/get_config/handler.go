package get_config

import (
	"encoding/json"
	"net/http"
)

// Handler は /api/config リクエストを処理する HTTP ハンドラーです。
type Handler struct {
	service *ConfigService
}

// NewHandler は新しい Handler インスタンスを生成します。
func NewHandler(service *ConfigService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	config := h.service.GetConfig()
	_ = json.NewEncoder(w).Encode(config)
}
