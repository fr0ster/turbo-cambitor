package futures_rest_api_test

import (
	"os"
	"testing"

	futures "github.com/fr0ster/turbo-cambitor/rest_api/binance/futures"
	"github.com/stretchr/testify/assert"
)

func allowLiveREST(t *testing.T) bool {
	t.Helper()
	// Увімкнути live REST тести, якщо RUN_LIVE_REST=1
	return os.Getenv("RUN_LIVE_REST") == "1"
}

func TestPing(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.Ping().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestExchangeInfo(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.ExchangeInfo().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestOrderBook(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.OrderBook("BTCUSDT").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestRecentTrades(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.RecentTrades("BTCUSDT").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestOldTradesLookup(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.OldTradesLookup("BTCUSDT").SetAPIKey().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestAggregateTrades(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.AggregateTrades("BTCUSDT").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestKlines(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.Klines("BTCUSDT", "1m").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestContinuousKlines(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.ContinuousKlines("BTCUSDT", "1m", "PERPETUAL").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestIndexPriceKlines(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.IndexPriceKlines("BTCUSDT", "1m").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestMarkPriceKlines(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.MarkPriceKlines("BTCUSDT", "1m").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestFundingRate(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.FundingRate().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestFundingInfo(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.FundingInfo().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestTicker24hr(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.Ticker24hr().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestTickerPrice(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.TickerPrice().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestTickerPriceV2(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.TickerPriceV2().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestBookTicker(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.BookTicker().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestOpenInterest(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.OpenInterest("BTCUSDT").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestCompositeIndexSymbol(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.CompositeIndexSymbol().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestMultiAssetsModeAssetIndex(t *testing.T) {
	t.Parallel()
	if !assert.True(t, allowLiveREST(t), "set RUN_LIVE_REST=1 to run live REST tests") {
		return
	}
	ra := futures.New(sign, true)
	response, err := ra.MultiAssetsModeAssetIndex().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
