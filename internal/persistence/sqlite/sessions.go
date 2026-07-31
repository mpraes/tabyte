package sqlite

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/mpraes/tabyte/internal/application"
	"github.com/mpraes/tabyte/internal/domain"
)

func (d *DB) UpsertSession(session domain.AnalysisSession) error {
	tx, err := d.sql.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	now := time.Now().UTC().Format(time.RFC3339)
	hash := ddlHash(session.DDLText)
	var total any
	if session.EstimatedTotalBytes != nil {
		total = *session.EstimatedTotalBytes
	}

	_, err = tx.Exec(`
		INSERT INTO analysis_sessions (
			id, engine, source_name, ddl_text, ddl_hash, parser_status,
			total_tables, warning_count, error_count, estimated_total_bytes,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			engine = excluded.engine,
			source_name = excluded.source_name,
			ddl_text = excluded.ddl_text,
			ddl_hash = excluded.ddl_hash,
			parser_status = excluded.parser_status,
			total_tables = excluded.total_tables,
			warning_count = excluded.warning_count,
			error_count = excluded.error_count,
			estimated_total_bytes = excluded.estimated_total_bytes,
			updated_at = excluded.updated_at
	`, session.ID, string(session.Engine), session.SourceName, session.DDLText, hash, session.Status,
		len(session.Tables), len(session.Warnings)+len(session.Signals), total, now, now)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM analysis_tables WHERE session_id = ?`, session.ID); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM analysis_warnings WHERE session_id = ?`, session.ID); err != nil {
		return err
	}

	warnByTable := map[string]int{}
	for _, w := range session.Warnings {
		warnByTable[strings.ToLower(w.Table)]++
	}
	for _, s := range session.Signals {
		warnByTable[strings.ToLower(s.Table)]++
	}

	for _, t := range session.Tables {
		tableID := session.ID + ":" + t.Name
		var rowSize, tableBytes, indexBytes any
		if t.EstimatedRowBytes != nil {
			rowSize = *t.EstimatedRowBytes
		}
		if t.EstimatedTableBytes != nil {
			tableBytes = *t.EstimatedTableBytes
		}
		if t.Calculation != nil {
			indexBytes = t.Calculation.IndexBytes
		}
		var growthValue any
		var growthUnit any
		if t.GrowthRowsPerPeriod > 0 {
			growthValue = float64(t.GrowthRowsPerPeriod)
			growthUnit = fmt.Sprintf("%s:%d", t.GrowthPeriod, t.GrowthHorizon)
		}
		_, err = tx.Exec(`
			INSERT INTO analysis_tables (
				id, session_id, schema_name, table_name, row_size_bytes, row_count_assumed,
				growth_rate_value, growth_rate_unit, estimated_table_bytes, index_bytes,
				warning_count, created_at, updated_at
			) VALUES (?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, tableID, session.ID, t.Name, rowSize, t.AssumedRowCount,
			growthValue, growthUnit, tableBytes, indexBytes,
			warnByTable[strings.ToLower(t.Name)], now, now)
		if err != nil {
			return err
		}
	}

	n := 0
	for _, w := range session.Warnings {
		n++
		id := fmt.Sprintf("%s:w:%d", session.ID, n)
		tableID := nullableTableID(session.ID, w.Table)
		_, err = tx.Exec(`
			INSERT INTO analysis_warnings (
				id, session_id, table_id, column_id, code, severity, category, title, message, created_at
			) VALUES (?, ?, ?, NULL, ?, 'warning', 'warning', ?, ?, ?)
		`, id, session.ID, tableID, w.Code, w.Code, w.Message, now)
		if err != nil {
			return err
		}
	}
	for _, s := range session.Signals {
		n++
		id := fmt.Sprintf("%s:s:%d", session.ID, n)
		tableID := nullableTableID(session.ID, s.Table)
		_, err = tx.Exec(`
			INSERT INTO analysis_warnings (
				id, session_id, table_id, column_id, code, severity, category, title, message, created_at
			) VALUES (?, ?, ?, NULL, ?, 'info', 'signal', ?, ?, ?)
		`, id, session.ID, tableID, s.Code, s.Code, s.Message, now)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *DB) DeleteSession(id string) error {
	_, err := d.sql.Exec(`DELETE FROM analysis_sessions WHERE id = ?`, id)
	return err
}

func (d *DB) LoadAll() ([]application.PersistedSession, error) {
	rows, err := d.sql.Query(`
		SELECT id, engine, source_name, ddl_text, parser_status
		FROM analysis_sessions
		ORDER BY created_at
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []application.PersistedSession
	for rows.Next() {
		var s application.PersistedSession
		var source sql.NullString
		if err := rows.Scan(&s.ID, &s.Engine, &source, &s.DDLText, &s.Status); err != nil {
			return nil, err
		}
		s.SourceName = source.String
		sessions = append(sessions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range sessions {
		tables, err := d.loadTables(sessions[i].ID)
		if err != nil {
			return nil, err
		}
		sessions[i].Tables = tables
	}
	return sessions, nil
}

func (d *DB) loadTables(sessionID string) ([]application.PersistedTable, error) {
	rows, err := d.sql.Query(`
		SELECT table_name, row_count_assumed, growth_rate_value, growth_rate_unit
		FROM analysis_tables
		WHERE session_id = ?
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []application.PersistedTable
	for rows.Next() {
		var t application.PersistedTable
		var growthVal sql.NullFloat64
		var growthUnit sql.NullString
		if err := rows.Scan(&t.Name, &t.AssumedRowCount, &growthVal, &growthUnit); err != nil {
			return nil, err
		}
		if growthVal.Valid && growthUnit.Valid {
			t.GrowthRowsPerPeriod = int64(growthVal.Float64)
			period, horizon := parseGrowthUnit(growthUnit.String)
			t.GrowthPeriod = period
			t.GrowthHorizon = horizon
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func ddlHash(ddl string) string {
	sum := sha256.Sum256([]byte(ddl))
	return hex.EncodeToString(sum[:])
}

func nullableTableID(sessionID, tableName string) any {
	if strings.TrimSpace(tableName) == "" {
		return nil
	}
	return sessionID + ":" + tableName
}

func parseGrowthUnit(s string) (period string, horizon int64) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return s, 0
	}
	period = parts[0]
	fmt.Sscanf(parts[1], "%d", &horizon)
	return period, horizon
}
