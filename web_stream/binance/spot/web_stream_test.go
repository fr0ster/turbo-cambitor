package spot_web_stream_test

import (
	"os"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	spot_rest "github.com/fr0ster/turbo-cambitor/rest_api/binance/spot"
	common "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/spot"
	"github.com/fr0ster/turbo-restler/web_socket"
	signature "github.com/fr0ster/turbo-signer/v2/signature"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	timeOut = 1 * time.Second
)

// Mock handler for WebSocket messages
func mockHandler(message web_socket.LogRecord) {
	js, err := simplejson.NewJson(message.Body)
	if err != nil {
		logrus.Errorf("Error parsing JSON: %v", err)
		return
	}
	logrus.Infof("Received message: %+v", js)
}

// Mock error handler for WebSocket errors
func mockErrHandler(err error) error {
	logrus.Errorf("Error: %v", err)
	return err
}

// API keys for endpoints that require auth (UserData listenKey)
var (
	apiKey = os.Getenv("SPOT_TEST_BINANCE_API_KEY")
	secret = os.Getenv("SPOT_TEST_BINANCE_SECRET_KEY")
	sign   = signature.NewSignHMAC(signature.PublicKey(apiKey), signature.SecretKey(secret))
)

func requireSpotAPIKeys(t *testing.T) bool {
	t.Helper()
	return apiKey != "" && secret != ""
}

func TestKlines(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.Klines("1m").SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestContinuousKlines(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.ContinuousKlines("1m", "BTCUSDT").SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestPartialBookDepths(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.PartialBookDepths(common.DepthStreamLevel5, common.DepthStreamRate100ms).SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestDiffBookDepths(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.DiffBookDepths(common.DepthStreamRate100ms).SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestAggTrades(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.AggTrades().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestTrades(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.Trades().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestBookTickers(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.BookTickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestTickers(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.Tickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestMiniTickers(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.MiniTickers().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestUserData(t *testing.T) {
	t.Parallel()
	if !assert.True(t, requireSpotAPIKeys(t), "SPOT_TEST_BINANCE_API_KEY/SECRET_KEY must be set") {
		return
	}
	ra := spot_rest.New(sign, true)
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
		_ = ra.CloseListenKey(listenKey)
	}()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestMarkPrice(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.MarkPrice().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestLiquidationOrder(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.LiquidationOrder().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestContractInfo(t *testing.T) {
	t.Parallel()
	stream := web_stream.NewDefault(true)
	wrapper, err := stream.ContractInfo().SetSymbol("BTCUSDT").SetMessageLogger(mockHandler).Connect()
	assert.NoError(t, err)
	if err != nil || wrapper == nil {
		return
	}
	wrapper.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			_ = mockErrHandler(evt.Error)
		}
	})
	defer wrapper.Disconnect()
	assert.NotNil(t, wrapper)
	time.Sleep(timeOut)
}

func TestStream(t *testing.T) {
	t.Parallel()
	stream, err := web_stream.
		NewDefault(true).
		Stream().Connect()
	defer stream.Disconnect()
	assert.NoError(t, err)
	err = stream.Subscribe(func(me web_socket.MessageEvent) {
		js, err := simplejson.NewJson(me.Body)
		if err != nil {
			logrus.Errorf("Error parsing JSON: %v", err)
			return
		}
		logrus.Infof("Received message: %+v", js)
	}, "btcusdt@aggTrade")
	assert.NoError(t, err)
	subscribes, err := stream.ListOfSubscriptions()
	assert.NoError(t, err)
	assert.NotNil(t, subscribes)
	time.Sleep(timeOut)
	stream.Unsubscribe("btcusdt@aggTrade")
	<-time.After(time.Second * 2)
}
