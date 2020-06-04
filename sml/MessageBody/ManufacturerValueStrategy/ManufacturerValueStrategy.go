package ManufacturerValueStrategy

import (
	"bytes"
	"encoding/hex"
	"goSML/sml/MessageBody/Util"
)

type ManufacturerValueStrategy struct {
}

func New() ManufacturerValueStrategy {
	return ManufacturerValueStrategy{}
}

var OBIST_T_ID = []byte{0x81, 0x81, 0xc7, 0x82, 0x03, 0xff}

func (m ManufacturerValueStrategy) Responsible(currentByteInArray int, message []byte) bool {
	currentByteInArray++ // first byte is to indicate list length 77
	currentByteInArray++ // next byte is to indicate length of obist-t-id (07 8181c78203ff -> 07)
	return bytes.Equal(message[currentByteInArray:currentByteInArray+len(OBIST_T_ID)], OBIST_T_ID)
}

func (m ManufacturerValueStrategy) ExtractStringValue(currentByteInArray int, message []byte) string {
	currentByteInArray++                      // first byte is to indicate list length 77
	currentByteInArray += len(OBIST_T_ID) + 1 // next byte is to indicate length of obist-t-id (07 8181c78203ff -> 07)
	currentByteInArray++                      // status
	currentByteInArray++                      // valTime
	currentByteInArray++                      // unit
	currentByteInArray++                      // scaler
	byteCount := Util.GetNoForHex(message[currentByteInArray])
	currentByteInArray++
	return hex.EncodeToString(message[currentByteInArray : currentByteInArray+byteCount])
}

func (m ManufacturerValueStrategy) ExtractIntValue(currentByteInArray int, message []byte) uint64 {
	//No int value
	return 0
}
