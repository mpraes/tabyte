# Tabyte — Local HTTP API v0.3

## Objective

Contract for the **local HTTP API** served by `tabyte serve`. Base URL:

```text
http://127.0.0.1:8787/api/v1
```

Default bind is loopback only. UI is served outside the API prefix (`GET /`, `GET /assets/…`).

## Conventions

### Format

- Request/response bodies use `application/json` unless noted (export CSV is an exception).
- Consistent envelope on JSON responses.

### Success envelope

```json
{
  "data": {},
  "meta": {
    "request_id": "req_local",
    "timestamp": "2026-07-30T11:30:00Z"
  },
  "error": null
}
```

### Error envelope

```json
{
  "data": null,
  "meta": {
    "request_id": "req_local",
    "timestamp": "2026-07-30T11:30:00Z"
  },
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "invalid JSON body",
    "details": []
  }
}
```

### Status codes in use

| Situation | Status |
|---|---|
| Successful read / update | `200 OK` |
| Session created | `201 Created` |
| Delete with no body | `204 No Content` |
| Invalid payload | `400 Bad Request` |
| Missing resource / persistence off | `404 Not Found` |
| Unexpected failure | `500 Internal Server Error` |

## Implemented routes

### System

#### `GET /api/v1/health`

Basic process health.

```json
{
  "data": {
    "status": "ok",
    "app": "tabyte",
    "version": "0.0.1"
  }
}
```

#### `GET /api/v1/info`

Instance capabilities.

```json
{
  "data": {
    "app": "tabyte",
    "version": "0.0.1",
    "mode": "local",
    "bind": "127.0.0.1:8787",
    "persistence": false,
    "external_required": false,
    "ai_insights": false,
    "engines": ["sqlserver", "postgres"]
  }
}
```

### Settings (requires `--persist`)

Without persistence, both return `404` with message `persistence is not enabled`.

#### `GET /api/v1/settings`

```json
{
  "data": {
    "settings": [
      { "key": "default_engine", "value": "postgres", "value_type": "string" }
    ]
  }
}
```

#### `PUT /api/v1/settings`

Upsert one setting.

```json
{
  "key": "default_engine",
  "value": "postgres",
  "value_type": "string"
}
```

### Analysis sessions

#### `POST /api/v1/analysis-sessions`

Create and analyze.

Request:

```json
{
  "engine": "postgres",
  "source_name": "paste.sql",
  "ddl_text": "CREATE TABLE users (id INT, name VARCHAR(100));"
}
```

`engine` values: `postgres` | `sqlserver`.

Response `201` includes session id, engine, status, `tables` (with columns, indexes, calculation, growth fields when set), totals, warnings, and signals.

#### `GET /api/v1/analysis-sessions`

List sessions (summary items suitable for the sidebar).

#### `GET /api/v1/analysis-sessions/{sessionId}`

Full session: `ddl_text`, `tables` with calculation detail, warnings, signals, human-readable totals when present.

#### `DELETE /api/v1/analysis-sessions/{sessionId}`

Removes the session (`204`).

#### `PATCH /api/v1/analysis-sessions/{sessionId}`

Reprocess with another engine.

```json
{ "engine": "sqlserver" }
```

#### `PATCH /api/v1/analysis-sessions/{sessionId}/tables/{tableName}`

Update assumed row count and recalculate volumes.

```json
{ "assumed_row_count": 5000 }
```

#### `PATCH /api/v1/analysis-sessions/{sessionId}/tables/{tableName}/growth`

Set simple growth projection.

```json
{
  "rows_per_period": 100,
  "period": "day",
  "horizon": 30
}
```

`period`: `hour` | `day` | `month`.

#### `GET /api/v1/analysis-sessions/{sessionId}/export?format=json|csv`

Downloads analysis export. Default format is `json`. CSV returns a downloadable CSV body (not the JSON envelope).

#### `GET /api/v1/analysis-sessions/{sessionId}/insights`

AI extension hook. Default provider is disabled:

```json
{
  "data": {
    "enabled": false,
    "insights": []
  }
}
```

### UI (non-API)

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/` | Embedded web UI |
| `GET` | `/assets/*` | CSS/JS assets |

## Not implemented (deferred)

These appeared in earlier drafts and remain out of the current mux:

- `GET /analysis-sessions/{id}/summary`
- `GET /analysis-sessions/{id}/tables` (and nested column/index/warning routes)
- `POST /imports/sql`
- `PATCH /settings` (bulk)
- `PUT /settings/{key}`

Table/column/warning detail is returned nested in session JSON instead.

## Session result shape (high level)

Create/get `data` typically includes:

- `id`, `engine`, `status`, `source_name`, `ddl_text` (get)
- `estimated_total_bytes`, optional `estimated_total_human`
- optional `projected_total_bytes` / `projected_total_human`
- `tables[]` with `name`, sizes, `assumed_row_count`, `calculation`, `columns`, `indexes`, growth fields
- `warnings[]`, `warning_count`, `signals[]`, `signal_count`
