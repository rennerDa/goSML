package Message

import (
	"bytes"
	"errors"
	MessageBody "goSML/sml/MessageBody"
)

type Message struct {
	currentByteInArray int
	message            []byte

	//StartMessage string
	MessageBody MessageBody.MessageBody
	//EndMessage   string
}

var START_ESCAPE_SEQUENCE = []byte{0x1b, 0x1b, 0x1b, 0x1b}
var START_TRANSMISSION_SEQUENCE = []byte{0x01, 0x01, 0x01, 0x01}

func New(currentByteInArray int, rawMessage []byte) (*Message, int, error) {
	message := Message{currentByteInArray: currentByteInArray, message: rawMessage}

	if !message.isStartOfMessageValid() {
		return nil, message.currentByteInArray, errors.New("Start of Message is not valid!")
	}

	//skip start message
	message.currentByteInArray = 51

	messageBody, currentByte, err := MessageBody.New(message.currentByteInArray, message.message)
	message.currentByteInArray = currentByte
	if err != nil {
		return nil, message.currentByteInArray, errors.New("Error creating message body!")
	}
	message.MessageBody = *messageBody

	return &message, message.currentByteInArray, nil
}

func (message Message) isStartOfMessageValid() bool {
	startOfEscapeSequence := message.currentByteInArray
	endOfEscapeSequence := startOfEscapeSequence + len(START_ESCAPE_SEQUENCE)
	endOfTransmissionSequence := endOfEscapeSequence + len(START_TRANSMISSION_SEQUENCE)
	message.currentByteInArray = endOfTransmissionSequence

	if !bytes.Equal(message.message[startOfEscapeSequence:endOfEscapeSequence], START_ESCAPE_SEQUENCE) {
		return false
	}
	if !bytes.Equal(message.message[endOfEscapeSequence:endOfTransmissionSequence], START_TRANSMISSION_SEQUENCE) {
		return false
	}
	return true
}
