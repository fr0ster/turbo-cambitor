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
	stream.low_stream.SetTimeOut(timeout)
	return stream
}

func New(
	host web_socket.WsHost,
	path web_socket.WsPath,
	scheme web_socket.WsScheme,
	silent bool,
	timeOut ...time.Duration) *StreamWrapper {
	timeOutVar := 10 * time.Second
	if len(timeOut) > 0 {
		timeOutVar = timeOut[0]
	}
	stream, err := web_socket.New(host, path, scheme, web_socket.TextMessage, silent)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	}
	return &StreamWrapper{
		wsHost:     host,
		wsPath:     path,
		low_stream: stream,
		timeOut:    timeOutVar,
	}
}
