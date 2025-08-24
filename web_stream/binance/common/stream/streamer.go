// streamer.go: базовий StreamWrapper
//
// Відповідає за низькорівневий менеджмент WebSocket-з'єднання, health socket, таймаути, auto-reconnect.
// Ловить помилки коннекта (1000/1006), може ресетити з'єднання або повертати помилку консьюмеру.
// Не має власного health monitor, не керує режимом роботи (sync/async) — це роблять wrappers.
// Має метод HealthPing(timeout) для ручного пінгу сервера.
// ARCHITECTURE OVERVIEW
//
// StreamWrapper (base):
//   - Відповідає за низькорівневий менеджмент WebSocket-з'єднання.
//   - Ловить помилки коннекта (1000/1006), може або ресетити з'єднання, або повертати помилку консьюмеру.
//   - Не має власного health monitor, не керує режимом роботи (sync/async).
//
// Sync/Async wrappers:
//   - НЕ ловлять помилки самостійно, не роблять reconnect.
//   - Керуються зовнішнім health monitor-ом, який пінгує сервер і керує станом з'єднання.
//   - Всі readiness/health перевірки виконуються wrapper-ом, а не базовим StreamWrapper.
//
// Health monitor:
//   - Окремий моніторинговий коннектор, який періодично пінгує сервер через dedicated health socket.
//   - Якщо unhealthy — закриває основне з'єднання, якщо healthy — відкриває/відновлює.
//   - Керує paused-станом для async/sync wrappers.
//
// Пінг-позапланово:
//   - StreamWrapper має метод HealthPing(timeout), який дозволяє вручну пінгувати сервер поза графіком монітору.
//   - Це потрібно для ручних readiness/health перевірок у тестах чи при ініціалізації.
//
// Всі ці принципи обов'язково враховувати при розробці та тестуванні wrappers.
package streamer

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/google/uuid"

	"github.com/fr0ster/turbo-cambitor/common"
	"github.com/fr0ster/turbo-restler/web_socket"
)

type StreamWrapper struct {
	socket web_socket.WebSocketCommonInterface
	// Dedicated health-check socket to probe server state independently of business traffic
	healthSocket web_socket.WebSocketCommonInterface
	factory      func() (web_socket.WebSocketCommonInterface, error)
	wsScheme     common.WsScheme
	wsHost       common.WsHost
	wsEndpoint   common.WsEndpoint
	timeOut      time.Duration
	mu           sync.Mutex

	// Separate mutex for health socket operations to reduce contention with main socket
	healthMu sync.Mutex

	readTimeout  *time.Duration
	writeTimeout *time.Duration

	// Налаштування, які можуть бути задані до Connect()
	pendingLogger     func(message web_socket.LogRecord)
	pendingPingHandle func(string) error
	pendingPongHandle func(string) error

	autoReconnect        bool
	maxReconnectAttempts int
	reconnectInterval    time.Duration

	// Health monitoring config/state
	healthInterval time.Duration
	healthTimeout  time.Duration
	healthStopChan chan struct{}
	healthRunning  bool
	healthy        bool

	// Pause flag controlled by health monitor
	paused bool

	// Default sync timeout used by helpers/tests
	defaultSyncTimeout time.Duration

	handlers map[string]int
}

var (
	// ErrPaused is returned in async mode when the connection is temporarily paused by the health monitor.
	ErrPaused = errors.New("connection paused")
	// ErrNotConnected indicates that there is no active connection and no pause/recovery in progress.
	ErrNotConnected = errors.New("not connected")
)

// Compile-time assertion: base wrapper is the reference implementation of StreamInterface.
var _ CambitorInterface = (*StreamWrapper)(nil)

// NewStreamWrapper creates a new StreamWrapper with a socket factory
func NewStreamWrapper(
	factory func() (web_socket.WebSocketCommonInterface, error),
	wsScheme common.WsScheme,
	wsHost common.WsHost,
	wsEndpoint common.WsEndpoint,
	timeOut ...time.Duration) *StreamWrapper {
	if len(timeOut) == 0 {
		// Slightly longer default to accommodate public WS management calls (SUBSCRIBE/UNSUBSCRIBE/LIST)
		timeOut = append(timeOut, 3*time.Second)
	}
	return &StreamWrapper{
		socket:     nil, // defer dialing until Connect()
		factory:    factory,
		wsScheme:   wsScheme,
		wsHost:     wsHost,
		wsEndpoint: wsEndpoint,
		timeOut:    timeOut[0],
		handlers:   make(map[string]int),
		// Sync wait will use the same timeout as Call by default
		defaultSyncTimeout: timeOut[0],
	}
}

