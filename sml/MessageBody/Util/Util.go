package Util

import "encoding/binary"

func GetNoForHex(number byte) int {
	return int(
		binary.BigEndian.Uint16(
			ExtendToMatchInt(
				[]byte{number}, 2)))
}

func ExtendToMatchInt(hex []byte, noBytes int) []byte {
	for len(hex) < noBytes { //uint64 = 8bytes
		hex = append([]byte{0}, hex...)
	}
	return hex
}
