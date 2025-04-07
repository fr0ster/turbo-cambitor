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

func (sw *StreamWrapper) EnableAutoReconnect(interval ...time.Duration) *StreamWrapper {
	sw.autoReconnect = true
	if interval == nil {
		interval = append(interval, sw.reconnectInterval)
	}
	sw.reconnectInterval = interval[0]
	sw.reconnectStopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-sw.reconnectStopChan:
				return
			default:
				if sw.low_stream == nil || !sw.low_stream.GetLoopStarted() {
					logrus.Warn("🔁 Attempting reconnect...")
					if err := sw.Reconnect(); err != nil {
						logrus.Warnf("Reconnect failed: %v", err)
						time.Sleep(sw.reconnectInterval)
						continue
					}
					logrus.Info("✅ Reconnected successfully")
				}
				time.Sleep(sw.reconnectInterval)
			}
		}
	}()
	return sw
}

func (sw *StreamWrapper) DisableAutoReconnect() *StreamWrapper {
	if sw.reconnectStopChan != nil {
		close(sw.reconnectStopChan)
	}
	sw.autoReconnect = false
	return sw
}
