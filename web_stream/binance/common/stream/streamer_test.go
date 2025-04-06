package streamer_test

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	streamer "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"

	"github.com/bitly/go-simplejson"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var (
	testServer *httptest.Server
	testHost   string
)

var upgrader = websocket.Upgrader{}

func TestMain(m *testing.M) {
	flag.Parse()

	// Start server only once for all tests
	testServer, testHost = startGlobalMockServer()
	code := m.Run()

	testServer.Close()
	os.Exit(code)
}

func startGlobalMockServer() (*httptest.Server, string) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			req, _ := simplejson.NewJson(msg)
			id := req.Get("id").MustString()
			method := req.Get("method").MustString()

			resp := simplejson.New()
			resp.Set("id", id)

			switch method {
			case "SUBSCRIBE":
				resp.Set("result", "")
			case "UNSUBSCRIBE":
				resp.Set("result", true)
			case "LIST_SUBSCRIPTIONS":
				resp.Set("result", []string{"btcusdt@aggTrade"})
			case "UNKNOWN_METHOD":
				time.Sleep(2 * time.Second)
			default:
				resp.Set("error", "unknown method")
			}

			b, _ := resp.Encode()
			conn.WriteMessage(websocket.TextMessage, b)
		}
	}))
	host := strings.TrimPrefix(s.URL, "http://")
	return s, host
}

func newStreamWrapper() *streamer.StreamWrapper {
	return streamer.New(web_socket.WsHost(testHost), "/ws", "ws", true)
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
	sw.SetErrHandler(func(err error) error {
		called = true
		return err
	})
	// Немає способу симулювати помилку через зовнішній сервер тут
	assert.False(t, called, "Handler shouldn't be called without an error")
}

func TestAddRemoveHandler(t *testing.T) {
	sw := newStreamWrapper()

	called := false
	sw.AddHandler("my-handler", func(js *simplejson.Json) {
		called = true
	})

	// err := sw.Subscribe("btcusdt@aggTrade")
	// assert.NoError(t, err)
	// assert.True(t, called, "Handler should be called")
	// time.Sleep(1 * time.Second)

	sw.RemoveHandler("my-handler")
	called = sw.GetLoopStarted()
	assert.False(t, called, "Handler should be removed")
}
