package streamer

import (
	"time"

	web_api "github.com/fr0ster/turbo-restler/web_api"
	web_stream "github.com/fr0ster/turbo-restler/web_stream"
)

type (
	Stream struct {
		wsHost       web_api.WsHost
		wsPath       web_api.WsPath
		wsHandler    web_stream.WsHandler
		wsErrHandler web_stream.ErrHandler
	}
)

func (stream *Stream) Start() (
	doneC chan struct{},
	stopC chan struct{},
	err error) {
	streamer, err := web_stream.New(
		stream.wsHost,
		stream.wsPath,
		stream.wsHandler,
		stream.wsErrHandler,
		true,
		60*time.Second)
	if err != nil {
		return
	}
	doneC, stopC, err = streamer.Start()
	if err != nil {
		return
	}
	return
}

func (stream *Stream) SetHandler(
	handler web_stream.WsHandler) *Stream {
	stream.wsHandler = handler
	return stream
}

func (stream *Stream) SetErrHandler(
	errHandler web_stream.ErrHandler) *Stream {
	stream.wsErrHandler = errHandler
	return stream
}

func New(
	host web_api.WsHost,
	path web_api.WsPath) *Stream {
	return &Stream{
		wsHost: host,
		wsPath: path,
	}
}
