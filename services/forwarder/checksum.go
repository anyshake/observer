package forwarder

import "unsafe"

func (a *ForwarderService) getChecksum(arr []int32) (checksum uint8) {
	for _, data := range arr {
		bytes := (*[4]byte)(unsafe.Pointer(&data))[:]
		for j := 0; j < int(unsafe.Sizeof(int32(0))); j++ {
			checksum ^= bytes[j]
		}
	}

	return checksum
}
