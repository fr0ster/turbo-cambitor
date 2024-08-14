package spot_web_api_test

import (
	"os"
	"testing"

	web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/spot"
	signature "github.com/fr0ster/turbo-signer/signature"
	"github.com/stretchr/testify/assert"
)

var (
	apiKey = os.Getenv("SPOT_TEST_BINANCE_API_KEY")
	secret = os.Getenv("SPOT_TEST_BINANCE_SECRET_KEY")
	sign   = signature.NewSignHMAC(signature.PublicKey(apiKey), signature.SecretKey(secret))
)

// Test 3: Account Information
func TestAccountInformation(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.AccountInformation())
	result, err := wa.AccountInformation().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 7: Exchange Info
func TestExchangeInfo(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.ExchangeInfo())
	result, err := wa.ExchangeInfo().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 9: Logon
func TestLogon(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.Logon())
	_, err := wa.Logon().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Equal(t, `error request: map[code:-2028 msg:HMAC-SHA-256 API key is not supported.]`, err.Error())
}

// Test 10: Logout
func TestLogout(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.Logout())
	result, err := wa.Logout().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 11: Order Book
func TestOrderBook(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.OrderBook())
	result, err := wa.OrderBook().Set("symbol", "BTCUSDT").Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 12: Ping
func TestPing(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.Ping())
	result, err := wa.Ping().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 13: Place, Query And Cancel Order
func TestPlaceAndQueryAndCancelOrder(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.PlaceOrder())
	result, err := wa.PlaceOrder().
		SetAPIKey().
		Set("symbol", "BTCUSDT").
		Set("side", "BUY").
		Set("type", "LIMIT").
		Set("quantity", "0.007").
		Set("price", "10000.0").
		Set("timeInForce", "GTC").
		SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, wa.QueryOrder())
	result, err = wa.QueryOrder().SetAPIKey().Set("orderId", result.Get("orderId").MustInt()).SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, wa.CancelOrder())
	result, err = wa.CancelOrder().SetAPIKey().Set("orderId", result.Get("orderId").MustInt()).SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 6: Place, Query And CancelReplace Order
func TestPlaceAndQueryAndCancelReplaceOrder(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.PlaceOrder())
	result, err := wa.PlaceOrder().
		SetAPIKey().
		Set("symbol", "BTCUSDT").
		Set("side", "BUY").
		Set("type", "LIMIT").
		Set("quantity", "0.01").
		Set("price", "10000").
		Set("timeInForce", "GTC").
		SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, wa.QueryOrder())
	result, err = wa.QueryOrder().SetAPIKey().Set("orderId", result.Get("orderId").MustInt()).SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, wa.CancelReplaceOrder())
	result, err = wa.CancelReplaceOrder().
		SetAPIKey().
		Set("orderId", result.Get("orderId").MustInt()).
		Set("side", "BUY").
		Set("quantity", "0.02").
		Set("price", "10001").SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 14: Query All Orders
func TestQueryAllOrders(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.QueryAllOrders())
	result, err := wa.QueryAllOrders().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 15: Query Open Orders
func TestQueryOpenOrders(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.QueryOpenOrders())
	result, err := wa.QueryOpenOrders().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 16: Query Order
func TestQueryOrder(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.QueryOrder())
	result, err := wa.QueryOrder().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 19: Status
func TestStatus(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.Status())
	result, err := wa.Status().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 20: Symbol Book Ticker
func TestSymbolBookTicker(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.SymbolBookTicker())
	result, err := wa.SymbolBookTicker().Set("symbols", []string{"BTCUSDT"}).Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 21: Symbol Price Ticker
func TestSymbolPriceTicker(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.SymbolPriceTicker())
	result, err := wa.SymbolPriceTicker().Set("symbols", []string{"BTCUSDT"}).Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 22: Time
func TestTime(t *testing.T) {
	wa := web_api.New(sign, true)
	assert.NotNil(t, wa.Time())
	result, err := wa.Time().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
