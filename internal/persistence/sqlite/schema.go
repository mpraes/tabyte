package sqlite

import "database/sql"

const schemaSQL = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS app_settings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	key TEXT NOT NULL UNIQUE,
	value TEXT NOT NULL,
	value_type TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS analysis_sessions (
	id TEXT PRIMARY KEY,
	engine TEXT NOT NULL,
	source_name TEXT,
	ddl_text TEXT NOT NULL,
	ddl_hash TEXT NOT NULL,
	parser_status TEXT NOT NULL,
	total_tables INTEGER NOT NULL,
	warning_count INTEGER NOT NULL,
	error_count INTEGER NOT NULL,
	estimated_total_bytes INTEGER,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS analysis_tables (
	id TEXT PRIMARY KEY,
	session_id TEXT NOT NULL REFERENCES analysis_sessions(id) ON DELETE CASCADE,
	schema_name TEXT,
	table_name TEXT NOT NULL,
	row_size_bytes INTEGER,
	row_count_assumed INTEGER NOT NULL,
	growth_rate_value REAL,
	growth_rate_unit TEXT,
	estimated_table_bytes INTEGER,
	index_bytes INTEGER,
	warning_count INTEGER NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_analysis_tables_session ON analysis_tables(session_id);

CREATE TABLE IF NOT EXISTS analysis_warnings (
	id TEXT PRIMARY KEY,
	session_id TEXT NOT NULL REFERENCES analysis_sessions(id) ON DELETE CASCADE,
	table_id TEXT,
	column_id TEXT,
	code TEXT NOT NULL,
	severity TEXT NOT NULL,
	category TEXT NOT NULL,
	title TEXT NOT NULL,
	message TEXT NOT NULL,
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_analysis_warnings_session ON analysis_warnings(session_id);
`

func Migrate(db *sql.DB) error {
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}
	_, err := db.Exec(schemaSQL)
	return err
}
