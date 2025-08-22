# turbo-cambitor
Low-level wrappers for working with crypto exchanges over their REST/WebSocket APIs.

## Highlights
- Binance WebStream with auto-reconnect and health-driven readiness.
- Dual-connection architecture: a dedicated service socket performs Ping/Pong health checks and supervises the main business socket.
- Two operation modes:
  - Async (default): fail fast on Unhealthy.
  - Sync: wait for Healthy before I/O up to a per-call timeout.

## Quick start (WebStream)
```go
sw := web_stream.
    New(true).
    Stream().
    SetSymbol("BTCUSDT")

// Choose mode
sw.SetModeAsync()           // fail-fast
// or
// sw.SetModeSync(2 * time.Second) // wait up to 2s before each I/O

// Connect and wait for health if needed
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
if err := sw.WaitHealthy(ctx); err != nil {
    // handle not healthy
}

// Now safe to call
resp, err := sw.ListSubscriptions()
_ = resp; _ = err
```

## Docs
- Dual-connection design (v0.3.0): `docs/DUAL_CONNECTION_DESIGN_v0.3.0_EN.md`, `docs/DUAL_CONNECTION_DESIGN_v0.3.0_UA.md`
- Roadmap and tests (v0.3.0): `docs/DUAL_CONNECTION_ROADMAP_v0.3.0_EN.md`, `docs/DUAL_CONNECTION_ROADMAP_v0.3.0_UA.md`
- Legacy drafts: `docs/DUAL_CONNECTION_DESIGN_EN.md`, `docs/DUAL_CONNECTION_DESIGN_UA.md`, `docs/DUAL_CONNECTION_ROADMAP_EN.md`, `docs/DUAL_CONNECTION_ROADMAP_UA.md`

## Notes
- Tests that require exchange credentials are gated by API key presence only.
