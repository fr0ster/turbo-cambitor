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
		stream       *web_stream.WebStream
		requestID    uint64
		doneC        chan struct{}
		stopC        chan struct{}
	}
)

func (stream *Stream) createStream() (err error) {
	if stream.stream == nil {
		stream.stream, err = web_stream.New(
			stream.wsHost,
			stream.wsPath,
			stream.wsHandler,
			stream.wsErrHandler,
			true,
			60*time.Second)
	}
	return
}

func (stream *Stream) Start() (
	doneC chan struct{},
	stopC chan struct{},
	err error) {
	if stream.stopC != nil && stream.doneC != nil {
		return stream.doneC, stream.stopC, nil
	}
	stream.createStream()
	doneC, stopC, err = stream.stream.Start()
	if err != nil {
		return
	}
	stream.doneC = doneC
	stream.stopC = stopC
	return
}

func (stream *Stream) Stop() (
	doneC chan struct{},
	stopC chan struct{}) {
	if stream.stream != nil {
		close(stream.stopC)
		stream.stopC = nil
		stream.doneC = nil
	}
	return nil, nil
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

func (wa *Stream) Socket() *web_api.WebApi {
	wa.createStream()
	return wa.stream.Socket()
}

func New(
	host web_api.WsHost,
	path web_api.WsPath) *Stream {
	return &Stream{
		wsHost:    host,
		wsPath:    path,
		requestID: 0,
	}
}