func (sw *StreamWrapper) Connect() (*StreamWrapper, error) {
	var sock web_socket.WebSocketCommonInterface

	// Initialize socket on first Connect (configure under lock, open outside)
	sw.mu.Lock()
	if sw.socket == nil {
		s, err := sw.factory()
		if err != nil {
			sw.mu.Unlock()
			return nil, fmt.Errorf("failed to create websocket client: %w", err)
		}
		// Apply deferred settings before opening
		if sw.pendingLogger != nil {
			s.SetMessageLogger(sw.pendingLogger)
		}
		if sw.readTimeout != nil {
			s.SetReadTimeout(*sw.readTimeout)
		}
		if sw.writeTimeout != nil {
			s.SetWriteTimeout(*sw.writeTimeout)
		}
		if sw.pendingPingHandle != nil {
			s.SetPingHandler(sw.pendingPingHandle)
		}
		if sw.pendingPongHandle != nil {
			s.SetPongHandler(sw.pendingPongHandle)
		}
		// Install state handlers (they may acquire sw.mu later during callbacks)
		s.SetConnectedHandler(func() {
			sw.mu.Lock()
			sw.paused = false
			sw.mu.Unlock()
		})
		s.SetDisconnectHandler(func() {
			sw.mu.Lock()
			if sw.autoReconnect {
				sw.paused = true
			}
			sw.mu.Unlock()
		})
		sw.socket = s
	}
	sock = sw.socket
	sw.mu.Unlock()

	// Open outside the lock to avoid deadlocks with handlers
	sock.Open()
	_ = sock.WaitStarted()
	return sw, nil
}

func (sw *StreamWrapper) Reconnect(maxAttempts int, delay time.Duration) error {
	for i := 0; i < maxAttempts; i++ {
		// Snapshot and close old socket outside of lock
		var old web_socket.WebSocketCommonInterface
		sw.mu.Lock()
		old = sw.socket
		sw.mu.Unlock()

		if old != nil {
			_ = old.Close()
			_ = old.WaitStopped()
		}

		// Create and configure a new socket
		socket, err := sw.factory()
		if err == nil {
			if sw.pendingLogger != nil {
				socket.SetMessageLogger(sw.pendingLogger)
			}
			if sw.readTimeout != nil {
				socket.SetReadTimeout(*sw.readTimeout)
			}
			if sw.writeTimeout != nil {
				socket.SetWriteTimeout(*sw.writeTimeout)
			}
			if sw.pendingPingHandle != nil {
				socket.SetPingHandler(sw.pendingPingHandle)
			}
			if sw.pendingPongHandle != nil {
				socket.SetPongHandler(sw.pendingPongHandle)
			}
			socket.SetConnectedHandler(func() {
				sw.mu.Lock()
				sw.paused = false
				sw.mu.Unlock()
			})
			socket.SetDisconnectHandler(func() {
				sw.mu.Lock()
				if sw.autoReconnect {
					sw.paused = true
				}
				sw.mu.Unlock()
			})

			// Open outside of sw.mu to avoid handler deadlocks
			socket.Open()
			_ = socket.WaitStarted()

			sw.mu.Lock()
			sw.socket = socket
			sw.mu.Unlock()
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("failed to reconnect after %d attempts", maxAttempts)
}

func (sw *StreamWrapper) Disconnect() {
	// Detach main socket under lock
	sw.mu.Lock()
	sock := sw.socket
	sw.socket = nil
	// Stop monitor while holding the lock (uses *_Locked)
	sw.stopHealthMonitorLocked()
	sw.mu.Unlock()

	// Close main socket outside lock to avoid deadlocks with callbacks
	if sock != nil {
		_ = sock.Close()
		_ = sock.WaitStopped()
	}

	// Close health socket under its own lock, no sw.mu held
	sw.healthMu.Lock()
	if sw.healthSocket != nil {
		_ = sw.healthSocket.Close()
		sw.healthSocket = nil
	}
	sw.healthMu.Unlock()
}

func (sw *StreamWrapper) Call(rq *simplejson.Json) (*simplejson.Json, error) {
	id := rq.Get("id").MustString()
	if rq.Get("id").MustString() == "" {
		id = uuid.New().String()
		rq.Set("id", id)
	}

	resultC := make(chan *simplejson.Json, 1)
	errC := make(chan error, 1)

	// Subscribe on the current socket, but be prepared to re-subscribe on retries
	sw.mu.Lock()
	curSock := sw.socket
	sw.mu.Unlock()

	subID := curSock.Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
			// Propagate close/timeout or any other errors to the caller.
			errC <- evt.Error
			return
		}
		resp, err := simplejson.NewJson(evt.Body)
		if err != nil {
			errC <- err
			return
		}
		if resp.Get("id").MustString() == id {
			resultC <- resp
		}
	})
	// Ensure we clean up the subscription on whatever socket we last used
	defer func() {
		if curSock != nil && subID != 0 {
			curSock.Unsubscribe(subID)
		}
	}()

	// No re-subscribe/send retries here; base stays neutral.

	jsonBytes, err := rq.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	if err := curSock.Send(web_socket.WriteEvent{Body: jsonBytes}); err != nil {
		// Surface send error as-is; variants decide on retries.
		return nil, fmt.Errorf("send error: %w", err)
	}

	timer := time.NewTimer(sw.timeOut)
	defer timer.Stop()
	for {
		select {
		case resp := <-resultC:
			return resp, nil
		case err := <-errC:
			return nil, err
		case <-timer.C:
			return nil, fmt.Errorf("timeout")
		}
	}
}

