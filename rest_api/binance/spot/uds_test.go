package spot_rest_api_test

import (
	"testing"

	futures "github.com/fr0ster/turbo-cambitor/rest_api/binance/spot"
	"github.com/stretchr/testify/assert"
)

func TestListenKey(t *testing.T) {
	if !assert.True(t, requireSpotAPIKeys(t), "SPOT_TEST_BINANCE_API_KEY/SECRET_KEY must be set") {
		return
	}
	ra := futures.New(sign, true)
	listenKey, err := ra.ListenKey()
	assert.Nil(t, err)
	assert.NotEmpty(t, listenKey)
}

func TestKeepAliveListenKey(t *testing.T) {
	if !assert.True(t, requireSpotAPIKeys(t), "SPOT_TEST_BINANCE_API_KEY/SECRET_KEY must be set") {
		return
	}
	ra := futures.New(sign, true)
	listenKey, err := ra.ListenKey()
	assert.Nil(t, err)
	err = ra.KeepAliveListenKey(listenKey)
	assert.Nil(t, err)
}
