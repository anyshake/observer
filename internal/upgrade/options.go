package upgrade

func (h *Helper) GetVersionCheckDomain() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.versionCheckDomain
}

func (h *Helper) SetVersionCheckDomain(domain string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.versionCheckDomain = domain
}

func (h *Helper) GetReleaseFetchUrl() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.releaseFetchUrl
}

func (h *Helper) SetReleaseFetchUrl(url string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.releaseFetchUrl = url
}