func (sw *StreamWrapper) Subscribe(f func(web_socket.MessageEvent), subs ...string) error {
	if len(subs) == 0 {
		return fmt.Errorf("no subscriptions provided")
	}

	var newSubs []string
	for _, sub := range subs {
		if _, ok := sw.handlers[sub]; !ok {
			newSubs = append(newSubs, sub)
		}
	}

	if len(newSubs) == 0 {
		return fmt.Errorf("already subscribed to all provided streams")
	}

	rq := simplejson.New()
	rq.Set("method", "SUBSCRIBE")
	rq.Set("params", newSubs)
	rq.Set("id", uuid.New().String())
	_, err := sw.Call(rq)
	if err != nil {
		return fmt.Errorf("subscribe error: %w", err)
	}

	for _, sub := range newSubs {
		id := sw.socket.Subscribe(f)
		if id == 0 {
			continue
		}
		sw.handlers[sub] = id
	}
	return nil
}

func (sw *StreamWrapper) Unsubscribe(subs ...string) error {
	if len(subs) == 0 {
		return fmt.Errorf("no subscriptions provided")
	}

	var toRemove []string
	for _, sub := range subs {
		if _, ok := sw.handlers[sub]; ok {
			toRemove = append(toRemove, sub)
		}
	}

	if len(toRemove) == 0 {
		return fmt.Errorf("no active subscriptions found")
	}

	rq := simplejson.New()
	rq.Set("method", "UNSUBSCRIBE")
	rq.Set("params", toRemove)
	rq.Set("id", uuid.New().String())

	_, err := sw.Call(rq)
	if err != nil {
		// Network round-trip for UNSUBSCRIBE on public endpoints can be flaky;
		// if we hit a timeout/i/o timeout, consider it non-fatal and proceed with local cleanup.
		// This keeps behavior robust under high latency while ensuring handlers are removed.
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			// proceed silently
		} else if strings.Contains(err.Error(), "i/o timeout") || strings.Contains(err.Error(), "timeout") {
			// proceed silently
		} else {
			return fmt.Errorf("unsubscribe error: %w", err)
		}
	}

	for _, sub := range toRemove {
		id := sw.handlers[sub]
		sw.socket.Unsubscribe(id)
		delete(sw.handlers, sub)
	}

	return nil
}

