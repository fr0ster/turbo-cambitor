package futures_web_stream_test

import (
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	common "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/futures"
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

	doneC, stopC, err :=
		stream.Symbol("BTCUSDT").Klines("1m").
			SetHandler(
				func(message *simplejson.Json) {
					mockHandler(message)
				}).
			SetErrHandler(
				func(err error) {
					mockErrHandler(err)
				}).
			Start()
	assert.NoError(t, err)
	assert.NotNil(t, doneC)
	assert.NotNil(t, stopC)

	// Stop the streamer after some time
	time.Sleep(timeOut)
	close(stopC)
}

func TestKlines(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").Klines("1m").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestContinuousKlines(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").ContinuousKlines("1m", "BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestPartialBookDepths(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").PartialBookDepths(common.DepthStreamLevel5, common.DepthStreamRate100ms).SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestDiffBookDepths(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").DiffBookDepths(common.DepthStreamRate100ms).SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestAggTrades(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").AggTrades().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestTrades(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").Trades().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestBookTickers(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").BookTickers().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestTickers(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").Tickers().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestMiniTickers(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").MiniTickers().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestUserData(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").UserData("listenKey").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestMarkPrice(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").MarkPrice().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestLiquidationOrder(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").LiquidationOrder().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestContractInfo(t *testing.T) {
	stream := web_stream.New(true)
	stopC, _, err := stream.Symbol("BTCUSDT").ContractInfo().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestStreamStartAndStop(t *testing.T) {
	stream := web_stream.New(true).Symbol("BTCUSDT").MarkPrice().SetHandler(mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, stream)
	doneC, stopC, err := stream.Start()
	assert.NoError(t, err)
	assert.NotNil(t, doneC)
	assert.NotNil(t, stopC)
	time.Sleep(timeOut)
	doneC, stopC = stream.Stop()
	assert.Nil(t, doneC)
	assert.Nil(t, stopC)
}
