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
	result, err := wa.AccountBalance().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 2: Account Balance V2
func TestAccountBalanceV2(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.AccountBalanceV2().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 3: Account Information
func TestAccountInformation(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.AccountInformation().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 4: Account Information V2
func TestAccountInformationV2(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.AccountInformationV2().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 5: Account Positions
func TestAccountPositions(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.AccountPositions().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 6: Account Positions V2
func TestAccountPositionsV2(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.AccountPositionsV2().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 9: Logon
func TestLogon(t *testing.T) {
	wa := web_api.New(sign, true)
	_, err := wa.Logon().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Equal(t, `error request: map[code:-4056 msg:HMAC_SHA256 API key is not supported.]`, err.Error())
}

// Test 11: Order Book
func TestOrderBook(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.OrderBook().SetAPIKey().Set("symbol", "BTCUSDT").Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 12: Ping
func TestPing(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.Ping().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 13: Place, Query And Cancel Order
func TestPlaceAndQueryAndCancelOrder(t *testing.T) {
	wa := web_api.New(sign, true)
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
	result, err = wa.QueryOrder().SetAPIKey().Set("symbol", "BTCUSDT").Set("orderId", result.Get("orderId").MustInt()).SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	result, err = wa.CancelOrder().SetAPIKey().Set("symbol", "BTCUSDT").Set("orderId", result.Get("orderId").MustInt()).SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 6: Place, Query And Modify Order
func TestPlaceAndQueryAndModifyOrder(t *testing.T) {
	wa := web_api.New(sign, true)
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
	result, err = wa.QueryOrder().SetAPIKey().Set("symbol", "BTCUSDT").Set("orderId", result.Get("orderId").MustInt()).SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	result, err = wa.ModifyOrder().
		SetAPIKey().
		Set("symbol", "BTCUSDT").
		Set("orderId", result.Get("orderId").MustInt()).
		Set("side", "BUY").
		Set("quantity", "0.02").
		Set("price", "10001").SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 17: Query Position
func TestQueryPosition(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.QueryPosition().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 18: Query Position V2
func TestQueryPositionV2(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.QueryPositionV2().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 19: Status
func TestStatus(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.Status().SetAPIKey().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 20: Symbol Book Ticker
func TestSymbolBookTicker(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.SymbolBookTicker().SetAPIKey().Set("symbol", "BTCUSDT").SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 21: Symbol Price Ticker
func TestSymbolPriceTicker(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.SymbolPriceTicker().SetAPIKey().Set("symbol", "BTCUSDT").SetTimestamp().SetSignature().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// Test 22: Time
func TestTime(t *testing.T) {
	wa := web_api.New(sign, true)
	result, err := wa.Time().Do()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
