# Typing Publish Participant Reservation Cap

Status: Development source candidate.

The process-local typing anti-flood guard already caps the complete active reservation map at 4096 entries and expires entries after the existing 10-second typing TTL. This continuation adds a second fairness boundary: one participant may hold at most 64 active conversation/user publish reservations at a time.

The participant cap is checked only when a new conversation/user key would be created. An existing key may still refresh after the existing 250 ms minimum interval, even while that participant is at the per-participant ceiling. A successful idle signal clears that conversation/user reservation, and normal TTL cleanup runs before both the global and participant capacity checks.

The cap uses only the same process-local conversation/user reservation keys and UTC timestamps already retained by the anti-flood guard. It does not store message content, draft/editor state, keystrokes, device identifiers, IP addresses, client metadata, presence history, attachment metadata, or telemetry identifiers.

When the participant or global reservation boundary rejects a new `typing=true` signal, the existing `ErrTypingRateLimited` path remains authoritative and the HTTP API continues to return the same low-detail 429 response. Idle signals are not delayed by this limiter.

This is not distributed production abuse protection. The counters are process-local, disappear with the process, and do not coordinate across replicas. Production readiness still requires separately reviewed distributed controls, identity/session acceptance, Wardveil Security and Privacy Shield acceptance, operational observability, deployment evidence, and representative client validation.
