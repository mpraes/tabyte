package httpapi

import (
	"net/http"

	"github.com/mpraes/tabyte/internal/application"
)

func NewMux(store *application.SessionStore, settings application.SettingsRepository, persistence bool) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", HandleHealth)
	mux.HandleFunc("GET /api/v1/info", HandleInfo(persistence))
	mux.HandleFunc("POST /api/v1/analysis-sessions", HandleCreateAnalysisSession(store))
	mux.HandleFunc("GET /api/v1/analysis-sessions/{sessionId}", HandleGetAnalysisSession(store))
	mux.HandleFunc("GET /api/v1/analysis-sessions", HandleListAnalysisSessions(store))
	mux.HandleFunc("DELETE /api/v1/analysis-sessions/{sessionId}", HandleDeleteAnalysisSession(store))
	mux.HandleFunc("PATCH /api/v1/analysis-sessions/{sessionId}", HandleReprocessAnalysisSession(store))
	mux.HandleFunc(
		"PATCH /api/v1/analysis-sessions/{sessionId}/tables/{tableName}",
		HandleUpdateTableRowCount(store),
	)
	mux.HandleFunc("GET /", HandleUI)
	mux.HandleFunc(
		"PATCH /api/v1/analysis-sessions/{sessionId}/tables/{tableName}/growth",
		HandleUpdateTableGrowth(store),
	)
	mux.HandleFunc(
		"GET /api/v1/analysis-sessions/{sessionId}/export",
		HandleExportAnalysisSession(store),
	)
	mux.HandleFunc("GET /api/v1/settings", HandleGetSettings(settings))
	mux.HandleFunc("PUT /api/v1/settings", HandlePutSettings(settings))
	return mux
}
