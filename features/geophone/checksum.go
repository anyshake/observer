package geophone

import (
	"math"
)

func GetChecksum(arr []float64) byte {
	var sum uint16 = 0

	for _, num := range arr {
		bytes := math.Float64bits(num)
		byteArr := make([]byte, 8)
		for i := 0; i < 8; i++ {
			byteArr[i] = byte(bytes >> uint(i*8))
		}

		for _, b := range byteArr {
			sum += uint16(b)
		}
	}

	return byte(sum % 256)
}

func CompareChecksum(geophone *Geophone) bool {
	var (
		checksum_vertical   = GetChecksum(geophone.Vertical[:])
		checksum_eastwest   = GetChecksum(geophone.EastWest[:])
		checksum_northsouth = GetChecksum(geophone.NorthSouth[:])
	)

	if checksum_vertical != geophone.Checksum[0] ||
		checksum_eastwest != geophone.Checksum[1] ||
		checksum_northsouth != geophone.Checksum[2] {
		return false
	}
	return true
}
