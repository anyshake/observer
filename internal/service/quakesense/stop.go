package quakesense

func (s *QuakeSenseServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.SetStoppedAt(s.timeSource.Get())
	s.status.SetIsRunning(false)

	_ = s.hardwareDev.Unsubscribe(ID)
	if s.mqttClient != nil {
		s.mqttClient.Disconnect(100)
	}
	s.channelBuffer.Reset()
	s.prevSamplerate = 0
	s.filterKernel = nil
	s.cancelFn()
	s.wg.Wait()

	return nil
}
