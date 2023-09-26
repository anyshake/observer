package geophone

func (g *Geophone) getSize(packetLen, checksumLen int) int {
	// channelLen*packetLen*int32 + checksumLen + 1
	return checksumLen*packetLen*4 + checksumLen + 1
}