func (sw *StreamWrapper) ListOfSubscriptions() ([]string, error) {
	rq := simplejson.New()
	rq.Set("method", "LIST_SUBSCRIPTIONS")
	rq.Set("id", uuid.New().String())
	resp, err := sw.Call(rq)
	if err != nil {
		return nil, err
	}

	// Always return a non-nil slice for friendlier callers/tests
	result := make([]string, 0)
	for _, v := range resp.Get("result").MustArray() {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result, nil
}

func (sw *StreamWrapper) GetConnection() web_socket.WebSocketCommonInterface {
	return sw.socket
}

func (sw *StreamWrapper) SetMaxReconnectAttempts(n int) CambitorInterface {
	sw.maxReconnectAttempts = n
	return sw
}

func (sw *StreamWrapper) SetReadTimeout(timeout time.Duration) CambitorInterface {
	sw.readTimeout = &timeout
	if sw.socket != nil {
		sw.socket.SetReadTimeout(timeout)
	}
	return sw
}

func (sw *StreamWrapper) SetWriteTimeout(timeout time.Duration) CambitorInterface {
	sw.writeTimeout = &timeout
	if sw.socket != nil {
		sw.socket.SetWriteTimeout(timeout)
	}
	return sw
}

func (sw *StreamWrapper) SetReconnectInterval(interval time.Duration) CambitorInterface {
	sw.reconnectInterval = interval
	return sw
}

func (sw *StreamWrapper) EnableAutoReconnect() CambitorInterface {
	sw.autoReconnect = true
	// Default intervals if not set
	if sw.reconnectInterval <= 0 {
		sw.reconnectInterval = 500 * time.Millisecond
	}
	if sw.healthInterval <= 0 {
		sw.healthInterval = sw.reconnectInterval
	}
	if sw.healthTimeout <= 0 {
		sw.healthTimeout = 700 * time.Millisecond
	}
	sw.startHealthMonitor()
	return sw
}

func (sw *StreamWrapper) DisableAutoReconnect() {
	sw.autoReconnect = false
	sw.stopHealthMonitor()
}

func (sw *StreamWrapper) SetMessageLogger(logger func(message web_socket.LogRecord)) CambitorInterface {
	// Зберегти й застосувати, якщо з'єднання вже є
	sw.pendingLogger = logger
	if sw.socket != nil {
		sw.socket.SetMessageLogger(logger)
	}
	return sw
}

func (sw *StreamWrapper) SetPingHandler(handler func(string) error) {
	sw.pendingPingHandle = handler
	if sw.socket != nil {
		sw.socket.SetPingHandler(handler)
	}
}

func (sw *StreamWrapper) SetPongHandler(handler func(string) error) {
	sw.pendingPongHandle = handler
	if sw.socket != nil {
		sw.socket.SetPongHandler(handler)
	}
}

// WaitReady checks socket readiness by sending Ping and waiting for Pong without changing any subscription state.
// It polls until ctx is done, using probeInterval between attempts. Returns nil when a Pong is received.
func (sw *StreamWrapper) WaitReady(ctx context.Context, probeInterval time.Duration) error {
	if probeInterval <= 0 {
		probeInterval = 200 * time.Millisecond
	}

	// Preserve previously set pong handler (if any) through StreamWrapper API
	prev := sw.pendingPongHandle

	readyCh := make(chan struct{}, 1)
	// Temporary pong handler that also forwards to previous
	sw.SetPongHandler(func(s string) error {
		select {
		case readyCh <- struct{}{}:
		default:
		}
		if prev != nil {
			return prev(s)
		}
		return nil
	})
	defer sw.SetPongHandler(prev)

	ticker := time.NewTicker(probeInterval)
	defer ticker.Stop()

	for {
		// Check context first
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Snapshot socket
		sw.mu.Lock()
		socket := sw.socket
		sw.mu.Unlock()
		if socket == nil || !socket.IsStarted() {
			// Wait for next probe
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(probeInterval):
				continue
			}
		}

		deadline := time.Now().Add(probeInterval)
		// Send Ping (opcode 0x9) without depending on gorilla/websocket here
		const pingOpcode = 0x9
		_ = socket.GetControl().WriteControl(pingOpcode, []byte("ping"), deadline)

		// Wait for Pong or next tick / context cancel
		select {
		case <-readyCh:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// retry
		}
	}
}

// WaitServerListening actively checks that the TCP port from wsHost is accepting connections.
// It returns true as soon as a TCP dial succeeds before the timeout, otherwise false.
func (sw *StreamWrapper) WaitServerListening(timeout time.Duration) bool {
	// Extract host:port from wsHost. It should already be in host:port form.
	addr := string(sw.wsHost)
	if strings.HasPrefix(addr, "ws://") || strings.HasPrefix(addr, "wss://") {
		// Defensive: strip scheme if accidentally included
		addr = strings.TrimPrefix(strings.TrimPrefix(addr, "ws://"), "wss://")
	}

	deadline := time.Now().Add(timeout)
	for {
		if time.Now().After(deadline) {
			return false
		}
		conn, err := net.DialTimeout("tcp", addr, 250*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return true
		}
		time.Sleep(150 * time.Millisecond)
	}
}

