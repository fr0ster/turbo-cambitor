package spot_web_api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	spot_web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/spot"
	signature "github.com/fr0ster/turbo-signer/signature"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func startMockWebSocketServer(t *testing.T, path string, response string) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		assert.NoError(t, err)
		defer conn.Close()
		_, _, err = conn.ReadMessage()
		assert.NoError(t, err)
		err = conn.WriteMessage(websocket.TextMessage, []byte(response))
		assert.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

func TestAccountInformation_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/account", `{"makerCommission":10}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/account", "ws", sign)
	response, err := wa.Call(wa.AccountInformation().SetAPIKey().SetTimestamp().SetSignature().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestExchangeInfo_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/exchangeInfo", `{"symbols":["BTCUSDT"]}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/exchangeInfo", "ws", sign)
	response, err := wa.Call(wa.ExchangeInfo().Set("symbols", []string{"BTCUSDT"}).Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestLogout_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/logout", `{"msg":"logged out"}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/logout", "ws", sign)
	response, err := wa.Call(wa.Logout().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestOrderBook_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/depth", `{"lastUpdateId":1}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/depth", "ws", sign)
	response, err := wa.Call(wa.OrderBook().Set("symbol", "BTCUSDT").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestPing_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/ping", `{"pong":"ok"}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/ping", "ws", sign)
	response, err := wa.Call(wa.Ping().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestPlaceOrder_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/order", `{"orderId": 123456}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/order", "ws", sign)
	response, err := wa.Call(wa.PlaceOrder().SetAPIKey().Set("symbol", "BTCUSDT").Set("side", "SELL").Set("type", "LIMIT").Set("quantity", "0.1").Set("price", "100000.0").Set("timeInForce", "GTC").SetTimestamp().SetSignature().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryOrder_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/order", `{"orderId": 123456}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/order", "ws", sign)
	response, err := wa.Call(wa.QueryOrder().SetAPIKey().Set("symbol", "BTCUSDT").Set("orderId", 123456).SetTimestamp().SetSignature().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestCancelOrder_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/order", `{"orderId": 123456}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/order", "ws", sign)
	response, err := wa.Call(wa.CancelOrder().SetAPIKey().Set("symbol", "BTCUSDT").Set("orderId", 123456).SetTimestamp().SetSignature().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestCancelReplaceOrder_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/order/cancelReplace", `{"orderId": 654321}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/order/cancelReplace", "ws", sign)
	response, err := wa.Call(wa.CancelReplaceOrder().SetAPIKey().Set("symbol", "BTCUSDT").Set("cancelOrderId", 123456).Set("side", "BUY").Set("type", "LIMIT").Set("quantity", "0.02").Set("price", "100001").Set("cancelReplaceMode", "STOP_ON_FAILURE").Set("timeInForce", "GTC").SetTimestamp().SetSignature().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryOpenOrders_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/openOrders", `[{"orderId":123}]`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/openOrders", "ws", sign)
	response, err := wa.Call(wa.QueryOpenOrders().SetAPIKey().SetTimestamp().SetSignature().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestStatus_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/system/status", `{"status": "normal"}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/system/status", "ws", sign)
	response, err := wa.Call(wa.Status().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestSymbolBookTicker_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/ticker/bookTicker", `{"symbol":"BTCUSDT"}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/ticker/bookTicker", "ws", sign)
	response, err := wa.Call(wa.SymbolBookTicker().Set("symbols", []string{"BTCUSDT"}).Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestSymbolPriceTicker_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/ticker/price", `{"symbol":"BTCUSDT"}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/ticker/price", "ws", sign)
	response, err := wa.Call(wa.SymbolPriceTicker().Set("symbols", []string{"BTCUSDT"}).Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestTime_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/api/v3/time", `{"serverTime":123456789}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := spot_web_api.New(wsHost, "/api/v3/time", "ws", sign)
	response, err := wa.Call(wa.Time().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
