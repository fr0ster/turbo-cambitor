package common_web_api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	common "github.com/fr0ster/turbo-cambitor/common"
	common_web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/common/web_api"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func runEchoServer(t *testing.T, expectedPath string) *httptest.Server {
	upgrader := websocket.Upgrader{}
	mux := http.NewServeMux()

	mux.HandleFunc(expectedPath, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		assert.NoError(t, err)

		go func(c *websocket.Conn) {
			defer c.Close()

			for {
				_ = c.SetReadDeadline(time.Now().Add(5 * time.Second))
				msgType, msg, err := c.ReadMessage()
				if err != nil {
					t.Logf("read error: %v", err)
					return
				}
				t.Logf("Echo server received: %s", string(msg))
				_ = c.WriteMessage(msgType, msg)
			}
		}(conn)
	})

	return httptest.NewServer(mux)
}

func TestWebApiWrapper_Echo_ReadWrite(t *testing.T) {
	const endpointPath = "/ws"

	server := runEchoServer(t, endpointPath)
	defer server.Close()

	host := strings.TrimPrefix(server.URL, "http://")

	wa := common_web_api.New(
		common.WsHost(host),
		common.WsEndpoint(endpointPath),
		common.WsScheme("ws"),
		nil, // No signature needed for echo test
	)

	js := simplejson.New()
	js.Set("method", "echo")
	js.Set("params", map[string]interface{}{
		"message": "hello echo",
	})

	// ⏩ Надсилаємо запит та читаємо відповідь
	resp, err := wa.Call(js)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// 🔍 Перевіряємо, що відповідь збігається з тим, що відправили
	sent, _ := js.Encode()
	received, _ := resp.Encode()

	assert.JSONEq(t, string(sent), string(received), "Echo response should match the sent message")
}
