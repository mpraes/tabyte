package sqlite_test

import (
	"path/filepath"
	"testing"

	"github.com/mpraes/tabyte/internal/application"
	"github.com/mpraes/tabyte/internal/domain"
	"github.com/mpraes/tabyte/internal/persistence/sqlite"
)

func TestSettingsAndSessionRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tabyte.db")

	db, err := sqlite.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	if err := db.UpsertSetting("theme", "dark", "string"); err != nil {
		t.Fatalf("upsert setting: %v", err)
	}
	settings, err := db.ListSettings()
	if err != nil {
		t.Fatalf("list settings: %v", err)
	}
	if len(settings) != 1 || settings[0].Key != "theme" || settings[0].Value != "dark" {
		t.Fatalf("unexpected settings: %+v", settings)
	}

	row := int64(28)
	tableBytes := int64(28000)
	total := int64(28000)
	session := domain.AnalysisSession{
		ID:         "as_test",
		Engine:     domain.EnginePostgres,
		SourceName: "a.sql",
		DDLText:    "CREATE TABLE a (id INT);",
		Status:     "created",
		Tables: []domain.Table{{
			Name:                "a",
			AssumedRowCount:     1000,
			EstimatedRowBytes:   &row,
			EstimatedTableBytes: &tableBytes,
			Calculation:         &domain.RowCalculation{EstimatedRowBytes: 28},
		}},
		EstimatedTotalBytes: &total,
		Warnings: []domain.Warning{{
			Code:    "WIDE_VARCHAR",
			Message: "wide",
			Table:   "a",
			Column:  "name",
		}},
		Signals: []domain.Signal{{
			Code:    "WIDE_ROW",
			Message: "wide row",
			Table:   "a",
		}},
	}
	if err := db.UpsertSession(session); err != nil {
		t.Fatalf("upsert session: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	db2, err := sqlite.Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer db2.Close()

	loaded, err := db2.LoadAll()
	if err != nil {
		t.Fatalf("load all: %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("want 1 session, got %d", len(loaded))
	}
	got := loaded[0]
	if got.ID != "as_test" || got.Engine != "postgres" || got.DDLText == "" {
		t.Fatalf("unexpected session: %+v", got)
	}
	if len(got.Tables) != 1 || got.Tables[0].Name != "a" || got.Tables[0].AssumedRowCount != 1000 {
		t.Fatalf("unexpected tables: %+v", got.Tables)
	}

	rebuilt, err := application.RebuildSession(got)
	if err != nil {
		t.Fatalf("rebuild: %v", err)
	}
	if rebuilt.EstimatedTotalBytes == nil || *rebuilt.EstimatedTotalBytes != 28000 {
		t.Fatalf("unexpected rebuilt total: %+v", rebuilt.EstimatedTotalBytes)
	}
	if len(rebuilt.Tables) != 1 || rebuilt.Tables[0].AssumedRowCount != 1000 {
		t.Fatalf("unexpected rebuilt tables: %+v", rebuilt.Tables)
	}

	if err := db2.DeleteSession("as_test"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	loaded, err = db2.LoadAll()
	if err != nil {
		t.Fatalf("load after delete: %v", err)
	}
	if len(loaded) != 0 {
		t.Fatalf("want 0 sessions after delete, got %d", len(loaded))
	}
}

func TestHydratePreservesGrowth(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tabyte.db")
	db, err := sqlite.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	row := int64(28)
	tableBytes := int64(140000)
	total := int64(140000)
	projRows := int64(8000)
	projBytes := int64(224000)
	session := domain.AnalysisSession{
		ID:         "as_growth",
		Engine:     domain.EnginePostgres,
		SourceName: "g.sql",
		DDLText:    "CREATE TABLE a (id INT);",
		Status:     "created",
		Tables: []domain.Table{{
			Name:                "a",
			AssumedRowCount:     5000,
			EstimatedRowBytes:   &row,
			EstimatedTableBytes: &tableBytes,
			GrowthRowsPerPeriod: 100,
			GrowthPeriod:        "day",
			GrowthHorizon:       30,
			ProjectedRowCount:   &projRows,
			ProjectedTableBytes: &projBytes,
			Calculation:         &domain.RowCalculation{EstimatedRowBytes: 28},
		}},
		EstimatedTotalBytes: &total,
		ProjectedTotalBytes: &projBytes,
	}
	if err := db.UpsertSession(session); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	loaded, err := db.LoadAll()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	rebuilt, err := application.RebuildSession(loaded[0])
	if err != nil {
		t.Fatalf("rebuild: %v", err)
	}
	t0 := rebuilt.Tables[0]
	if t0.AssumedRowCount != 5000 || t0.GrowthRowsPerPeriod != 100 || t0.GrowthPeriod != "day" || t0.GrowthHorizon != 30 {
		t.Fatalf("growth not restored: %+v", t0)
	}
	if t0.ProjectedRowCount == nil || *t0.ProjectedRowCount != 8000 {
		t.Fatalf("projected rows: %+v", t0.ProjectedRowCount)
	}
}
