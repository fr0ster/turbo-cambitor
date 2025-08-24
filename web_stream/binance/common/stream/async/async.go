// МОНІТОР КЕРУЄ СТАНОМ ВРАППЕРА
// async перевіряє стан і якщо коннект на паузі бо пінгів нема, то одразу вертає помилку
// sync чекає якийсь час і якщо монітор не піднімає флаг що сервер доступний, то через той час повертає помилку
// а якщо монітор каже сервер відновився і пінги йдуть, то виконує запит
// sync/async не ловлять помилок самі, вони блять керуются моніторінговим коннектором який пінгує сервер.
// А ще в нього має бути метод який надасть можливість пінгувати сервер поза графіком моніторінгу
// async.go: асинхронний wrapper
//
// НЕ ловить помилки самостійно, не робить reconnect.
package async

import (
	"strings"
	"time"

	"github.com/bitly/go-simplejson"

	common "github.com/fr0ster/turbo-cambitor/common"
	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"
	"github.com/fr0ster/turbo-restler/web_socket"
)

// StreamWrapper is an auto-reconnect async-mode variant (fail-fast when not ready).
type StreamWrapper struct{ *stream.StreamWrapper }

// NewStreamWrapper constructs an auto-reconnect wrapper in async mode.
func NewStreamWrapper(
	factory func() (web_socket.WebSocketCommonInterface, error),
	wsScheme common.WsScheme,
	wsHost common.WsHost,
	wsEndpoint common.WsEndpoint,
	timeOut ...time.Duration,
) *StreamWrapper {
	sw := stream.NewStreamWrapper(factory, wsScheme, wsHost, wsEndpoint, timeOut...)
	sw.EnableAutoReconnect()
	return &StreamWrapper{sw}
}

// compile-time check
var _ stream.CambitorInterface = (*StreamWrapper)(nil)

// Enforce fixed behavior: auto-reconnect enabled, async mode.
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

// Keep async mode; ignore attempts to switch to sync.
func (w *StreamWrapper) SetModeAsync() *stream.StreamWrapper { return w.StreamWrapper }
func (w *StreamWrapper) SetModeSync(defaultTimeout time.Duration) *stream.StreamWrapper {
	return w.StreamWrapper
}

// ------------------------------
// Client-facing overrides (Async semantics: fail-fast if not ready)
// ------------------------------

func (w *StreamWrapper) ensureReadyAsync() error {
	if w.StreamWrapper == nil {
		return stream.ErrNotConnected
	}
	// Fail fast if connection is paused by health monitor
	if w.StreamWrapper.IsPaused() {
		return stream.ErrPaused
	}
	if !w.StreamWrapper.IsConnected() || w.StreamWrapper.HealthStatus() != stream.HealthHealthy {
		return stream.ErrNotConnected
	}
	return nil
}

func (w *StreamWrapper) Call(rq *simplejson.Json) (*simplejson.Json, error) {
	if err := w.ensureReadyAsync(); err != nil {
		return nil, err
	}
	resp, err := w.StreamWrapper.Call(rq)
	if err != nil {
		msg := err.Error()
		// Async variant: якщо монітор уже підняв connection після рестарту — зробити один повтор
		if w.StreamWrapper.IsAutoReconnectEnabled() && w.StreamWrapper.IsConnected() && w.StreamWrapper.HealthStatus() == stream.HealthHealthy && (strings.Contains(msg, "close") || strings.Contains(msg, "broken pipe") || strings.Contains(msg, "timeout")) {
			return w.StreamWrapper.Call(rq)
		}
	}
	return resp, err
}

func (w *StreamWrapper) Subscribe(f func(web_socket.MessageEvent), subs ...string) error {
	if err := w.ensureReadyAsync(); err != nil {
		return err
	}
	return w.StreamWrapper.Subscribe(f, subs...)
}

func (w *StreamWrapper) Unsubscribe(subs ...string) error {
	if err := w.ensureReadyAsync(); err != nil {
		return err
	}
	return w.StreamWrapper.Unsubscribe(subs...)
}

func (w *StreamWrapper) ListOfSubscriptions() ([]string, error) {
	if err := w.ensureReadyAsync(); err != nil {
		return nil, err
	}
	return w.StreamWrapper.ListOfSubscriptions()
}
