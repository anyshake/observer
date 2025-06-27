package frp_client

import "fmt"

func (s *FrpClientServiceImpl) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	serverAddr, err := (&frpClientConfigServerAddrImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.serverAddr = serverAddr.(string)

	serverPort, err := (&frpClientConfigServerPortImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.serverPort = serverPort.(int64)

	disableCustomTLSFirstByte, err := (&frpClientConfigDisableCustomTLSFirstByteImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.disableCustomTLSFirstByte = disableCustomTLSFirstByte.(bool)

	keyFile, err := (&frpClientConfigKeyFileImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.keyFile = keyFile.(string)

	certFile, err := (&frpClientConfigCertFileImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.certFile = certFile.(string)

	trustedCaFile, err := (&frpClientConfigTrustedCaFileImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.trustedCaFile = trustedCaFile.(string)

	tlsServerName, err := (&frpClientConfigTlsServerNameImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.tlsServerName = tlsServerName.(string)

	user, err := (&frpClientConfigUserImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.user = user.(string)

	authToken, err := (&frpClientConfigTokenImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.authToken = authToken.(string)

	connPoolCount, err := (&frpClientConfigPoolCountImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.connPoolCount = connPoolCount.(int64)

	tcpMux, err := (&frpClientConfigTcpMuxImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.tcpMux = tcpMux.(bool)

	tlsEnable, err := (&frpClientConfigTlsEnableImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.tlsEnable = tlsEnable.(bool)

	transportProtocol, err := (&frpClientConfigProtocolImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.transportProtocol = transportProtocol.(string)

	proxyName, err := (&frpClientConfigProxyNameImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.proxyName = proxyName.(string)

	useEncryption, err := (&frpClientConfigUseEncryptionImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.useEncryption = useEncryption.(bool)

	useCompression, err := (&frpClientConfigUseCompressionImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.useCompression = useCompression.(bool)

	useDomainAccess, err := (&frpClientConfigUseDomainAccessImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.useDomainAccess = useDomainAccess.(bool)

	remoteOutboundPort, err := (&frpClientConfigRemoteOutboundPortImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.remoteOutboundPort = remoteOutboundPort.(int64)

	subdomain, err := (&frpClientConfigSubdomainImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.subdomain = subdomain.(string)

	customDomains, err := (&frpClientConfigCustomDomainsImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.customDomains = customDomains.([]string)

	return nil
}
