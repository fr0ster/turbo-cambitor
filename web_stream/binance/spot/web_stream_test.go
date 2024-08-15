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

	doneC, stopC, err :=
		stream.Klines("1m").
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
	assert.NotNil(t, stream.Klines("1m"))
	stopC, _, err := stream.Klines("1m").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestContinuousKlines(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.ContinuousKlines("1m", "BTCUSDT"))
	stopC, _, err := stream.ContinuousKlines("1m", "BTCUSDT").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestPartialBookDepths(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.PartialBookDepths(common.DepthStreamLevel5, common.DepthStreamRate100ms))
	stopC, _, err := stream.PartialBookDepths(common.DepthStreamLevel5, common.DepthStreamRate100ms).SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestDiffBookDepths(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.DiffBookDepths(common.DepthStreamRate100ms))
	stopC, _, err := stream.DiffBookDepths(common.DepthStreamRate100ms).SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestAggTrades(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.AggTrades())
	stopC, _, err := stream.AggTrades().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestTrades(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.Trades())
	stopC, _, err := stream.Trades().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestBookTickers(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.BookTickers())
	stopC, _, err := stream.BookTickers().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestTickers(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.Tickers())
	stopC, _, err := stream.Tickers().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestMiniTickers(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.MiniTickers())
	stopC, _, err := stream.MiniTickers().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestUserData(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.UserData("listenKey"))
	stopC, _, err := stream.UserData("listenKey").SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestMarkPrice(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.MarkPrice())
	stopC, _, err := stream.MarkPrice().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestLiquidationOrder(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.LiquidationOrder())
	stopC, _, err := stream.LiquidationOrder().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestContractInfo(t *testing.T) {
	stream := web_stream.New(true)
	assert.NotNil(t, stream.ContractInfo())
	stopC, _, err := stream.ContractInfo().SetHandler(mockHandler).SetErrHandler(mockErrHandler).Start()
	assert.NoError(t, err)
	assert.NotNil(t, stopC)
	close(stopC)
}

func TestSubscribing(t *testing.T) {
	stream := web_stream.New(true).DynamicStream()
	assert.NotNil(t, stream)
	response, err := stream.Subscribe([]string{"btcusdt@aggTrade"})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	response, err = stream.ListOfSubscriptions()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	response, err = stream.Unsubscribe([]string{"btcusdt@aggTrade"})
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
