package streamer

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (stream *StreamWrapper) Close() {
	if stream.low_stream == nil {
		return
	}
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
	logicErrC := make(chan error, 1)
	userErrC := ws.low_stream.GetErrorC()

	ws.AddHandler(id, func(response *simplejson.Json) {
		if response == nil {
			select {
			case userErrC <- fmt.Errorf("received nil response"):
			default:
			}
			return
		}

		if errStr := response.Get("error").MustString(); errStr != "" {
			select {
			case logicErrC <- fmt.Errorf("call error: %s", errStr):
			default:
			}
			return
		}

		if response.Get("id").MustString() != id {
			return
		}

		select {
		case resultC <- response:
		default:
		}
	})
	defer ws.RemoveHandler(id)

	if err := ws.low_stream.Send(rq); err != nil {
		ws.triggerAutoReconnectIfNeeded(err) // 👈 помилка при надсиланні
		return nil, fmt.Errorf("send error: %w", err)
	}

	select {
	case err := <-userErrC:
		ws.triggerAutoReconnectIfNeeded(err) // 👈 помилка від сокету
		return nil, fmt.Errorf("call error: %w", err)

	case err := <-logicErrC:
		ws.low_stream.ErrorHandler()(err) // 👈 викликаємо твій кастомний логічний хендлер
		return nil, err

	case resp := <-resultC:
		return resp, nil

	case <-time.After(ws.timeOut):
		timeoutErr := fmt.Errorf("timeout")
		ws.triggerAutoReconnectIfNeeded(timeoutErr)
		return nil, timeoutErr
	}
}

func (ws *StreamWrapper) triggerAutoReconnectIfNeeded(err error) {
	if !ws.autoReconnect {
		return
	}

	go func() {
		logrus.Warnf("🔁 Auto-reconnect triggered due to: %v", err)

		// ✅ Ця пауза обов'язкова: перед першою спробою!
		time.Sleep(ws.reconnectInterval + 200*time.Millisecond)

		attempts := 0

		for {
			select {
			case <-ws.reconnectStopChan:
				return
			default:
				if ws.low_stream == nil || !ws.low_stream.GetLoopStarted() {
					logrus.Warnf("🔁 Attempting reconnect (attempt %d)...", attempts+1)

					if err := ws.Reconnect(); err != nil {
						logrus.Warnf("Reconnect failed: %v", err)
						attempts++
						if attempts >= ws.maxReconnectAttempts {
							logrus.Errorf("❌ Reconnect failed after %d attempts — stopping auto-reconnect", attempts)
							ws.Close()
							return
						}
						time.Sleep(ws.reconnectInterval + time.Duration(attempts*100)*time.Millisecond)
						continue
					}

					logrus.Info("✅ Reconnected successfully")
					attempts = 0
				}

				time.Sleep(ws.reconnectInterval)
			}
		}
	}()
}

func (sw *StreamWrapper) Reconnect() error {
	sw.reconnectMu.Lock()
	defer sw.reconnectMu.Unlock()

	// Якщо вже успішно реконнектились — не робимо нічого
	if sw.reconnectedOnce {
		return nil
	}

	sw.Close()

	err := sw.Connect()
	if err != nil {
		logrus.Errorf("Reconnect error: %v", err)
		return err
	}

	sw.reconnectedOnce = true
	logrus.Info("✅ Reconnected successfully (marked as once)")
	return nil
}

func (sw *StreamWrapper) Connect() error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	var err error
	for i := 0; i < 10; i++ {
		var stream *web_socket.WebSocketWrapper
		stream, err = web_socket.New(sw.wsHost, sw.wsPath, sw.wsScheme, sw.messageType, sw.silent)
		if err == nil {
			sw.low_stream = stream
			return nil
		}
		logrus.Errorf("🔁 Connect failed: %v", err)
		time.Sleep(500 * time.Millisecond) // 🧘 чекати на старт сервера
	}
	return fmt.Errorf("Connect failed: %w", err)
}
