package streamer

import (
	"time"

	"github.com/fr0ster/turbo-restler/web_socket"
)

type (
	StreamWrapper struct {
		symbol     string
		wsHost     web_socket.WsHost
		wsPath     web_socket.WsPath
		low_stream *web_socket.WebSocketWrapper
		timeOut    time.Duration
	}
)
