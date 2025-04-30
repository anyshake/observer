package helicorder

func (s *HelicorderServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.dataProvider.queryCache.Clear()
	s.status.SetStoppedAt(s.timeSource.Get())
	s.status.SetIsRunning(false)

	s.cancelFn()
	s.wg.Wait()

	return nil
}
