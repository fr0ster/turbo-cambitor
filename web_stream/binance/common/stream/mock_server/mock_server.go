package mock_server

import (
	"fmt"
	"net"
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
	currentPort  int
)

// StartReusableMockServer запускає WebSocket mock сервер на заданому порту.
// Якщо port == 0, буде обрано вільний еферемний порт. Повертає фактичний порт.
func StartReusableMockServer(port int) int {
	mockServerMu.Lock()
	defer mockServerMu.Unlock()

	if currentSrv != nil {
		if currentPort > 0 {
			logrus.Infof("☑️ Mock server on port %d is already running", currentPort)
			return currentPort
		}
	}

	addr := fmt.Sprintf(":%d", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocketConnection(w, r, currentPort)
	})

	// Створюємо listener, щоб отримати фактичний порт (коли port == 0)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Errorf("Failed to start mock server listener on %s: %v", addr, err)
		return 0
	}
	actualPort := ln.Addr().(*net.TCPAddr).Port
	currentPort = actualPort
	currentSrv = &http.Server{Handler: mux}

	go func(p int) {
		logrus.Infof("🚀 Starting mock server on :%d", p)
		err := currentSrv.Serve(ln)
		if err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Server error: %v", err)
		}
		mockServerMu.Lock()
		currentSrv = nil
		currentPort = 0
		mockServerMu.Unlock()
	}(actualPort)

	return actualPort
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request, port int) {
	upgrader := websocket.Upgrader{
		// Allow all origins for tests to avoid cross-origin issues
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()

	stopChan := make(chan struct{})
	pingPongStopChan := make(chan struct{})
	pingPongActive := false

	// ⬇️ Підписки клієнта
	subscriptions := make(map[string]chan struct{})

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
			streamsRaw := req.Get("params").MustArray()
			var subscribed []string
			var already []string

			for _, s := range streamsRaw {
				stream, ok := s.(string)
				if !ok || stream == "" {
					continue
				}
				if _, exists := subscriptions[stream]; exists {
					already = append(already, stream)
					continue
				}

				stop := make(chan struct{})
				subscriptions[stream] = stop

				go func(c *websocket.Conn, stream string, stop chan struct{}) {
					ticker := time.NewTicker(500 * time.Millisecond)
					defer ticker.Stop()

					for {
						select {
						case <-stop:
							logrus.Infof("🛑 Stream %s stopped", stream)
							return
						case t := <-ticker.C:
							msg := simplejson.New()
							msg.Set("type", "mock_data")
							msg.Set("stream", stream)
							msg.Set("time", t.Format(time.RFC3339))
							msg.Set("value", fmt.Sprintf("%.2f", 30000+1000*float64(time.Now().UnixNano()%1000)/1000.0))

							b, _ := msg.Encode()
							mockServerMu.Lock()
							err := c.WriteMessage(websocket.TextMessage, b)
							mockServerMu.Unlock()
							if err != nil {
								logrus.Warnf("❌ Failed to send mock data for %s: %v", stream, err)
								return
							}
						}
					}
				}(conn, stream, stop)

				subscribed = append(subscribed, stream)
			}

			if len(subscribed) > 0 {
				resp.Set("result", subscribed)
			} else {
				resp.Set("error", fmt.Sprintf("Already subscribed to: %v", already))
			}

		case "UNSUBSCRIBE":
			streamsRaw := req.Get("params").MustArray()
			var unsubscribed []string
			var notFound []string

			for _, s := range streamsRaw {
				stream, ok := s.(string)
				if !ok || stream == "" {
					continue
				}
				if stop, exists := subscriptions[stream]; exists {
					close(stop)
					delete(subscriptions, stream)
					unsubscribed = append(unsubscribed, stream)
				} else {
					notFound = append(notFound, stream)
				}
			}

			if len(unsubscribed) > 0 {
				resp.Set("result", unsubscribed)
			} else {
				resp.Set("error", fmt.Sprintf("Not subscribed to: %v", notFound))
			}

		case "LIST_SUBSCRIPTIONS":
			activeStreams := make([]string, 0, len(subscriptions))
			for k := range subscriptions {
				activeStreams = append(activeStreams, k)
			}
			resp.Set("result", activeStreams)

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

		case "PONG_CONTROL":
			timeout := req.Get("params").Get("timeout").MustInt(1000)
			if !pingPongActive {
				pingPongActive = true
				go startPingController(conn, timeout, pingPongStopChan)

				msg := simplejson.New()
				msg.Set("type", "ping")
				msg.Set("time", time.Now().Format(time.RFC3339))
				b, _ := msg.Encode()
				conn.WriteMessage(websocket.TextMessage, b)

				resp.Set("result", fmt.Sprintf("pong control started with timeout %d ms", timeout))
			} else {
				resp.Set("result", "pong control already active")
			}

		default:
			resp.Set("error", "unknown method")
		}

		b, _ := resp.Encode()
		mockServerMu.Lock()
		conn.WriteMessage(websocket.TextMessage, b)
		mockServerMu.Unlock()
	}

	// 🧹 Зупиняємо всі активні підписки
	for _, stop := range subscriptions {
		close(stop)
	}
	close(stopChan)
	close(pingPongStopChan)
}

