package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mpraes/tabyte/internal/application"
	"github.com/mpraes/tabyte/internal/domain"
)

type createSessionRequest struct {
	Engine     string `json:"engine"`
	SourceName string `json:"source_name"`
	DDLText    string `json:"ddl_text"`
}

type updateTableRowCountRequest struct {
	AssumedRowCount int64 `json:"assumed_row_count"`
}



func HandleCreateAnalysisSession(store *application.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createSessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON body")
			return
		}

		session, err := application.CreateSession(store, application.CreateSessionInput{
			Engine:     req.Engine,
			SourceName: req.SourceName,
			DDLText:    req.DDLText,
		})
		if err != nil {
			WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}

		tables := make([]map[string]any, 0, len(session.Tables))
		for _, t := range session.Tables {
			cols := make([]map[string]any, 0, len(t.Columns))
			for _, c := range t.Columns {
				cols = append(cols, map[string]any{
					"name":          c.Name,
					"original_type": c.OriginalType,
				})
			}
			tables = append(tables, map[string]any{
				"name":         t.Name,
				"column_count": len(t.Columns),
				"columns":      cols,
			})
		}
		WriteJSON(w, http.StatusCreated, map[string]any{
			"id":                    session.ID,
			"engine":                session.Engine,
			"status":                session.Status,
			"table_count":           len(session.Tables),
			"tables":                tablesJSONWithCalculation(session.Tables),
			"estimated_total_bytes": session.EstimatedTotalBytes,
			"warnings":              warningsJSON(session.Warnings),
			"warning_count":         len(session.Warnings),
		})
	}
}

func HandleGetAnalysisSession(store *application.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("sessionId")
		session, ok := application.GetSession(store, id)
		if !ok {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "session not found")
			return
		}
		tables := make([]map[string]any, 0, len(session.Tables))
		for _, t := range session.Tables {
			cols := make([]map[string]any, 0, len(t.Columns))
			for _, c := range t.Columns {
				cols = append(cols, map[string]any{
					"name":          c.Name,
					"original_type": c.OriginalType,
				})
			}
			tables = append(tables, map[string]any{
				"name":         t.Name,
				"column_count": len(t.Columns),
				"columns":      cols,
			})
		}
		WriteJSON(w, http.StatusOK, map[string]any{
			"id":          session.ID,
			"engine":      session.Engine,
			"source_name": session.SourceName,
			"status":      session.Status,
			"tables":      tablesJSONWithCalculation(session.Tables),
			"warnings":    warningsJSON(session.Warnings),
			"warning_count":  len(session.Warnings),
			"estimated_total_bytes": session.EstimatedTotalBytes,
		})
	}
}

func HandleListAnalysisSessions(store *application.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions := application.ListSessions(store)
		items := make([]map[string]any, 0, len(sessions))
		for _, s := range sessions {
			items = append(items, map[string]any{
				"id":          s.ID,
				"engine":      s.Engine,
				"source_name": s.SourceName,
				"status":      s.Status,
				"table_count": len(s.Tables),
			})
		}
		WriteJSON(w, http.StatusOK, items)
	}
}

func HandleDeleteAnalysisSession(store *application.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("sessionId")
		if !application.DeleteSession(store, id) {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "session not found")
			return
		}
		w.WriteHeader(http.StatusNoContent) // 204, no body
	}
}

func tablesJSONWithCalculation(tables []domain.Table) []map[string]any {
	out := make([]map[string]any, 0, len(tables))
	for _, t := range tables {
		cols := make([]map[string]any, 0, len(t.Columns))
		for _, c := range t.Columns {
			m := map[string]any{
				"name":            c.Name,
				"original_type":   c.OriginalType,
				"normalized_type": c.NormalizedType,
			}
			if c.Length != nil {
				m["length"] = *c.Length
			}
			if c.Precision != nil {
				m["precision"] = *c.Precision
			}
			if c.Scale != nil {
				m["scale"] = *c.Scale
			}
			if c.AssumedAvgLength != nil {
				m["assumed_avg_length"] = *c.AssumedAvgLength
			}
			if c.EstimatedBytes != nil {
				m["estimated_bytes"] = *c.EstimatedBytes
			}
			cols = append(cols, m)
		}

		item := map[string]any{
			"name":         t.Name,
			"column_count": len(t.Columns),
			"columns":      cols,
			"assumed_row_count": t.AssumedRowCount,
		}
		if t.EstimatedRowBytes != nil {
			item["estimated_row_bytes"] = *t.EstimatedRowBytes
		}
		if t.EstimatedTableBytes != nil {
			item["estimated_table_bytes"] = *t.EstimatedTableBytes
		}
		if t.Calculation != nil {
			item["calculation"] = map[string]any{
				"column_payload_bytes": t.Calculation.ColumnPayloadBytes,
				"row_header_bytes":     t.Calculation.RowHeaderBytes,
				"null_bitmap_bytes":    t.Calculation.NullBitmapBytes,
				"estimated_row_bytes":  t.Calculation.EstimatedRowBytes,
				// omit index_bytes for now
			}
		}
		out = append(out, item)
	}
	return out
}

func warningsJSON(warnings []domain.Warning) []map[string]any {
	out := make([]map[string]any, 0, len(warnings))
	for _, w := range warnings {
		item := map[string]any{
			"code":    w.Code,
			"message": w.Message,
			"table":   w.Table,
		}
		if w.Column != "" {
			item["column"] = w.Column
		}
		out = append(out, item)
	}
	return out
}

func HandleUpdateTableRowCount(store *application.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.PathValue("sessionId")
		tableName := r.PathValue("tableName")

		var req updateTableRowCountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON body")
			return
		}

		session, err := application.UpdateTableRowCount(store, sessionID, tableName, req.AssumedRowCount)
		if err != nil {
			switch {
			case errors.Is(err, application.ErrSessionNotFound),
				errors.Is(err, application.ErrTableNotFound):
				WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			default:
				WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			}
			return
		}

		WriteJSON(w, http.StatusOK, map[string]any{
			"id":          session.ID,
			"engine":      session.Engine,
			"status":      session.Status,
			"tables":      tablesJSONWithCalculation(session.Tables),
			"warnings":    warningsJSON(session.Warnings),
			"warning_count":  len(session.Warnings),
			"estimated_total_bytes": session.EstimatedTotalBytes,
		})
	}
}