package httpapi

import (
	"net/http"

	"github.com/mpraes/tabyte/internal/application"
)

func NewMux(store *application.SessionStore) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", HandleHealth)
	mux.HandleFunc("GET /api/v1/info", HandleInfo)
	mux.HandleFunc("POST /api/v1/analysis-sessions", HandleCreateAnalysisSession(store))
	mux.HandleFunc("GET /api/v1/analysis-sessions/{sessionId}", HandleGetAnalysisSession(store))
	mux.HandleFunc("GET /api/v1/analysis-sessions", HandleListAnalysisSessions(store))
	mux.HandleFunc("DELETE /api/v1/analysis-sessions/{sessionId}", HandleDeleteAnalysisSession(store))
	mux.HandleFunc(
		"PATCH /api/v1/analysis-sessions/{sessionId}/tables/{tableName}",
		HandleUpdateTableRowCount(store),
	)
	mux.HandleFunc("GET /", HandleUI)
	return mux
}