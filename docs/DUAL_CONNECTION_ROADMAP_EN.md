# Roadmap: Dual WebSocket Connections (Main + Service)

This document details phases, API, tasks, acceptance criteria, test plan, and risks for implementing a dual-connection architecture in Cambitor.

## 0) Scope
- Add a service (health) WebSocket for Ping/Pong health monitoring.
- Control the main connection via pause/resume states.
- Add operation modes: Async (fail-fast) and Sync (wait).
- Update public API and tests; remove dependency on `waitServerListening` and local Pong helpers.

## 1) Public API (signatures)
- Config/modes:
  - `func (sw *StreamWrapper) SetModeAsync() *StreamWrapper`
  - `func (sw *StreamWrapper) SetModeSync(defaultTimeout time.Duration) *StreamWrapper`
  - `type HealthOptions struct { Interval, Timeout, BackoffMin, BackoffMax time.Duration; FailThreshold, SuccessThreshold int; Jitter float64 }`
  - `func (sw *StreamWrapper) EnableHealthMonitor(opts HealthOptions) *StreamWrapper`
  - `func (sw *StreamWrapper) DisableHealthMonitor() *StreamWrapper`
- Health:
  - `type HealthState int` + `const (HealthConnecting HealthState = iota; HealthHealthy; HealthUnhealthy)`
  - `func (sw *StreamWrapper) HealthStatus() HealthState`
  - `func (sw *StreamWrapper) HealthPing(timeout time.Duration) bool`
  - `func (sw *StreamWrapper) WaitHealthy(ctx context.Context) error`
  - Back-compat: `WaitForPong(timeout)`/`WaitReady(ctx, interval)` routed via the service channel.

### Function-level tests (step-by-step scenarios)

Below are concrete test scenarios for each new/updated API.

1. `SetModeAsync()`
  - Pre: HealthMonitor enabled, server Unhealthy (mock drops Pong N times).
  - Act: call `sw.SetModeAsync()`, then `sw.Call(LIST_SUBSCRIPTIONS)`.
  - Expect: immediate `ErrNotConnected`; after recovery (mock resumes Pong), retry succeeds.

2. `SetModeSync(defaultTimeout)`
  - Pre: Unhealthy.
  - Act: `sw.SetModeSync(2s)`, invoke `sw.Call(...)` in a goroutine; after 1s, recover health (mock Pong/restart).
  - Expect: call blocks < 2s, then succeeds; if not recovered, returns `context deadline exceeded`.

3. `EnableHealthMonitor(opts)`
  - Pre: monitor off; set `opts.Interval=200ms`, `Timeout=150ms`, `FailThreshold=2`, `SuccessThreshold=1`.
  - Act: `EnableHealthMonitor(opts)`; force mock to miss 2 Pongs, then respond.
  - Expect: states `Connecting→Healthy→Unhealthy` (after two misses) → `Healthy` (after success); verify via `HealthStatus()` and timing.

4. `DisableHealthMonitor()`
  - Pre: monitor on and pinging.
  - Act: `DisableHealthMonitor()`; observe over ~2*Interval.
  - Expect: no new Pings (state does not degrade to Unhealthy despite induced misses; optionally count Pings on mock).

5. `HealthStatus()`
  - Pre: freshly started healthConn.
  - Act/Expect: initially `Connecting`; after first Pong — `Healthy`; on consecutive failures — `Unhealthy`.

6. `HealthPing(timeout)`
  - Case 1: server alive — returns `true` < timeout.
  - Case 2: server unresponsive — returns `false` ≈ timeout.
  - Case 3: during reconnect — ensure no panic and a correct boolean outcome.

7. `WaitHealthy(ctx)`
  - Case 1: server Healthy — returns quickly (<100ms with small intervals).
  - Case 2: server Unhealthy — with 500ms timeout returns `ctx.Err()`.
  - Case 3: Unhealthy→Healthy within the window — returns `nil` before deadline.

8. `IsPaused()` (indirect pause/resume test)
  - Pre: active subscriptions on mainConn.
  - Act: drive health to `Unhealthy` (misses ≥ FailThreshold).
  - Expect: `IsPaused()==true`; mainConn closed; `Subscribe/Call` in Async mode → `ErrNotConnected`, in Sync → blocks until `Healthy`.
  - After recovery: `IsPaused()==false`; subscriptions re-applied; mock emits data to client.

