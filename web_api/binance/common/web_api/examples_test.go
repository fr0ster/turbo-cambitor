package common_web_api_test

import (
	"time"

	common_web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/common/web_api"
	web_socket "github.com/fr0ster/turbo-restler/web_socket"
	"github.com/gorilla/websocket"
)

// Example: override the socket factory.
func ExampleNew_withFactory() {
	wa := common_web_api.New("ws-api.binance.com", "/ws-api/v3", "wss", nil,
		common_web_api.WithFactory(func() (web_socket.WebSocketCommonInterface, error) {
			// You can inject a custom dialer, URL, mocks, etc.
			return web_socket.NewWebSocketWrapper(websocket.DefaultDialer, "wss://ws-api.binance.com/ws-api/v3")
		}),
	)
	_ = wa
	// Output:
}

// Example: provide detailed socket config (timeouts, buffers, reconnects) via WithWebSocketConfig.
func ExampleNew_withWebSocketConfig() {
	wa := common_web_api.New("ws-api.binance.com", "/ws-api/v3", "wss", nil,
		common_web_api.WithWebSocketConfig(web_socket.WebSocketConfig{
			// Leave URL empty to auto-use fields from New(...)
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		}),
	)
	_ = wa
	// Output:
}
