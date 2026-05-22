# arnada

A lightweight, high-performance HTTP reverse proxy gateway built from scratch in Go. 

## Why This Exists

`arnada` is not built to replace production enterprise giants like Caddy or Traefik. Instead, it was created as a deep-dive engineering exercise to understand the mechanics under the hood of modern edge routers: connection streaming, interface wrapping, and request lifecycle management. 

Building it from the ground up allows for deep architectural customization and a playground for exploring low-level networking primitives in Go.

## What It Achieves

* **Zero-Allocation Context Pipeline:** Cooperatively builds a single `RequestContext` object across middlewares without context map sprawl, ensuring a massive reduction in memory thrashing.
* **Preserved Interface Streaming:** A custom `StatWriter` layer that accurately tracks metrics while safely forwarding low-level `http.Flusher` and `http.Hijacker` interfaces to prevent streaming hangs.
* **Single-Point Observability:** Consolidates multi-step gateway lifecycles into exactly **one** comprehensive, structured JSON log entry per complete request.
* **Concurrency-Safe Routing:** Dynamically pairs incoming Host headers to target destinations using a pre-cached, thread-safe proxy manager map.

## Current Architecture

```text
Incoming Request 
      │
      ▼
 1. Logger Middleware      ──► Injects StatWriter, registers top-level defer safety net
      │
      ▼
 2. Context Initializer   ──► Spans a blank RequestContext pointer into the writer & context
      │
      ▼
 3. Request ID Middleware  ──► Resolves/assigns X-Request-ID; populates initial HTTP metadata
      │
      ▼
 4. Reverse Proxy Handler  ──► Matches Host to upstream target and pipes traffic via httputil
      │
      ▼
 Response Unwinding        ──► StatWriter captures final status/bytes; deferred Logger fires ONE JSON log
 ```

 Future Roadmap

    [ ] Path-Based Routing: Move beyond basic Host header matching to support exact and prefix path routing (e.g., /api/v1/*).

    [ ] Persistent Log Pipeline: Transition from standard out streaming to a dedicated database storage layer.

    [ ] Admin Web Interface: Build a lightweight, real-time GUI dashboard to search, filter, and inspect live request streams visually.