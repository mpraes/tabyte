# Tabyte — Functional Requirements v0.3

## Product vision

**Tabyte** is an open-source tool for storage estimation and structural performance-impact analysis from relational DDLs. The first supported engines are **SQL Server** and **PostgreSQL**, with engine-specific rules for data types, table structure, and indexing.

The product runs locally through a **cross-platform CLI** that starts a local HTTP server and opens the embedded web UI in the default browser. The same flow works on Windows PowerShell and Ubuntu/Linux, with a local-first distribution model.

## Functional objective

The system accepts DDLs, interprets tables and relevant schema objects, applies estimation rules for the selected engine, and presents understandable results. Estimates compose base table structure with index overhead where applicable; they are not live measurements from a database instance.

## First-version scope

The product covers:

- DDL analysis by pasted text
- Engine selection (SQL Server / PostgreSQL)
- Estimated size per column, row, table, and schema
- Basic structural warnings and performance signals
- Simple growth projection
- Optional local SQLite persistence for history and settings
- Embedded web UI for the full analyze workflow

It does **not** require an external database server. SQLite is used only as an optional internal store when persistence is enabled (`tabyte serve --persist`).

Out of scope for this version: real benchmarks, execution-plan reading, mandatory remote DB connections, engines beyond SQL Server and PostgreSQL, advanced automated tuning, and SQL file upload.

## Actors

### Analyst user

Pastes or edits DDL and wants estimated storage impact plus structural performance signals.

### Architect / DBA user

Compares schema decisions, reviews types, estimates growth, and evaluates index, row-width, and modeling trade-offs.

### CLI user

Starts the app from a terminal (`tabyte serve`), gets a local URL, and uses the UI in the browser on Windows or Linux with the same flow.

## Functional requirements

### RF-01 — Start via CLI

The system provides a primary command such as `tabyte serve`. On success it prints the active local URL and, unless `--no-open` is set, opens the default browser.

### RF-02 — Operate on localhost

The HTTP server binds by default to `127.0.0.1` / localhost only, avoiding unnecessary external network exposure.

### RF-03 — Provide DDL as text

The user can paste or edit DDL in a text area for analysis. Content stays in the current session until replaced. File upload is out of scope for this delivery; `source_name` remains optional session metadata (API/persistence), not a required UI field. The UI currently sends a fixed placeholder such as `paste.sql`.

### RF-04 — Select database engine

The user explicitly chooses **SQL Server** or **PostgreSQL**. Parsing, type mapping, and estimation follow the selected engine.

### RF-05 — Validate minimal input

The system validates that input is sufficient before processing. Empty, incomplete, or semantically insufficient DDL yields a clear error message.

### RF-06 — Interpret table structures

The system identifies at least tables, columns, data types, nullability, primary keys, and indexes supported by the initial parser. Unsupported constructs are signaled without aborting the whole analysis when possible.

### RF-07 — Normalize types per engine

Textual types from the DDL are converted to an internal normalized model, preserving length, precision, scale, and other attributes needed for estimation.

### RF-08 — Estimate size per column

The system estimates storage per column from the normalized type, parameters, and engine rules, shown in estimated bytes (or a derived unit).

### RF-09 — Estimate size per row

The system estimates row size including column payload and modeled engine overhead (header, null bitmap, and related factors). For SQL Server, base structure is estimated before adding nonclustered index contribution.

### RF-10 — Estimate volume per table

Table volume is row size × assumed row count. A configurable default cardinality applies when the user has not set a value.

### RF-11 — Estimate total schema volume

The system consolidates estimated schema volume across identified tables and engine-considered components, shown in a readable unit (KB, MB, GB).

### RF-12 — Set row count per table

The user can set or change assumed rows per table. Recalculation happens via explicit apply (UI) / `PATCH` on the table resource (API).

### RF-13 — Project growth

The user can project growth per table with a simple rate (rows per hour/day/month and a horizon). This supports early volumetry scenarios, not advanced temporal modeling.

### RF-14 — Show per-table detail

