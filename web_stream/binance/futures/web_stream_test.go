package futures_web_stream_test

import (
	"os"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	futures_rest "github.com/fr0ster/turbo-cambitor/rest_api/binance/futures"
	common "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/futures"
	"github.com/fr0ster/turbo-restler/web_socket"
	signature "github.com/fr0ster/turbo-signer/v2/signature"
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

// API keys for endpoints that require auth (e.g., creating listenKey for UserData stream)
var (
	apiKey = os.Getenv("FUTURE_TEST_BINANCE_API_KEY")
	secret = os.Getenv("FUTURE_TEST_BINANCE_SECRET_KEY")
	sign   = signature.NewSignHMAC(signature.PublicKey(apiKey), signature.SecretKey(secret))
)

func requireFuturesAPIKeys(t *testing.T) bool {
	t.Helper()
	return apiKey != "" && secret != ""
}

func TestKlines(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.Klines("1m").SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	time.Sleep(2 * timeOut)
	assert.NotNil(t, wrapper)
}

func TestPartialBookDepths(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.PartialBookDepths(common.DepthStreamLevel5, common.DepthStreamRate100ms).SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestDiffBookDepths(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.DiffBookDepths(common.DepthStreamRate100ms).SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestAggTrades(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.AggTrades().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestTrades(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.Trades().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestBookTickers(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.BookTickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestTickers(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.Tickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestMiniTickers(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.MiniTickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestUserData(t *testing.T) {
	t.Parallel()
	// Require API keys to obtain a real listenKey
	if !assert.True(t, requireFuturesAPIKeys(t), "FUTURE_TEST_BINANCE_API_KEY/SECRET_KEY must be set") {
		return
	}
	ra := futures_rest.New(sign, true)
	listenKey, err := ra.ListenKey()
	assert.NoError(t, err)
	assert.NotEmpty(t, listenKey)

	stream := web_stream.NewDefault(true)
	wrapper, err := stream.UserData(listenKey).SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer func() {
		wrapper.Disconnect()
		_ = ra.CloseListenKey()
	}()
	assert.NotNil(t, wrapper)
}

func TestMarkPrice(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.MarkPrice().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestLiquidationOrder(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.LiquidationOrder().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
}

func TestContractInfo(t *testing.T) {
	t.Parallel()
	doneCh := make(chan struct{}, 1)
	stream, err := web_stream.NewDefault(true).
		Stream().
		SetMessageLogger(mockHandler).
		Connect()
	// SetMessageLogger(mockHandler).
	// Connect()
	assert.NoError(t, err)
	assert.NotNil(t, stream)
	stream.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
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
	time.Sleep(2 * timeOut)
	select {
	case <-doneCh:
		t.Log("Received contract info message")
	case <-time.After(3 * timeOut):
		t.Error("Timeout waiting for contract info message")
		return
	}
	err = stream.Unsubscribe("btcusdt@contractInfo")
	assert.NoError(t, err)
}

func TestStream(t *testing.T) {
	t.Parallel()
	stream, err := web_stream.
		NewDefault(true).
		Stream().Connect()
	defer stream.Disconnect()
	assert.NoError(t, err)
	stream.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
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
	t.Parallel()
	stream, err := web_stream.
		New("fstream.binance.com/ws", "btcusdt@aggTrade", "wss").
		Stream().SetMessageLogger(mockHandler).
		SetReadTimeout(5 * time.Second).
		SetWriteTimeout(5 * time.Second).
		Connect()
	defer stream.Disconnect()
	assert.NoError(t, err)
	stream.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	// Ensure the connection is responsive before querying subscriptions (best effort)
	_ = stream.WaitForPong(2 * time.Second)
	// Give the public stream a moment to establish before querying subscriptions
	time.Sleep(2 * timeOut)
	assert.NoError(t, err)
	subscribes, err := stream.ListOfSubscriptions()
	if err != nil {
		// Retry once if the first management call races with start-up
		time.Sleep(1500 * time.Millisecond)
		subscribes, err = stream.ListOfSubscriptions()
	}
	assert.NoError(t, err)
	assert.NotNil(t, subscribes)
	time.Sleep(2 * timeOut)
	stream.Unsubscribe("btcusdt@aggTrade")
}
