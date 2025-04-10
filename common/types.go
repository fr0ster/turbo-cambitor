package common

// WsHost represents the WebSocket host (e.g., "stream.binance.com:9443")
type WsHost string

// WsEndpoint represents the WebSocket endpoint (e.g., "/ws")
type WsEndpoint string

// WsScheme represents the WebSocket scheme (e.g., "wss")
type WsScheme string

// Suffix returns the raw path as string
func (p WsScheme) Suffix() string {
	return string(p)
}

const (
	// WsSchemeWSS represents the WebSocket Secure scheme
	WsSchemeWSS WsScheme = "wss"
	// WsSchemeWS represents the WebSocket scheme
	WsSchemeWS WsScheme = "ws"
	// WsSchemeHTTP WsScheme = "http"
	WsSchemeHTTP WsScheme = "http"
	// WsSchemeHTTPS WsScheme = "https"
	WsSchemeHTTPS WsScheme = "https"
	// WsSchemeGRPC WsScheme = "grpc"
	WsSchemeGRPC WsScheme = "grpc"
	// WsSchemeGRPCS WsScheme = "grpcs"
	WsSchemeGRPCS WsScheme = "grpcs"
	// WsSchemeGRPCWeb WsScheme = "mqtt"
	WsSchemeGRPCWeb WsScheme = "mqtt"
)
