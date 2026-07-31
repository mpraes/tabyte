package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mpraes/tabyte/internal/application"
)

type putSettingRequest struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	ValueType string `json:"value_type"`
}

func HandleGetSettings(settings application.SettingsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if settings == nil {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "persistence is not enabled")
			return
		}
		list, err := settings.ListSettings()
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		items := make([]map[string]any, 0, len(list))
		for _, s := range list {
			items = append(items, map[string]any{
				"key":        s.Key,
				"value":      s.Value,
				"value_type": s.ValueType,
			})
		}
		WriteJSON(w, http.StatusOK, map[string]any{"settings": items})
	}
}

func HandlePutSettings(settings application.SettingsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if settings == nil {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "persistence is not enabled")
			return
		}
		var req putSettingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON body")
			return
		}
		key := strings.TrimSpace(req.Key)
		if key == "" {
			WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "key is required")
			return
		}
		valueType := strings.TrimSpace(req.ValueType)
		if valueType == "" {
			valueType = "string"
		}
		if err := settings.UpsertSetting(key, req.Value, valueType); err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		WriteJSON(w, http.StatusOK, map[string]any{
			"key":        key,
			"value":      req.Value,
			"value_type": valueType,
		})
	}
}
