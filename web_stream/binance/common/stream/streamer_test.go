package streamer_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/fr0ster/turbo-cambitor/common"
	mock_server "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream/mock_server"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/google/uuid"

	streamer "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var mockPort int

func TestMain(m *testing.M) {
	flag.Parse()

	// Start server only once for all tests
	mockPort = mock_server.StartReusableMockServer(0)
	// run tests
	os.Exit(m.Run())
}

func newStreamWrapper() *streamer.StreamWrapper {
	scheme := "ws"
	host := fmt.Sprintf("localhost:%d", mockPort)
	endpoint := "/ws"
	factory := func() (web_socket.WebSocketCommonInterface, error) {
		url := fmt.Sprintf("%s://%s%s", scheme, host, endpoint)
		return web_socket.NewWebSocketWrapper(websocket.DefaultDialer, url)
	}

	sw := streamer.NewStreamWrapper(factory, common.WsScheme(scheme), common.WsHost(host), common.WsEndpoint(endpoint))

	// Loosen deadlines to avoid spurious i/o timeouts in CI
	sw.SetReadTimeout(3 * time.Second)
	sw.SetWriteTimeout(3 * time.Second)
	// Be a bit more patient on readiness in tests
	sw.SetModeSync(3 * time.Second)

	if _, err := sw.Connect(); err != nil {
		panic(err)
	}

	return sw
}

// ----------------------------
//            TESTS
// ----------------------------

// Test for Subscribe, ListOfSubscriptions, and Unsubscribe methods
func TestSubscribers(t *testing.T) {
	t.Parallel()
	sw := newStreamWrapper()
	defer func() {
		err := sw.Unsubscribe("btcusdt@aggTrade")
		assert.NoError(t, err)
	}()
	err := sw.Subscribe(func(me web_socket.MessageEvent) {
		if me.Error != nil {
			t.Logf("❌ Received error: %v", me.Error)
			return
		}
		t.Logf("📨 Received message: %s", string(me.Body))
	},
		"btcusdt@aggTrade")
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	subs, err := sw.ListOfSubscriptions()
	assert.NoError(t, err)
	assert.Equal(t, []string{"btcusdt@aggTrade"}, subs)
	time.Sleep(1 * time.Second)
}
func TestMultiSubscribers(t *testing.T) {
	t.Parallel()
	sw := newStreamWrapper()

	// Перелік стрімів
	streams := []string{
		"btcusdt@aggTrade",
		"btcusdt@depth",
		"btcusdt@kline_1m",
	}

	defer func() {
		err := sw.Unsubscribe(streams...)
		assert.NoError(t, err)
	}()

	received := make(map[string]bool)
	var mu sync.Mutex

	err := sw.Subscribe(func(me web_socket.MessageEvent) {
		if me.Error != nil {
			t.Logf("❌ Received error: %v", me.Error)
			return
		}

		js, err := simplejson.NewJson(me.Body)
		if err != nil {
			t.Logf("❌ JSON parse error: %v", err)
			t.Logf("↩️ Raw body: %s", string(me.Body))
			return
		}

		if e := js.Get("error").MustString(); e != "" {
			t.Logf("❌ Message has error: %s", e)
			return
		}

		stream := js.Get("stream").MustString()
		if stream == "" {
			t.Logf("⚠️ Received message without 'stream': %s", string(me.Body))
			return
		}

		t.Logf("📨 Received message from stream: %s", stream)
		mu.Lock()
		received[stream] = true
		mu.Unlock()
	}, streams...)

	assert.NoError(t, err)

	// Дочекайся повідомлень
	time.Sleep(1 * time.Second)

	subs, err := sw.ListOfSubscriptions()
	assert.NoError(t, err)
	assert.ElementsMatch(t, streams, subs)

	// Перевірка, що всі стріми щось отримали
	for _, s := range streams {
		mu.Lock()
		_, ok := received[s]
		mu.Unlock()
		assert.True(t, ok, "❌ No message received from stream: %s", s)
	}

	time.Sleep(1 * time.Second)
}

