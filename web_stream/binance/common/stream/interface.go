package streamer

import (
	"context"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/fr0ster/turbo-restler/web_socket"
)

// CambitorInterface mirrors the full public API of StreamWrapper so all variants implement the same contract.
type CambitorInterface interface {
	// Connection lifecycle
	Connect() (*StreamWrapper, error)
	Reconnect(maxAttempts int, delay time.Duration) error
	Disconnect()

	// Business operations
	Call(rq *simplejson.Json) (*simplejson.Json, error)
	Subscribe(f func(web_socket.MessageEvent), subscriptions ...string) error
	Unsubscribe(subscriptions ...string) error
	ListOfSubscriptions() ([]string, error)

	// Access to underlying socket
	GetConnection() web_socket.WebSocketCommonInterface

	// Configuration
	SetMaxReconnectAttempts(n int) CambitorInterface
	SetReadTimeout(timeout time.Duration) CambitorInterface
	SetWriteTimeout(timeout time.Duration) CambitorInterface
	SetReconnectInterval(interval time.Duration) CambitorInterface
	EnableAutoReconnect() CambitorInterface
	DisableAutoReconnect()
	SetMessageLogger(logger func(message web_socket.LogRecord)) CambitorInterface
	SetPingHandler(handler func(string) error)
	SetPongHandler(handler func(string) error)

	// Readiness/health helpers
	WaitReady(ctx context.Context, probeInterval time.Duration) error
	WaitServerListening(timeout time.Duration) bool
	WaitForPong(timeout time.Duration) bool
	IsConnected() bool
	HealthPing(timeout time.Duration) bool

	// Health state and readiness
	SetModeSync(defaultTimeout time.Duration) *StreamWrapper
	HealthStatus() HealthState
	WaitConnected(ctx context.Context) error
	WaitHealthy(ctx context.Context) error
	IsPaused() bool
	IsAutoReconnectEnabled() bool
}
