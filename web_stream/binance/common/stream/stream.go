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

func (rq *Stream) Start() (
	doneC chan struct{},
	stopC chan struct{},
	err error) {
	stream, err := web_stream.New(
		rq.wsHost,
		rq.wsPath,
		rq.wsHandler,
		rq.wsErrHandler,
		true,
		60*time.Second)
	if err != nil {
		return
	}
	doneC, stopC, err = stream.Start()
	if err != nil {
		return
	}
	return
}

func (rq *Stream) SetHandler(
	handler web_stream.WsHandler) *Stream {
	rq.wsHandler = handler
	return rq
}

func (rq *Stream) SetErrHandler(
	errHandler web_stream.ErrHandler) *Stream {
	rq.wsErrHandler = errHandler
	return rq
}

func New(
	host web_api.WsHost,
	path web_api.WsPath) *Stream {
	return &Stream{
		wsHost: host,
		wsPath: path,
	}
}
