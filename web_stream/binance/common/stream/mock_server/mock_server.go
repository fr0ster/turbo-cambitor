package mock_server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	mockServerMu sync.Mutex
	currentSrv   *http.Server
)

// StartReusableMockServer запускає WebSocket mock сервер на заданому порту
func StartReusableMockServer(port int) {
	mockServerMu.Lock()
	defer mockServerMu.Unlock()

	addr := fmt.Sprintf(":%d", port)

	// Якщо є попередній сервер — зупиняємо
	if currentSrv != nil {
		logrus.Warnf("🛑 Stopping existing server on port %d", port)
		_ = currentSrv.Close()
		time.Sleep(100 * time.Millisecond)
	}

	upgrader := websocket.Upgrader{}
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.Errorf("Upgrade error: %v", err)
			return
		}
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logrus.Warnf("Read error: %v", err)
				break
			}

			req, err := simplejson.NewJson(msg)
			if err != nil {
				logrus.Warn("Invalid JSON received")
				continue
			}

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
				logrus.Info("🤐 Simulating timeout: no response sent")
				time.Sleep(2 * time.Second)
				continue
			case "CLOSE_GRACEFUL":
				logrus.Info("🔌 Graceful close requested by client")
				time.Sleep(100 * time.Millisecond)
				conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
				return
			case "CLOSE_ABRUPT":
				logrus.Info("💥 Abrupt close requested by client")
				return
			case "ERROR":
				params := req.Get("params").MustArray()
				errorText := "generic error"
				if len(params) > 0 {
					if s, ok := params[0].(string); ok && s != "" {
						errorText = s
					}
				}
				resp.Set("error", errorText)
				resp.Set("id", id) // 👈 обов’язково

				b, _ := resp.Encode()
				conn.WriteMessage(websocket.TextMessage, b)

				// ❗️ НЕ закриваємо з'єднання
				continue
			default:
				resp.Set("error", "unknown method")
			}

			b, _ := resp.Encode()
			conn.WriteMessage(websocket.TextMessage, b)
		}
	})

	currentSrv = &http.Server{Addr: addr, Handler: mux}

	go func() {
		logrus.Infof("🚀 Starting mock server on %s", addr)
		err := currentSrv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server error: %v", err)
		}
	}()
}
