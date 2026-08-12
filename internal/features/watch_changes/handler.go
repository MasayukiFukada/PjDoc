package watch_changes

import (
	"fmt"
	"net/http"
)

// Handler は SSE 接続を受け取る HTTP ハンドラーです。
type Handler struct {
	Broadcaster *Broadcaster
}

// NewHandler は Handler のインスタンスを生成します。
func NewHandler(b *Broadcaster) *Handler {
	return &Handler{Broadcaster: b}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	messageChan := make(chan string, 10)
	h.Broadcaster.Register(messageChan)
	defer h.Broadcaster.Unregister(messageChan)

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-messageChan:
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
