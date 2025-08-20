package futures_web_api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	futures_web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/futures"
	signature "github.com/fr0ster/turbo-signer/v2/signature"
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

func TestPing_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v1/ping", `{"pong":"ok"}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v1/ping", "ws", sign)
	response, err := wa.Call(wa.Ping().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestTime_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v1/time", `{"serverTime": 123456789}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v1/time", "ws", sign)
	response, err := wa.Call(wa.Time().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestOrderBook_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v1/depth", `{"lastUpdateId":1}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v1/depth", "ws", sign)
	response, err := wa.Call(wa.OrderBook().Set("symbol", "BTCUSDT").Set("symbol", "BTCUSDT").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestAccountBalance_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v2/balance", `[{"asset":"USDT","balance":"1000.0"}]`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v2/balance", "ws", sign)
	response, err := wa.Call(wa.AccountBalance().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestAccountInformation_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v2/account", `{"makerCommission":10}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v2/account", "ws", sign)
	response, err := wa.Call(wa.AccountInformation().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestAccountPositions_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v2/positionRisk", `[{"symbol":"BTCUSDT","positionAmt":"0.1"}]`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v2/positionRisk", "ws", sign)
	response, err := wa.Call(wa.AccountPositions().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryPosition_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v1/positionRisk", `[{"symbol":"BTCUSDT","positionAmt":"0.2"}]`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v1/positionRisk", "ws", sign)
	response, err := wa.Call(wa.QueryPosition().Set("symbol", "BTCUSDT").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryPositionV2_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v2/positionRisk", `[{"symbol":"BTCUSDT","positionAmt":"0.3"}]`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v2/positionRisk", "ws", sign)
	response, err := wa.Call(wa.QueryPositionV2().Set("symbol", "BTCUSDT").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestStatus_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v1/leverageBracket", `[{"symbol":"BTCUSDT","brackets":[{"initialLeverage":10}]}]`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v1/leverageBracket", "ws", sign)
	response, err := wa.Call(wa.Status().Set("symbol", "BTCUSDT").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestSymbolBookTicker_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v1/ticker/bookTicker", `{"symbol":"BTCUSDT"}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v1/ticker/bookTicker", "ws", sign)
	response, err := wa.Call(wa.SymbolBookTicker().Set("symbol", "BTCUSDT").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestSymbolPriceTicker_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v1/ticker/price", `{"symbol":"BTCUSDT"}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v1/ticker/price", "ws", sign)
	response, err := wa.Call(wa.SymbolPriceTicker().Set("symbol", "BTCUSDT").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestPlaceOrder_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v1/order", `{"orderId": 123456}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v1/order", "ws", sign)
	response, err := wa.Call(wa.PlaceOrder().Set("symbol", "BTCUSDT").Set("side", "BUY").Set("type", "LIMIT").Set("quantity", "0.01").Set("price", "10000").Set("timeInForce", "GTC").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryOrder_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v1/order", `{"orderId": 123456}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v1/order", "ws", sign)
	response, err := wa.Call(wa.QueryOrder().Set("symbol", "BTCUSDT").Set("orderId", "123456").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestCancelOrder_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v1/order", `{"orderId": 123456}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v1/order", "ws", sign)
	response, err := wa.Call(wa.CancelOrder().Set("symbol", "BTCUSDT").Set("orderId", "123456").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestModifyOrder_Mock(t *testing.T) {
	server := startMockWebSocketServer(t, "/fapi/v1/order", `{"orderId": 123456}`)
	sign := signature.NewSignHMAC("test", "test")
	wsHost := strings.TrimPrefix(server.URL, "http://")
	wa := futures_web_api.New(wsHost, "/fapi/v1/order", "ws", sign)
	response, err := wa.Call(wa.ModifyOrder().Set("symbol", "BTCUSDT").Set("orderId", "123456").Set("quantity", "0.02").Set("price", "10001").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
