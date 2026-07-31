# Tabyte

Open-source local tool for **storage estimation** and **structural performance signals** from relational DDLs. First engines: **SQL Server** and **PostgreSQL**.

You run a CLI that starts a local HTTP API (`127.0.0.1`). A web UI is planned to sit on top of that API.

Requirements and design notes live under [`docs/`](docs/).

| Doc | Purpose |
|---|---|
| [`docs/requisitos_funcionais.md`](docs/requisitos_funcionais.md) | Functional requirements (RF-01…RF-24) |
| [`docs/requisitos_nao_funcionais.md`](docs/requisitos_nao_funcionais.md) | Non-functional requirements |
| [`docs/endpoints.md`](docs/endpoints.md) | Local HTTP API contract |
| [`docs/arquitetura.md`](docs/arquitetura.md) | Architecture decisions |
| [`docs/estimativas_postgres_mssql.md`](docs/estimativas_postgres_mssql.md) | Estimation rules reference |
| [`docs/stack_sugerida.md`](docs/stack_sugerida.md) | Suggested stack |
| [`docs/entidades_sqlite.md`](docs/entidades_sqlite.md) | Optional SQLite entity model |

---

## Status

**Backend core is usable via API** (local `127.0.0.1`), with embedded stub UI and optional SQLite persistence.

Version reported by `/api/v1/info`: `0.0.1`.

### Done

| Area | What works |
|---|---|
| CLI | `tabyte serve` (`--no-open`, `--persist`) |
| Localhost | Loopback-only bind (RF-23); opens browser by default |
| System API | `GET /api/v1/health`, `GET /api/v1/info` (`external_required: false`) |
| Sessions | Create, get, list, delete, reprocess engine, export JSON/CSV |
| Estimates | Columns, rows (`calculation`), tables, schema totals (+ human), indexes |
| Growth | `PATCH .../tables/{name}/growth` |
| Alerts | Structural warnings (RF-18) and signals (RF-19) |
| Persistence | Optional SQLite settings + session history (`--persist`) |
| AI hook | `GET .../insights` with disabled provider (RF-24; no external calls) |
| Tests | Unit tests + [`tests/smoke.sh`](tests/smoke.sh) |

Functional RFs **RF-01…RF-24** for the initial local product are covered at the API/core level. Remaining product work is mainly richer UI and optional real AI providers.

### Still open (product polish)

| Item | Notes |
|---|---|
| Richer web UI | Stub exists at `GET /`; full paste/analyze UX still thin |
| Dedicated result routes | `/summary`, `/tables`, `/warnings` (today nested in session JSON) |
| File import | `POST /imports/sql` |
| Live AI provider | Extension interface exists; default is disabled/noop |
---

## Quick start

Requires Go (see `go.mod`).

```bash
go run ./cmd/tabyte serve
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

See [`docs/requisitos_funcionais.md`](docs/requisitos_funcionais.md) for the full RF list and first-delivery priority.
