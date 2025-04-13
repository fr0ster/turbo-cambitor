package common_web_api

import (
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	common "github.com/fr0ster/turbo-cambitor/common"
	web_socket "github.com/fr0ster/turbo-restler/web_socket"
	signature "github.com/fr0ster/turbo-signer/signature"

	"github.com/gorilla/websocket"
)

func (wa *WebApiWrapper) Lock() {
	wa.mutex.Lock()
}

func (wa *WebApiWrapper) Unlock() {
	wa.mutex.Unlock()
}

func New(
	host common.WsHost,
	endpoint common.WsEndpoint,
	scheme common.WsScheme,
	sign signature.Sign) *WebApiWrapper {
	factory := func() (web_socket.WebSocketInterface, error) {
		url := string(scheme) + "://" + string(host) + string(endpoint)
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			return nil, err
		}
		return web_socket.NewWebSocketWrapper(conn), nil
	}
	wa, err := factory()
	if err != nil {
		return nil
	}
	return &WebApiWrapper{
		waScheme:   scheme,
		waHost:     host,
		waEndpoint: endpoint,
		mutex:      &sync.Mutex{},
		sign:       sign,
		connection: wa,
	}
}

func (wa *WebApiWrapper) SetTimeOut(timeout time.Duration) *WebApiWrapper {
	wa.timeout = &timeout
	return wa
}

func (wa *WebApiWrapper) Call(js *simplejson.Json) (result *simplejson.Json, err error) {
	wa.mutex.Lock()
	defer wa.mutex.Unlock()

	writer := wa.connection.GetWriter()
	rq, err := js.MarshalJSON()
	if err != nil {
		return nil, err
	}
	if wa.timeout != nil {
		wa.connection.SetWriteTimeout(*wa.timeout)
	}
	if err := writer.WriteMessage(websocket.TextMessage, rq); err != nil {
		return nil, err
	}

	reader := wa.connection.GetReader()
	if wa.timeout != nil {
		wa.connection.SetReadTimeout(*wa.timeout)
	}
	_, resp, err := reader.ReadMessage()
	if err != nil {
		return nil, err
	}
	result, err = simplejson.NewJson(resp)
	return
}
