package metrics

func (s *MetricsServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.SetStoppedAt(s.timeSource.Get())
	s.status.SetIsRunning(false)

	_ = s.oltpTracerProvider.Shutdown(s.oltpCtx)
	s.cancelFn()
	s.wg.Wait()

	return nil
}
