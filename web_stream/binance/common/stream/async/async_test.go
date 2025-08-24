// async_test.go: тести для async wrapper
//
// Перевіряє роботу health monitor, readiness, паузи, пінги, reconnect.
package async_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"

	common "github.com/fr0ster/turbo-cambitor/common"
	common_web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	async "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream/async"
	mock_server "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream/mock_server"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/stretchr/testify/assert"
)

var mockPort int

func TestMain(m *testing.M) {
	flag.Parse()
	mockPort = mock_server.StartReusableMockServer(0)
	os.Exit(m.Run())
}

func newStreamWrapper() *async.StreamWrapper {
	scheme := "ws"
	host := fmt.Sprintf("localhost:%d", mockPort)
	endpoint := "/ws"
	factory := func() (web_socket.WebSocketCommonInterface, error) {
		url := fmt.Sprintf("%s://%s%s", scheme, host, endpoint)
		return web_socket.NewWebSocketWrapper(websocket.DefaultDialer, url)
	}

	sw := async.NewStreamWrapper(factory, common.WsScheme(scheme), common.WsHost(host), common.WsEndpoint(endpoint))

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

// In async mode, operations fail fast when not connected/healthy.
func TestAsync_FailFast_WhenNotReady(t *testing.T) {
	scheme := "ws"
	host := fmt.Sprintf("localhost:%d", mockPort)
	endpoint := "ws"
	b := common_web_stream.NewStreamBuilder(common.WsScheme(scheme), common.WsHost(host), common.WsEndpoint(endpoint))
	s := b.StreamAsync()
	sw, err := s.Connect()
	assert.NoError(t, err)

	// Trigger abrupt close without restart to force disconnected state
	rq := simplejson.New()
	rq.Set("method", "CLOSE_ABRUPT")
	rq.Set("id", "async-failfast")
	_, _ = sw.Call(rq)

	// Immediately try to call again; expect an error due to not ready
	rq2 := simplejson.New()
	rq2.Set("method", "LIST_SUBSCRIPTIONS")
	rq2.Set("id", "async-after-close")
	resp2, err2 := sw.Call(rq2)
	assert.Error(t, err2)
	assert.Nil(t, resp2)
}

// Async mode with auto-reconnect should recover after server restart.
func TestAsync_Close_WithAutoReconnect(t *testing.T) {
	scheme := "ws"
	host := fmt.Sprintf("localhost:%d", mockPort)
	endpoint := "ws"
	b := common_web_stream.NewStreamBuilder(common.WsScheme(scheme), common.WsHost(host), common.WsEndpoint(endpoint))
	s := b.StreamAsync()
	sw, err := s.Connect()
	assert.NoError(t, err)

	// Виставляємо таймаути як у робочому тесті
	sw.SetReadTimeout(3 * time.Second)
	sw.SetWriteTimeout(3 * time.Second)
	sw.SetMaxReconnectAttempts(30).SetReconnectInterval(200 * time.Millisecond).EnableAutoReconnect()

	// 1. Відправляємо команду graceful close з рестартом
	rq := simplejson.New()
	rq.Set("method", "CLOSE_GRACEFUL")
	rq.Set("id", "async-auto")
	rq.SetPath([]string{"params", "restart"}, 500)
	_, _ = sw.Call(rq)

	// 2. Одразу після падіння — має бути помилка (fail-fast)
	rqFail := simplejson.New()
	rqFail.Set("method", "LIST_SUBSCRIPTIONS")
	rqFail.Set("id", "async-fail-after-close")
	respFail, errFail := sw.Call(rqFail)
	assert.Error(t, errFail)
	assert.Nil(t, respFail)

	// Дочекайся відновлення коннекту (через сервісний канал/монітор)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	assert.NoError(t, sw.WaitConnected(ctx), "should auto-recover after CLOSE_GRACEFUL with restart")

	// 3. Повторно викликаєм після відновлення коннекту
	rq2 := simplejson.New()
	rq2.Set("method", "LIST_SUBSCRIPTIONS")
	rq2.Set("id", "async-post-auto")
	resp2, err2 := sw.Call(rq2)
	assert.NoError(t, err2)
	assert.NotNil(t, resp2)
}

var closeMutex = &sync.Mutex{}

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
