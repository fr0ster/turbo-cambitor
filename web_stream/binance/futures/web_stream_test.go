package futures_web_stream_test

import (
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	common "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/futures"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	timeOut = 1 * time.Second
)

// Mock handler for WebSocket messages
func genMockHandler(id string) func(message *simplejson.Json) {
	return func(message *simplejson.Json) {
		if message.Get("e").MustString() == id {
			logrus.Infof("Received message: %+v", message)
		}
	}
}
func mockHandler(message web_socket.LogRecord) {
	logrus.Infof("Log record: %+v", message)
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
	doneCh := make(chan struct{})
	wrapper :=
		stream.
			ContractInfo().
			SetSymbol("BTCUSDT").
			SetMessageLogger(func(message web_socket.LogRecord) {
				logrus.Infof("Received message: %+v", message)
				// if message.Get("e").MustString() == "CONTRACT_INFO" {
				doneCh <- struct{}{}
				// }
			})
	assert.NotNil(t, wrapper)
	err := wrapper.Connect()
	defer wrapper.Disconnect()
	assert.NoError(t, err)
	err = wrapper.Subscribe(func(me web_socket.MessageEvent) {
		logrus.Infof("Received message: %+v", me)
		doneCh <- struct{}{}
	}, "btcusdt@contractInfo")
	assert.NoError(t, err)
	select {
	case <-doneCh:
		t.Log("Received contract info message")
	case <-time.After(3 * timeOut):
		t.Error("Timeout waiting for connection")
		return
	}
}

func TestStream(t *testing.T) {
	stream := web_stream.
		NewDefault(true).
		Stream().
		SetSymbol("BTCUSDT")
	err := stream.Connect()
	defer stream.Disconnect()
	assert.NoError(t, err)
	err = stream.Subscribe(func(me web_socket.MessageEvent) {
		js, err := simplejson.NewJson(me.Body)
		if err == nil {
			logrus.Infof("Received message: %+v", js)
		}
	}, "btcusdt@aggTrade")
	assert.NoError(t, err)
	subscribes, err := stream.ListOfSubscriptions()
	assert.NoError(t, err)
	assert.NotNil(t, subscribes)
	time.Sleep(timeOut)
	stream.Unsubscribe("btcusdt@aggTrade")
}
