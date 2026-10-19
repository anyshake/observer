package frp_client

import (
	"context"
	"fmt"
)

func (s *FrpClientServiceImpl) Init(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return fmt.Errorf("service %s cannot be initialized while running", ID)
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(ctx, s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	serverAddr, err := (&frpClientConfigServerAddrImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.serverAddr = serverAddr.(string)

	serverPort, err := (&frpClientConfigServerPortImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.serverPort = serverPort.(int64)

	disableCustomTLSFirstByte, err := (&frpClientConfigDisableCustomTLSFirstByteImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.disableCustomTLSFirstByte = disableCustomTLSFirstByte.(bool)

	keyFile, err := (&frpClientConfigKeyFileImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.keyFile = keyFile.(string)

	certFile, err := (&frpClientConfigCertFileImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.certFile = certFile.(string)

	trustedCaFile, err := (&frpClientConfigTrustedCaFileImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.trustedCaFile = trustedCaFile.(string)

	tlsServerName, err := (&frpClientConfigTlsServerNameImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.tlsServerName = tlsServerName.(string)

	user, err := (&frpClientConfigUserImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.user = user.(string)

	authToken, err := (&frpClientConfigTokenImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.authToken = authToken.(string)

	connPoolCount, err := (&frpClientConfigPoolCountImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.connPoolCount = connPoolCount.(int64)

	tcpMux, err := (&frpClientConfigTcpMuxImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.tcpMux = tcpMux.(bool)

	tlsEnable, err := (&frpClientConfigTlsEnableImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.tlsEnable = tlsEnable.(bool)

	transportProtocol, err := (&frpClientConfigProtocolImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.transportProtocol = transportProtocol.(string)

	proxyName, err := (&frpClientConfigProxyNameImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.proxyName = proxyName.(string)

	useEncryption, err := (&frpClientConfigUseEncryptionImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.useEncryption = useEncryption.(bool)

	useCompression, err := (&frpClientConfigUseCompressionImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.useCompression = useCompression.(bool)

	useDomainAccess, err := (&frpClientConfigUseDomainAccessImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.useDomainAccess = useDomainAccess.(bool)

	remoteOutboundPort, err := (&frpClientConfigRemoteOutboundPortImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.remoteOutboundPort = remoteOutboundPort.(int64)

	subdomain, err := (&frpClientConfigSubdomainImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.subdomain = subdomain.(string)

	customDomains, err := (&frpClientConfigCustomDomainsImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.customDomains = customDomains.([]string)

	return nil
}
