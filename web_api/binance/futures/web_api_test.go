package futures_web_api_test

import (
	"os"
	"testing"

	web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/futures"
	signature "github.com/fr0ster/turbo-signer/signature"
	"github.com/stretchr/testify/assert"
)

var (
	apiKey = os.Getenv("FUTURE_TEST_BINANCE_API_KEY")
	secret = os.Getenv("FUTURE_TEST_BINANCE_SECRET_KEY")
	sign   = signature.NewSignHMAC(signature.PublicKey(apiKey), signature.SecretKey(secret))
)

// Test 1: Account Balance
func TestAccountBalance(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.AccountBalance().SetAPIKey().SetTimestamp().SetSignature().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 3: Account Information
func TestAccountInformation(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.AccountInformation().SetAPIKey().SetTimestamp().SetSignature().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 5: Account Positions
func TestAccountPositions(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.AccountPositions().SetAPIKey().SetTimestamp().SetSignature().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 11: Order Book
func TestOrderBook(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.OrderBook().SetAPIKey().Set("symbol", "BTCUSDT").Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 12: Ping
func TestPing(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.Ping().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 13: Place, Query And Cancel Order
func TestPlaceAndQueryAndCancelOrder(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.PlaceOrder().
		SetAPIKey().
		Set("symbol", "BTCUSDT").
		Set("side", "BUY").
		Set("type", "LIMIT").
		Set("quantity", "0.01").
		Set("price", "10000").
		Set("timeInForce", "GTC").
		SetTimestamp().
		SetSignature().
		Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
	response, err = wa.Call(
		wa.QueryOrder().
			SetAPIKey().
			Set("symbol", "BTCUSDT").
			Set("orderId", response.Get("orderId").MustInt()).
			SetTimestamp().
			SetSignature().
			Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
	response, err = wa.Call(
		wa.CancelOrder().
			SetAPIKey().
			Set("symbol", "BTCUSDT").
			Set("orderId", response.Get("orderId").MustInt()).
			SetTimestamp().
			SetSignature().
			Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 6: Place, Query And Modify Order
func TestPlaceAndQueryAndModifyOrder(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.PlaceOrder().
		SetAPIKey().
		Set("symbol", "BTCUSDT").
		Set("side", "BUY").
		Set("type", "LIMIT").
		Set("quantity", "0.01").
		Set("price", "10000").
		Set("timeInForce", "GTC").
		SetTimestamp().SetSignature().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
	response, err = wa.Call(wa.QueryOrder().SetAPIKey().Set("symbol", "BTCUSDT").Set("orderId", response.Get("orderId").MustInt()).SetTimestamp().SetSignature().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
	response, err = wa.Call(wa.ModifyOrder().
		SetAPIKey().
		Set("symbol", "BTCUSDT").
		Set("orderId", response.Get("orderId").MustInt()).
		Set("side", "BUY").
		Set("quantity", "0.02").
		Set("price", "10001").SetTimestamp().SetSignature().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 17: Query Position
func TestQueryPosition(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.QueryPosition().Set("symbol", "BTCUSDT").SetAPIKey().SetTimestamp().SetSignature().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 18: Query Position V2
func TestQueryPositionV2(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.QueryPositionV2().SetAPIKey().SetTimestamp().SetSignature().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 19: Status
func TestStatus(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.Status().SetAPIKey().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 20: Symbol Book Ticker
func TestSymbolBookTicker(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.SymbolBookTicker().SetAPIKey().Set("symbol", "BTCUSDT").SetTimestamp().SetSignature().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 21: Symbol Price Ticker
func TestSymbolPriceTicker(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.SymbolPriceTicker().SetAPIKey().Set("symbol", "BTCUSDT").SetTimestamp().SetSignature().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// Test 22: Time
func TestTime(t *testing.T) {
	wa := web_api.New(sign, true)
	response, err := wa.Call(wa.Time().Do())
	assert.Nil(t, err)
	assert.NotNil(t, response)
}
