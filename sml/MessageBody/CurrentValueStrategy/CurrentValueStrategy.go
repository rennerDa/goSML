package CurrentValueStrategy

import (
	"bytes"
	"encoding/binary"
	"goSML/sml/MessageBody/Util"
)

type CurrentValueStrategy struct {
}

func New() CurrentValueStrategy {
	return CurrentValueStrategy{}
}

var OBIST_T_ID = []byte{0x01, 0x00, 0x10, 0x07, 0x00, 0xff}

func (m CurrentValueStrategy) Responsible(currentByteInArray int, message []byte) bool {
	currentByteInArray++ // first byte is to indicate list length 77
	currentByteInArray++ // next byte is to indicate length of obist-t-id (07 8181c78203ff -> 07)
	return bytes.Equal(message[currentByteInArray:currentByteInArray+len(OBIST_T_ID)], OBIST_T_ID)
}

func (m CurrentValueStrategy) ExtractStringValue(currentByteInArray int, message []byte) string {
	// No string value
	return ""
}

func (m CurrentValueStrategy) ExtractIntValue(currentByteInArray int, message []byte) uint64 {
	currentByteInArray++                      // first byte is to indicate list length 77
	currentByteInArray += len(OBIST_T_ID) + 1 // next byte is to indicate length of obist-t-id (07 8181c78203ff -> 07)
	currentByteInArray++                      // status
	currentByteInArray++                      // valTime
	currentByteInArray += 2                   // unit
	currentByteInArray += 2                   // scaler

	byteCount := Util.GetNoForHex(message[currentByteInArray] << 4 >> 4)
	currentByteInArray++
	totalValueHex := message[currentByteInArray : currentByteInArray+byteCount-1]
	totalValueHex = Util.ExtendToMatchInt(totalValueHex, 8)
	return binary.BigEndian.Uint64(totalValueHex)
}
