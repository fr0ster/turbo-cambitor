package streamer

import (
	"time"

	"github.com/fr0ster/turbo-restler/web_socket"
)

// StreamInterface describes only the business-level logic for working with streams
type StreamInterface interface {
	// Connect establishes a WebSocket connection and starts the stream
	Connect() error

	// Reconnect recreates the connection with retry attempts
	Reconnect(maxAttempts int, delay time.Duration) error

	// Disconnect closes the WebSocket connection
	Disconnect()

	// Subscribe sends a business-level subscription (e.g., to Binance)
	Subscribe(f func(web_socket.MessageEvent), subscriptions ...string) error

	// Unsubscribe removes a subscription from the business stream
	Unsubscribe(subscriptions ...string) error

	// ListOfSubscriptions returns the current business-level subscriptions
	ListOfSubscriptions() ([]string, error)

	// SetSymbol sets the symbol (e.g., BTCUSDT)
	SetSymbol(symbol string) StreamInterface

	// GetConnection returns the underlying WebSocketInterface, if it needs to be passed further
	GetConnection() web_socket.WebSocketInterface

	// SetMaxReconnectAttempts sets the maximum number of reconnect attempts
	SetMaxReconnectAttempts(n int) StreamInterface

	// SetReconnectInterval sets the interval between reconnect attempts
	SetReconnectInterval(interval time.Duration) StreamInterface

	// EnableAutoReconnect enables automatic reconnection
	EnableAutoReconnect() StreamInterface

	// DisableAutoReconnect disables automatic reconnection
	DisableAutoReconnect()

	// SetMessageLogger sets a function for logging messages
	SetMessageLogger(logger func(message web_socket.LogRecord)) StreamInterface
}