9. Back-compat: `WaitForPong(timeout)` / `WaitReady(ctx, interval)`
  - Ensure they use the service connection; return `true` on Pong and `false` on timeout.

10. Coordination with `EnableAutoReconnect()`
  - Pre: HealthMonitor enabled; AutoReconnect also enabled.
  - Act: trigger crash+restart in mock.
  - Expect: no duplicate reconnect/resubscribe; single pause→resume cycle; data resumes.

## 2) Internal structures/fields
- In `StreamWrapper` add:
  - `healthConn web_socket.WebSocketCommonInterface`
  - `healthState atomic.Value` (HealthState)
  - `paused atomic.Bool`
  - `mode enum` (Async/Sync) + `defaultSyncTimeout time.Duration`
  - `lifecycleMu sync.Mutex` (main lifecycle)
  - `healthMu sync.Mutex` (service lifecycle)
  - Contexts/cancels for loops, state channels
- Reuse existing: handler registration, timeouts, message logger, resubscription logic.

## 3) Phases

### Phase A: Health API skeleton (no loop)
- Add types, fields, `HealthStatus`, `HealthPing`, `WaitHealthy` (temporary control Ping).
- `EnableHealthMonitor` stores options only; no background monitor yet.
- Tests: switch to `HealthPing/WaitHealthy` without fully removing old helpers (mark deprecated).

### Phase B: Service socket & monitor
- Create/maintain `healthConn` via factory.
- Implement ping loop with thresholds (`FailThreshold/SuccessThreshold`).
- Implement reconnect with backoff + jitter.
- Publish state transitions.
- Metrics/logs: RTT, failure/success counts, backoff level.

### Phase C: Main pause/resume
- `pauseMain()`: set `paused=true`, close mainConn cleanly, stop loops.
- `resumeMain()`: bring up mainConn (Open → WaitStarted → WaitReady), reapply timeouts/handlers/logger, resubscribe.
- Serialize with `lifecycleMu`, ensure idempotency.

### Phase D: Async/Sync modes
- In public methods (Call/Subscribe/Unsubscribe/ListOfSubscriptions):
  - Async: if `paused` or `HealthStatus!=Healthy` → `ErrNotConnected`.
  - Sync: `WaitHealthy(ctx)` with per-call timeout before I/O.
- Configure defaults/timeouts.

### Phase E: Test cleanup and compatibility
- Comment out `waitServerListening` and local `waitForPong` in `streamer_test.go`.
- Rewrite reconnect scenarios using `HealthPing/WaitHealthy/HealthStatus`.
- Add tests for flapping, Sync-mode timeouts, and healthConn degradation.

### Phase F: Docs/Examples
- Update README and add usage snippets for Async/Sync.
- Include migration notes and behavior with AutoReconnect conflicts.

## 4) Acceptance Criteria
- All unit/integration tests pass stably, including under debugger.
- No panics or races (`-race` clean).
- No gorilla dependency in prod Cambitor (use restler API).
- Reconnect/pause/resume deterministic; subscriptions are restored automatically.
- Ping frequency within safe bounds; no server flooding.

## 5) Test Plan
- HealthConn:
  - Successful Ping/Pong; Pong timeout; N consecutive failures → Unhealthy; successes → Healthy.
  - Reconnect with backoff; jitter; monitoring persists after recovery.
- MainConn:
  - Pause on Unhealthy (APIs fail-fast or block based on mode).
  - Resume on Healthy: Connect → WaitStarted → WaitReady → resubscribe; data flows.
- Modes:
  - Async: `ErrNotConnected` during Unhealthy; success after `WaitHealthy`.
  - Sync: block until Healthy/timeout; verify wait duration.
- Flapping/stress:
  - Frequent state changes; no deadlocks/panics; no resubscription storms.
- Security:
  - No secrets in logs; ping caps; timeouts respected; structured logs.

## 6) Risks and Mitigation
- Connection limits — make service channel optional; warn in logs.
- Excessive pings — enforce bounds; conservative defaults.
- Dual reconnection logic — prioritize health monitor; disable or delegate AutoReconnect.
- State divergence — serialize transitions; idempotent pause/resume.

## 7) Rollout
- Phase 1 (A+B): API skeleton + healthConn with monitor behind a feature flag.
- Phase 2 (C+D): Pause/Resume + modes; test coverage; remove old helpers from tests.
- Phase 3 (E+F): Cleanup, docs, examples, integration validation.
