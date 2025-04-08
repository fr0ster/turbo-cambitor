package streamer

import (
	"sync"
	"time"

	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/sirupsen/logrus"
)

func (stream *StreamWrapper) SetSymbol(symbol string) *StreamWrapper {
	stream.symbol = symbol
	return stream
}

func (stream *StreamWrapper) SetTimeOut(timeout time.Duration) *StreamWrapper {
	stream.timeOut = timeout
	stream.low_stream.SetTimeOut(timeout)
	return stream
}

func New(
	host web_socket.WsHost,
	path web_socket.WsPath,
	scheme web_socket.WsScheme,
	silent bool,
	timeOut ...time.Duration) *StreamWrapper {
	timeOutVar := 10 * time.Second
	if len(timeOut) > 0 {
		timeOutVar = timeOut[0]
	}
	stream, err := web_socket.New(host, path, scheme, web_socket.TextMessage, silent)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	}
	return &StreamWrapper{
		wsHost:      host,
		wsPath:      path,
		wsScheme:    scheme,
		symbol:      "",
		messageType: web_socket.TextMessage,
		silent:      silent,
		low_stream:  stream,
		timeOut:     timeOutVar,
		mu:          sync.Mutex{},
	}
}

func (sw *StreamWrapper) SetMaxReconnectAttempts(n int) *StreamWrapper {
	sw.maxReconnectAttempts = n
	return sw
}

func (ws *StreamWrapper) EnableAutoReconnect(interval ...time.Duration) *StreamWrapper {
	ws.autoReconnect = true
	ws.stopOnce = sync.Once{} // ⚠️ обнуляємо, щоби `Disable...` можна було викликати знов

	if len(interval) > 0 {
		ws.reconnectInterval = interval[0]
	}

	ws.reconnectStopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-ws.reconnectStopChan:
				return
			default:
				logrus.Warn("🔁 Attempting reconnect...")
				if err := ws.Reconnect(); err != nil {
					logrus.Warnf("Reconnect failed: %v", err)
					time.Sleep(ws.reconnectInterval)
					continue
				}
				logrus.Info("✅ Reconnected successfully")
			}
			time.Sleep(ws.reconnectInterval)
		}
	}()

	return ws
}

func (ws *StreamWrapper) DisableAutoReconnect() {
	ws.autoReconnect = false
	ws.stopOnce.Do(func() {
		if ws.reconnectStopChan != nil {
			close(ws.reconnectStopChan)
		}
	})
}

func (sw *StreamWrapper) GetConnection() *web_socket.WebSocketWrapper {
	return sw.low_stream
}

func (sw *StreamWrapper) IsLoopStarted() bool {
	return sw.low_stream != nil && sw.low_stream.GetLoopStarted()
}

func (sw *StreamWrapper) SetPongHandler(handler ...func(appData string) error) *StreamWrapper {
	sw.low_stream.SetPongHandler(handler...)
	return sw
}

func (sw *StreamWrapper) SetPingHandler(handler ...func(appData string) error) *StreamWrapper {
	sw.low_stream.SetPingHandler(handler...)
	return sw
}
