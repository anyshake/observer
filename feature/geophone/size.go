package geophone

func (g *Geophone) getPacketSize(packetLen, checksumLen int) int {
	// channelLen*packetLen*int32 + checksumLen + 1
	return checksumLen*packetLen*4 + checksumLen + 1
}