// WaitForPong sends a Ping control frame and waits for a Pong within the given timeout.
// Returns true if Pong is received in time. This is a convenience wrapper for tests and readiness probes.
func (sw *StreamWrapper) WaitForPong(timeout time.Duration) bool {
	sw.mu.Lock()
	socket := sw.socket
	sw.mu.Unlock()

	if socket == nil || !socket.IsStarted() {
		return false
	}

	done := make(chan struct{}, 1)

	// Preserve previously set pong handler (if any)
	prev := sw.pendingPongHandle
	sw.SetPongHandler(func(s string) error {
		select {
		case done <- struct{}{}:
		default:
		}
		if prev != nil {
			return prev(s)
		}
		return nil
	})
	defer sw.SetPongHandler(prev)

	// Send Ping control (opcode 0x9)
	const pingOpcode = 0x9
	deadline := time.Now().Add(timeout)
	if err := socket.GetControl().WriteControl(pingOpcode, []byte("ping"), deadline); err != nil {
		return false
	}

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// ----------------------------
//      HEALTH MONITORING
// ----------------------------

// IsConnected returns true if the main socket exists and has started its loops.
func (sw *StreamWrapper) IsConnected() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.socket != nil && sw.socket.IsStarted() && !sw.socket.IsStopped()
}

// HealthPing sends a Ping using the dedicated health socket and waits for a Pong within the timeout.
// It creates and opens the health socket if needed.
func (sw *StreamWrapper) HealthPing(timeout time.Duration) bool {
	if timeout <= 0 {
		timeout = 700 * time.Millisecond
	}

	// Ensure health socket is open
	if !sw.ensureHealthSocket() {
		return false
	}

	sw.healthMu.Lock()
	socket := sw.healthSocket
	sw.healthMu.Unlock()
	if socket == nil || !socket.IsStarted() {
		return false
	}

	done := make(chan struct{}, 1)
	// Install temporary pong handler
	// We don't preserve external handlers on health socket; it's dedicated to probing
	socket.SetPongHandler(func(string) error {
		select {
		case done <- struct{}{}:
		default:
		}
		return nil
	})
	const pingOpcode = 0x9
	deadline := time.Now().Add(timeout)
	if err := socket.GetControl().WriteControl(pingOpcode, []byte("ping"), deadline); err != nil {
		return false
	}

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Start the health monitor goroutine (idempotent). It will:
// - periodically ping the server via a dedicated health socket;
// - if unhealthy => close the main connection;
// - if healthy again => (re)connect the main connection.
func (sw *StreamWrapper) startHealthMonitor() {
	sw.mu.Lock()
	alreadyRunning := sw.healthRunning
	if !sw.healthRunning {
		sw.healthRunning = true
		sw.healthStopChan = make(chan struct{})
	}
	interval := sw.healthInterval
	timeout := sw.healthTimeout
	sw.mu.Unlock()

	if alreadyRunning {
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-sw.healthStopChan:
				return
			case <-ticker.C:
				healthy := sw.HealthPing(timeout)

				if !healthy {
					// Mark unhealthy and ensure main socket is down
					sw.mu.Lock()
					wasHealthy := sw.healthy
					sw.healthy = false
					sw.paused = true
					socket := sw.socket
					sw.mu.Unlock()

					if wasHealthy {
						// Transition to unhealthy, close main
						if socket != nil {
							_ = socket.Close()
							_ = socket.WaitStopped()
						}
					}
					continue
				}

				// healthy == true
				sw.mu.Lock()
				wasHealthy := sw.healthy
				sw.healthy = true
				needReconnect := sw.socket == nil || sw.socket.IsStopped()
				sw.mu.Unlock()

				// If we were unhealthy -> healthy OR we are healthy but main is down, reconnect now.
				if !wasHealthy || needReconnect {
					// Try to (re)connect main socket using existing settings
					_ = sw.Reconnect(1, 0)
				}
				// Clear pause regardless; readiness gated by IsConnected
				sw.mu.Lock()
				sw.paused = false
				sw.mu.Unlock()
			}
		}
	}()
}

func (sw *StreamWrapper) stopHealthMonitor() {
	sw.mu.Lock()
	sw.stopHealthMonitorLocked()
	sw.mu.Unlock()
}