For each table the UI/API presents name, interpreted column count, estimated bytes per row, assumed rows, estimated volume, columns with individual sizes, indexes, and calculation breakdown. API v0 embeds this in create/get session `tables[]`; dedicated `/tables` collection routes are deferred.

### RF-15 — Show calculation memory

The system shows how the result was obtained: column payload, row header, null bitmap, and index contribution when applicable.

### RF-16 — Interpret supported indexes

Within the initial version, the system recognizes indexes defined in the DDL or enough structure to infer basics (e.g. primary keys and explicit indexes).

### RF-17 — Estimate structural index impact

The system provides a basic estimate of index storage impact. SQL Server separates base structure from nonclustered indexes; PostgreSQL treats indexes as additional cost tied to read performance.

### RF-18 — Emit structural modeling warnings

Rule-based warnings cover cases such as potentially wide columns, overly generic types, unnecessary precision, excess variable columns, and structures prone to volume or indexing cost growth.

### RF-19 — Emit performance-impact signals

Structural signals that may affect read, write, ordering, or index maintenance are reported without claiming real workload performance.

### RF-20 — Export results

Results can be exported as structured JSON and CSV for documentation, comparison, or external automation (UI download links and `GET .../export`).

### RF-21 — Reprocess after parameter changes

The user can change engine, row volume, growth, or other supported parameters and recalculate without reloading the whole app (`PATCH` session for engine; table PATCHes for rows/growth).

### RF-22 — Optional local settings persistence

Optional local persistence of settings and analysis history via embedded SQLite. Persistence is not required for basic operation; enable with `--persist`.

### RF-23 — No mandatory external integration

The product is fully usable without external APIs. Parsing, estimation, and rule-based alerts run locally (`external_required: false` on `/info`).

### RF-24 — Extension point for optional AI

An extension point exists for future textual AI insights (`GET .../insights`). Insights annotate analysis results and never replace the calculation engine. The default provider is disabled (no external calls).

## Web UI (current product surface)

The embedded UI at `GET /` provides:

| Area | Behavior |
|---|---|
| Workspace | Brand, engine selector, DDL textarea, **Analyze** |
| Sessions | List, open, delete, refresh (in-memory or SQLite when `--persist`) |
| Settings | View/upsert key-value settings when persistence is on; guidance when off |
| Results | Schema totals, per-table cards (rows, calculation, columns, indexes, growth), warnings, signals, insights stub, JSON/CSV export |

## Business rules

### RN-01 — Engine defines estimation semantics

Every analysis follows the selected engine. The same DDL text may produce different results under SQL Server vs PostgreSQL.

### RN-02 — Calculation is not a live measurement

Results are estimates, not exact physical storage on a real instance. The UI states this explicitly.

### RN-03 — Indexes are benefit and cost

When indexes are considered, the system communicates both potential read/search/order benefits and extra storage/maintenance cost.

### RN-04 — Partial failures do not void the whole result

If part of the DDL cannot be interpreted, the system analyzes the rest when global consistency allows, and surfaces limitations.

## Main use cases

### UC-01 — Start locally

User runs `tabyte serve`. System starts the local server, prints the URL, and opens the UI.

### UC-02 — Analyze a table DDL

User pastes DDL, selects engine, clicks Analyze. System interprets structure, estimates sizes, and shows table detail.

### UC-03 — Analyze multiple tables

User provides multiple `CREATE TABLE` statements. System identifies tables individually and shows per-table and schema totals.

### UC-04 — Project volumetry

User sets assumed rows and growth parameters. System recalculates totals and projections.

### UC-05 — Review structural alerts

User reviews warnings and performance signals to see which schema decisions deserve attention.

### UC-06 — Export result

After analysis, user exports JSON or CSV for docs or comparison.

### UC-07 — Reopen a session

User opens a prior session from the sidebar; DDL, engine, and results are restored in the workspace.

## Priority for first delivery

Critical requirements for the first delivery: RF-01–RF-12, RF-14, RF-15, RF-18, RF-19, RF-20, RF-23, plus the embedded UI surface above. Optional persistence and AI can land incrementally as long as the release already delivers explainable calculation, simple local operation, and defensible structural alerts.
