package common_web_api

import (
	"fmt"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/fr0ster/turbo-restler/web_socket"
	signature "github.com/fr0ster/turbo-signer/signature"
)

func (wa *WebApiWrapper) Lock() {
	wa.mutex.Lock()
}

func (wa *WebApiWrapper) Unlock() {
	wa.mutex.Unlock()
}

func (wa *WebApiWrapper) SetPingHandler(handler ...func(appData string) error) *WebApiWrapper {
	wa.connection.SetPingHandler(handler...)
	return wa
}

func (wa *WebApiWrapper) SetTimeOut(d time.Duration) *WebApiWrapper {
	wa.timeOut = d
	return wa
}

func (wa *WebApiWrapper) Call(rq *simplejson.Json) (result *simplejson.Json, err error) {
	resultC := make(chan *simplejson.Json)
	wa.connection.AddHandler(rq.Get("id").MustString(), func(data *simplejson.Json) {
		if data.Get("id").MustString() == rq.Get("id").MustString() {
			resultC <- data
		}
	})
	wa.connection.Start()
	wa.connection.Send(rq)
	var ok bool
	select {
	case <-time.After(wa.timeOut):
		err = fmt.Errorf("timeout")
	case result, ok = <-resultC:
		if !ok {
			err = fmt.Errorf("error")
		}
	}
	wa.connection.Stop()
	return
}

func New(host web_socket.WsHost, path web_socket.WsPath, sign signature.Sign) *WebApiWrapper {
	wa, err := web_socket.New(host, path)
	if err != nil {
		return nil
	}
	return &WebApiWrapper{
		waHost:        host,
		waPath:        path,
		mutex:         &sync.Mutex{},
		sign:          sign,
		connection:    wa,
		wsHandlerMap:  make(web_socket.WsHandlerMap),
		timeOut:       5 * time.Second,
		isLoopStarted: false,
	}
}
