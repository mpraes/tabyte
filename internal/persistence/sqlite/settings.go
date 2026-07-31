package sqlite

import (
	"time"

	"github.com/mpraes/tabyte/internal/application"
)

func (d *DB) ListSettings() ([]application.Setting, error) {
	rows, err := d.sql.Query(`SELECT key, value, value_type FROM app_settings ORDER BY key`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []application.Setting
	for rows.Next() {
		var s application.Setting
		if err := rows.Scan(&s.Key, &s.Value, &s.ValueType); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (d *DB) UpsertSetting(key, value, valueType string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := d.sql.Exec(`
		INSERT INTO app_settings (key, value, value_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			value_type = excluded.value_type,
			updated_at = excluded.updated_at
	`, key, value, valueType, now, now)
	return err
}
