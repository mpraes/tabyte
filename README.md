# Tabyte

Open-source local tool for **storage estimation** and **structural performance signals** from relational DDLs. First engines: **SQL Server** and **PostgreSQL**.

A CLI starts a local HTTP API on `127.0.0.1` and serves an embedded web UI for paste/analyze workflows.

Requirements and design notes live under [`docs/`](docs/).

| Doc | Purpose |
|---|---|
| [`docs/functional_requirements.md`](docs/functional_requirements.md) | Functional requirements (RF-01…RF-24) |
| [`docs/non_functional_requirements.md`](docs/non_functional_requirements.md) | Non-functional requirements |
| [`docs/endpoints.md`](docs/endpoints.md) | Local HTTP API contract |
| [`docs/architecture.md`](docs/architecture.md) | Architecture decisions |
| [`docs/storage_estimates.md`](docs/storage_estimates.md) | Estimation rules reference |
| [`docs/suggested_stack.md`](docs/suggested_stack.md) | Suggested stack |
| [`docs/sqlite_entities.md`](docs/sqlite_entities.md) | Optional SQLite entity model |

---

## Status

**Local product is usable**: API + embedded UI, optional SQLite persistence.

Version reported by `/api/v1/info`: `0.0.1`.

### Done

| Area | What works |
|---|---|
| CLI | `tabyte serve` (`--no-open`, `--persist`) |
| Localhost | Loopback-only bind; opens browser by default |
| Web UI | Engine select, paste DDL, Analyze, session history, settings, table detail, growth, export, warnings/signals/insights |
| System API | `GET /api/v1/health`, `GET /api/v1/info` (`external_required: false`) |
| Sessions | Create, get, list, delete, reprocess engine, export JSON/CSV |
| Estimates | Columns, rows (`calculation`), tables, schema totals (+ human), indexes |
| Growth | `PATCH .../tables/{name}/growth` |
| Alerts | Structural warnings (RF-18) and signals (RF-19) |
| Persistence | Optional SQLite settings + session history (`--persist`) |
| AI hook | `GET .../insights` with disabled provider (RF-24; no external calls) |
| Tests | Unit tests + [`tests/smoke.sh`](tests/smoke.sh) |

### Still open (product polish)

| Item | Notes |
|---|---|
| Dedicated result routes | `/summary`, `/tables`, `/warnings` (today nested in session JSON) |
| File import | `POST /imports/sql` |
| Live AI provider | Extension interface exists; default is disabled/noop |

---

## Quick start

Requires Go (see `go.mod`).

```bash
go run ./cmd/tabyte serve
```

With optional persistence:

```bash
go run ./cmd/tabyte serve --persist
```

Server prints:

```text
Tabyte listening on http://127.0.0.1:8787
```

Smoke test (server must already be running):

```bash
./tests/smoke.sh
```

### Example: create a session

```bash
curl -s -X POST http://127.0.0.1:8787/api/v1/analysis-sessions \
  -H 'Content-Type: application/json' \
  -d '{
    "engine": "postgres",
    "source_name": "demo.sql",
    "ddl_text": "CREATE TABLE users (id INT, name VARCHAR(100));"
  }'
```

Useful follow-ups:

```bash
# list sessions
curl -s http://127.0.0.1:8787/api/v1/analysis-sessions

# update assumed rows for a table
curl -s -X PATCH http://127.0.0.1:8787/api/v1/analysis-sessions/{id}/tables/users \
  -H 'Content-Type: application/json' \
  -d '{"assumed_row_count": 5000}'
```

---

## Layout

```text
cmd/tabyte/                 CLI entrypoint
internal/runtime/           process bootstrap (serve)
internal/httpapi/           HTTP handlers / JSON envelope
internal/application/       session orchestration, volume updates
internal/domain/            core models (session, table, column, engine)
internal/parser/            DDL structure extraction
internal/engine/postgres/   Postgres normalize + estimate
internal/engine/sqlserver/  SQL Server normalize + estimate
internal/persistence/sqlite/  optional local store (--persist)
internal/platform/          OS helpers (browser open)
web/                        embedded UI assets
docs/                       product and API docs
tests/smoke.sh              end-to-end API smoke checks
```

---

## Design notes

- Estimates are **not** physical measurements on a live instance (RN-02).
- Engine choice changes semantics for types, row overhead, and indexes (RN-01).
- Optional AI insights annotate results only; they never replace the calculation engine (RF-24).

See [`docs/functional_requirements.md`](docs/functional_requirements.md) for the full RF list and first-delivery priority.
