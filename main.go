package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/influxdata/influxdb-client-go"
	"github.com/tarm/serial"
	"goSML/sml/Message"
	"log"
	"strconv"
	"time"
)

const (
	DB       = ""
	username = ""
	password = ""
)

const listItemByte = 112   // x70 -> int
const messageTypeByte = 96 // x60 -> int

func main() {
	log.Println("Starting smart meter watch!")
	s := openConnection()

	//startSequence := "1b1b1b1b0101010176"
	startSequence := []byte{0x1b, 0x1b, 0x1b, 0x1b, 0x01, 0x01, 0x01, 0x01, 0x76}
	//endSequence := "1b1b1b1b1a"
	endSequence := []byte{0x1b, 0x1b, 0x1b, 0x1b, 0x1a}
	falseSequence := []byte{0xb1, 0xb1, 0xb1}

	buf := make([]byte, 128)
	cache := []byte{}
	for {
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		cache = append(cache, buf[:n]...)
		packetFound := GetBytesInBetween(cache, startSequence, endSequence)
		if len(packetFound) > 0 {
			//log.Print("found packet")
			//log.Print(hex.EncodeToString(packetFound))
			index := bytes.Index(cache, endSequence)
			cache = cache[index+len(endSequence):]
			go decode(packetFound)
		} else if bytes.Index(cache, falseSequence) != -1 || len(cache) > 5000 {
			log.Println("WARN - Wrong byte-order detected or cache is way to big. Reset!")
			log.Println(len(cache))
			cache = []byte{}
			err := s.Close()
			if err != nil {
				log.Println("ERROR - Error closing open connection to serial port")
				log.Fatal(err)
			}
			s = openConnection()
		}
	}
}

func openConnection() *serial.Port {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Println("ERROR - Could not connect to serial port!")
		log.Fatal(err)
	}
	return s
}

func decode(byteMessage []byte) {
	message, _, err := Message.New(0, byteMessage)

	if err != nil || message == nil {
		log.Fatal("ERROR - ")
	}

	go commitToDb(*message)
}

func commitToDb(message Message.Message) {
	client := influxdb2.NewClient("http://localhost:8086", "")
	writeApi := client.WriteApiBlocking("", "renner_metrics")

	year, month, day := time.Now().Date()
	weekday := time.Now().Weekday()

	tags := map[string]string{
		"counter":    "1",
		"year":       strconv.Itoa(year),
		"month":      strconv.Itoa(int(month)),
		"dayOfMonth": strconv.Itoa(day),
		"weekday":    strconv.Itoa(int(weekday)),
	}
	fields := map[string]interface{}{
		"messageBodyTransactionId": int(message.MessageBody.TransactionId),
		"totalValue":               int(message.MessageBody.TotalValue),
		"currentValue":             int(message.MessageBody.CurrentValue),
		"binary":                   hex.EncodeToString(message.MessageBody.Message)}

	p := influxdb2.NewPoint("power_consumption",
		tags,
		fields,
		time.Now())
	// Write data
	err := writeApi.WritePoint(context.Background(), p)
	if err != nil {
		client.Close()
		fmt.Printf("Write error: %s\n", err.Error())
	}
	client.Close()
}

func GetBytesInBetween(cache []byte, start []byte, end []byte) (result []byte) {
	s := bytes.Index(cache, start)
	if s == -1 {
		//log.Println("No start found")
		return
	}
	e := bytes.Index(cache, end)
	if e == -1 {
		//log.Println("No end found")
		return
	}
	numberOfCRCAfterEscapeSequence := 3 // 3 bytes -> fillbyte + 2 CRC bytes
	endIndexOfMessage := e + len(end) + numberOfCRCAfterEscapeSequence
	if endIndexOfMessage > len(cache) {
		//log.Println("End is not finished yet")
		return
	}
	return cache[s:endIndexOfMessage]
}
