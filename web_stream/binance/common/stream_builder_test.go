package common_web_stream_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"sync/atomic"

	common "github.com/fr0ster/turbo-cambitor/common"
	common_web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	web_socket "github.com/fr0ster/turbo-restler/web_socket"
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

	_, err := stream.Stream().Connect()
	if err != nil {
		t.Fatalf("failed to connect to mock websocket: %v", err)
	}
}

func TestStreamBuilder_WithDialer(t *testing.T) {
	// Mock WS server
	server := httptest.NewServer(http.HandlerFunc(mockWSHandler))
	defer server.Close()

	hostOnly := strings.TrimPrefix(server.URL, "http://")
	symbol := "btcusdt"

	// Custom dialer that toggles a flag via NetDialContext
	var used atomic.Bool
	d := *websocket.DefaultDialer
	d.NetDialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		used.Store(true)
		var nd net.Dialer
		return nd.DialContext(ctx, network, addr)
	}

	builder := common_web_stream.NewStreamBuilder(
		common.WsScheme("ws"),
		common.WsHost(hostOnly+"/ws"),
		common.WsEndpoint(""),
		symbol,
	).WithDialer(&d)

	// Build and connect
	_, err := builder.AggTrades().Stream().Connect()
	if err != nil {
		t.Fatalf("failed to connect with custom dialer: %v", err)
	}
	if !used.Load() {
		t.Fatalf("expected custom dialer to be used, but it wasn't")
	}
}

func TestStreamBuilder_WithFactory(t *testing.T) {
	// Mock WS server
	server := httptest.NewServer(http.HandlerFunc(mockWSHandler))
	defer server.Close()

	hostOnly := strings.TrimPrefix(server.URL, "http://")
	symbol := "btcusdt"

	var factoryUsed atomic.Bool
	factory := func(url string) (web_socket.WebSocketCommonInterface, error) {
		factoryUsed.Store(true)
		return web_socket.NewWebSocketWrapper(websocket.DefaultDialer, url)
	}

	builder := common_web_stream.NewStreamBuilder(
		common.WsScheme("ws"),
		common.WsHost(hostOnly+"/ws"),
		common.WsEndpoint(""),
		symbol,
	).WithFactory(factory)

	// Ensure that if both factory and dialer are set, factory wins
	d := *websocket.DefaultDialer
	d.HandshakeTimeout = 2 * time.Second
	builder = builder.WithDialer(&d)

	_, err := builder.AggTrades().Stream().Connect()
	if err != nil {
		t.Fatalf("failed to connect with custom factory: %v", err)
	}
	if !factoryUsed.Load() {
		t.Fatalf("expected custom factory to be used, but it wasn't")
	}
}