func TestParallelSubscribers(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup
	clientCount := 5
	streams := []string{"btcusdt@aggTrade", "ethusdt@depth", "bnbusdt@kline_1m"}

	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			sw := newStreamWrapper()
			defer func() {
				for _, s := range streams {
					_ = sw.Unsubscribe(s)
				}
				sw.Disconnect()
			}()

			// Buffer only expected unique stream events and avoid blocking
			received := make(chan string, len(streams))
			seen := make(map[string]bool)
			var seenMu sync.Mutex

			handler := func(me web_socket.MessageEvent) {
				if me.Error != nil {
					t.Errorf("❌ Client %d received error: %v", clientID, me.Error)
					return
				}

				var msg map[string]interface{}
				err := json.Unmarshal(me.Body, &msg)
				if err != nil {
					t.Errorf("❌ Client %d received invalid JSON: %s", clientID, me.Body)
					return
				}

				// (optional) handle only mock_data messages
				if msg["type"] != "mock_data" {
					t.Logf("ℹ️ Client %d received non-data message: %s", clientID, me.Body)
					return
				}

				streamRaw, ok := msg["stream"]
				if !ok {
					t.Errorf("⚠️ Client %d: message has no 'stream': %s", clientID, me.Body)
					return
				}
				stream, ok := streamRaw.(string)
				if !ok {
					t.Errorf("⚠️ Client %d: 'stream' is not a string: %v", clientID, streamRaw)
					return
				}
				// Deduplicate and send non-blocking to avoid stalling the reader
				seenMu.Lock()
				if !seen[stream] {
					seen[stream] = true
					select {
					case received <- stream:
					default:
						// drop if buffer full
					}
				}
				seenMu.Unlock()
			}

			// Subscribe to each stream
			for _, s := range streams {
				err := sw.Subscribe(handler, s)
				assert.NoError(t, err, "Client %d failed to subscribe to %s", clientID, s)
			}

			time.Sleep(1 * time.Second)

			// Check active subscriptions
			subs, err := sw.ListOfSubscriptions()
			assert.NoError(t, err, "Client %d failed to list subscriptions", clientID)
			assert.ElementsMatch(t, streams, subs, "Client %d subscriptions mismatch", clientID)

			// Wait for some messages
			timeout := time.After(3 * time.Second)
			streamsReceived := make(map[string]bool)

		loop:
			for {
				select {
				case s := <-received:
					t.Logf("📨 Client %d received stream: %s", clientID, s)
					streamsReceived[s] = true
					if len(streamsReceived) == len(streams) {
						break loop
					}
				case <-timeout:
					t.Errorf("⏱ Client %d did not receive all streams in time", clientID)
					break loop
				}
			}

			for _, s := range streams {
				assert.True(t, streamsReceived[s], "Client %d missing stream %s", clientID, s)
			}
		}(i)
	}

	wg.Wait()
}

// testListOfSubscriptions uses a short-lived connection to avoid interfering with active streams
func testListOfSubscriptions(t *testing.T, factory func() (web_socket.WebSocketCommonInterface, error), timeout time.Duration) ([]string, error) {
	t.Helper()
	tmp, err := factory()
	if err != nil {
		return nil, err
	}
	tmp.Open()
	defer tmp.Close()
	if timeout > 0 {
		tmp.SetReadTimeout(timeout)
		tmp.SetWriteTimeout(timeout)
	}
	_ = tmp.WaitStarted()

	rq := simplejson.New()
	rq.Set("method", "LIST_SUBSCRIPTIONS")
	rq.Set("id", uuid.New().String())
	body, _ := rq.MarshalJSON()
	if err := tmp.Send(web_socket.WriteEvent{Body: body}); err != nil {
		return nil, err
	}
	ch := make(chan []byte, 1)
	sid := tmp.Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Kind == web_socket.KindData {
			select {
			case ch <- evt.Body:
			default:
			}
		}
	})
	defer tmp.Unsubscribe(sid)

	var resp []byte
	select {
	case resp = <-ch:
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout")
	}
	js, err := simplejson.NewJson(resp)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	for _, v := range js.Get("result").MustArray() {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out, nil
}

func TestCall(t *testing.T) {
	t.Parallel()
	sw := newStreamWrapper()
	rq := simplejson.New()
	rq.Set("method", "LIST_SUBSCRIPTIONS")
	rq.Set("id", "test-id")

	resp, err := sw.Call(rq)
	assert.NoError(t, err)
	assert.Equal(t, "test-id", resp.Get("id").MustString())
	assert.Equal(t, []interface{}{}, resp.Get("result").MustArray())
}

