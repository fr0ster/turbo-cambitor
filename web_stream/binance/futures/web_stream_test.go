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

func mockHandler(message web_socket.LogRecord) {
	jr, err := simplejson.NewJson(message.Body)
	if err != nil {
		logrus.Errorf("Error parsing JSON: %v", err)
		return
	}
	logrus.Infof("Log record: %+v", jr)
}

// Mock error handler for WebSocket errors
func mockErrHandler(err error) error {
	logrus.Errorf("Error: %v", err)
	return err
}

func TestKlines(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.Klines("1m").SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	time.Sleep(timeOut)
	assert.NotNil(t, wrapper)
}

func TestContinuousKlines(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.ContinuousKlines("1m", "BTCUSDT").SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestPartialBookDepths(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.PartialBookDepths(common.DepthStreamLevel5, common.DepthStreamRate100ms).SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestDiffBookDepths(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.DiffBookDepths(common.DepthStreamRate100ms).SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestAggTrades(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.AggTrades().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestTrades(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.Trades().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestBookTickers(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.BookTickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestTickers(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.Tickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestMiniTickers(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.MiniTickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestUserData(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.UserData("listenKey").SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestMarkPrice(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.MarkPrice().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestLiquidationOrder(t *testing.T) {
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.LiquidationOrder().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestContractInfo(t *testing.T) {
	doneCh := make(chan struct{}, 1)
	stream, err := web_stream.NewDefault(true).
		ContractInfo().
		SetSymbol("BTCUSDT").
		SetMessageLogger(func(message web_socket.LogRecord) {
			if message.Body == nil {
				return
			}
			js, err := simplejson.NewJson(message.Body)
			if err != nil {
				logrus.Errorf("Error parsing JSON: %v", err)
				logrus.Errorf("Message: %s", message.Body)
				return
			}
			kind := web_socket.MessageKind(js.Get("Kind").MustInt())
			if kind == web_socket.KindData {
				b := js.Get("Body").MustString()
				if b == "" {
					return
				}
				logrus.Infof("Received message: %s", b)
				// if message.Get("e").MustString() == "CONTRACT_INFO" {
				select {
				case doneCh <- struct{}{}:
				default:
				}
				// }
			}
		}).Connect()
	// SetMessageLogger(mockHandler).
	// Connect()
	assert.NoError(t, err)
	assert.NotNil(t, stream)
	defer stream.Disconnect()
	assert.NoError(t, err)
	err = stream.Subscribe(func(me web_socket.MessageEvent) {
		js, err := simplejson.NewJson(me.Body)
		if err != nil {
			logrus.Errorf("Error parsing JSON: %v", err)
			return
		}
		logrus.Infof("Received message: %+v", js)
		doneCh <- struct{}{}
	}, "btcusdt@contractInfo")
	// err = stream.Subscribe(nil, "btcusdt@contractInfo")
	assert.NoError(t, err)
	subs, err := stream.ListOfSubscriptions()
	assert.NoError(t, err)
	assert.NotNil(t, subs)
	time.Sleep(timeOut)
	select {
	case <-doneCh:
		t.Log("Received contract info message")
	case <-time.After(3 * timeOut):
		t.Error("Timeout waiting for connection")
		return
	}
	err = stream.Unsubscribe("btcusdt@contractInfo")
	assert.NoError(t, err)
}

func TestStream(t *testing.T) {
	stream, err := web_stream.
		NewDefault(true).
		Stream().Connect()
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

func TestStaticStream(t *testing.T) {
	stream, err := web_stream.
		New("fstream.binance.com/ws", "btcusdt@aggTrade", "wss").
		Stream().SetMessageLogger(mockHandler).Connect()
	defer stream.Disconnect()
	assert.NoError(t, err)
	time.Sleep(timeOut)
	assert.NoError(t, err)
	subscribes, err := stream.ListOfSubscriptions()
	assert.NoError(t, err)
	assert.NotNil(t, subscribes)
	time.Sleep(timeOut)
	stream.Unsubscribe("btcusdt@aggTrade")
}
