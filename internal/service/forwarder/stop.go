package forwarder

func (s *ForwarderServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.SetStoppedAt(s.timeSource.Get())
	s.status.SetIsRunning(false)

	_ = s.hardwareDev.Unsubscribe(ID)
	if s.listener != nil {
		_ = s.listener.Close()
	}
	s.cancelFn()
	s.wg.Wait()

	return nil
}
