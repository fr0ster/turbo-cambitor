package common_web_api

import (
	web_socket "github.com/fr0ster/turbo-restler/web_socket"
	"github.com/gorilla/websocket"
)

// Option applies a configuration to WebApiWrapper during construction.
type Option func(*WebApiWrapper)

// WithFactory overrides the default WebSocket factory.
func WithFactory(f func() (web_socket.WebSocketCommonInterface, error)) Option {
	return func(wa *WebApiWrapper) {
		if f != nil {
			wa.factory = f
		}
	}
}

// WithWebSocketConfig builds a factory using NewWebSocketWrapperWithConfig.
func WithWebSocketConfig(config web_socket.WebSocketConfig) Option {
	return func(wa *WebApiWrapper) {
		// Ensure URL is constructed if missing
		if config.URL == "" {
			config.URL = string(wa.waScheme) + "://" + string(wa.waHost) + string(wa.waEndpoint)
		}
		if config.Dialer == nil {
			config.Dialer = websocket.DefaultDialer
		}
		wa.factory = func() (web_socket.WebSocketCommonInterface, error) {
			return web_socket.NewWebSocketWrapperWithConfig(config)
		}
	}
}
