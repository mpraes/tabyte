package httpapi

import (
	"github.com/mpraes/tabyte/internal/application"
	"github.com/mpraes/tabyte/internal/domain"
)

func humanTotalBytes(session domain.AnalysisSession) map[string]any {
	out := map[string]any{}
	if session.EstimatedTotalBytes != nil {
		out["estimated_total_human"] = application.FormatBytes(*session.EstimatedTotalBytes)
	}
	if session.ProjectedTotalBytes != nil {
		out["projected_total_human"] = application.FormatBytes(*session.ProjectedTotalBytes)
	}
	return out
}