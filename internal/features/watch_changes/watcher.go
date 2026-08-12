package watch_changes

import (
	"sync"
)

// Broadcaster は接続されたクライアントにホットリロードイベントを通知する構造体です。
type Broadcaster struct {
	mu      sync.Mutex
	clients map[chan string]bool
}

// NewBroadcaster は Broadcaster インスタンスを作成します。
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[chan string]bool),
	}
}

// Register は新しいクライアントチャネルを登録します。
func (b *Broadcaster) Register(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[ch] = true
}

// Unregister はクライアントチャネルの登録を解除します。
func (b *Broadcaster) Unregister(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clients, ch)
}

// NotifyReload は全接続クライアントへリロードシグナルを送信します。
func (b *Broadcaster) NotifyReload() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.clients {
		select {
		case ch <- "reload":
		default:
		}
	}
}
