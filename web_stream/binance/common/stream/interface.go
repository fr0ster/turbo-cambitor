package streamer

import (
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/fr0ster/turbo-restler/web_socket"
)

// StreamInterface описує лише бізнес-рівень логіки роботи зі стрімами
type StreamInterface interface {
	// Connect встановлює WebSocket-зʼєднання і запускає стрім
	Connect(...bool) error

	// Reconnect перестворює зʼєднання з повторними спробами
	Reconnect(maxAttempts int, delay time.Duration) error

	// Call виконує запит і чекає відповідь у стилі RPC
	Call(rq *simplejson.Json) (*simplejson.Json, error)

	// Subscribe надсилає бізнес-рівневу підписку (наприклад, на Binance)
	Subscribe(subscriptions ...string) error

	// Unsubscribe знімає підписку на бізнес-потік
	Unsubscribe(subscriptions ...string) error

	// ListOfSubscriptions повертає поточні бізнес-підписки
	ListOfSubscriptions() ([]string, error)

	// SetSymbol задає символ (наприклад, BTCUSDT)
	SetSymbol(symbol string) StreamInterface

	// GetConnection повертає underlying WebSocketInterface, якщо потрібно передати далі
	GetConnection() web_socket.WebSocketInterface

	// SetMaxReconnectAttempts встановлює максимальну кількість спроб реконекту
	SetMaxReconnectAttempts(n int) StreamInterface

	// SetReconnectInterval задає інтервал між спробами реконекту
	SetReconnectInterval(interval time.Duration) StreamInterface

	// EnableAutoReconnect вмикає автоматичний реконект
	EnableAutoReconnect() StreamInterface

	// DisableAutoReconnect вимикає автоматичний реконект
	DisableAutoReconnect()
}
