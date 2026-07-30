package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/mpraes/tabyte/internal/application"
	"github.com/mpraes/tabyte/internal/domain"
)

type createSessionRequest struct {
	Engine     string `json:"engine"`
	SourceName string `json:"source_name"`
	DDLText    string `json:"ddl_text"`
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
			"id":          session.ID,
			"engine":      session.Engine,
			"status":      session.Status,
			"table_count": len(session.Tables),
			"tables":      tablesJSON(session.Tables),
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
			"tables":      tablesJSON(session.Tables),
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

func tablesJSON(tables []domain.Table) []map[string]any {
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
		}
		if t.EstimatedRowBytes != nil {
			item["estimated_row_bytes"] = *t.EstimatedRowBytes
		}
		out = append(out, item)
	}
	return out
}