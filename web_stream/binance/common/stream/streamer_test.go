package streamer_test

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	mock_server "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream/mock_server"
	"github.com/sirupsen/logrus"

	streamer "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"

	"github.com/bitly/go-simplejson"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flag.Parse()

	// Start server only once for all tests
	mock_server.StartReusableMockServer(8080)
	os.Exit(m.Run())
}

func newStreamWrapper() *streamer.StreamWrapper {
	sw := streamer.New("localhost:8080", "/ws", "ws", false)
	err := sw.Connect()
	if err != nil {
		panic(err)
	}
	return sw
}

// ----------------------------
//            TESTS
// ----------------------------

func TestSubscribe(t *testing.T) {
	sw := newStreamWrapper()
	err := sw.Subscribe("btcusdt@aggTrade")
	assert.NoError(t, err)
}

func TestUnsubscribe(t *testing.T) {
	sw := newStreamWrapper()
	err := sw.Unsubscribe("btcusdt@aggTrade")
	assert.NoError(t, err)
}

func TestListOfSubscriptions(t *testing.T) {
	sw := newStreamWrapper()
	subs, err := sw.ListOfSubscriptions()
	assert.NoError(t, err)
	assert.Equal(t, []string{"btcusdt@aggTrade"}, subs)
}

func TestCall(t *testing.T) {
	sw := newStreamWrapper()
	rq := simplejson.New()
	rq.Set("method", "LIST_SUBSCRIPTIONS")
	rq.Set("id", "test-id")

	resp, err := sw.Call(rq)
	assert.NoError(t, err)
	assert.Equal(t, "test-id", resp.Get("id").MustString())
	assert.Equal(t, []interface{}{"btcusdt@aggTrade"}, resp.Get("result").MustArray())
}

func TestCallTimeout(t *testing.T) {
	sw := newStreamWrapper()

	rq := simplejson.New()
	rq.Set("method", "UNKNOWN_METHOD")
	rq.Set("id", "timeout-test")

	// Встановлюємо короткий таймаут
	sw.SetTimeOut(400 * time.Millisecond)

	start := time.Now()
	resp, err := sw.Call(rq)
	duration := time.Since(start)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.LessOrEqual(t, duration.Milliseconds(), int64(1000), "should timeout fast")
}

func TestClose(t *testing.T) {
	sw := newStreamWrapper()
	assert.NotPanics(t, func() {
		sw.Close()
	})
}

func TestSetErrHandler(t *testing.T) {
	sw := newStreamWrapper()

	called := false
	errC := make(chan error, 1)

	sw.SetErrHandler(func(err error) error {
		called = true
		select {
		case errC <- err:
		default:
			t.Logf("⚠️ errC full, dropping error: %v", err)
		}
		return err
	})

	rq := simplejson.New()
	rq.Set("method", "ERROR") // 👈 сервер обробляє цей метод
	rq.Set("id", "err-test-id")
	rq.Set("params", []interface{}{"custom-server-error"}) // 🧠 додаємо потрібну помилку

	resp, err := sw.Call(rq)

	select {
	case receivedErr := <-errC:
		assert.True(t, called, "Error handler should have been called")
		assert.Error(t, receivedErr)
		assert.Contains(t, receivedErr.Error(), "custom-server-error")
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for error handler to be called")
	}

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "custom-server-error")
}

func TestAddRemoveHandler(t *testing.T) {
	sw := newStreamWrapper()

	called := false
	sw.AddHandler("my-handler", func(js *simplejson.Json) {
		called = true
	})

	err := sw.Subscribe("btcusdt@aggTrade")
	assert.NoError(t, err)
	assert.True(t, called, "Handler should be called")
	time.Sleep(1 * time.Second)

	sw.RemoveHandler("my-handler")
	called = sw.GetLoopStarted()
	assert.False(t, called, "Handler should be removed")
}

func TestGracefulCloseStreamer(t *testing.T) {
	sw := newStreamWrapper()

	rq := simplejson.New()
	rq.Set("method", "CLOSE_GRACEFUL")
	rq.Set("id", "graceful-close")

	resp, err := sw.Call(rq)

	assert.Error(t, err, "should return error due to closed connection")
	assert.Nil(t, resp, "no response expected on graceful close")
	assert.Contains(t, err.Error(), "close", "error should mention closed connection")
}

func TestAbruptCloseStreamer(t *testing.T) {
	sw := newStreamWrapper()

	rq := simplejson.New()
	rq.Set("method", "CLOSE_ABRUPT")
	rq.Set("id", "abrupt-close")

	resp, err := sw.Call(rq)

	assert.Error(t, err, "should return error due to abrupt close")
	assert.Nil(t, resp, "no response expected on abrupt close")
	assert.Contains(t, err.Error(), "close", "error should mention closed connection")
}

