package timesync

func (s *TimeSyncServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.SetStoppedAt(s.timeSource.Get())
	s.status.SetIsRunning(false)

	s.cancelFn()
	s.wg.Wait()

	return nil
}
