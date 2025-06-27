package frp_client

import (
	"context"
	"sync"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/fatedier/frp/client"
)

const ID = "service_frp_client"

type FrpClientServiceImpl struct {
	mu     sync.Mutex
	status service.Status

	wg              sync.WaitGroup
	ctx             context.Context
	cancelFn        context.CancelFunc
	timeSource      *timesource.Source
	actionHandler   *action.Handler
	localServerAddr string

	// Server address to FRP server
	serverAddr string
	// Server port to FRP server
	serverPort int64
	// Disable custom TLS first byte to avoid traffic being sniffed.
	// See https://github.com/fatedier/frp/issues/3193#issuecomment-1332547301
	disableCustomTLSFirstByte bool
	// Specifies the path of the secret key file that client will load.
	keyFile string
	// Specifies the path of the cert file that client will load.
	certFile string
	// Specifies the path of the trusted ca file that will load.
	trustedCaFile string
	// specifies the custom server name of tls certificate.
	// By default, server name if same to ServerAddr.
	tlsServerName string
	// The token for authentication when connecting to the server
	authToken string
	// A prefix in proxy name to distinguish different users
	user string
	// Count of connection pool
	connPoolCount int64
	// TCP multiplexing availablity, should be same as FRP server config
	tcpMux bool
	// Enable TLS encryption for connection to FRP server
	tlsEnable bool
	// Transport protocol when connecting to FRP server
	transportProtocol string
	// Custom proxy name for AnyShake Observer web interface
	proxyName string
	// Whether to encrypt traffic when proxing
	useEncryption bool
	// Whether to compress traffic when proxing
	useCompression bool
	// Specify which proxy type to use in ProxyBaseConfig (tcp or http),
	// if true, proxyConfigurer will be v1.HTTPProxyConfig,
	// else v1.TCPProxyConfig will be applied.
	useDomainAccess bool
	// Remote outbound port, this field is available only when useDomainAccess is disabled
	remoteOutboundPort int64
	// Subdomain to bind for this proxy, this field is available only when useDomainAccess is enabled
	subdomain string
	// List of custom domains, this field is available only when useDomainAccess is enabled
	customDomains []string

	clientObj *client.Service
}
