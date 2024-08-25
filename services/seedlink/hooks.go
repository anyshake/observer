package seedlink

import (
	"github.com/anyshake/observer/utils/logger"
	"github.com/bclswl0827/slgo/handlers"
)

type hooks struct {
	serviceName string
}

func (h *hooks) OnConnection(client *handlers.SeedLinkClient) {
	logger.GetLogger(h.serviceName).Infof("client %v connected", client.RemoteAddr())
}

func (h *hooks) OnData(client *handlers.SeedLinkClient, data []byte) {}

func (h *hooks) OnClose(client *handlers.SeedLinkClient) {
	logger.GetLogger(h.serviceName).Infof("client %v disconnected", client.RemoteAddr())
}

func (h *hooks) OnCommand(client *handlers.SeedLinkClient, command []string) {
	logger.GetLogger(h.serviceName).Infof("client %v issued command %v", client.RemoteAddr(), command)
}
