package httpapi

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mpraes/tabyte/internal/application"
	"github.com/mpraes/tabyte/internal/domain"
)

func HandleExportAnalysisSession(store *application.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("sessionId")
		session, ok := application.GetSession(store, id)
		if !ok {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "session not found")
			return
		}

		format := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("format")))
		if format == "" {
			format = "json"
		}

		switch format {
		case "json":
			writeExportJSON(w, session)
		case "csv":
			writeExportCSV(w, session)
		default:
			WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "format must be json or csv")
		}
	}
}

func writeExportJSON(w http.ResponseWriter, session domain.AnalysisSession) {
	payload := map[string]any{
		"id":                    session.ID,
		"engine":                session.Engine,
		"source_name":           session.SourceName,
		"status":                session.Status,
		"tables":                tablesJSONWithCalculation(session.Tables),
		"estimated_total_bytes": session.EstimatedTotalBytes,
		"projected_total_bytes": session.ProjectedTotalBytes,
		"warnings":              warningsJSON(session.Warnings),
		"signals":               signalsJSON(session.Signals),
	}
	for k, v := range humanTotalBytes(session) {
		payload[k] = v
	}

	filename := fmt.Sprintf("tabyte-%s.json", session.ID)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeExportCSV(w http.ResponseWriter, session domain.AnalysisSession) {
	filename := fmt.Sprintf("tabyte-%s.csv", session.ID)
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)

	cw := csv.NewWriter(w)
	_ = cw.Write([]string{
		"session_id", "engine", "table", "column_count",
		"estimated_row_bytes", "assumed_row_count", "estimated_table_bytes",
		"index_bytes", "estimated_total_bytes", "estimated_total_human",
	})

	total := ""
	human := ""
	if session.EstimatedTotalBytes != nil {
		total = fmt.Sprintf("%d", *session.EstimatedTotalBytes)
		human = application.FormatBytes(*session.EstimatedTotalBytes)
	}

	for _, t := range session.Tables {
		rowBytes, tableBytes, indexBytes := "", "", "0"
		if t.EstimatedRowBytes != nil {
			rowBytes = fmt.Sprintf("%d", *t.EstimatedRowBytes)
		}
		if t.EstimatedTableBytes != nil {
			tableBytes = fmt.Sprintf("%d", *t.EstimatedTableBytes)
		}
		if t.Calculation != nil {
			indexBytes = fmt.Sprintf("%d", t.Calculation.IndexBytes)
		}
		_ = cw.Write([]string{
			session.ID,
			string(session.Engine),
			t.Name,
			fmt.Sprintf("%d", len(t.Columns)),
			rowBytes,
			fmt.Sprintf("%d", t.AssumedRowCount),
			tableBytes,
			indexBytes,
			total,
			human,
		})
	}
	cw.Flush()
}