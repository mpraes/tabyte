package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/mpraes/tabyte/internal/application"
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

		WriteJSON(w, http.StatusCreated, map[string]any{
			"id":     session.ID,
			"engine": session.Engine,
			"status": session.Status,
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
		WriteJSON(w, http.StatusOK, map[string]any{
			"id":          session.ID,
			"engine":      session.Engine,
			"source_name": session.SourceName,
			"status":      session.Status,
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