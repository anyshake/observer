package upgrade

func (u *Helper) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.closed = true
	if u.dnsClient == nil {
		return nil
	}
	return u.dnsClient.Close()
}
