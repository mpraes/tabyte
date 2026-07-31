# Tabyte — Storage Estimation Reference (PostgreSQL & SQL Server) v0.3

## Objective

Document how Tabyte estimates **column** and **row** storage for **PostgreSQL** and **SQL Server**. Estimates are **documented, reproducible, and engine-specific**. They are not exact physical measurements on a live instance.

Implementation lives in:

- `internal/engine/postgres/` (`EstimateColumn`, row helpers)
- `internal/engine/sqlserver/` (`EstimateColumn`, row helpers)
- Domain field `AssumedAvgLength` on `domain.Column` for variable-length types

## Principle

Schema metadata alone cannot yield absolute physical accuracy for variable types without assumptions about stored values. Tabyte therefore:

- applies fixed sizes where the engine documents them
- uses an **assumed average length** for variable types (not always the declared max)
- layers row overhead (header, null bitmap, index contribution) outside pure column sizing

## SQL Server — column rules

| Normalized type | Estimate |
|---|---|
| `smallint` | 2 |
| `int` | 4 |
| `bigint` | 8 |
| `tinyint` | 1 |
| `boolean` (maps from `bit`) | 1 per column (row-level bit packing is a later refinement) |
| `uuid` | 16 |
| `date` | **3** |
| `datetime` | 8 |
| `smalldatetime` | 4 |
| `datetime2` | 6 / 7 / 8 by fractional precision (scale) |
| `timestamp` | 8 (normalized fallback) |
| `float` | 8 (default double) |
| `char(n)` | `n` |
| `nchar(n)` | `n * 2` |
| `varchar` | `assumedAvg + 2` |
| `nvarchar` | `assumedAvg * 2 + 2` |
| `text` | `256 + 2` (assumed payload heuristic) |
| `numeric` / `decimal` | Official precision bands: 1–9→5, 10–19→9, 20–28→13, 29–38→17 |

### Assumed average length

For variable types, if `AssumedAvgLength` is unset:

- use `Length` when below the fallback (64)
- otherwise use `Length / 2`
- otherwise fallback `64`

## PostgreSQL — column rules

| Normalized type | Estimate |
|---|---|
| `smallint` | 2 |
| `int` | 4 |
| `bigint` | 8 |
| `boolean` | 1 |
| `uuid` | 16 |
| `date` | **4** |
| `timestamp` | 8 |
| `float` | 8 |
| `char` | varlena size of declared length |
| `varchar` | varlena size of assumed average |
| `text` | varlena size of assumed 256 |
| `numeric` | groups of 4 digits × 2 bytes + overhead (~4) |

### Varlena heuristic

```text
if dataBytes < 127 → dataBytes + 1
else → dataBytes + 4
```

## Row and table layer

Column bytes alone understate storage. Engines also model:

- row header / fixed overhead
- null bitmap
- index contribution
- table volume = estimated row size × assumed row count
- optional growth projection (UI/API) on top of current assumed volume

See RF-09–RF-13 and the `calculation` object returned on each table in the API/UI.

## Known limitations

- Variable types depend on assumed content length, not live data samples.
- SQL Server `bit` packing is approximate at column level.
- LOB / TOAST / out-of-row storage is not fully modeled.
- Page fill, alignment, and compression are out of scope for v0.

## Acceptance criteria (engine modules)

- SQL Server and PostgreSQL do **not** share one generic `numeric` / `varchar` heuristic.
- SQL Server `date` is 3 bytes; PostgreSQL `date` is 4 bytes.
- PostgreSQL does not treat `varchar(n)` as fixed `n` bytes.
- Code and docs distinguish column-level vs row-level contributions.
