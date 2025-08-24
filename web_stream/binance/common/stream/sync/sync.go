// МОНІТОР КЕРУЄ СТАНОМ ВРАППЕРА
// sync перевіряє стан через монітор, і якщо сервер непінгується — wrapper деякий час чекає, періодично перевіряє стан через монітор.
// Якщо монітор не піднімає флаг, що сервер доступний — тільки тоді повертає помилку.
// Якщо монітор каже, що сервер відновився і пінги йдуть — виконує запит.
// sync не повертає помилку одразу, а чекає певний час і пробує перепідключитись, поки монітор не підтвердить доступність.
// async повертає помилку одразу, якщо монітор каже що сервер недоступний.
// sync/async не ловлять помилок самі, вони керуються монітором, який пінгує сервер.
// А ще є метод для позаграфікової перевірки стану серверу.
package sync

import (
	"time"

	"context"

	"github.com/bitly/go-simplejson"

	common "github.com/fr0ster/turbo-cambitor/common"
	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"
	"github.com/fr0ster/turbo-restler/web_socket"
)

// StreamWrapper is an auto-reconnect sync-mode variant.
type StreamWrapper struct {
	*stream.StreamWrapper
	syncTimeout time.Duration
}

// NewStreamWrapper constructs an auto-reconnect wrapper in sync mode.
func NewStreamWrapper(
	factory func() (web_socket.WebSocketCommonInterface, error),
	wsScheme common.WsScheme,
	wsHost common.WsHost,
	wsEndpoint common.WsEndpoint,
	timeOut ...time.Duration,
) *StreamWrapper {
	sw := stream.NewStreamWrapper(factory, wsScheme, wsHost, wsEndpoint, timeOut...)
	sw.EnableAutoReconnect()
	def := 3 * time.Second
	if len(timeOut) > 0 && timeOut[0] > 0 {
		def = timeOut[0]
	}
	sw.SetModeSync(def)
	return &StreamWrapper{StreamWrapper: sw, syncTimeout: def}
}

// compile-time check
var _ stream.CambitorInterface = (*StreamWrapper)(nil)

// Enforce fixed behavior: auto-reconnect enabled, sync mode.
func (w *StreamWrapper) EnableAutoReconnect() stream.CambitorInterface { return w }
func (w *StreamWrapper) DisableAutoReconnect()                         {}

func (w *StreamWrapper) SetMaxReconnectAttempts(n int) stream.CambitorInterface {
	w.StreamWrapper.SetMaxReconnectAttempts(n)
	return w
}
func (w *StreamWrapper) SetReconnectInterval(interval time.Duration) stream.CambitorInterface {
	w.StreamWrapper.SetReconnectInterval(interval)
	return w
}
func (w *StreamWrapper) SetMessageLogger(logger func(message web_socket.LogRecord)) stream.CambitorInterface {
	w.StreamWrapper.SetMessageLogger(logger)
	return w
}

func (w *StreamWrapper) SetReadTimeout(timeout time.Duration) stream.CambitorInterface {
	w.StreamWrapper.SetReadTimeout(timeout)
	return w
}

func (w *StreamWrapper) SetWriteTimeout(timeout time.Duration) stream.CambitorInterface {
	w.StreamWrapper.SetWriteTimeout(timeout)
	return w
}

// Keep sync mode but allow updating default timeout.
func (w *StreamWrapper) SetModeAsync() *stream.StreamWrapper { return w.StreamWrapper }
func (w *StreamWrapper) SetModeSync(defaultTimeout time.Duration) *stream.StreamWrapper {
	if defaultTimeout > 0 {
		w.syncTimeout = defaultTimeout
	}
	return w.StreamWrapper.SetModeSync(defaultTimeout)
}

// ------------------------------
// Client-facing overrides (Sync semantics: wait for connect within default timeout)
// ------------------------------

func (w *StreamWrapper) waitReady() error {
	if w.StreamWrapper == nil {
		return stream.ErrNotConnected
	}
	// In sync mode, we wait for connection to be ready using configured timeout
	to := w.syncTimeout
	if to <= 0 {
		to = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), to)
	defer cancel()
	if err := w.StreamWrapper.WaitConnected(ctx); err != nil {
		return stream.ErrNotConnected
	}
	if !w.StreamWrapper.IsConnected() || w.StreamWrapper.HealthStatus() != stream.HealthHealthy {
		return stream.ErrNotConnected
	}
	return nil
}

func (w *StreamWrapper) Call(rq *simplejson.Json) (*simplejson.Json, error) {
	if err := w.waitReady(); err != nil {
		return nil, err
	}
	return w.StreamWrapper.Call(rq)
}

func (w *StreamWrapper) Subscribe(f func(web_socket.MessageEvent), subs ...string) error {
	if err := w.waitReady(); err != nil {
		return err
	}
	return w.StreamWrapper.Subscribe(f, subs...)
}

func (w *StreamWrapper) Unsubscribe(subs ...string) error {
	if err := w.waitReady(); err != nil {
		return err
	}
	return w.StreamWrapper.Unsubscribe(subs...)
}

func (w *StreamWrapper) ListOfSubscriptions() ([]string, error) {
	if err := w.waitReady(); err != nil {
		return nil, err
	}
	return w.StreamWrapper.ListOfSubscriptions()
}
