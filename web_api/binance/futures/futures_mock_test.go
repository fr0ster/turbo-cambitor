package futures_web_api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	futures_web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/futures"
	signature "github.com/fr0ster/turbo-signer/signature"
	"github.com/stretchr/testify/assert"
)

func startMockServer(t *testing.T, path string, response string, statusCode int) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		_, err := w.Write([]byte(response))
		assert.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

func TestPing_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/ping", `{}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.Ping().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestTime_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/time", `{"serverTime": 123456789}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.Time().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestOrderBook_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/depth?symbol=BTCUSDT", `{"lastUpdateId":1,"bids":[["50000.0","1.0"]],"asks":[["51000.0","1.5"]]}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.OrderBook().Set("symbol", "BTCUSDT").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestAccountBalance_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v2/balance", `[{"asset":"USDT","balance":"1000.0","availableBalance":"900.0"}]`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.AccountBalance().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestAccountInformation_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v2/account", `{"makerCommission":10}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.AccountInformation().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestAccountPositions_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v2/positionRisk", `[{"symbol":"BTCUSDT","positionAmt":"0.1"}]`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.AccountPositions().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryPosition_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/positionRisk?symbol=BTCUSDT", `[{"symbol":"BTCUSDT","positionAmt":"0.2"}]`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.QueryPosition().Set("symbol", "BTCUSDT").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryPositionV2_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v2/positionRisk", `[{"symbol":"BTCUSDT","positionAmt":"0.3"}]`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.QueryPositionV2().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestStatus_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/leverageBracket", `[{"symbol":"BTCUSDT","brackets":[{"initialLeverage":10}]}]`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.Status().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestSymbolBookTicker_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/ticker/bookTicker?symbol=BTCUSDT", `{"symbol":"BTCUSDT","bidPrice":"50000.0","askPrice":"51000.0"}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.SymbolBookTicker().Set("symbol", "BTCUSDT").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestSymbolPriceTicker_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/ticker/price?symbol=BTCUSDT", `{"symbol":"BTCUSDT","price":"60000.00"}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.SymbolPriceTicker().Set("symbol", "BTCUSDT").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestCancelOrder_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/order", `{"orderId":123456}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.CancelOrder().Set("symbol", "BTCUSDT").Set("orderId", "123456").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestModifyOrder_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/order", `{"orderId":123456}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.ModifyOrder().Set("symbol", "BTCUSDT").Set("orderId", "123456").Set("quantity", "0.2").Set("price", "60000").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestPlaceOrder_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/order", `{"orderId":123456}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.PlaceOrder().Set("symbol", "BTCUSDT").Set("side", "BUY").Set("type", "LIMIT").Set("quantity", "0.01").Set("price", "50000").Set("timeInForce", "GTC").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryOrder_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/order", `{"orderId":123456}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.QueryOrder().Set("symbol", "BTCUSDT").Set("orderId", "123456").Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestLogon_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/logon", `{"success":true}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.Logon().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestLogout_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/logout", `{"success":true}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.Logout().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestUserDataStreamStart_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/listenKey", `{"listenKey":"abc123"}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.UserDataStreamStart().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestUserDataStreamPing_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/listenKey", `{}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.UserDataStreamPing().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestUserDataStreamStop_Mock(t *testing.T) {
	server := startMockServer(t, "/fapi/v1/listenKey", `{}`, http.StatusOK)
	sign := signature.NewSignHMAC("test", "test")
	wa := futures_web_api.New(server.Listener.Addr().String(), "", "http", sign)
	response, err := wa.Call(wa.UserDataStreamStop().Do())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
