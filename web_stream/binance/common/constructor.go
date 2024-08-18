package common_web_stream

import (
	"sync"

	"github.com/fr0ster/turbo-restler/web_socket"
)

func (wa *WebStream) Lock() {
	wa.mutex.Lock()
}

func (wa *WebStream) Unlock() {
	wa.mutex.Unlock()
}

func New(host web_socket.WsHost) *WebStream {
	return &WebStream{
		waHost: host,
		mutex:  &sync.Mutex{},
	}
}
