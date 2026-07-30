package httpapi

import (
	"net/http"

	webui "github.com/mpraes/tabyte/web"
)

func HandleUI(w http.ResponseWriter, r *http.Request) {
	data, err := webui.Assets.ReadFile("index.html")
	if err != nil {
		http.Error(w, "ui unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}