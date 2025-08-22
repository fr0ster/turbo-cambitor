package streamer

import (
	"context"
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

	// Operation mode and defaults
	mode               OperationMode
	defaultSyncTimeout time.Duration

	handlers map[string]int
}

// NewStreamWrapper creates a new StreamWrapper with a socket factory
func NewStreamWrapper(
	factory func() (web_socket.WebSocketCommonInterface, error),
	wsScheme common.WsScheme,
	wsHost common.WsHost,
	wsEndpoint common.WsEndpoint,
	timeOut ...time.Duration) *StreamWrapper {
	if len(timeOut) == 0 {
		timeOut = append(timeOut, time.Second)
	}
	return &StreamWrapper{
		socket:             nil, // defer dialing until Connect()
		factory:            factory,
		wsScheme:           wsScheme,
		wsHost:             wsHost,
		wsEndpoint:         wsEndpoint,
		timeOut:            timeOut[0],
		handlers:           make(map[string]int),
		mode:               ModeAsync,
		defaultSyncTimeout: 5 * time.Second,
	}
}

func (sw *StreamWrapper) Connect() (*StreamWrapper, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	// Initialize socket on first Connect
	if sw.socket == nil {
		socket, err := sw.factory()
		if err != nil {
			return nil, fmt.Errorf("failed to create websocket client: %w", err)
		}
		sw.socket = socket

		// Застосувати відкладені налаштування (логер/хендлери/таймаути)
		if sw.pendingLogger != nil {
			sw.socket.SetMessageLogger(sw.pendingLogger)
		}
		if sw.readTimeout != nil {
			sw.socket.SetReadTimeout(*sw.readTimeout)
		}
		if sw.writeTimeout != nil {
			sw.socket.SetWriteTimeout(*sw.writeTimeout)
		}
		if sw.pendingPingHandle != nil {
			sw.socket.SetPingHandler(sw.pendingPingHandle)
		}
		if sw.pendingPongHandle != nil {
			sw.socket.SetPongHandler(sw.pendingPongHandle)
		}
	}
	sw.socket.Open()
	// Дочекайся старту лупів, щоб уникнути гонок ініціалізації в тестах/дебазі
	_ = sw.socket.WaitStarted()
	return sw, nil
}

