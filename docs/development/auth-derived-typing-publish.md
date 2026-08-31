# Auth-Derived Typing Publication — Development

Typing publication no longer accepts a caller-authored `user_id` field in the HTTP request body. The shared authenticated runtime is the only source of the publishing user identity, and strict JSON rejects attempts to reintroduce a body-level identity field.

The request body is minimized to `sequence` and `state`. The service still enforces conversation membership, publish privacy preference, sequence monotonicity, and short-lived typing state. Observation continues to expose only the minimized active-typing projection.

This strengthens the trust boundary without changing the underlying Development memory-backed typing store or privacy preference persistence.

This slice does not establish production GoreeCloud Identity sessions, durable Privacy Shield preference backing, production presence fan-out, native client typing UI, deployment acceptance, or Stable qualification.
