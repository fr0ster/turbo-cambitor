package common_web_stream_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fr0ster/turbo-cambitor/web_stream/binance/common/common_web_stream"
	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"
	"github.com/gorilla/websocket"
)

// ...

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // ⚠️ дозволяє все
}

func mockWSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// не можна використовувати t.Logf тут, бо t не передається — просто лог
		println("Upgrade error:", err.Error())
		return
	}
	defer conn.Close()

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if err := conn.WriteMessage(mt, msg); err != nil {
			break
		}
	}
}

func TestStreamBuilder_AggTrades(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockWSHandler))
	defer server.Close()

	// Отримаємо хост без http://
	wsHost := strings.TrimPrefix(server.URL, "http://")

	builder := common_web_stream.NewStreamBuilder("ws", stream.WsHost(wsHost), "btcusdt")
	stream := builder.AggTrades()

	err := stream.Connect()
	if err != nil {
		t.Fatalf("failed to connect to mock websocket: %v", err)
	}
}
