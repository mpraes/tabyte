package application

import "github.com/mpraes/tabyte/internal/domain"

// InsightProvider generates optional textual insights from an analysis result.
// Implementations must not alter estimates; they only annotate.
type InsightProvider interface {
	Enabled() bool
	Insights(session domain.AnalysisSession) ([]domain.Insight, error)
}

// DisabledInsightProvider is the default: core stays fully local with no AI calls.
type DisabledInsightProvider struct{}

func (DisabledInsightProvider) Enabled() bool { return false }

func (DisabledInsightProvider) Insights(domain.AnalysisSession) ([]domain.Insight, error) {
	return []domain.Insight{}, nil
}

func ListInsights(store *SessionStore, sessionID string, provider InsightProvider) ([]domain.Insight, bool, error) {
	session, ok := store.Get(sessionID)
	if !ok {
		return nil, false, ErrSessionNotFound
	}
	if provider == nil {
		provider = DisabledInsightProvider{}
	}
	insights, err := provider.Insights(session)
	if err != nil {
		return nil, provider.Enabled(), err
	}
	if insights == nil {
		insights = []domain.Insight{}
	}
	return insights, provider.Enabled(), nil
}