func (sw *StreamWrapper) Reconnect(maxAttempts int, delay time.Duration) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	for i := 0; i < maxAttempts; i++ {
		if sw.socket != nil {
			// Close current socket and ensure loops are fully stopped.
			// Use WaitStopped() to avoid blocking on a stale Stopped() channel
			// in cases where Close/Halt already waited and consumed the signal.
			_ = sw.socket.Close()
			_ = sw.socket.WaitStopped()
		}
		socket, err := sw.factory()
		if err == nil {
			// Apply the same pending settings as in Connect()
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
			socket.Open()
			// Гарантуємо, що лупи запущені перед поверненням нового сокета
			_ = socket.WaitStarted()
			sw.socket = socket
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("failed to reconnect after %d attempts", maxAttempts)
}

func (sw *StreamWrapper) Disconnect() {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	if sw.socket != nil {
		sw.socket.Close()
		sw.socket = nil
	}
	// Also stop and close health socket
	sw.stopHealthMonitorLocked()
	sw.healthMu.Lock()
	if sw.healthSocket != nil {
		_ = sw.healthSocket.Close()
		sw.healthSocket = nil
	}
	sw.healthMu.Unlock()
}

func (sw *StreamWrapper) Call(rq *simplejson.Json) (*simplejson.Json, error) {
	// Gate by mode/health
	if err := sw.ensureReadyForOperation(); err != nil {
		return nil, err
	}
	id := rq.Get("id").MustString()
	if rq.Get("id").MustString() == "" {
		id = uuid.New().String()
		rq.Set("id", id)
	}

	resultC := make(chan *simplejson.Json, 1)
	errC := make(chan error, 1)

	subID := sw.socket.Subscribe(func(evt web_socket.MessageEvent) {
		if evt.Error != nil {
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
	defer sw.socket.Unsubscribe(subID)

	jsonBytes, err := rq.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	if err := sw.socket.Send(web_socket.WriteEvent{Body: jsonBytes}); err != nil {
		return nil, fmt.Errorf("send error: %w", err)
	}

	select {
	case resp := <-resultC:
		return resp, nil
	case err := <-errC:
		return nil, err
	case <-time.After(sw.timeOut):
		return nil, fmt.Errorf("timeout")
	}
}

func (sw *StreamWrapper) Subscribe(f func(web_socket.MessageEvent), subs ...string) error {
	if err := sw.ensureReadyForOperation(); err != nil {
		return err
	}
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
	if err := sw.ensureReadyForOperation(); err != nil {
		return err
	}
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
		return fmt.Errorf("unsubscribe error: %w", err)
	}

	for _, sub := range toRemove {
		id := sw.handlers[sub]
		sw.socket.Unsubscribe(id)
		delete(sw.handlers, sub)
	}

	return nil
}

func (sw *StreamWrapper) ListOfSubscriptions() ([]string, error) {
	if err := sw.ensureReadyForOperation(); err != nil {
		return nil, err
	}
	rq := simplejson.New()
	rq.Set("method", "LIST_SUBSCRIPTIONS")
	rq.Set("id", uuid.New().String())
	resp, err := sw.Call(rq)
	if err != nil {
		return nil, err
	}

	var result []string
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

func (sw *StreamWrapper) SetMaxReconnectAttempts(n int) StreamInterface {
	sw.maxReconnectAttempts = n
	return sw
}

func (sw *StreamWrapper) SetReadTimeout(timeout time.Duration) StreamInterface {
	sw.readTimeout = &timeout
	if sw.socket != nil {
		sw.socket.SetReadTimeout(timeout)
	}
	return sw
}

func (sw *StreamWrapper) SetWriteTimeout(timeout time.Duration) StreamInterface {
	sw.writeTimeout = &timeout
	if sw.socket != nil {
		sw.socket.SetWriteTimeout(timeout)
	}
	return sw
}

func (sw *StreamWrapper) SetReconnectInterval(interval time.Duration) StreamInterface {
	sw.reconnectInterval = interval
	return sw
}

func (sw *StreamWrapper) EnableAutoReconnect() StreamInterface {
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

func (sw *StreamWrapper) SetMessageLogger(logger func(message web_socket.LogRecord)) StreamInterface {
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

				if !wasHealthy {
					// Transition back to healthy
					if needReconnect {
						// Try to (re)connect main socket using existing settings
						_ = sw.Reconnect(1, 0)
					}
					// Clear pause regardless; readiness gated by IsConnected
					sw.mu.Lock()
					sw.paused = false
					sw.mu.Unlock()
				}
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
//    HEALTH/READY API (NEW)
// ----------------------------

// OperationMode defines how public methods behave under unhealthy state.
type OperationMode int

const (
	// ModeAsync returns ErrNotConnected immediately when unhealthy/paused.
	ModeAsync OperationMode = iota
	// ModeSync blocks in WaitHealthy (with default timeout) before proceeding.
	ModeSync
)

// HealthState represents current health as seen by the monitor.
type HealthState int

const (
	HealthConnecting HealthState = iota
	HealthHealthy
	HealthUnhealthy
)

// SetModeAsync switches Cambitor to fail-fast mode.
func (sw *StreamWrapper) SetModeAsync() *StreamWrapper {
	sw.mu.Lock()
	sw.mode = ModeAsync
	sw.mu.Unlock()
	return sw
}

// SetModeSync switches Cambitor to wait mode with a default per-call timeout.
func (sw *StreamWrapper) SetModeSync(defaultTimeout time.Duration) *StreamWrapper {
	if defaultTimeout <= 0 {
		defaultTimeout = 5 * time.Second
	}
	sw.mu.Lock()
	sw.mode = ModeSync
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

// WaitHealthy blocks until health is Healthy and main connection is up, or ctx is done.
func (sw *StreamWrapper) WaitHealthy(ctx context.Context) error {
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
		// If monitor disabled, probe directly
		if !sw.isHealthMonitorRunning() {
			_ = sw.HealthPing(700 * time.Millisecond)
		}
		if sw.HealthStatus() == HealthHealthy {
			if sw.IsConnected() {
				return nil
			}
			// Try a lightweight reconnect attempt if healthy but main is down
			_ = sw.Reconnect(1, 0)
			if sw.IsConnected() {
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

// IsPaused reports whether main operations are currently paused by health monitor.
func (sw *StreamWrapper) IsPaused() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.paused
}

// Helper: whether monitor is running (without racing callers).
func (sw *StreamWrapper) isHealthMonitorRunning() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.healthRunning
}

// ensureReadyForOperation enforces async/sync behavior before I/O.
func (sw *StreamWrapper) ensureReadyForOperation() error {
	sw.mu.Lock()
	mode := sw.mode
	sw.mu.Unlock()

	if mode == ModeAsync {
		if !sw.IsConnected() || sw.HealthStatus() != HealthHealthy {
			return fmt.Errorf("not connected")
		}
		return nil
	}

	// ModeSync: wait with default timeout
	ctx, cancel := context.WithTimeout(context.Background(), sw.defaultSyncTimeout)
	defer cancel()
	if err := sw.WaitHealthy(ctx); err != nil {
		return fmt.Errorf("not connected: %w", err)
	}
	if !sw.IsConnected() {
		return fmt.Errorf("not connected")
	}
	return nil
}
