package streamer

import (
	web_api "github.com/fr0ster/turbo-restler/web_api"
	web_stream "github.com/fr0ster/turbo-restler/web_stream"

	"github.com/sirupsen/logrus"
)

type (
	Stream struct {
		symbol       string
		wsHost       web_api.WsHost
		wsPath       web_api.WsPath
		wsHandler    web_stream.WsHandler
		wsErrHandler web_stream.ErrHandler
		low_stream   *web_stream.WebStream
	}
)

func (stream *Stream) Start() (err error) {
	return stream.low_stream.Start()
}

func (stream *Stream) Stop() {
	stream.low_stream.Stop()
}

func (stream *Stream) Close() {
	stream.low_stream.Close()
}

func (stream *Stream) SetHandler(handler web_stream.WsHandler) *Stream {
	stream.wsHandler = handler
	return stream
}

func (stream *Stream) SetErrHandler(errHandler web_stream.ErrHandler) *Stream {
	stream.wsErrHandler = errHandler
	return stream
}

func (stream *Stream) AddSubscription(handler web_stream.WsHandler, subscription string) {
	stream.low_stream.AddHandler(subscription, handler)
	stream.low_stream.Subscribe(subscription)
}

func (stream *Stream) ListOfSubscriptions(handler web_stream.WsHandler) error {
	return stream.low_stream.ListOfSubscriptions(handler)
}

func (stream *Stream) RemoveSubscriptions(subscription string) {
	stream.low_stream.RemoveHandler(subscription)
	stream.low_stream.Unsubscribe(subscription)
}

func (stream *Stream) SetSymbol(symbol string) *Stream {
	stream.symbol = symbol
	return stream
}

func New(
	host web_api.WsHost,
	path web_api.WsPath,
	scheme ...web_api.WsScheme) *Stream {
	stream, err := web_stream.New(host, path, scheme...)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	}
	return &Stream{
		wsHost:     host,
		wsPath:     path,
		low_stream: stream,
	}
}
