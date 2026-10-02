package tailscale

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os"

	"tailscale.com/ipn"
	"tailscale.com/tsnet"
)

type nodeConfig struct {
	Hostname   string
	AuthKey    string
	ControlURL string
	Store      ipn.StateStore
	UserLogf   func(format string, args ...any)
}

type embeddedNode interface {
	Start() error
	Up(context.Context) error
	Listen(network, addr string) (net.Listener, error)
	ListenPacket(network, addr string) (net.PacketConn, error)
	TailscaleIPs() (netip.Addr, netip.Addr)
	Close() error
}

type nodeFactory func(nodeConfig) embeddedNode

type tsnetNode struct {
	server     *tsnet.Server
	runtimeDir string
}

func (n *tsnetNode) Start() error {
	runtimeDir, err := os.MkdirTemp("", "anyshake-observer-tsnet-")
	if err != nil {
		return fmt.Errorf("failed to create tsnet runtime directory: %w", err)
	}
	n.runtimeDir = runtimeDir
	n.server.Dir = runtimeDir

	if err := n.server.Start(); err != nil {
		return errors.Join(err, n.cleanupRuntimeDir())
	}
	return nil
}

func (n *tsnetNode) Up(ctx context.Context) error {
	_, err := n.server.Up(ctx)
	return err
}

func (n *tsnetNode) Listen(network, addr string) (net.Listener, error) {
	return n.server.Listen(network, addr)
}

func (n *tsnetNode) ListenPacket(network, addr string) (net.PacketConn, error) {
	return n.server.ListenPacket(network, addr)
}

func (n *tsnetNode) TailscaleIPs() (netip.Addr, netip.Addr) {
	return n.server.TailscaleIPs()
}

func (n *tsnetNode) Close() error {
	return errors.Join(n.server.Close(), n.cleanupRuntimeDir())
}

func (n *tsnetNode) cleanupRuntimeDir() error {
	if n.runtimeDir == "" {
		return nil
	}
	runtimeDir := n.runtimeDir
	n.runtimeDir = ""
	if err := os.RemoveAll(runtimeDir); err != nil {
		return fmt.Errorf("failed to remove tsnet runtime directory: %w", err)
	}
	return nil
}

func newTSNetNode(config nodeConfig) embeddedNode {
	return &tsnetNode{server: &tsnet.Server{
		Hostname:   config.Hostname,
		AuthKey:    config.AuthKey,
		ControlURL: config.ControlURL,
		Store:      config.Store,
		UserLogf:   config.UserLogf,
		Logf:       func(string, ...any) {},
	}}
}
