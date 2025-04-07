package streamer

import (
	"sync"
	"time"

	"github.com/fr0ster/turbo-restler/web_socket"
)

type (
	StreamWrapper struct {
		symbol               string
		wsHost               web_socket.WsHost
		wsPath               web_socket.WsPath
		wsScheme             web_socket.WsScheme
		messageType          web_socket.MessageType
		silent               bool
		low_stream           *web_socket.WebSocketWrapper
		mu                   sync.Mutex
		timeOut              time.Duration
		autoReconnect        bool
		reconnectInterval    time.Duration
		reconnectStopChan    chan struct{}
		reconnectedOnce      bool
		reconnectMu          sync.Mutex
		maxReconnectAttempts int
		stopOnce             sync.Once
	}
)
