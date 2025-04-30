package metrics

func (s *MetricsServiceImpl) IsEnabled() bool {
	enable, err := (&metricsConfigEnabledImpl{}).Get(s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
