// sync_test.go: тести для sync wrapper
//
// Перевіряє роботу health monitor, readiness, паузи, пінги, reconnect.
package sync_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"

	common "github.com/fr0ster/turbo-cambitor/common"
	common_web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	mock_server "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream/mock_server"
	"github.com/stretchr/testify/assert"
)

var mockPort int

func TestMain(m *testing.M) {
	flag.Parse()
	mockPort = mock_server.StartReusableMockServer(0)
	os.Exit(m.Run())
}

// Graceful close with server restart, then manual reconnect
func TestSync_GracefulClose_WithRestart(t *testing.T) {
	t.Parallel()
	scheme := "ws"
	host := fmt.Sprintf("localhost:%d", mockPort)
	endpoint := "ws"
	b := common_web_stream.NewStreamBuilder(common.WsScheme(scheme), common.WsHost(host), common.WsEndpoint(endpoint))
	s := b.StreamSync()
	sw, err := s.Connect()
	assert.NoError(t, err)

	sw.SetReadTimeout(3 * time.Second)
	sw.SetWriteTimeout(3 * time.Second)

	rq := simplejson.New()
	rq.Set("method", "CLOSE_GRACEFUL")
	rq.Set("id", "sync-graceful")
	rq.SetPath([]string{"params", "restart"}, 400)
	resp, callErr := sw.Call(rq)
	assert.Error(t, callErr)
	assert.Nil(t, resp)

	err = sw.Reconnect(10, 200*time.Millisecond)
	assert.NoError(t, err)

	rq2 := simplejson.New()
	rq2.Set("method", "LIST_SUBSCRIPTIONS")
	rq2.Set("id", "sync-post")
	resp2, err2 := sw.Call(rq2)
	assert.NoError(t, err2)
	assert.NotNil(t, resp2)
}

// Abrupt close with server restart, then manual reconnect
func TestSync_AbruptClose_WithRestart(t *testing.T) {
	t.Parallel()
	scheme := "ws"
	host := fmt.Sprintf("localhost:%d", mockPort)
	endpoint := "ws"
	b := common_web_stream.NewStreamBuilder(common.WsScheme(scheme), common.WsHost(host), common.WsEndpoint(endpoint))
	s := b.StreamSync()
	sw, err := s.Connect()
	assert.NoError(t, err)

	sw.SetReadTimeout(3 * time.Second)
	sw.SetWriteTimeout(3 * time.Second)

	rq := simplejson.New()
	rq.Set("method", "CLOSE_ABRUPT")
	rq.Set("id", "sync-abrupt")
	rq.SetPath([]string{"params", "restart"}, 400)
	resp, callErr := sw.Call(rq)
	assert.Error(t, callErr)
	assert.Nil(t, resp)

	err = sw.Reconnect(10, 200*time.Millisecond)
	assert.NoError(t, err)

	rq2 := simplejson.New()
	rq2.Set("method", "LIST_SUBSCRIPTIONS")
	rq2.Set("id", "sync-post")
	resp2, err2 := sw.Call(rq2)
	assert.NoError(t, err2)
	assert.NotNil(t, resp2)
}

// Graceful close + auto-reconnect enabled
func TestSync_GracefulClose_WithAutoReconnect(t *testing.T) {
	t.Parallel()
	scheme := "ws"
	host := fmt.Sprintf("localhost:%d", mockPort)
	endpoint := "ws"
	b := common_web_stream.NewStreamBuilder(common.WsScheme(scheme), common.WsHost(host), common.WsEndpoint(endpoint))
	s := b.StreamSync()
	sw, err := s.Connect()
	assert.NoError(t, err)

	sw.SetReadTimeout(3 * time.Second)
	sw.SetWriteTimeout(3 * time.Second)

	sw.SetMaxReconnectAttempts(30).SetReconnectInterval(200 * time.Millisecond).EnableAutoReconnect()

	rq := simplejson.New()
	rq.Set("method", "CLOSE_GRACEFUL")
	rq.Set("id", "sync-graceful-auto")
	rq.SetPath([]string{"params", "restart"}, 500)
	_, _ = sw.Call(rq)

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	assert.NoError(t, sw.WaitConnected(ctx))

	rq2 := simplejson.New()
	rq2.Set("method", "LIST_SUBSCRIPTIONS")
	rq2.Set("id", "sync-post-auto")
	resp2, err2 := sw.Call(rq2)
	assert.NoError(t, err2)
	assert.NotNil(t, resp2)
}

// Abrupt close + auto-reconnect enabled
func TestSync_AbruptClose_WithAutoReconnect(t *testing.T) {
	t.Parallel()
	scheme := "ws"
	host := fmt.Sprintf("localhost:%d", mockPort)
	endpoint := "ws"
	b := common_web_stream.NewStreamBuilder(common.WsScheme(scheme), common.WsHost(host), common.WsEndpoint(endpoint))
	s := b.StreamSync()
	sw, err := s.Connect()
	assert.NoError(t, err)

	sw.SetReadTimeout(3 * time.Second)
	sw.SetWriteTimeout(3 * time.Second)

	sw.SetMaxReconnectAttempts(10).SetReconnectInterval(200 * time.Millisecond).EnableAutoReconnect()

	rq := simplejson.New()
	rq.Set("method", "CLOSE_ABRUPT")
	rq.Set("id", "sync-abrupt-auto")
	rq.SetPath([]string{"params", "restart"}, 500)
	_, _ = sw.Call(rq)

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	assert.NoError(t, sw.WaitConnected(ctx))

	rq2 := simplejson.New()
	rq2.Set("method", "LIST_SUBSCRIPTIONS")
	rq2.Set("id", "sync-post-auto")
	resp2, err2 := sw.Call(rq2)
	assert.NoError(t, err2)
	assert.NotNil(t, resp2)
}
