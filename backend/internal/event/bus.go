package event

import (
	"context"
	"log"
	"sync"
)

// Event mewakili suatu kejadian dalam sistem
type Event struct {
	Type string
	Data interface{}
}

// Handler adalah fungsi yang menangani event
type Handler func(ctx context.Context, event Event) error

// Bus adalah event bus internal dengan worker pool
type Bus struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
	eventCh  chan Event
	workers  int
	wg       sync.WaitGroup
	closeCh  chan struct{}
	closed   bool
	closeMu  sync.Mutex
}

// NewBus membuat event bus baru dengan jumlah worker tertentu
func NewBus(workers int) *Bus {
	b := &Bus{
		handlers: make(map[string][]Handler),
		eventCh:  make(chan Event, 100),
		workers:  workers,
		closeCh:  make(chan struct{}),
	}
	b.startWorkers()
	return b
}

// Subscribe mendaftarkan handler untuk tipe event tertentu
func (b *Bus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish mengirim event ke bus (non-blocking, akan di-drop jika channel penuh)
func (b *Bus) Publish(ctx context.Context, event Event) {
	b.closeMu.Lock()
	closed := b.closed
	b.closeMu.Unlock()
	if closed {
		log.Printf("event bus closed, dropping event: %s", event.Type)
		return
	}
	select {
	case b.eventCh <- event:
	default:
		log.Printf("event channel full, dropping event: %s", event.Type)
	}
}

// startWorkers menjalankan worker goroutine
func (b *Bus) startWorkers() {
	for i := 0; i < b.workers; i++ {
		b.wg.Add(1)
		go b.worker()
	}
}

// worker mengambil event dari channel dan menjalankan handler yang terdaftar
func (b *Bus) worker() {
	defer b.wg.Done()
	for {
		select {
		case event := <-b.eventCh:
			b.mu.RLock()
			handlers := b.handlers[event.Type]
			b.mu.RUnlock()
			for _, h := range handlers {
				// Jalankan handler secara asinkron agar tidak memblokir worker lain
				go func(handler Handler, ev Event) {
					if err := handler(context.Background(), ev); err != nil {
						log.Printf("event handler error: %v", err)
					}
				}(h, event)
			}
		case <-b.closeCh:
			return
		}
	}
}

// Stop menghentikan semua worker dan menutup bus. Tidak ada event baru yang diproses.
func (b *Bus) Stop() {
	b.closeMu.Lock()
	if b.closed {
		b.closeMu.Unlock()
		return
	}
	b.closed = true
	b.closeMu.Unlock()
	close(b.closeCh)
	b.wg.Wait()
	close(b.eventCh)
	log.Println("event bus stopped")
}