package streamer

import (
	"fmt"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/google/uuid"

	"github.com/fr0ster/turbo-cambitor/common"
	"github.com/fr0ster/turbo-restler/web_socket"
)

type StreamWrapper struct {
	socket     web_socket.WebSocketInterface
	factory    func() (web_socket.WebSocketInterface, error)
	wsScheme   common.WsScheme
	wsHost     common.WsHost
	wsEndpoint common.WsEndpoint
	timeOut    time.Duration
	mu         sync.Mutex

	readTimeout  *time.Duration
	writeTimeout *time.Duration

	autoReconnect        bool
	reconnectStopChan    chan struct{}
	maxReconnectAttempts int
	reconnectInterval    time.Duration

	handlers map[string]int
}

// NewStreamWrapper creates a new StreamWrapper with a socket factory
func NewStreamWrapper(
	factory func() (web_socket.WebSocketInterface, error),
	wsScheme common.WsScheme,
	wsHost common.WsHost,
	wsEndpoint common.WsEndpoint,
	timeOut ...time.Duration) *StreamWrapper {
	if len(timeOut) == 0 {
		timeOut = append(timeOut, time.Second)
	}

	socket, err := factory()
	if err != nil {
		return nil
	}
	return &StreamWrapper{
		socket:     socket,
		factory:    factory,
		wsScheme:   wsScheme,
		wsHost:     wsHost,
		wsEndpoint: wsEndpoint,
		timeOut:    timeOut[0],
		handlers:   make(map[string]int),
	}
}

func (sw *StreamWrapper) Connect() (*StreamWrapper, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if sw.socket == nil {
		return nil, fmt.Errorf("socket is nil")
	}
	sw.socket.Open()
	return sw, nil
}

func (sw *StreamWrapper) Reconnect(maxAttempts int, delay time.Duration) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	for i := 0; i < maxAttempts; i++ {
		if sw.socket != nil {
			sw.socket.Close()
		}
		socket, err := sw.factory()
		if err == nil {
			socket.Open()
			sw.socket = socket
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("failed to reconnect after %d attempts", maxAttempts)
}

func (sw *StreamWrapper) Disconnect() {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	if sw.socket != nil {
		sw.socket.Close()
		sw.socket = nil
	}
}

func (sw *StreamWrapper) Call(rq *simplejson.Json) (*simplejson.Json, error) {
	id := rq.Get("id").MustString()
	if rq.Get("id").MustString() == "" {
		id = uuid.New().String()
		rq.Set("id", id)
	}

	resultC := make(chan *simplejson.Json, 1)
	errC := make(chan error, 1)

	subID, err := sw.socket.Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			errC <- evt.Error
			return
		}
		resp, err := simplejson.NewJson(evt.Body)
		if err != nil {
			errC <- err
			return
		}
		if resp.Get("id").MustString() == id {
			resultC <- resp
		}
	})
	if err != nil {
		return nil, fmt.Errorf("subscribe error: %w", err)
	}
	defer sw.socket.Unsubscribe(subID)

	jsonBytes, err := rq.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	if err := sw.socket.Send(web_socket.WriteEvent{Body: jsonBytes}); err != nil {
		return nil, fmt.Errorf("send error: %w", err)
	}

	select {
	case resp := <-resultC:
		return resp, nil
	case err := <-errC:
		return nil, err
	case <-time.After(sw.timeOut):
		return nil, fmt.Errorf("timeout")
	}
}

func (sw *StreamWrapper) Subscribe(f func(web_socket.MessageEvent), subs ...string) error {
	if len(subs) == 0 {
		return fmt.Errorf("no subscriptions provided")
	}

	var newSubs []string
	for _, sub := range subs {
		if _, ok := sw.handlers[sub]; !ok {
			newSubs = append(newSubs, sub)
		}
	}

	if len(newSubs) == 0 {
		return fmt.Errorf("already subscribed to all provided streams")
	}

	rq := simplejson.New()
	rq.Set("method", "SUBSCRIBE")
	rq.Set("params", newSubs)
	rq.Set("id", uuid.New().String())
	_, err := sw.Call(rq)
	if err != nil {
		return fmt.Errorf("subscribe error: %w", err)
	}

	for _, sub := range newSubs {
		id, err := sw.socket.Subscribe(f)
		if err != nil {
			continue
		}
		sw.handlers[sub] = id
	}
	return nil
}

