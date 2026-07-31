# Tabyte — Non-Functional Requirements v0.3

## Objective

This document defines **only** Tabyte’s non-functional requirements: quality attributes and technical constraints, not business capabilities. Categories follow quality-characteristic thinking similar to ISO/IEC 25010 (performance, compatibility, usability, reliability, security, maintainability, portability).

Assumed context: a **local-first** app distributed as a **cross-platform CLI**, with a browser UI on `localhost`, implemented in **Go**, with assets embedded in the binary and optional **SQLite** for internal persistence.

## Scope

This document does not define what Tabyte does for the business (analyze DDL, estimate volume, expose a given endpoint) except where that directly affects a quality attribute.

## Performance and efficiency

### RNF-01 — Startup time

The app should start quickly on a typical local machine: load config, start the local server, and serve the web UI with a utility-tool feel.

### RNF-02 — Interactive response time

Interactive operations should feel fluid for small and medium schemas. Normal local use should not feel frozen.

### RNF-03 — Resource efficiency

CPU and memory use should suit common developer workstations and laptops. The product must not assume heavy infrastructure for local use.

### RNF-04 — Low operational dependency

Final-user environments should need as few external dependencies as possible. Prefer a self-contained binary.

## Compatibility and portability

### RNF-05 — Cross-platform parity

Main local flows should behave equivalently on Windows and Ubuntu/Linux. OS-specific differences should be minimized and documented when unavoidable.

### RNF-06 — Simple installation

Getting the binary should be simple for a technical user, with a short, predictable, documented entry path.

### RNF-07 — Portable builds

The project should build and package for Windows and Linux without an overly complex pipeline.

### RNF-08 — No external database server

Normal operation must not require an external database server. Local persistence, when used, is embedded and transparent (application file format).

## Technical usability

### RNF-09 — Clear operability

Start, use, and stop should be understandable for a technical user. Status, error, and success messages should be short and useful.

### RNF-10 — Interface consistency

The local web UI should stay consistent in visuals, terminology, and structure across panels and messages.

### RNF-11 — Error prevention

Validations, clear messages, and immediate feedback should reduce avoidable operator errors on incomplete or invalid input.

### RNF-12 — Basic accessibility

The web UI should meet basic accessibility expectations: keyboard use, visible focus, sufficient contrast, and semantic markup.

## Reliability and resilience

### RNF-13 — Deterministic behavior

Under the same input, configuration, and version, behavior and results should be consistent and reproducible.

### RNF-14 — Partial-failure tolerance

Localized errors should degrade in a controlled way. Prefer clear messages over abrupt process death when isolation is possible.

### RNF-15 — Operational recoverability

After a non-critical failure or unexpected exit, restart should be low effort. Local persistence must not become a frequent corruption/manual-recovery hotspot.

### RNF-16 — Local observability

Logs should be sufficient for local diagnosis without external observability infrastructure.

## Security and privacy

### RNF-17 — Restricted local bind

By default, accept connections only on `localhost` / `127.0.0.1`. Any relaxation must be explicit and non-default.

### RNF-18 — Privacy of analyzed content

User content stays local by default. Nothing submitted is sent to external services without explicit user action or conscious configuration.

### RNF-19 — Minimal necessary persistence

When persistence is on, store only what the experience needs. Prefer simplicity, predictability, and local control.

### RNF-20 — Transparency of external integration

Future external integrations (e.g. AI providers) must clearly communicate privacy and data-traffic effects.

## Maintainability

### RNF-21 — Internal modularity

Clear module boundaries among UI, CLI runtime, domain, persistence, and integrations. Modularity is central to safe evolution.

### RNF-22 — Testability

Architecture must allow automated tests at unit, integration, and local-runtime smoke levels.

### RNF-23 — Modifiability

Incremental changes in persistence, UI, or future integrations should not require rewriting the whole core.

### RNF-24 — Analyzability

Codebase structure, naming, and separation of domain vs infrastructure should be understandable to contributors.

### RNF-25 — Technical traceability

Relevant decisions on local runtime, packaging, persistence, and quality should be traceable in project docs.

## Technology constraints

### RNF-26 — Primary language

Primary implementation language is **Go**.

### RNF-27 — Assets embedded in the binary

Static web UI assets are packaged with Go `embed`.

### RNF-28 — Optional persistence via SQLite

When local persistence is needed, the preferred mechanism is **SQLite**.

### RNF-29 — Prefer no-CGO when viable

Prefer components that keep cross-platform build/runtime simple; avoid CGO requirements when practical (e.g. pure-Go SQLite driver).

## Local deployment and operation

### RNF-30 — Single CLI runtime

Primary operational model: CLI starts the local server; UI is used in the browser. Same pattern on Windows and Linux.

### RNF-31 — No mandatory auxiliary services

Default install must not require companion processes, resident agents, or parallel local infrastructure.

### RNF-32 — Clean shutdown

Shutdown should not leave orphan processes, stuck ports, or unnecessary temp files.

## Priority

Most critical for first delivery: RNF-04, RNF-05, RNF-06, RNF-08, RNF-09, RNF-13, RNF-17, RNF-21, RNF-22, RNF-26, RNF-27, RNF-28, RNF-30.
