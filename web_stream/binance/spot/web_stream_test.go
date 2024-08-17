package spot_web_stream_test

import (
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	common "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/spot"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	timeOut = 1 * time.Second
)

// Mock handler for WebSocket messages
func mockHandler(message *simplejson.Json) {
	logrus.Infof("Received message: %+v", message)
}

// Mock error handler for WebSocket errors
func mockErrHandler(err error) {
	logrus.Errorf("Error: %v", err)
}

func TestStartBinanceStreamer(t *testing.T) {
	// Test 2: Binance stream
	stream := web_stream.New(true)
	assert.NotNil(t, stream)

	streamer :=
		stream.Klines("1m").SetSymbol("BTCUSDT").
			SetHandler(
				func(message *simplejson.Json) {
					mockHandler(message)
				}).
			SetErrHandler(
				func(err error) {
					mockErrHandler(err)
				})
	err := streamer.Start()
	assert.NoError(t, err)

	// Stop the streamer after some time
	time.Sleep(timeOut)
	streamer.Stop()
}

func TestKlines(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.Klines("1m").SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestContinuousKlines(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.ContinuousKlines("1m", "BTCUSDT").SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestPartialBookDepths(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.PartialBookDepths(common.DepthStreamLevel5, common.DepthStreamRate100ms).SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestDiffBookDepths(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.DiffBookDepths(common.DepthStreamRate100ms).SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestAggTrades(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.AggTrades().SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestTrades(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.Trades().SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestBookTickers(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.BookTickers().SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestTickers(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.Tickers().SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestMiniTickers(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.MiniTickers().SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestUserData(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.UserData("listenKey").SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestMarkPrice(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.MarkPrice().SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestLiquidationOrder(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.LiquidationOrder().SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestContractInfo(t *testing.T) {
	stream := web_stream.New(true)
	err := stream.ContractInfo().SetSymbol("BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
}

func TestStream(t *testing.T) {
	stream := web_stream.
		New(true).
		Stream().
		SetSymbol("BTCUSDT").
		SetErrHandler(mockErrHandler)
	err := stream.Start()
	assert.NoError(t, err)
	stream.AddSubscription(mockHandler, "btcusdt@aggTrade")
	time.Sleep(timeOut)
	stream.RemoveSubscriptions("btcusdt@aggTrade")
	time.Sleep(timeOut)
	stream.Stop()
}
