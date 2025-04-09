package streamer

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/google/uuid"
)

type StreamWrapper struct {
	socket  web_socket.WebSocketInterface
	factory func() (web_socket.WebSocketInterface, error)
	symbol  string
	wsPath  WsPath
	timeOut time.Duration
	mu      sync.Mutex
}

// NewStreamWrapper створює новий StreamWrapper з фабрикою сокета
func NewStreamWrapper(factory func() (web_socket.WebSocketInterface, error), wsPath WsPath) *StreamWrapper {
	return &StreamWrapper{
		factory: factory,
		wsPath:  wsPath,
		timeOut: 10 * time.Second,
	}
}

func (sw *StreamWrapper) Connect() error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	socket, err := sw.factory()
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	socket.Open()
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
	id := uuid.New().String()
	rq.Set("id", id)

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
		sw.wsPath = WsPath("/" + strings.ToLower(symbol) + sw.wsPath.Suffix())
	}
	return sw
}

func (sw *StreamWrapper) GetConnection() web_socket.WebSocketInterface {
	return sw.socket
}
