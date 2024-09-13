package futures_rest_api_test

import (
	"os"
	"testing"

	futures "github.com/fr0ster/turbo-cambitor/rest_api/binance/futures"
	signature "github.com/fr0ster/turbo-signer/signature"
	"github.com/stretchr/testify/assert"
)

var (
	apiKey = os.Getenv("FUTURE_TEST_BINANCE_API_KEY")
	secret = os.Getenv("FUTURE_TEST_BINANCE_SECRET_KEY")
	sign   = signature.NewSignHMAC(signature.PublicKey(apiKey), signature.SecretKey(secret))
)

func TestCallRestAPI(t *testing.T) {
	ra := futures.New(sign, true)
	listenkey, err := ra.ListenKey()
	assert.NoError(t, err)
	assert.NotEmpty(t, listenkey)

	rq, err := ra.Account().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
}

func TestAccount(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.Account().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestBalance(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.Balance().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestPositionRisk(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.PositionRisk().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestUserCommissionRate(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.UserCommissionRate().Set("symbol", "BTCUSDT").SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestAccountConfiguration(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.AccountConfiguration().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestSymbolConfiguration(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.SymbolConfiguration().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestForceOrders(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.ForceOrders().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestNotionalLeverageBrackets(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.NotionalLeverageBrackets().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestMultiAssetsMargin(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.MultiAssetsMargin().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestPositionMode(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.PositionMode().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestIncome(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.Income().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestBNBBurnStatus(t *testing.T) {
	ra := futures.New(sign, true)
	rq, err := ra.BNBBurnStatus().SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, rq)
	response, err := ra.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
