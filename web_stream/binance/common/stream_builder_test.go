package common_web_stream_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	common "github.com/fr0ster/turbo-cambitor/common"
	common_web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // дозволяє будь-яке джерело
}

func mockWSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
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
	// Запускаємо мок-сервер
	server := httptest.NewServer(http.HandlerFunc(mockWSHandler))
	defer server.Close()

	// Отримаємо хост без "http://"
	hostOnly := strings.TrimPrefix(server.URL, "http://") // localhost:12345
	symbol := "btcusdt"

	// Конструктор StreamBuilder (тут wsEndpoint буде "")
	builder := common_web_stream.NewStreamBuilder(
		common.WsScheme("ws"),
		common.WsHost(hostOnly+"/ws"),
		common.WsEndpoint(""),
		symbol,
	)

	stream := builder.AggTrades()

	err := stream.Stream().Connect()
	if err != nil {
		t.Fatalf("failed to connect to mock websocket: %v", err)
	}
}
