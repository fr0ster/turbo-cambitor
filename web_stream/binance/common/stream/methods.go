package streamer

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (stream *StreamWrapper) Start() (err error) {
	return stream.low_stream.Start()
}

func (stream *StreamWrapper) Stop() {
	stream.low_stream.Stop()
}

func (stream *StreamWrapper) Close() {
	stream.low_stream.Close()
}

func (stream *StreamWrapper) SetErrHandler(errHandler web_socket.ErrHandler) *StreamWrapper {
	stream.low_stream.SetErrHandler(errHandler)
	return stream
}

func (ws *StreamWrapper) AddHandler(handlerId string, handler web_socket.WsHandler) *StreamWrapper {
	ws.low_stream.AddHandler(handlerId, handler)
	return ws
}

func (ws *StreamWrapper) RemoveHandler(handlerId string) *StreamWrapper {
	ws.low_stream.RemoveHandler(handlerId)
	return ws
}

func (ws *StreamWrapper) Subscribe(subscriptions ...string) (err error) {
	if len(subscriptions) == 0 {
		err = fmt.Errorf("no subscriptions")
		return
	}
	// Send subscription request
	rq := simplejson.New()
	rq.Set("method", "SUBSCRIBE")
	rq.Set("id", uuid.New().String())
	rq.Set("params", subscriptions)
	response, err := ws.syncWebApiCall(rq)
	if err != nil {
		return
	}
	if !(response.Get("result").MustString() == "" && response.Get("id").MustString() == rq.Get("id").MustString()) {
		err = fmt.Errorf("something went wrong")
	}
	return
}

func (ws *StreamWrapper) ListOfSubscriptions(handler web_socket.WsHandler) (resultOut []string, err error) {
	rq := simplejson.New()
	rq.Set("method", "LIST_SUBSCRIPTIONS")
	rq.Set("id", uuid.New().String())
	response, err := ws.syncWebApiCall(rq)
	result := response.Get("result").MustArray()
	if len(result) == 0 && response.Get("id").MustString() != rq.Get("id").MustString() {
		err = fmt.Errorf("something went wrong")
		return
	}
	for _, v := range result {
		resultOut = append(resultOut, v.(string))
	}
	return
}

func (ws *StreamWrapper) Unsubscribe(subscriptions ...string) (err error) {
	if len(subscriptions) == 0 {
		err = fmt.Errorf("no subscriptions")
		return
	}
	// Send unsubscribe request
	rq := simplejson.New()
	rq.Set("method", "UNSUBSCRIBE")
	rq.Set("id", uuid.New().String())
	rq.Set("params", subscriptions)
	response, err := ws.syncWebApiCall(rq)
	if err != nil {
		return
	}
	if !(response.Get("result").MustBool() && response.Get("id").MustString() == rq.Get("id").MustString()) {
		err = fmt.Errorf("something went wrong")
	}
	return
}

func (ws *StreamWrapper) syncWebApiCall(rq *simplejson.Json) (responseOut *simplejson.Json, err error) {
	// Створюємо канал для отримання відповіді
	resultC := make(chan *simplejson.Json, 1)

	// Додаємо обробник відповіді
	ws.AddHandler(rq.Get("id").MustString(), func(response *simplejson.Json) {
		if response.Get("id").MustString() == rq.Get("id").MustString() {
			resultC <- response
		}
	})

	// Відправляємо запит
	err = ws.low_stream.Send(rq)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	}
	select {
	case <-time.After(ws.timeOut):
		err = fmt.Errorf("timeout")
		return
	case response, ok := <-resultC:
		if ok {
			if response.Get("id").MustString() == rq.Get("id").MustString() {
				responseOut = response
			}
		} else {
			err = fmt.Errorf("error: %v", response.Get("error"))
		}
	}
	// Видаляємо обробник відповіді
	ws.RemoveHandler(rq.Get("id").MustString())
	return
}
