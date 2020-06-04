package MessageBody

import (
	"encoding/binary"
	"goSML/sml/MessageBody/CurrentValueStrategy"
	"goSML/sml/MessageBody/DeviceIdentificationValueStrategy"
	"goSML/sml/MessageBody/ManufacturerValueStrategy"
	"goSML/sml/MessageBody/TotalValueStrategy"
	"goSML/sml/MessageBody/Util"
	"reflect"
)

type MessageBody struct {
	TransactionId uint64
	TotalValue    uint64
	CurrentValue  uint64
	Message       []byte

	valueEntryStrategies []ValueEntryStrategy
	currentByteInArray   int
}

const listItemByte = 112 // x70 -> int

func New(currentByteInArray int, message []byte) (*MessageBody, int, error) {
	messageBody := MessageBody{currentByteInArray: currentByteInArray, Message: message}
	initializeStrategies(&messageBody)

	messageBody.currentByteInArray++ // TODO check first byte
	messageBody.setTransactionId()

	messageBody.currentByteInArray += 4  // groupNo + abortOnError
	messageBody.currentByteInArray += 1  // messageBody (List with 2 items)
	messageBody.currentByteInArray += 3  // TODO Message Type (getListResponse) -> use Implementation to detect correct type
	messageBody.currentByteInArray += 1  // list with 7 entries
	messageBody.currentByteInArray += 1  // clientId
	messageBody.currentByteInArray += 11 // TODO serverId
	messageBody.currentByteInArray += 7  // TODO reqFileId
	messageBody.currentByteInArray += 8  // secList (1 + 2 + 5)

	valueCount, newPosition := getListCount(messageBody.Message, messageBody.currentByteInArray)
	messageBody.currentByteInArray = newPosition

	for listEntry := 0; listEntry < valueCount; listEntry++ {
		for _, strategy := range messageBody.valueEntryStrategies {
			if strategy.Responsible(messageBody.currentByteInArray, messageBody.Message) {
				if reflect.TypeOf(strategy) == reflect.TypeOf(TotalValueStrategy.TotalValueStrategy{}) {
					messageBody.TotalValue = strategy.ExtractIntValue(messageBody.currentByteInArray, messageBody.Message)
				}
				if reflect.TypeOf(strategy) == reflect.TypeOf(CurrentValueStrategy.CurrentValueStrategy{}) {
					messageBody.CurrentValue = strategy.ExtractIntValue(messageBody.currentByteInArray, messageBody.Message)
				}
				break
			}
		}

		messageBody.currentByteInArray++ // Skip list with 7 entries
		messageBody.moveBytePointerToNextListEntry()
	}

	return &messageBody, messageBody.currentByteInArray, nil
}

func (messageBody *MessageBody) moveBytePointerToNextListEntry() {
	for i := 0; i < 7; i++ {
		if messageBody.Message[messageBody.currentByteInArray] == 0x83 {
			messageBody.currentByteInArray += 50
			//Public key needs 2 bytes to acquire size => _3 02 __ => 50bytes
			continue
		}
		messageBody.currentByteInArray += Util.GetNoForHex(messageBody.Message[messageBody.currentByteInArray] << 4 >> 4) // shifting bits to get only second hex digit (f. e. B5 => 05)
	}
}

func (messageBody *MessageBody) setTransactionId() {
	transactionIdLength := int(messageBody.Message[messageBody.currentByteInArray])
	messageBody.currentByteInArray++

	transactionIdHex := messageBody.Message[messageBody.currentByteInArray : messageBody.currentByteInArray+transactionIdLength-1]
	transactionIdHex = extendToMatchInt(transactionIdHex, 8)
	transactionId := binary.BigEndian.Uint64(transactionIdHex)

	messageBody.currentByteInArray += transactionIdLength - 1
	messageBody.TransactionId = transactionId
}

func extendToMatchInt(transactionIdHex []byte, noBytes int) []byte {
	for len(transactionIdHex) < noBytes { //uint64 = 8bytes
		transactionIdHex = append([]byte{0}, transactionIdHex...)
	}
	return transactionIdHex
}

func getListCount(message []byte, currentPosition int) (int, int) {
	hexValueCount := message[currentPosition] - listItemByte
	valueCount := Util.GetNoForHex(hexValueCount)
	currentPosition++
	return valueCount, currentPosition
}

func initializeStrategies(messageBody *MessageBody) {
	messageBody.valueEntryStrategies = append(messageBody.valueEntryStrategies, ManufacturerValueStrategy.New())
	messageBody.valueEntryStrategies = append(messageBody.valueEntryStrategies, DeviceIdentificationValueStrategy.New())
	messageBody.valueEntryStrategies = append(messageBody.valueEntryStrategies, TotalValueStrategy.New())
	messageBody.valueEntryStrategies = append(messageBody.valueEntryStrategies, CurrentValueStrategy.New())
}
