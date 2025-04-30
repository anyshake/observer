package archiver

func (s *ArchiverServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.SetStoppedAt(s.timeSource.Get())
	s.status.SetIsRunning(false)

	_ = s.hardwareDev.Unsubscribe(ID)
	s.cancelFn()
	s.wg.Wait()

	return nil
}
