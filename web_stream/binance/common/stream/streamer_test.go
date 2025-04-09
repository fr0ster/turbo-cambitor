package streamer_test

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	mock_server "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream/mock_server"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/gorilla/websocket"
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
	factory := func() (web_socket.WebSocketInterface, error) {
		url := "ws://localhost:8080/ws"
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			return nil, err
		}
		return web_socket.NewWebSocketWrapper(conn), nil
	}

	sw := streamer.NewStreamWrapper(factory, "/ws").SetSymbol("btcusdt")

	if err := sw.Connect(); err != nil {
		panic(err)
	}

	return sw.(*streamer.StreamWrapper)
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

var closeMutex = &sync.Mutex{}

func TestGracefulCloseStreamer_WithRestart(t *testing.T) {
	closeMutex.Lock()
	defer closeMutex.Unlock()
	func() {
		sw := newStreamWrapper()

		rq := simplejson.New()
		rq.Set("method", "CLOSE_GRACEFUL")
		rq.Set("id", "graceful-close")
		rq.SetPath([]string{"params", "restart"}, 500)

		resp, err := sw.Call(rq)
		assert.Error(t, err, "should return error due to graceful close")
		assert.Nil(t, resp, "no response expected")

		err = sw.Reconnect(10, 1000*time.Millisecond)
		assert.True(t, err == nil, "client should reconnect after CLOSE_GRACEFUL with restart")
	}()
	func() {
		sw := newStreamWrapper()
		rq := simplejson.New()
		rq.Set("method", "LIST_SUBSCRIPTIONS")
		rq.Set("id", "test-id")

		resp, err := sw.Call(rq)
		assert.NoError(t, err)
		assert.Equal(t, "test-id", resp.Get("id").MustString())
		assert.Equal(t, []interface{}{"btcusdt@aggTrade"}, resp.Get("result").MustArray())
	}()
}

func TestAbruptCloseStreamer_WithRestart(t *testing.T) {
	closeMutex.Lock()
	defer closeMutex.Unlock()
	func() {
		sw := newStreamWrapper()

		rq := simplejson.New()
		rq.Set("method", "CLOSE_ABRUPT")
		rq.Set("id", "abrupt-close")
		rq.SetPath([]string{"params", "restart"}, 500)

		resp, err := sw.Call(rq)
		assert.Error(t, err, "should return error due to abrupt close")
		assert.Nil(t, resp, "no response expected")

		err = sw.Reconnect(10, 200*time.Millisecond)
		assert.True(t, err == nil, "client should reconnect after CLOSE_ABRUPT with restart")
	}()
	func() {
		sw := newStreamWrapper()
		rq := simplejson.New()
		rq.Set("method", "LIST_SUBSCRIPTIONS")
		rq.Set("id", "test-id")

		resp, err := sw.Call(rq)
		assert.NoError(t, err)
		assert.Equal(t, "test-id", resp.Get("id").MustString())
		assert.Equal(t, []interface{}{"btcusdt@aggTrade"}, resp.Get("result").MustArray())
	}()
}

func TestGracefulCloseStreamer_WithAutoReconnect(t *testing.T) {
	closeMutex.Lock()
	defer closeMutex.Unlock()
	sw := newStreamWrapper()
	sw.SetMaxReconnectAttempts(3).SetReconnectInterval(200 * time.Millisecond).EnableAutoReconnect()

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
	closeMutex.Lock()
	defer closeMutex.Unlock()
	sw := newStreamWrapper()
	sw.SetMaxReconnectAttempts(3).SetReconnectInterval(200 * time.Millisecond).EnableAutoReconnect()

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

func TestPingPongHandlingWithActiveStream(t *testing.T) {
	sw := newStreamWrapper()

	pingCalled := false
	messageReceived := false

	// Встановлюємо Ping handler
	sw.GetConnection().SetPingHandler(func(appData string, ctrl web_socket.ControlWriter) error {
		pingCalled = true
		t.Logf("📡 Ping received: %s", appData)
		time.Sleep(100 * time.Millisecond)
		t.Log("✅ Simulated business logic")
		return ctrl.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(1*time.Second))
	})

	// Підписка на WebSocket-повідомлення
	sw.GetConnection().Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			t.Logf("❌ Received error: %v", evt.Error)
			return
		}
		messageReceived = true
		t.Logf("📨 Received message: %s", string(evt.Body))
	})

	// Стартуємо ping/pong сценарій
	rq := simplejson.New()
	rq.Set("method", "PONG_CONTROL")
	rq.Set("id", "ping_control")
	rq.SetPath([]string{"params", "timeout"}, 100)

	resp, err := sw.Call(rq)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Створюємо бізнес-підписку
	err = sw.Subscribe("btcusdt@aggTrade")
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	assert.True(t, pingCalled, "Ping handler should have been called")
	assert.True(t, messageReceived, "WebSocket message should have been received")
}
