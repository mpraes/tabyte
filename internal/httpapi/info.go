package httpapi

import "net/http"

func HandleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{
		"app":         "tabyte",
		"version":     "0.0.1",
		"mode":        "local",
		"bind":        "127.0.0.1:8787",
		"persistence": false,
		"engines":     []string{"sqlserver", "postgres"},
	})
}