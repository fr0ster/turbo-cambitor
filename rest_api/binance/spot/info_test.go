package spot_rest_api_test

import (
	"testing"

	spot "github.com/fr0ster/turbo-cambitor/rest_api/binance/spot"
	"github.com/stretchr/testify/assert"
)

func TestExchangeInfo(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.ExchangeInfo().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestPing(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.Ping().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestOrderBook(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.OrderBook("BTCUSDT").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestRecentTrades(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.RecentTrades("BTCUSDT").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestOldTradesLookup(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.OldTradesLookup("BTCUSDT").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestAggregateTrades(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.AggregateTrades("BTCUSDT").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestKlines(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.Klines("BTCUSDT", "1m").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestUIKlines(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.UIKlines("BTCUSDT", "1m").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestCurrentAvgPrice(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.CurrentAvgPrice("BTCUSDT").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestTicker24hr(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.Ticker24hr("BTCUSDT").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestTradingDay(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.TradingDay().Set("symbol", "BTCUSDT").Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestTickerPrice(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.TickerPrice().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestTickerBookTicker(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.TickerBookTicker().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestTickerV3(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.TickerV3().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