var closeNormalMutex = &sync.Mutex{}
var abruptCloseMutex = &sync.Mutex{}

func TestGracefulCloseStreamer_WithRestart(t *testing.T) {
	closeNormalMutex.Lock()
	defer closeNormalMutex.Unlock()
	sw := newStreamWrapper()

	rq := simplejson.New()
	rq.Set("method", "CLOSE_GRACEFUL")
	rq.Set("id", "graceful-close")
	rq.SetPath([]string{"params", "restart"}, 500)

	resp, err := sw.Call(rq)
	assert.Error(t, err, "should return error due to graceful close")
	assert.Nil(t, resp, "no response expected")

	ok := retryConnect(sw, 10, 1000*time.Millisecond)
	assert.True(t, ok, "client should reconnect after CLOSE_GRACEFUL with restart")
}

func TestAbruptCloseStreamer_WithRestart(t *testing.T) {
	abruptCloseMutex.Lock()
	defer abruptCloseMutex.Unlock()
	sw := newStreamWrapper()

	rq := simplejson.New()
	rq.Set("method", "CLOSE_ABRUPT")
	rq.Set("id", "abrupt-close")
	rq.SetPath([]string{"params", "restart"}, 500)

	resp, err := sw.Call(rq)
	assert.Error(t, err, "should return error due to abrupt close")
	assert.Nil(t, resp, "no response expected")

	ok := retryConnect(sw, 10, 200*time.Millisecond)
	assert.True(t, ok, "client should reconnect after CLOSE_ABRUPT with restart")
}

func retryConnect(sw *streamer.StreamWrapper, maxAttempts int, delay time.Duration) bool {
	for i := 0; i < maxAttempts; i++ {
		time.Sleep(delay)

		logrus.Infof("🔁 Attempting manual reconnect (%d/%d)", i+1, maxAttempts)
		if err := sw.Reconnect(); err == nil {
			time.Sleep(300 * time.Millisecond) // 🧘 дати з'єднанню стабілізуватись
			rq := simplejson.New()
			rq.Set("method", "LIST_SUBSCRIPTIONS")
			rq.Set("id", fmt.Sprintf("recheck-%d", i))
			if _, err := sw.Call(rq); err == nil {
				return true
			}
		}
	}
	return false
}

func TestGracefulCloseStreamer_WithAutoReconnect(t *testing.T) {
	closeNormalMutex.Lock()
	defer closeNormalMutex.Unlock()
	sw := newStreamWrapper()
	sw.EnableAutoReconnect(1000 * time.Millisecond)

	defer sw.DisableAutoReconnect()

	rq := simplejson.New()
	rq.Set("method", "CLOSE_GRACEFUL")
	rq.Set("id", "graceful-close")
	rq.SetPath([]string{"params", "restart"}, 500)

	resp, err := sw.Call(rq)
	assert.Error(t, err, "should return error due to graceful close")
	assert.Nil(t, resp)

	// 🧘‍♂️ Дати серверу стартануть після рестарту
	time.Sleep(500 * time.Millisecond)

	ok := waitUntilConnected(sw, 10, 300*time.Millisecond)
	assert.True(t, ok, "should auto-reconnect after CLOSE_GRACEFUL with restart")
}

func TestAbruptCloseStreamer_WithAutoReconnect(t *testing.T) {
	abruptCloseMutex.Lock()
	defer abruptCloseMutex.Unlock()
	sw := newStreamWrapper()
	sw.EnableAutoReconnect(200 * time.Millisecond)

	defer sw.DisableAutoReconnect()

	rq := simplejson.New()
	rq.Set("method", "CLOSE_ABRUPT")
	rq.Set("id", "abrupt-close")
	rq.SetPath([]string{"params", "restart"}, 500)

	resp, err := sw.Call(rq)
	assert.Error(t, err, "should return error due to abrupt close")
	assert.Nil(t, resp)

	ok := waitUntilConnected(sw, 10, 300*time.Millisecond)
	assert.True(t, ok, "should auto-reconnect after CLOSE_ABRUPT with restart")
}

func waitUntilConnected(sw *streamer.StreamWrapper, maxAttempts int, delay time.Duration) bool {
	for i := 0; i < maxAttempts; i++ {
		time.Sleep(delay)

		// if !sw.IsLoopStarted() {
		// 	logrus.Warnf("⏳ loop not started on attempt %d", i)
		// 	continue
		// }

		rq := simplejson.New()
		rq.Set("method", "LIST_SUBSCRIPTIONS")
		rq.Set("id", fmt.Sprintf("auto-recheck-%d", i))

		if _, err := sw.Call(rq); err == nil {
			logrus.Infof("✅ Connected and received subscriptions at attempt %d", i)
			sw.DisableAutoReconnect()
			return true
		}

		logrus.Warnf("❌ Call failed at attempt %d", i)
	}
	return false
}
