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

func TestKlines(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.Klines("1m").SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestContinuousKlines(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.ContinuousKlines("1m", "BTCUSDT").SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestPartialBookDepths(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.PartialBookDepths(common.DepthStreamLevel5, common.DepthStreamRate100ms).SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestDiffBookDepths(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.DiffBookDepths(common.DepthStreamRate100ms).SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestAggTrades(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.AggTrades().SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestTrades(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.Trades().SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestBookTickers(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.BookTickers().SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestTickers(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.Tickers().SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestMiniTickers(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.MiniTickers().SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestUserData(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.UserData("listenKey").SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestMarkPrice(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.MarkPrice().SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestLiquidationOrder(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.LiquidationOrder().SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestContractInfo(t *testing.T) {
	stream := web_stream.New(true)
	wrapper := stream.ContractInfo().SetSymbol("BTCUSDT").AddHandler("default", mockHandler).SetErrHandler(mockErrHandler)
	assert.NotNil(t, wrapper)
}

func TestStream(t *testing.T) {
	stream := web_stream.
		New(true).
		Stream().
		SetSymbol("BTCUSDT").
		SetErrHandler(mockErrHandler)
	stream.AddHandler("btcusdt@aggTrade", mockHandler)
	err := stream.Subscribe("btcusdt@aggTrade")
	assert.NoError(t, err)
	time.Sleep(timeOut)
	stream.Unsubscribe("btcusdt@aggTrade")
	stream.RemoveHandler("btcusdt@aggTrade")
}
