package httpapi

import (
	"errors"
	"net/http"

	"github.com/mpraes/tabyte/internal/application"
	"github.com/mpraes/tabyte/internal/domain"
)

func HandleListInsights(store *application.SessionStore, provider application.InsightProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.PathValue("sessionId")
		insights, enabled, err := application.ListInsights(store, sessionID, provider)
		if err != nil {
			if errors.Is(err, application.ErrSessionNotFound) {
				WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
				return
			}
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		WriteJSON(w, http.StatusOK, map[string]any{
			"session_id": sessionID,
			"enabled":    enabled,
			"insights":   insightsJSON(insights),
			"count":      len(insights),
		})
	}
}

func insightsJSON(insights []domain.Insight) []map[string]any {
	out := make([]map[string]any, 0, len(insights))
	for _, in := range insights {
		out = append(out, map[string]any{
			"provider": in.Provider,
			"category": in.Category,
			"text":     in.Text,
		})
	}
	return out
}
