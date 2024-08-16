package streamer

import (
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"

	web_api "github.com/fr0ster/turbo-restler/web_api"

	"github.com/sirupsen/logrus"
)

type (
	// WsCallBackMap map of callback functions
	WsHandlerMap map[string]WsHandler

	// WsHandler handles messages
	WsHandler func(*simplejson.Json)

	// ErrHandler handles errors
	ErrHandler func(err error)
	Stream     struct {
		symbol          string
		wsHost          web_api.WsHost
		wsPath          web_api.WsPath
		callStack       WsHandlerMap
		wsErrHandler    ErrHandler
		socket          *web_api.WebApi
		quit            chan struct{}
		timeOut         time.Duration
		streamIsStarted bool
		responseC       chan *simplejson.Json
	}
)

func (stream *Stream) Start() (err error) {
	if !stream.streamIsStarted {
		stream.socket.Socket().SetReadLimit(655350)
		go func() {
			for {
				select {
				case <-stream.quit:
					return
				// case <-time.After(100 * time.Microsecond):
				default:
					response, err := stream.socket.Read()
					if err != nil {
						stream.wsErrHandler(err)
					}
					for _, handler := range stream.callStack {
						handler(response)
					}
				}
			}
		}()
		stream.streamIsStarted = true
	}
	return
}

func (stream *Stream) Stop() {
	close(stream.quit)
	stream.streamIsStarted = false
}

func (stream *Stream) Close() {
	stream.socket.Close()
}

func (stream *Stream) SetHandler(handler WsHandler) *Stream {
	if stream.symbol != "" {
		genHandler := func(symbol string) func(message *simplejson.Json) {
			return func(message *simplejson.Json) {
				if message.Get("s").MustString() == symbol {
					handler(message)
				}
			}
		}
		stream.callStack["defaultHandler"] = genHandler(stream.symbol)
	}
	return stream
}

func (stream *Stream) SetErrHandler(errHandler ErrHandler) *Stream {
	stream.wsErrHandler = errHandler
	return stream
}

func (stream *Stream) AddSubscriptions(handler WsHandler, params ...string) {
	if len(params) > 0 {
		const id = 1
		// Add handlers for each subscription
		for _, subscription := range params {
			parts := strings.Split(subscription, "@")
			stream.callStack[subscription] = func(message *simplejson.Json) {
				if message.Get("e").MustString() == parts[1] && strings.EqualFold(message.Get("s").MustString(), parts[0]) {
					handler(message)
				}
			}
		}
		// Send subscription request
		rq := simplejson.New()
		rq.Set("method", "SUBSCRIBE")
		rq.Set("id", id)
		rq.Set("params", params)
		err := stream.socket.Send(rq)
		if err != nil {
			logrus.Fatalf("Error: %v", err)
		}
	}
}

func (stream *Stream) GetResponse() *simplejson.Json {
	return <-stream.responseC
}

func (stream *Stream) ListSubscriptions(async ...bool) (response *simplejson.Json, err error) {
	if len(async) == 0 {
		async = append(async, false)
	}
	err = func() (err error) {
		const id = 2

		if stream.callStack["listResponseHandler"] == nil {
			stream.callStack["listResponseHandler"] = func(message *simplejson.Json) {
				if message.Get("id").MustInt() == id {
					stream.responseC <- message
				}
			}
		}
		// Send subscription request
		rq := simplejson.New()
		rq.Set("method", "LIST_SUBSCRIPTIONS")
		rq.Set("id", id)
		err = stream.socket.Send(rq)
		return
	}()
	if !async[0] {
		response = <-stream.responseC
	}
	return
}

func (stream *Stream) RemoveSubscriptions(params ...string) {
	const id = 3
	// Add handlers for each subscription
	for _, subscription := range params {
		stream.callStack[subscription] = nil
		// Видалення запису з мапи
		delete(stream.callStack, subscription)
	}
	// Send subscription request
	rq := simplejson.New()
	rq.Set("method", "UNSUBSCRIBE")
	rq.Set("id", id)
	rq.Set("params", params)
	err := stream.socket.Send(rq)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	}
}

func (stream *Stream) SetSymbol(symbol string) *Stream {
	stream.symbol = symbol
	return stream
}

func (stream *Stream) SetTimeOut(timeout time.Duration) *Stream {
	stream.timeOut = timeout
	return stream
}

func New(
	host web_api.WsHost,
	path web_api.WsPath,
	scheme ...web_api.WsScheme) *Stream {
	stream, err := web_api.New(host, path, scheme...)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	}
	// Встановлення обробника для ping повідомлень
	stream.Socket().SetPingHandler(func(appData string) error {
		logrus.Debug("Received ping:", appData)
		err := stream.Socket().WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
		if err != nil {
			logrus.Debug("Error sending pong:", err)
		}
		return nil
	})
	return &Stream{
		wsHost:          host,
		wsPath:          path,
		socket:          stream,
		quit:            make(chan struct{}),
		callStack:       make(WsHandlerMap),
		timeOut:         1 * time.Second,
		streamIsStarted: false,
		responseC:       make(chan *simplejson.Json),
	}
}
