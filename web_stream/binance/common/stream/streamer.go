package streamer

import (
	"fmt"
	"strings"
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
	symbol     string
	wsScheme   common.WsScheme
	wsHost     common.WsHost
	wsEndpoint common.WsEndpoint
	timeOut    time.Duration
	mu         sync.Mutex

	autoReconnect        bool
	reconnectStopChan    chan struct{}
	maxReconnectAttempts int
	reconnectInterval    time.Duration

	// для контролю паралельного доступу
	reconnectMu sync.Mutex
}

// NewStreamWrapper створює новий StreamWrapper з фабрикою сокета
func NewStreamWrapper(
	factory func() (web_socket.WebSocketInterface, error),
	wsScheme common.WsScheme,
	wsHost common.WsHost,
	wsEndpoint common.WsEndpoint,
	timeOut ...time.Duration) *StreamWrapper {
	if len(timeOut) == 0 {
		timeOut = append(timeOut, 10*time.Second)
	}
	return &StreamWrapper{
		factory:    factory,
		wsScheme:   wsScheme,
		wsHost:     wsHost,
		wsEndpoint: wsEndpoint,
		timeOut:    timeOut[0],
	}
}

func (sw *StreamWrapper) Connect(start ...bool) error {
	if len(start) == 0 {
		start = append(start, true)
	}
	sw.mu.Lock()
	defer sw.mu.Unlock()

	socket, err := sw.factory()
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	if start[0] {
		socket.Open()
	}
	sw.socket = socket
	return nil
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

func (sw *StreamWrapper) Call(rq *simplejson.Json) (*simplejson.Json, error) {
	id := rq.Get("id").MustString()
	if rq.Get("id").MustString() == "" {
		id = uuid.New().String()
		rq.Set("id", id)
	}

	resultC := make(chan *simplejson.Json, 1)
	errC := make(chan error, 1)

	subID := sw.socket.Subscribe(func(evt web_socket.MessageEvent) {
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
	defer sw.socket.Unsubscribe(subID)

	jsonBytes, err := rq.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	if err := sw.socket.Send(jsonBytes); err != nil {
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

func (sw *StreamWrapper) Subscribe(subs ...string) error {
	rq := simplejson.New()
	rq.Set("method", "SUBSCRIBE")
	rq.Set("params", subs)
	rq.Set("id", uuid.New().String())
	_, err := sw.Call(rq)
	return err
}

func (sw *StreamWrapper) Unsubscribe(subs ...string) error {
	rq := simplejson.New()
	rq.Set("method", "UNSUBSCRIBE")
	rq.Set("params", subs)
	rq.Set("id", uuid.New().String())
	_, err := sw.Call(rq)
	return err
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

func (sw *StreamWrapper) SetSymbol(symbol string) StreamInterface {
	sw.symbol = symbol
	if symbol != "" {
		sw.wsEndpoint = common.WsEndpoint(fmt.Sprintf("/%s/%s", sw.wsEndpoint, strings.ToLower(symbol)))
	}
	return sw
}

func (sw *StreamWrapper) GetConnection() web_socket.WebSocketInterface {
	return sw.socket
}

func (sw *StreamWrapper) SetMaxReconnectAttempts(n int) StreamInterface {
	sw.maxReconnectAttempts = n
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
