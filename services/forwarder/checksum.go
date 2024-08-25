package forwarder

import "unsafe"

func (a *ForwarderService) getChecksum(arr []int32) uint8 {
	checksum := uint8(0)

	for i := 0; i < len(arr); i++ {
		bytes := (*[4]byte)(unsafe.Pointer(&arr[i]))[:]

		for j := 0; j < int(unsafe.Sizeof(int32(0))); j++ {
			checksum ^= bytes[j]
		}
	}

	return checksum
}
