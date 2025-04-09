package common_web_stream

import (
	"sync"

	"github.com/gorilla/websocket"
)

func (wa *WebStream) Lock() {
	wa.mutex.Lock()
}

func (wa *WebStream) Unlock() {
	wa.mutex.Unlock()
}

func New(conn *websocket.Conn) *WebStream {
	return &WebStream{
		symbol: "",
		waHost: host,
		mutex:  &sync.Mutex{},
		silent: silent[0],
	}
}
