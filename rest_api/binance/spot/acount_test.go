package spot_rest_api_test

import (
	"os"
	"testing"

	spot "github.com/fr0ster/turbo-cambitor/rest_api/binance/spot"
	signature "github.com/fr0ster/turbo-signer/v2/signature"
	"github.com/stretchr/testify/assert"
)

var (
	apiKey = os.Getenv("SPOT_TEST_BINANCE_API_KEY")
	secret = os.Getenv("SPOT_TEST_BINANCE_SECRET_KEY")
	sign   = signature.NewSignHMAC(signature.PublicKey(apiKey), signature.SecretKey(secret))
)

func requireSpotAPIKeys(t *testing.T) bool {
	t.Helper()
	return apiKey != "" && secret != ""
}

func TestAccount(t *testing.T) {
	if !assert.True(t, requireSpotAPIKeys(t), "SPOT_TEST_BINANCE_API_KEY/SECRET_KEY must be set") {
		return
	}
	ra := spot.New(sign, true)
	response, err := ra.Account().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryCurrentOrderCountUsage(t *testing.T) {
	if !assert.True(t, requireSpotAPIKeys(t), "SPOT_TEST_BINANCE_API_KEY/SECRET_KEY must be set") {
		return
	}
	ra := spot.New(sign, true)
	response, err := ra.QueryCurrentOrderCountUsage().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryAllocations(t *testing.T) {
	if !assert.True(t, requireSpotAPIKeys(t), "SPOT_TEST_BINANCE_API_KEY/SECRET_KEY must be set") {
		return
	}
	ra := spot.New(sign, true)
	response, err := ra.QueryAllocations("BTCUSDT").SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
