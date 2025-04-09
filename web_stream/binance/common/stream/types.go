package streamer

// WsHost represents the WebSocket host (e.g., "stream.binance.com:9443")
type WsHost string

// WsPath represents the WebSocket path (e.g., "/btcusdt@aggTrade")
type WsPath string

// WsScheme represents the WebSocket scheme (e.g., "wss")
type WsScheme string

// Suffix returns the raw path as string
func (p WsPath) Suffix() string {
	return string(p)
}
