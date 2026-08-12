package watch_changes

import (
	"sync"
)

// Broadcaster は接続されたクライアントにホットリロードイベントを通知する仕組みでございます。
type Broadcaster struct {
	mu      sync.Mutex
	clients map[chan string]bool
}

// NewBroadcaster は Broadcaster インスタンスを作成いたしますわ。
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[chan string]bool),
	}
}

// Register は新しいクライアントチャネルを登録いたします。
func (b *Broadcaster) Register(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[ch] = true
}

// Unregister はクライアントチャネルの登録を解除いたしますわ。
func (b *Broadcaster) Unregister(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clients, ch)
}

// NotifyReload は全接続クライアントへリロードシグナルを美しく送信いたします。
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