func TestGracefulCloseStreamer(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
		assert.Equal(t, []interface{}{}, resp.Get("result").MustArray())
	}()
}

func TestAbruptCloseStreamer_WithRestart(t *testing.T) {
	t.Parallel()
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
		assert.Equal(t, []interface{}{}, resp.Get("result").MustArray())
	}()
}

func TestGracefulCloseStreamer_WithAutoReconnect(t *testing.T) {
	t.Parallel()
	closeMutex.Lock()
	defer closeMutex.Unlock()
	sw := newStreamWrapper()
	sw.SetMaxReconnectAttempts(30).SetReconnectInterval(200 * time.Millisecond).EnableAutoReconnect()

	defer sw.DisableAutoReconnect()

	rq := simplejson.New()
	rq.Set("method", "CLOSE_GRACEFUL")
	rq.Set("id", "graceful-close")
	rq.SetPath([]string{"params", "restart"}, 500)

	resp, err := sw.Call(rq)
	assert.Error(t, err, "should return error due to graceful close")
	assert.Nil(t, resp)

	// Дочекайся відновлення коннекту (через сервісний канал/монітор)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	assert.NoError(t, sw.WaitConnected(ctx), "should auto-recover after CLOSE_GRACEFUL with restart")
	rq2 := simplejson.New()
	rq2.Set("method", "LIST_SUBSCRIPTIONS")
	rq2.Set("id", "post-recover")
	resp2, err2 := sw.Call(rq2)
	assert.NoError(t, err2)
	assert.NotNil(t, resp2)
}

func TestAbruptCloseStreamer_WithAutoReconnect(t *testing.T) {
	t.Parallel()
	closeMutex.Lock()
	defer closeMutex.Unlock()
	sw := newStreamWrapper()
	// Give more room for abrupt close + server restart timing
	sw.SetMaxReconnectAttempts(10).SetReconnectInterval(200 * time.Millisecond).EnableAutoReconnect()

	rq := simplejson.New()
	rq.Set("method", "CLOSE_ABRUPT")
	rq.Set("id", "abrupt-close")
	rq.SetPath([]string{"params", "restart"}, 500)

	resp, err := sw.Call(rq)
	assert.Error(t, err, "should return error due to abrupt close")
	assert.Nil(t, resp)

	// Дочекайся відновлення коннекту (через сервісний канал/монітор)
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	assert.NoError(t, sw.WaitConnected(ctx), "should auto-recover after CLOSE_ABRUPT with restart")
	// Після відновлення можна вимкнути авто-реконнект і виконати простий Call
	sw.DisableAutoReconnect()
	rq2 := simplejson.New()
	rq2.Set("method", "LIST_SUBSCRIPTIONS")
	rq2.Set("id", "post-recover")
	resp2, err2 := sw.Call(rq2)
	assert.NoError(t, err2)
	assert.NotNil(t, resp2)
}

// removed waitUntilConnected: tests now rely on cambitor's WaitHealthy

// waitServerListening активно перевіряє, що TCP-порт слухає з'єднання
// func waitServerListening(addr string, timeout time.Duration) bool {
//     deadline := time.Now().Add(timeout)
//     for {
//         if time.Now().After(deadline) {
//             return false
//         }
//         conn, err := net.DialTimeout("tcp", addr, 250*time.Millisecond)
//         if err == nil {
//             _ = conn.Close()
//             return true
//         }
//         time.Sleep(150 * time.Millisecond)
//     }
// }

// waitForPong надсилає Ping і очікує Pong за короткий таймаут
// func waitForPong(conn web_socket.WebSocketCommonInterface, timeout time.Duration) bool {
//     if conn == nil || !conn.IsStarted() {
//         return false
//     }
//     done := make(chan struct{}, 1)
//     // Тимчасовий pong handler
//     conn.SetPongHandler(func(string) error {
//         select {
//         case done <- struct{}{}:
//         default:
//         }
//         return nil
//     })
//     // Надсилаємо Ping control
//     deadline := time.Now().Add(timeout)
//     if err := conn.GetControl().WriteControl(websocket.PingMessage, []byte("ping"), deadline); err != nil {
//         return false
//     }
//     select {
//     case <-done:
//         return true
//     case <-time.After(timeout):
//         return false
//     }
// }

