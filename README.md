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

**Backend core is usable via API.** Sessions are in-memory. There is no web UI yet, no SQLite persistence, and no index/alert/export features.

Version reported by `/api/v1/info`: `0.0.1`.

### Done

| Area | What works |
|---|---|
| CLI | `tabyte serve` starts the local server |
| Localhost | Binds to `127.0.0.1:8787` by default |
| System API | `GET /api/v1/health`, `GET /api/v1/info` |
| Sessions | Create, get, list, delete analysis sessions |
| DDL input | DDL text on create; empty / no-table / bad-engine validation |
| Engines | `sqlserver` and `postgres` |
| Parser | `CREATE TABLE` (+ optional schema / `IF NOT EXISTS`); columns and types |
| Normalization | Engine-specific type mapping (length, precision, scale) |
| Column estimate | Per-column estimated bytes (common scalar / string / numeric types) |
| Row estimate | Row header + null bitmap + column payload; exposed as `calculation` |
| Table / schema volume | Default `assumed_row_count` (1000); schema `estimated_total_bytes` |
| Row count update | `PATCH /api/v1/analysis-sessions/{id}/tables/{tableName}` |
| Multi-table DDL | Multiple `CREATE TABLE` in one session |
| Offline | No external services required for the core path |
| Tests | Unit tests for parser/normalize/estimate; [`tests/smoke.sh`](tests/smoke.sh) |

Rough RF coverage already in place: **RF-02, RF-04–RF-12, RF-14, RF-15 (partial), RF-23**. RF-01 is partial (serve works; browser auto-open not yet).

### Not done yet

Prioritized for a usable first delivery (from the RF doc):

| Priority | Item | Notes |
|---|---|---|
| High | Web UI | Paste DDL, pick engine, show tables / breakdown / calculation |
| High | Open browser on serve | RF-01 |
| High | Structural alerts | RF-18 (wide columns, generic types, etc.) |
| High | Performance signals | RF-19 (structural only, not real workload claims) |
| Medium | Indexes in parser | PK / explicit indexes (RF-16) |
| Medium | Index storage estimate | Include in `calculation` (RF-17); currently omitted |
| Medium | Growth projection | Simple rate per hour/day/month (RF-13) |
| Medium | Reprocess params | Change engine / options without full app reload (RF-21 partial today) |
| Later | Export JSON/CSV | RF-20 |
| Later | Optional SQLite | Settings + history (RF-22); package stub only |
| Later | Dedicated result routes | `/summary`, `/tables`, `/warnings` (today nested in session JSON) |
| Later | File import | `POST /imports/sql` |
| Later | AI extension hook | RF-24 |
| Later | Richer parser | Nullability, unsupported-structure warnings without failing whole DDL |

Also still scaffold-only: `web/` (embed UI), `internal/persistence/sqlite`, `internal/platform` (browser/OS helpers).

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
internal/persistence/sqlite/  (planned) optional local store
internal/platform/          (planned) OS helpers (browser, paths)
web/                        (planned) embedded UI assets
docs/                       product and API docs
tests/smoke.sh              end-to-end API smoke checks
```

---

## Design notes

- Estimates are **not** physical measurements on a live instance (RN-02).
- Engine choice changes semantics for types, row overhead, and (later) indexes (RN-01).
- Index contribution to storage is intentionally out of the current `calculation` payload.

See [`docs/requisitos_funcionais.md`](docs/requisitos_funcionais.md) for the full RF list and first-delivery priority.
