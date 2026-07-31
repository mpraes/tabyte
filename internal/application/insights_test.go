package application_test

import (
	"testing"

	"github.com/mpraes/tabyte/internal/application"
	"github.com/mpraes/tabyte/internal/domain"
)

func TestListInsightsDisabled(t *testing.T) {
	store := application.NewSessionStore(nil)
	session := domain.AnalysisSession{
		ID:     "as_1",
		Engine: domain.EnginePostgres,
		Status: "created",
	}
	store.Save(session)

	insights, enabled, err := application.ListInsights(store, "as_1", application.DisabledInsightProvider{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if enabled {
		t.Fatal("expected disabled")
	}
	if len(insights) != 0 {
		t.Fatalf("want empty insights, got %#v", insights)
	}

	_, _, err = application.ListInsights(store, "missing", application.DisabledInsightProvider{})
	if err != application.ErrSessionNotFound {
		t.Fatalf("want ErrSessionNotFound, got %v", err)
	}
}
