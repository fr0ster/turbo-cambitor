package common_web_api_test

import (
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	common "github.com/fr0ster/turbo-cambitor/common"
	common_web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/common/web_api"
	web_socket "github.com/fr0ster/turbo-restler/web_socket"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// Note: runEchoServer is defined in web_api_wrapper_test.go within the same package.

func TestWebApiWrapper_WithFactory_Override(t *testing.T) {
	const endpointPath = "/ws"

	server := runEchoServer(t, endpointPath)
	defer server.Close()

	host := strings.TrimPrefix(server.URL, "http://")

	var called atomic.Int32
	factory := func() (web_socket.WebSocketCommonInterface, error) {
		called.Add(1)
		url := "ws://" + host + endpointPath
		return web_socket.NewWebSocketWrapper(websocket.DefaultDialer, url)
	}

	// Intentionally pass bogus host to ensure default factory would fail if used
	wa := common_web_api.New(
		common.WsHost("invalid.local:0"),
		common.WsEndpoint("/broken"),
		common.WsScheme("ws"),
		nil,
		common_web_api.WithFactory(factory),
	)

	js := simplejson.New()
	js.Set("method", "echo")
	js.Set("params", map[string]any{"message": "hello"})

	resp, err := wa.SetTimeOut(1500 * time.Millisecond).Call(js)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.GreaterOrEqual(t, int(called.Load()), 1, "custom factory must be used")

	sent, _ := js.Encode()
	received, _ := resp.Encode()
	assert.JSONEq(t, string(sent), string(received))
}

func TestWebApiWrapper_WithWebSocketConfig_ConstructURL(t *testing.T) {
	const endpointPath = "/ws"

	server := runEchoServer(t, endpointPath)
	defer server.Close()

	host := strings.TrimPrefix(server.URL, "http://")

	wa := common_web_api.New(
		common.WsHost(host),
		common.WsEndpoint(endpointPath),
		common.WsScheme("ws"),
		nil,
		common_web_api.WithWebSocketConfig(web_socket.WebSocketConfig{
			// URL left empty on purpose to test auto-construction from wa fields
			BufferSize:   16,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		}),
	)

	js := simplejson.New()
	js.Set("method", "echo")
	js.Set("params", map[string]any{"message": "world"})

	resp, err := wa.SetTimeOut(2 * time.Second).Call(js)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestWebApiWrapper_WithWebSocketConfig_ExplicitURL_Wins(t *testing.T) {
	const endpointPath = "/ws"

	server := runEchoServer(t, endpointPath)
	defer server.Close()

	host := strings.TrimPrefix(server.URL, "http://")
	url := "ws://" + host + endpointPath

	wa := common_web_api.New(
		// Bogus fields – should be ignored because explicit URL is provided in config
		common.WsHost("invalid.local:1234"),
		common.WsEndpoint("/nope"),
		common.WsScheme("wss"),
		nil,
		common_web_api.WithWebSocketConfig(web_socket.WebSocketConfig{
			URL:          url,
			Dialer:       websocket.DefaultDialer,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		}),
	)

	js := simplejson.New()
	js.Set("method", "echo")
	js.Set("params", map[string]any{"message": "explicit"})

	resp, err := wa.SetTimeOut(2 * time.Second).Call(js)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
