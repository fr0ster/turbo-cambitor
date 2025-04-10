package spot_web_stream_test

import (
	"testing"
	"time"

	common "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/spot"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	timeOut = 1 * time.Second
)

// Mock handler for WebSocket messages
func mockHandler(message web_socket.LogRecord) {
	logrus.Infof("Received message: %+v", message)
}

// Mock error handler for WebSocket errors
func mockErrHandler(err error) error {
	logrus.Errorf("Error: %v", err)
	return err
}

func TestKlines(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.Klines("1m").SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestContinuousKlines(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.ContinuousKlines("1m", "BTCUSDT").SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestPartialBookDepths(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.PartialBookDepths(common.DepthStreamLevel5, common.DepthStreamRate100ms).SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestDiffBookDepths(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.DiffBookDepths(common.DepthStreamRate100ms).SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestAggTrades(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.AggTrades().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestTrades(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.Trades().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestBookTickers(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.BookTickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestTickers(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.Tickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestMiniTickers(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.MiniTickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestUserData(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.UserData("listenKey").SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestMarkPrice(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.MarkPrice().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestLiquidationOrder(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.LiquidationOrder().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestContractInfo(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper := stream.ContractInfo().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler)
	assert.NotNil(t, wrapper)
}

func TestStream(t *testing.T) {
	stream := web_stream.
		NewDefault(true).
		Stream().
		SetSymbol("BTCUSDT").
		SetMessageLogger(mockHandler)
	err := stream.Connect()
	defer stream.Disconnect()
	assert.NoError(t, err)
	err = stream.Subscribe(func(me web_socket.MessageEvent) {
		logrus.Infof("Received message: %+v", me)
	}, "btcusdt@aggTrade")
	assert.NoError(t, err)
	time.Sleep(timeOut)
	stream.Unsubscribe("btcusdt@aggTrade")
}
