package streamer

import (
	"sync/atomic"
	"time"

	"github.com/bitly/go-simplejson"
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
	stream.createStream()
	doneC, stopC, err = stream.stream.Start()
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

func (wa *Stream) Subscribe(streams []string) (response *simplejson.Json, err error) {
	wa.wsPath = web_api.WsPath("")
	wa.createStream()
	socket := wa.stream.Socket()
	rq := simplejson.New()
	rq.Set("method", "SUBSCRIBE")
	rq.Set("id", atomic.AddUint64(&wa.requestID, 1))
	rq.Set("params", streams)
	return socket.Call(rq)
}

func (wa *Stream) ListOfSubscriptions() (response *simplejson.Json, err error) {
	wa.wsPath = web_api.WsPath("")
	wa.createStream()
	socket := wa.stream.Socket()
	rq := simplejson.New()
	rq.Set("method", "LIST_SUBSCRIPTIONS")
	rq.Set("id", atomic.AddUint64(&wa.requestID, 1))
	return socket.Call(rq)
}

func (wa *Stream) Unsubscribe(streams []string) (response *simplejson.Json, err error) {
	wa.wsPath = web_api.WsPath("")
	wa.createStream()
	socket := wa.stream.Socket()
	rq := simplejson.New()
	rq.Set("method", "UNSUBSCRIBE")
	rq.Set("id", atomic.AddUint64(&wa.requestID, 1))
	rq.Set("params", streams)
	return socket.Call(rq)
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
