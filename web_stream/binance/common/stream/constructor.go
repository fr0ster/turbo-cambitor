package streamer

import (
	"time"

	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/sirupsen/logrus"
)

func (stream *StreamWrapper) SetSymbol(symbol string) *StreamWrapper {
	stream.symbol = symbol
	return stream
}

func (stream *StreamWrapper) SetTimeOut(timeout time.Duration) *StreamWrapper {
	stream.timeOut = timeout
	return stream
}

func New(
	host web_socket.WsHost,
	path web_socket.WsPath,
	scheme web_socket.WsScheme) *StreamWrapper {
	stream, err := web_socket.New(host, path, scheme)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	}
	return &StreamWrapper{
		wsHost:     host,
		wsPath:     path,
		low_stream: stream,
		timeOut:    10 * time.Second,
	}
}
