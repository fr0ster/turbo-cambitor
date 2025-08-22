# Cambitor: Dual WebSocket Connections (Main + Service) — Design and Behavior

This document introduces a dual-connection architecture for Cambitor with a main WebSocket (business operations) and a dedicated service WebSocket (health monitoring and lifecycle control). It defines async and sync operation modes, public API, and testing changes.

## Goals
- Robust handling of server crashes and network instability.
- Promptly pause the main connection on missing Pong and resume when recovered.
- Offer client choice between fail-fast and blocking (wait) models.
- Remove test dependency on artificial server readiness checks (waitServerListening) in favor of Ping/Pong-driven health.

## Architecture
- Service connection (healthConn):
  - Dedicated long-lived WebSocket to the same endpoint.
  - Health monitor loop: at `healthInterval`, send Ping and wait for Pong up to `healthTimeout`.
  - Debounce: `failThreshold` misses → Unhealthy; `successThreshold` successes → Healthy.
  - Self-reconnects with backoff (min/max + jitter) and keeps probing.
  - Public API:
    - `HealthStatus() enum {Connecting, Healthy, Unhealthy}`
    - `HealthPing(timeout) bool` — one-shot on-demand ping
    - `WaitHealthy(ctx) error` — block until Healthy or context timeout
    - `EnableHealthMonitor(opts) / DisableHealthMonitor()`

- Main connection (mainConn):
  - Executes business methods: `Call`, `Subscribe`, `Unsubscribe`, `ListOfSubscriptions`.
  - On Unhealthy: enter "paused" (recommended: close socket to avoid half-open state).
  - On Healthy: ensure main is up (Connect → WaitStarted → WaitReady), reapply subscriptions/handlers/timeouts.
  - Internal AutoReconnect is subordinated to the health monitor (or disabled).

## Operation modes
- Async (fail-fast):
  - Public methods check health/paused; if not Healthy — return `ErrNotConnected` immediately.
  - Client uses `WaitHealthy`/`HealthPing`/`HealthStatus` and retries.

- Sync (wait):
  - Before I/O, methods call `WaitHealthy(ctx)` with a per-call timeout; they block until Healthy or timeout.
  - Recommended: one-thread-per-cambitor for simpler coordination.

## Public API (additions)
- `SetModeAsync()` / `SetModeSync(defaultTimeout)`
- `EnableHealthMonitor(opts{ healthInterval, healthTimeout, failThreshold, successThreshold, backoff })`
- `DisableHealthMonitor()`
- `HealthStatus() enum`
- `HealthPing(timeout time.Duration) bool`
- `WaitHealthy(ctx context.Context) error`
- `IsPaused() bool`
- Backward-compat: `WaitForPong/WaitReady` routed via service connection (no gorilla in prod).

## State machines
- Health: `Connecting → Healthy → Unhealthy → (reconnect) → Connecting → Healthy …`
- Main: `Active → Paused → Reconnecting → Active …`
- Triggers:
  - N Ping failures → Unhealthy → `pauseMain()` (close/stop loops).
  - M Pong successes → Healthy → `resumeMain()` (Connect → WaitStarted → WaitReady → reapply).

## Concurrency
- Atomic flags for state; transitions serialized via a small mutex/control goroutine.
- All control/data writes serialized (restler writeMu is already in place).
- Main lifecycle under a lifecycle mutex to avoid concurrent reconnects.

## Defaults
- `healthInterval`: 1s (tests), 5–10s (prod)
- `healthTimeout`: 500–1000ms (tests), 2–3s (prod)
- `failThreshold`: 2–3; `successThreshold`: 1–2
- Service reconnect backoff: 200ms → 2s + jitter
- Sync per-call timeout: 5s in tests; configurable

## Test changes
- Remove `waitServerListening` and local `waitForPong`.
- Auto-reconnect scenarios:
  - Async: after crash+restart, `Call` returns `ErrNotConnected` until `HealthStatus == Healthy`; after `WaitHealthy`, `Call` succeeds.
  - Sync: `Call` blocks until Healthy or timeout; ensure it unblocks after restart.
- Readiness via `HealthPing`/`WaitHealthy` (Cambitor API).

## Edge cases
- If the exchange limits connections, make the service connection optional and document trade-offs.
- Reapply subscriptions after resume automatically (as today).
- Anti-flood: cap Ping frequency.

---

# Implementation Roadmap

1. Types and config in `streamer.go`
   - Modes (Async/Sync), health options, atomics, mutexes, state channels.
   - API: `HealthStatus`, `HealthPing`, `WaitHealthy`, `SetMode*`, `Enable/DisableHealthMonitor`.

2. Service socket
   - Create/maintain a dedicated WS.
   - Ping loop with thresholds (fail/success).
   - Reconnect with backoff + jitter.
   - Publish state and drive transitions.

3. Main integration
   - `pauseMain()` (close/stop), `resumeMain()` (Connect → WaitStarted → WaitReady → reapply configs/handlers/subscriptions).
   - Gate public methods (fail-fast vs wait) based on mode.

4. AutoReconnect refactor
   - Subordinate to health monitor or disable when health monitor is enabled.

5. Tests
   - Update scenarios (Async/Sync), remove `waitServerListening`/local Pong helpers.
   - Add checks using `HealthPing/WaitHealthy/HealthStatus`.

6. Docs
   - Modes, options, migration notes.

---

## Pros and Cons of the Dual-Connection Design

### Pros
- Faster, deterministic health detection independent of main data flow.
- Clean separation of concerns: control-plane (health) vs data-plane (business).
- Smoother auto-pause/resume with fewer half-open states (close-on-pause policy).
- Predictable tests and debugging without external server readiness hacks.
- Backward compatible surface; readiness helpers route via the service channel.

### Cons
- Additional socket per client increases resource usage and may hit exchange connection limits.
- Higher complexity: extra state machine, lifecycle mutexes, and reconciliation logic.
- Tuning required (intervals, thresholds, backoff) to avoid flapping or slow detection.
- Risk of state divergence (service says Healthy while main not yet resumed) without strict serialization.
- Potential reconnect storms across many instances if not jittered and rate-limited.

## Security Considerations

- Transport security:
  - Use WSS in production; validate certificates and consider pinning where feasible.
  - Enforce sane read/write timeouts to avoid resource hangs and Slowloris-like effects.

- Authentication and secrets:
  - Do not log API keys, auth headers, or signed payloads. Scrub logs by default.
  - On auth failures (4xx/invalid key), do not auto-retry indefinitely; trip a permanent-fail state.

- DoS and rate limiting:
  - Cap health ping frequency; respect exchange heartbeat policies. Default to conservative intervals.
  - Add jittered backoff on service reconnect to avoid thundering herds.
  - Gate resubscription storms after resume (batch or stagger re-subscribe with small delays).

- Concurrency safety:
  - Serialize all control/data writes (use underlying write mutex); avoid concurrent close/send.
  - Protect lifecycle with a dedicated mutex; ensure idempotent pause/resume.
  - Prevent goroutine/leak on state transitions; tie loops to contexts and WaitStopped.

- Integrity and replay:
  - On resume, re-fetch exchange info/time if required and verify stream alignment.
  - Consider sequence numbers or last update IDs where applicable to avoid gaps.

- Configuration hardening:
  - Provide upper/lower bounds for intervals/timeouts to prevent unsafe configs.
  - Feature-flag the service channel for venues with strict connection caps.

- Observability and auditing:
  - Emit structured events on state transitions (Connecting/Healthy/Unhealthy, pause/resume).
  - Expose metrics (health RTT, failures, resumes, backoff level) for alerting.