func startPingController(conn *websocket.Conn, timeoutMs int, stopChan chan struct{}) {
	var lastPong time.Time
	var lastPongMu sync.Mutex

	timeout := time.Duration(timeoutMs) * time.Millisecond
	ticker := time.NewTicker(timeout / 2)
	defer ticker.Stop()

	conn.SetPongHandler(func(appData string) error {
		logrus.Infof("✅ Pong received: %s", appData)
		lastPongMu.Lock()
		lastPong = time.Now()
		lastPongMu.Unlock()
		return nil
	})

	lastPong = time.Now()
	logrus.Infof("🏓 Native ping-pong control started with timeout %d ms", timeoutMs)

	for {
		select {
		case <-ticker.C:
			// надсилаємо ping
			mockServerMu.Lock()
			err := conn.WriteControl(websocket.PingMessage, []byte(fmt.Sprintf("ping-%d", time.Now().UnixNano())), time.Now().Add(time.Second))
			mockServerMu.Unlock()

			if err != nil {
				logrus.Warnf("❌ Failed to send ping: %v", err)
				return
			}

			// надсилаємо службовий ping-пакет (JSON)
			msg := simplejson.New()
			msg.Set("type", "ping")
			msg.Set("time", time.Now().Format(time.RFC3339))
			b, _ := msg.Encode()

			mockServerMu.Lock()
			_ = conn.WriteMessage(websocket.TextMessage, b)
			mockServerMu.Unlock()

			// перевіряємо таймаут
			lastPongMu.Lock()
			since := time.Since(lastPong)
			lastPongMu.Unlock()

			if since > timeout {
				logrus.Warn("⏱ Pong timeout reached, closing connection")
				notify := simplejson.New()
				notify.Set("type", "closing")
				notify.Set("reason", "pong timeout")
				notify.Set("time", time.Now().Format(time.RFC3339))
				nb, _ := notify.Encode()

				mockServerMu.Lock()
				_ = conn.WriteMessage(websocket.TextMessage, nb)
				_ = conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(1008, "Pong timeout"),
					time.Now().Add(time.Second))
				mockServerMu.Unlock()
				return
			}
		case <-stopChan:
			logrus.Infof("🛑 Ping controller stopped")
			return
		}
	}
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
	_ = conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(code, reason))
	mockServerMu.Unlock()

	close(stopChan)

	if restartAfter > 0 {
		go func() {
			time.Sleep(time.Duration(restartAfter) * time.Millisecond)

			logrus.Infof("🧘 Waiting before restarting server on :%d", port)

			mockServerMu.Lock()
			if currentSrv != nil {
				logrus.Warnf("🛑 Closing server on port %d", port)
				_ = currentSrv.Close()
				currentSrv = nil
			}
			mockServerMu.Unlock()

			time.Sleep(200 * time.Millisecond)
			p := port
			if p == 0 {
				p = currentPort
			}
			StartReusableMockServer(p)
		}()
	}
}