func (sw *StreamWrapper) Unsubscribe(subs ...string) error {
	if len(subs) == 0 {
		return fmt.Errorf("no subscriptions provided")
	}

	var toRemove []string
	for _, sub := range subs {
		if _, ok := sw.handlers[sub]; ok {
			toRemove = append(toRemove, sub)
		}
	}

	if len(toRemove) == 0 {
		return fmt.Errorf("no active subscriptions found")
	}

	rq := simplejson.New()
	rq.Set("method", "UNSUBSCRIBE")
	rq.Set("params", toRemove)
	rq.Set("id", uuid.New().String())

	_, err := sw.Call(rq)
	if err != nil {
		return fmt.Errorf("unsubscribe error: %w", err)
	}

	for _, sub := range toRemove {
		id := sw.handlers[sub]
		sw.socket.Unsubscribe(id)
		delete(sw.handlers, sub)
	}

	return nil
}

func (sw *StreamWrapper) ListOfSubscriptions() ([]string, error) {
	rq := simplejson.New()
	rq.Set("method", "LIST_SUBSCRIPTIONS")
	rq.Set("id", uuid.New().String())

	resp, err := sw.Call(rq)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, v := range resp.Get("result").MustArray() {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result, nil
}

func (sw *StreamWrapper) GetConnection() web_socket.WebSocketInterface {
	return sw.socket
}

func (sw *StreamWrapper) SetMaxReconnectAttempts(n int) StreamInterface {
	sw.maxReconnectAttempts = n
	return sw
}

func (sw *StreamWrapper) SetReadTimeout(timeout time.Duration) StreamInterface {
	sw.readTimeout = &timeout
	if sw.socket != nil {
		sw.socket.SetReadTimeout(timeout)
	}
	return sw
}

func (sw *StreamWrapper) SetWriteTimeout(timeout time.Duration) StreamInterface {
	sw.writeTimeout = &timeout
	if sw.socket != nil {
		sw.socket.SetWriteTimeout(timeout)
	}
	return sw
}

func (sw *StreamWrapper) SetReconnectInterval(interval time.Duration) StreamInterface {
	sw.reconnectInterval = interval
	return sw
}

func (sw *StreamWrapper) EnableAutoReconnect() StreamInterface {
	sw.autoReconnect = true
	sw.reconnectStopChan = make(chan struct{})

	go func() {
		attempts := 0
		for {
			select {
			case <-sw.reconnectStopChan:
				return
			default:
				// Перевірка на nil з'єднання
				if sw.socket == nil {
					goto tryReconnect
				}

				// Неблокуюча перевірка на закриття WebSocket
				select {
				case <-sw.socket.Done():
					goto tryReconnect
				default:
					// з'єднання активне, чекаємо і продовжуємо
					time.Sleep(sw.reconnectInterval)
					continue
				}

			tryReconnect:
				if attempts >= sw.maxReconnectAttempts {
					return
				}
				if err := sw.Reconnect(1, sw.reconnectInterval); err == nil {
					attempts = 0
				} else {
					attempts++
				}
				time.Sleep(sw.reconnectInterval)
			}
		}
	}()

	return sw
}

func (sw *StreamWrapper) DisableAutoReconnect() {
	sw.autoReconnect = false
	if sw.reconnectStopChan != nil {
		close(sw.reconnectStopChan)
		sw.reconnectStopChan = nil
	}
}

func (sw *StreamWrapper) SetMessageLogger(logger func(message web_socket.LogRecord)) StreamInterface {
	if sw.socket != nil {
		sw.socket.SetMessageLogger(logger)
	}
	return sw
}
