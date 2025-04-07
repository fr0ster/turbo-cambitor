package streamer

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/google/uuid"
)

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

func (ws *StreamWrapper) GetLoopStarted() bool {
	return ws.low_stream.GetLoopStarted()
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
	response, err := ws.Call(rq)
	if err != nil {
		return
	}
	if !(response.Get("result").MustString() == "" && response.Get("id").MustString() == rq.Get("id").MustString()) {
		err = fmt.Errorf("something went wrong")
	}
	return
}

func (ws *StreamWrapper) ListOfSubscriptions() (resultOut []string, err error) {
	rq := simplejson.New()
	rq.Set("method", "LIST_SUBSCRIPTIONS")
	rq.Set("id", uuid.New().String())
	response, err := ws.Call(rq)
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
	response, err := ws.Call(rq)
	if err != nil {
		return
	}
	if !(response.Get("result").MustBool() && response.Get("id").MustString() == rq.Get("id").MustString()) {
		err = fmt.Errorf("something went wrong")
	}
	return
}

func (ws *StreamWrapper) Call(rq *simplejson.Json) (*simplejson.Json, error) {
	idRaw := rq.Get("id").Interface()
	id, ok := idRaw.(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("invalid or missing id in request")
	}

	resultC := make(chan *simplejson.Json, 1)
	errorC := make(chan error, 1)

	// Тимчасовий handler на помилки з сокета (EOF тощо)
	// ws.low_stream.SetErrHandler(func(err error) error {
	// 	select {
	// 	case errorC <- err:
	// 	default:
	// 	}
	// 	return err
	// })

	ws.AddHandler(id, func(response *simplejson.Json) {
		if response == nil {
			err := fmt.Errorf("received nil response")
			ws.low_stream.ErrorHandler()(err)
			select {
			case errorC <- err:
			default:
			}
			return
		}

		respIDRaw := response.Get("id").Interface()
		if respIDStr, ok := respIDRaw.(string); ok && respIDStr == id {
			if errRaw := response.Get("error").Interface(); errRaw != nil {
				if errStr, ok := errRaw.(string); ok && errStr != "" {
					err := fmt.Errorf("%v", errStr)
					ws.low_stream.ErrorHandler()(err)
					select {
					case errorC <- err:
					default:
					}
					return
				}
			}

			select {
			case resultC <- response:
			default:
			}
		}
	})
	defer ws.RemoveHandler(id)

	// Відправка запиту
	if err := ws.low_stream.Send(rq); err != nil {
		return nil, fmt.Errorf("send error: %w", err)
	}

	// Очікування результату, помилки або тайм-ауту
	select {
	case err := <-errorC:
		return nil, fmt.Errorf("call error: %w", err)
	case resp := <-resultC:
		return resp, nil
	case <-time.After(ws.timeOut):
		return nil, fmt.Errorf("timeout")
	}
}