func TestPingPongHandlingWithActiveStream(t *testing.T) {
	t.Parallel()
	sw := newStreamWrapper()

	// pingCalled := false
	messageReceived := false

	// Встановлюємо Ping handler
	// sw.GetConnection().SetPingHandler(func(appData string, ctrl web_socket.ControlWriter) error {
	// 	pingCalled = true
	// 	t.Logf("📡 Ping received: %s", appData)
	// 	time.Sleep(100 * time.Millisecond)
	// 	t.Log("✅ Simulated business logic")
	// 	return ctrl.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(1*time.Second))
	// })

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
	err = sw.Subscribe(func(me web_socket.MessageEvent) {
		if me.Error != nil {
			t.Logf("❌ Received error: %v", me.Error)
			return
		}
		messageReceived = true
		t.Logf("📨 Received message: %s", string(me.Body))
	}, "btcusdt@aggTrade")
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	// assert.True(t, pingCalled, "Ping handler should have been called")
	assert.True(t, messageReceived, "WebSocket message should have been received")
}

func TestStreamWrapperLoggerCalled(t *testing.T) {
	t.Parallel()
	var logCalled bool
	var logErr error

	logChan := make(chan web_socket.LogRecord, 1)

	sw := newStreamWrapper()
	sw.SetMessageLogger(func(r web_socket.LogRecord) {
		logCalled = true
		logErr = r.Err
		select {
		case logChan <- r:
		default:
			// Don't block if the channel is full
		}
	})

	// Відправляємо щось, щоб викликати логування на receive
	err := sw.Subscribe(func(me web_socket.MessageEvent) {
		if me.Error != nil {
			t.Logf("❌ Received error: %v", me.Error)
			return
		}
		t.Logf("📨 Received message: %s", string(me.Body))
	}, "btcusdt@aggTrade")
	assert.NoError(t, err)

	select {
	case rec := <-logChan:
		t.Logf("📘 Received log: op=%s body=%s err=%v", rec.Op, string(rec.Body), rec.Err)
		assert.Equal(t, web_socket.OpSend, rec.Op, "should log a send operation")
		assert.NotNil(t, string(rec.Body), "subscribe", "should include subscribe payload")
		assert.Nil(t, rec.Err)
	case <-time.After(2 * time.Second):
		t.Fatal("expected log record not received")
	}

	assert.True(t, logCalled, "logger should have been called")
	assert.Nil(t, logErr)

	sw.Disconnect()
}

func TestReadWriteTimeout(t *testing.T) {
	t.Parallel()
	sw := newStreamWrapper()
	// Встановлюємо таймаут на читання
	sw.SetReadTimeout(1 * time.Second)
	// Встановлюємо таймаут на запис
	sw.SetWriteTimeout(1 * time.Second)

	rq := simplejson.New()
	rq.Set("method", "LIST_SUBSCRIPTIONS")
	rq.Set("id", "test-id")

	resp, err := sw.Call(rq)
	assert.NoError(t, err)
	assert.Equal(t, "test-id", resp.Get("id").MustString())
	assert.Equal(t, []interface{}{}, resp.Get("result").MustArray())
}

// Ensure helper testListOfSubscriptions is used
func Test_ListOfSubscriptions_Helper(t *testing.T) {
	// Не запускаємо паралельно; синхронізуємося з тестами, що перезапускають сервер
	closeMutex.Lock()
	defer closeMutex.Unlock()
	scheme := "ws"
	host := fmt.Sprintf("localhost:%d", mockPort)
	endpoint := "/ws"
	factory := func() (web_socket.WebSocketCommonInterface, error) {
		url := fmt.Sprintf("%s://%s%s", scheme, host, endpoint)
		return web_socket.NewWebSocketWrapper(websocket.DefaultDialer, url)
	}
	subs, err := testListOfSubscriptions(t, factory, 2500*time.Millisecond)
	assert.NoError(t, err)
	assert.NotNil(t, subs)
}
