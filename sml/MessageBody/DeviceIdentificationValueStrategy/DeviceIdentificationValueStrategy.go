package DeviceIdentificationValueStrategy

import (
	"bytes"
	"encoding/hex"
	"goSML/sml/MessageBody/Util"
)

type DeviceIdentificationValueStrategy struct {
}

func New() DeviceIdentificationValueStrategy {
	return DeviceIdentificationValueStrategy{}
}

var OBIST_T_ID = []byte{0x01, 0x00, 0x00, 0x00, 0x09, 0xff}

func (m DeviceIdentificationValueStrategy) Responsible(currentByteInArray int, message []byte) bool {
	currentByteInArray++ // first byte is to indicate list length 77
	currentByteInArray++ // next byte is to indicate length of obist-t-id (07 8181c78203ff -> 07)
	return bytes.Equal(message[currentByteInArray:currentByteInArray+len(OBIST_T_ID)], OBIST_T_ID)
}

func (m DeviceIdentificationValueStrategy) ExtractStringValue(currentByteInArray int, message []byte) string {
	currentByteInArray++                      // first byte is to indicate list length 77
	currentByteInArray += len(OBIST_T_ID) + 1 // next byte is to indicate length of obist-t-id (07 8181c78203ff -> 07)
	currentByteInArray++                      // status
	currentByteInArray++                      // valTime
	currentByteInArray++                      // unit
	currentByteInArray++                      // scaler
	byteCount := Util.GetNoForHex(message[currentByteInArray])
	currentByteInArray++
	return hex.EncodeToString(message[currentByteInArray : currentByteInArray+byteCount-1])
}

func (m DeviceIdentificationValueStrategy) ExtractIntValue(currentByteInArray int, message []byte) uint64 {
	//No int value
	return 0
}
