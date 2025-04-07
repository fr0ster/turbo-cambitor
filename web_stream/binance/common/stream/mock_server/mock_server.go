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

	if currentSrv != nil {
		logrus.Warnf("🖑 Stopping existing server on port %d", port)
		_ = currentSrv.Close()
		time.Sleep(200 * time.Millisecond)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocketConnection(w, r, port)
	})

	currentSrv = &http.Server{Addr: addr, Handler: mux}

	go func() {
		logrus.Infof("🚀 Starting mock server on %s", addr)
		err := currentSrv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Server error: %v", err)
		}
	}()
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request, port int) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()

	stopChan := make(chan struct{})
	go startPinger(conn, stopChan)

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
		restartAfter := req.Get("params").Get("restart").MustInt()

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
			handleClose(conn, stopChan, restartAfter, port, "graceful shutdown", websocket.CloseNormalClosure)
			return
		case "CLOSE_ABRUPT":
			handleClose(conn, stopChan, restartAfter, port, "abrupt shutdown", websocket.CloseAbnormalClosure)
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
			b, _ := resp.Encode()
			mockServerMu.Lock()
			conn.WriteMessage(websocket.TextMessage, b)
			mockServerMu.Unlock()
			continue
		default:
			resp.Set("error", "unknown method")
		}

		b, _ := resp.Encode()
		mockServerMu.Lock()
		conn.WriteMessage(websocket.TextMessage, b)
		mockServerMu.Unlock()
	}

	close(stopChan)
}

func handleClose(conn *websocket.Conn, stopChan chan struct{}, restartAfter, port int, reason string, code int) {
	logrus.Infof("Closing connection: %s", reason)

	notify := simplejson.New()
	notify.Set("type", "closing")
	notify.Set("reason", reason)
	notify.Set("time", time.Now().Format(time.RFC3339))
	b, _ := notify.Encode()
	mockServerMu.Lock()
	_ = conn.WriteMessage(websocket.TextMessage, b)
	mockServerMu.Unlock()

	time.Sleep(100 * time.Millisecond)
	mockServerMu.Lock()
	_ = conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(code, reason))
	mockServerMu.Unlock()

	close(stopChan)

	if restartAfter > 0 {
		go func() {
			time.Sleep(time.Duration(restartAfter) * time.Millisecond)

			logrus.Infof("🧘 Waiting before restarting server on :%d", port)

			mockServerMu.Lock()

			logrus.Warnf("🛑 Closing server on port %d", port)
			if currentSrv != nil {
				_ = currentSrv.Close()
				currentSrv = nil
			}
			mockServerMu.Unlock()

			time.Sleep(200 * time.Millisecond) // дати порту TIME_WAIT

			StartReusableMockServer(port)
		}()
	}
}

func startPinger(conn *websocket.Conn, stopChan chan struct{}) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			msg := simplejson.New()
			msg.Set("type", "ping")
			msg.Set("time", time.Now().Format(time.RFC3339))
			b, _ := msg.Encode()
			mockServerMu.Lock()
			conn.WriteMessage(websocket.TextMessage, b)
			mockServerMu.Unlock()
		case <-stopChan:
			return
		}
	}
}
