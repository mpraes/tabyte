package httpapi

import "net/http"

func HandleInfo(persistence bool, aiInsights bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]any{
			"app":                "tabyte",
			"version":            "0.0.1",
			"mode":               "local",
			"bind":               "127.0.0.1:8787",
			"persistence":        persistence,
			"external_required":  false,
			"ai_insights":        aiInsights,
			"engines":            []string{"sqlserver", "postgres"},
		})
	}
}
