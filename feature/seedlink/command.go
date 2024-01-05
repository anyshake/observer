package seedlink

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/anyshake/observer/driver/seedlink"
	"github.com/anyshake/observer/feature"
)

func (s *SeedLink) handleCommand(options *feature.FeatureOptions, slGlobal *seedlink.SeedLinkGlobal, slClient *seedlink.SeedLinkClient, conn net.Conn) {
	// Builtin seedlink list of SeedLink Protocol
	SeedLinkCommands := map[string]seedlink.SeedLinkCommand{
		"END":          {HasExtraArgs: false, SeedLinkCommandCallback: &seedlink.END{}},
		"DATA":         {HasExtraArgs: true, SeedLinkCommandCallback: &seedlink.DATA{}},
		"TIME":         {HasExtraArgs: true, SeedLinkCommandCallback: &seedlink.TIME{}},
		"INFO":         {HasExtraArgs: true, SeedLinkCommandCallback: &seedlink.INFO{}},
		"HELLO":        {HasExtraArgs: false, SeedLinkCommandCallback: &seedlink.HELLO{}},
		"SELECT":       {HasExtraArgs: true, SeedLinkCommandCallback: &seedlink.SELECT{}},
		"STATION":      {HasExtraArgs: true, SeedLinkCommandCallback: &seedlink.STATION{}},
		"CAPABILITIES": {HasExtraArgs: true, SeedLinkCommandCallback: &seedlink.CAPABILITIES{}},
	}

	// Dispatch connection events
	s.OnStart(options, "user connected ", conn.RemoteAddr().String())
	defer s.OnStop(options, "user disconnected ", conn.RemoteAddr().String())

	// Create a new reader
	reader := bufio.NewReader(conn)
	defer conn.Close()

	for {
		// Read client message
		clientMessage, err := reader.ReadString('\r')
		if err != nil {
			return
		} else {
			// Remove '\n' & '\r' from message and convert to uppercase
			trimmedMessage := strings.ReplaceAll(clientMessage, "\n", "")
			clientMessage = strings.ToUpper(strings.TrimSuffix(trimmedMessage, "\r"))
		}

		// Ignore empty message
		if len(clientMessage) == 0 {
			continue
		}

		// Disconnect if BYE received
		if clientMessage == "BYE" {
			return
		}

		// Exit from stream mode
		if clientMessage != "END" {
			slClient.StreamMode = false
		}

		// Check if command is whitelisted
		var (
			isCommandValid = true
			argumentList   = strings.Split(clientMessage, " ")
			mainArgument   = argumentList[0]
		)

		// Get command details from command list
		cmd, ok := SeedLinkCommands[mainArgument]
		if !ok {
			isCommandValid = false
		}

		// Send error if command is invalid
		if !isCommandValid {
			conn.Write([]byte(seedlink.RES_ERR))
			s.OnError(options, fmt.Errorf("RECV ERR: %s <%s>", conn.RemoteAddr().String(), clientMessage))
		} else {
			s.OnReady(options, fmt.Sprintf("RECV OK: %s <%s>", conn.RemoteAddr().String(), clientMessage))
		}

		// Check for extra arguments
		if isCommandValid && cmd.HasExtraArgs {
			if len(argumentList) == 0 {
				conn.Write([]byte(seedlink.RES_ERR))
			} else {
				err = cmd.Callback(slGlobal, slClient, options, s.handleMessage, conn, argumentList[1:]...)
				if err != nil {
					cmd.Fallback(slGlobal, slClient, options, conn, argumentList[1:]...)
				}
			}
		} else if isCommandValid {
			err = cmd.Callback(slGlobal, slClient, options, s.handleMessage, conn)
			if err != nil {
				cmd.Fallback(slGlobal, slClient, options, conn)
			}
		}

		// Clear selected channels
		if clientMessage == "END" {
			slClient.Channels = []string{}
		}
	}
}