func (sw *StreamWrapper) stopHealthMonitorLocked() {
	if sw.healthRunning {
		sw.healthRunning = false
		if sw.healthStopChan != nil {
			close(sw.healthStopChan)
			sw.healthStopChan = nil
		}
	}
}

// ensureHealthSocket makes sure the health socket is created and started.
func (sw *StreamWrapper) ensureHealthSocket() bool {
	sw.healthMu.Lock()
	defer sw.healthMu.Unlock()
	if sw.healthSocket != nil && sw.healthSocket.IsStarted() && !sw.healthSocket.IsStopped() {
		return true
	}
	// Close stale
	if sw.healthSocket != nil {
		_ = sw.healthSocket.Close()
		_ = sw.healthSocket.WaitStopped()
		sw.healthSocket = nil
	}
	// Create
	socket, err := sw.factory()
	if err != nil {
		return false
	}
	// Use concise timeouts if configured
	if sw.readTimeout != nil {
		socket.SetReadTimeout(*sw.readTimeout)
	}
	if sw.writeTimeout != nil {
		socket.SetWriteTimeout(*sw.writeTimeout)
	}
	socket.Open()
	_ = socket.WaitStarted()
	sw.healthSocket = socket
	return true
}

// ----------------------------
//    HEALTH/READY API
// ----------------------------

// HealthState represents current health as seen by the monitor.
type HealthState int

const (
	HealthConnecting HealthState = iota
	HealthHealthy
	HealthUnhealthy
)

// SetModeSync sets default wait timeout used by WaitConnected helpers; kept for compatibility.
func (sw *StreamWrapper) SetModeSync(defaultTimeout time.Duration) *StreamWrapper {
	if defaultTimeout <= 0 {
		defaultTimeout = 5 * time.Second
	}
	sw.mu.Lock()
	sw.defaultSyncTimeout = defaultTimeout
	sw.mu.Unlock()
	return sw
}

// HealthStatus returns the current health state.
func (sw *StreamWrapper) HealthStatus() HealthState {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	if sw.healthRunning {
		if sw.healthy {
			return HealthHealthy
		}
		return HealthUnhealthy
	}
	// No monitor: infer from main connection
	if sw.socket != nil && sw.socket.IsStarted() && !sw.socket.IsStopped() {
		return HealthHealthy
	}
	return HealthConnecting
}

// WaitConnected blocks until the main connection is established (IsConnected) and the endpoint is reachable
// when a health monitor is enabled. If no monitor is running, it will proactively probe and attempt one
// lightweight reconnect.
func (sw *StreamWrapper) WaitConnected(ctx context.Context) error {
	// Fast-path
	if sw.HealthStatus() == HealthHealthy && sw.IsConnected() {
		return nil
	}
	// Polling loop with modest sleep to avoid busy-wait
	tick := time.NewTicker(150 * time.Millisecond)
	defer tick.Stop()
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		monitorRunning := sw.isHealthMonitorRunning()
		if !monitorRunning {
			// If monitor disabled, probe directly and drive a reconnect if needed
			_ = sw.HealthPing(700 * time.Millisecond)
			if sw.HealthStatus() == HealthHealthy {
				if sw.IsConnected() {
					return nil
				}
				// Only attempt reconnect ourselves when no monitor is running to avoid races
				_ = sw.Reconnect(1, 0)
				if sw.IsConnected() {
					return nil
				}
			}
		} else {
			// Monitor is active: just wait until it establishes the main connection
			if sw.HealthStatus() == HealthHealthy && sw.IsConnected() {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
		}
	}
}

// WaitHealthy is kept for backward compatibility. Prefer WaitConnected.
func (sw *StreamWrapper) WaitHealthy(ctx context.Context) error { return sw.WaitConnected(ctx) }

// IsPaused reports whether main operations are currently paused by health monitor.
func (sw *StreamWrapper) IsPaused() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.paused
}

// Mode returns current OperationMode (sync/async) in a threadsafe way.
// Mode concept removed from base; async/sync variants gate externally.

// IsAutoReconnectEnabled reports whether auto-reconnect/health monitor is enabled.
func (sw *StreamWrapper) IsAutoReconnectEnabled() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.autoReconnect
}

// Helper: whether monitor is running (without racing callers).
func (sw *StreamWrapper) isHealthMonitorRunning() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.healthRunning
}

// ensureReadyForOperation enforces async/sync behavior before I/O.
// ensureReadyForOperation removed from base; variants control readiness.
