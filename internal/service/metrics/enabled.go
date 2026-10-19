package metrics

import "context"

func (s *MetricsServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&metricsConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
